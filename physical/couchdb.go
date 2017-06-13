package physical

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/mgutz/logxi/v1"
)

// CouchDBBackend allows the management of couchdb users
type CouchDBBackend struct {
	logger log.Logger
	client *couchDBClient
}

type couchDBClient struct {
	endpoint string
	username string
	password string
	*http.Client
}

type couchDBListItem struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value struct {
		Revision string
	} `json:"value"`
}

type couchDBList struct {
	TotalRows int               `json:"total_rows"`
	Offset    int               `json:"offset"`
	Rows      []couchDBListItem `json:"rows"`
}

func (m *couchDBClient) rev(key string) (string, error) {
	req, err := http.NewRequest("HEAD", fmt.Sprintf("%s/%s", m.endpoint, key), nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(m.username, m.password)

	resp, err := m.Client.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil
	}
	etag := resp.Header.Get("Etag")
	if len(etag) < 2 {
		return "", nil
	}
	return etag[1 : len(etag)-1], nil
}

func (m *couchDBClient) put(e couchDBEntry) error {
	bs, err := json.Marshal(e)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", m.endpoint, e.ID), bytes.NewReader(bs))
	if err != nil {
		return err
	}
	req.SetBasicAuth(m.username, m.password)
	_, err = m.Client.Do(req)

	return err
}

func (m *couchDBClient) get(key string) (*Entry, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", m.endpoint, url.PathEscape(key)), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(m.username, m.password)
	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET returned %s", resp.Status)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	entry := couchDBEntry{}
	if err := json.Unmarshal(bs, &entry); err != nil {
		return nil, err
	}
	return entry.Entry, nil
}

func (m *couchDBClient) list(prefix string) ([]couchDBListItem, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/_all_docs", m.endpoint), nil)
	req.SetBasicAuth(m.username, m.password)
	values := req.URL.Query()
	values.Set("skip", "0")
	values.Set("limit", "100")
	values.Set("include_docs", "false")
	if prefix != "" {
		values.Set("startkey", fmt.Sprintf("%q", prefix))
		values.Set("endkey", fmt.Sprintf("%q", prefix+"{}"))
	}
	req.URL.RawQuery = values.Encode()

	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	results := couchDBList{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results.Rows, nil
}

func newCouchDBBackend(conf map[string]string, logger log.Logger) (Backend, error) {
	endpoint := os.Getenv("COUCHDB_ENDPOINT")
	if endpoint == "" {
		endpoint = conf["endpoint"]
	}
	if endpoint == "" {
		return nil, fmt.Errorf("missing endpoint")
	}

	username := os.Getenv("COUCHDB_USERNAME")
	if username == "" {
		username = conf["username"]
	}

	password := os.Getenv("COUCHDB_PASSWORD")
	if password == "" {
		password = conf["password"]
	}

	return &CouchDBBackend{
		client: &couchDBClient{
			endpoint: endpoint,
			username: username,
			password: password,
			Client:   &http.Client{},
		},
		logger: logger,
	}, nil
}

type couchDBEntry struct {
	Entry   *Entry `json:"entry"`
	Rev     string `json:"_rev,omitempty"`
	ID      string `json:"_id"`
	Deleted *bool  `json:"_deleted,omitempty"`
}

// Put is used to insert or update an entry
func (m *CouchDBBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"couchdb", "put"}, time.Now())

	revision, _ := m.client.rev(url.PathEscape(entry.Key))

	return m.client.put(couchDBEntry{
		Entry: entry,
		Rev:   revision,
		ID:    url.PathEscape(entry.Key),
	})
}

// Get is used to fetch an entry
func (m *CouchDBBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"couchdb", "get"}, time.Now())

	return m.client.get(key)
}

// Delete is used to permanently delete an entry
func (m *CouchDBBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"couchdb", "delete"}, time.Now())

	revision, _ := m.client.rev(url.PathEscape(key))
	deleted := true
	return m.client.put(couchDBEntry{
		ID:      url.PathEscape(key),
		Rev:     revision,
		Deleted: &deleted,
	})
}

// List is used to list all the keys under a given prefix
func (m *CouchDBBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"couchdb", "list"}, time.Now())

	items, err := m.client.list(prefix)
	if err != nil {
		return nil, err
	}

	var out []string
	seen := make(map[string]interface{})
	for _, result := range items {
		trimmed := strings.TrimPrefix(result.ID, prefix)
		sep := strings.Index(trimmed, "/")
		if sep == -1 {
			out = append(out, trimmed)
		} else {
			trimmed = trimmed[:sep+1]
			if _, ok := seen[trimmed]; !ok {
				out = append(out, trimmed)
				seen[trimmed] = struct{}{}
			}
		}
	}
	return out, nil
}
