package physical

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/strutil"
	log "github.com/mgutz/logxi/v1"

	// CockroachDB uses the Postgres SQL driver
	_ "github.com/lib/pq"
)

// CockroachDBBackend Backend is a physical backend that stores data
// within a CockroachDB database.
type CockroachDBBackend struct {
	table         string
	client        *sql.DB
	rawStatements map[string]string
	statements    map[string]*sql.Stmt
	logger        log.Logger
	permitPool    *PermitPool
}

// newCockroachDBBackend constructs a CockroachDB backend using the given
// API client, server address, credentials, and database.
func newCockroachDBBackend(conf map[string]string, logger log.Logger) (Backend, error) {
	// Get the CockroachDB credentials to perform read/write operations.
	connURL, ok := conf["connection_url"]
	if !ok || connURL == "" {
		return nil, fmt.Errorf("missing connection_url")
	}

	dbTable, ok := conf["table"]
	if !ok {
		dbTable = "vault_kv_store"
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	var err error
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("cockroachdb: max_parallel set", "max_parallel", maxParInt)
		}
	}

	// Create CockroachDB handle for the database.
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cockroachdb: %v", err)
	}

	// Create the required table if it doesn't exists.
	createQuery := "CREATE TABLE IF NOT EXISTS " + dbTable +
		" (path STRING, value BYTES, PRIMARY KEY (path))"
	if _, err := db.Exec(createQuery); err != nil {
		return nil, fmt.Errorf("failed to create mysql table: %v", err)
	}

	// Setup the backend
	c := &CockroachDBBackend{
		table:  dbTable,
		client: db,
		rawStatements: map[string]string{
			"put": "INSERT INTO " + dbTable + " VALUES($1, $2)" +
				" ON CONFLICT (path) DO " +
				" UPDATE SET (path, value) = ($1, $2)",
			"get":    "SELECT value FROM " + dbTable + " WHERE path = $1",
			"delete": "DELETE FROM " + dbTable + " WHERE path = $1",
			"list":   "SELECT path FROM " + dbTable + " WHERE path LIKE $1",
		},
		statements: make(map[string]*sql.Stmt),
		logger:     logger,
		permitPool: NewPermitPool(maxParInt),
	}

	// Prepare all the statements required
	for name, query := range c.rawStatements {
		if err := c.prepare(name, query); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// prepare is a helper to prepare a query for future execution
func (c *CockroachDBBackend) prepare(name, query string) error {
	stmt, err := c.client.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare '%s': %v", name, err)
	}
	c.statements[name] = stmt
	return nil
}

// Put is used to insert or update an entry.
func (c *CockroachDBBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"cockroachdb", "put"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	_, err := c.statements["put"].Exec(entry.Key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (c *CockroachDBBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"cockroachdb", "get"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	var result []byte
	err := c.statements["get"].QueryRow(key).Scan(&result)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ent := &Entry{
		Key:   key,
		Value: result,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (c *CockroachDBBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"cockroachdb", "delete"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	_, err := c.statements["delete"].Exec(key)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (c *CockroachDBBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"cockroachdb", "list"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	likePrefix := prefix + "%"
	rows, err := c.statements["list"].Query(likePrefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %v", err)
		}

		key = strings.TrimPrefix(key, prefix)
		if i := strings.Index(key, "/"); i == -1 {
			// Add objects only from the current 'folder'
			keys = append(keys, key)
		} else if i != -1 {
			// Add truncated 'folder' paths
			keys = strutil.AppendIfMissing(keys, string(key[:i+1]))
		}
	}

	sort.Strings(keys)
	return keys, nil
}

// Transaction is used to run multiple entries via a transaction
func (c *CockroachDBBackend) Transaction(txns []TxnEntry) error {
	defer metrics.MeasureSince([]string{"cockroachdb", "transaction"}, time.Now())
	if len(txns) == 0 {
		return nil
	}

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	return crdb.ExecuteTx(c.client, func(tx *sql.Tx) error {
		return c.transaction(tx, txns)
	})
}

func (c *CockroachDBBackend) transaction(tx *sql.Tx, txns []TxnEntry) error {
	deleteStmt, err := tx.Prepare(c.rawStatements["delete"])
	if err != nil {
		return err
	}
	putStmt, err := tx.Prepare(c.rawStatements["put"])
	if err != nil {
		return err
	}

	for _, op := range txns {
		switch op.Operation {
		case DeleteOperation:
			_, err = deleteStmt.Exec(op.Entry.Key)
		case PutOperation:
			_, err = putStmt.Exec(op.Entry.Key, op.Entry.Value)
		default:
			return fmt.Errorf("%q is not a supported transaction operation", op.Operation)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
