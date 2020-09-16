package cassandra

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/cassandra"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

func TestBackend_basic(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, hostname := cassandra.PrepareTestContainer(t, "latest")
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, hostname),
			testAccStepRole(t),
			testAccStepReadCreds(t, "test"),
		},
	})
}

func TestBackend_roleCrud(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, hostname := cassandra.PrepareTestContainer(t, "latest")
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, hostname),
			testAccStepRole(t),
			testAccStepRoleWithOptions(t),
			testAccStepReadRole(t, "test", testRole),
			testAccStepReadRole(t, "test2", testRole),
			testAccStepDeleteRole(t, "test"),
			testAccStepDeleteRole(t, "test2"),
			testAccStepReadRole(t, "test", ""),
			testAccStepReadRole(t, "test2", ""),
		},
	})
}

func testAccStepConfig(t *testing.T, hostname string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Data: map[string]interface{}{
			"hosts":            hostname,
			"username":         "cassandra",
			"password":         "cassandra",
			"protocol_version": 3,
		},
	}
}

func testAccStepRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/test",
		Data: map[string]interface{}{
			"creation_cql": testRole,
		},
	}
}

func testAccStepRoleWithOptions(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/test2",
		Data: map[string]interface{}{
			"creation_cql": testRole,
			"lease":        "30s",
			"consistency":  "All",
		},
	}
}

func testAccStepDeleteRole(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + n,
	}
}

func testAccStepReadCreds(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "creds/" + name,
		Check: func(resp *logical.Response) error {
			var d struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[WARN] Generated credentials: %v", d)

			return nil
		},
	}
}

func testAccStepReadRole(t *testing.T, name string, cql string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if cql == "" {
					return nil
				}

				return fmt.Errorf("response is nil")
			}

			var d struct {
				CreationCQL string `mapstructure:"creation_cql"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.CreationCQL != cql {
				return fmt.Errorf("bad: %#v\n%#v\n%#v\n", resp, cql, d.CreationCQL)
			}

			return nil
		},
	}
}

const testRole = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;
GRANT ALL PERMISSIONS ON ALL KEYSPACES TO {{username}};`
