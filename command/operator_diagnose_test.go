// +build !race

package command

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/vault/diagnose"
	"github.com/mitchellh/cli"
)

func testOperatorDiagnoseCommand(tb testing.TB) *OperatorDiagnoseCommand {
	tb.Helper()

	ui := cli.NewMockUi()
	return &OperatorDiagnoseCommand{
		diagnose: diagnose.New(ioutil.Discard),
		BaseCommand: &BaseCommand{
			UI: ui,
		},
		skipEndEnd: true,
	}
}

func TestOperatorDiagnoseCommand_Run(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		args     []string
		expected []*diagnose.Result
	}{
		{
			"diagnose_ok",
			[]string{
				"-config", "./server/test-fixtures/config_diagnose_ok.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "Parse Configuration",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Start Listeners",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Listener TLS Checks",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"TLS is disabled in a listener config stanza.",
							},
						},
					},
				},
				{
					Name:   "Storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Consul TLS Checks",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Consul Storage Direct Server Access",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "Service Discovery",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Consul Service Discovery TLS Checks",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Consul Service Discovery Direct Server Access Checks",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "Create Vault Server Configuration Seals",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Create Core Configuration",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Initialize Randomness for Core",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "HA Storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create HA Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Consul HA Storage Direct Server Access",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Consul TLS Checks",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "Determine Redirect Address",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Cluster Address Checks",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Core Creation Checks",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Start Listeners",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Listener TLS Checks",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"TLS is disabled in a listener config stanza.",
							},
						},
					},
				},
				{
					Name:    "Barrier Encryption Checks",
					Status:  diagnose.ErrorStatus,
					Message: "Diagnose could not create a barrier seal object",
				},
				{
					Name:   "Server Runtime Checks",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "shamir Seal Finalization",
					Status: diagnose.OkStatus,
				},
			},
		},
		{
			"diagnose_raft_problems",
			[]string{
				"-config", "./server/test-fixtures/config_raft.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "Storage",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:    "Raft Folder Permission Checks",
							Status:  diagnose.WarningStatus,
							Message: "too many permissions",
						},
						{
							Name:    "Raft Quorum Checks",
							Status:  diagnose.WarningStatus,
							Message: "0 voters found",
						},
					},
				},
			},
		},
		{
			"diagnose_invalid_storage",
			[]string{
				"-config", "./server/test-fixtures/nostore_config.hcl",
			},
			[]*diagnose.Result{
				{
					Name:    "Storage",
					Status:  diagnose.ErrorStatus,
					Message: "No storage stanza in Vault Server Configuration.",
				},
			},
		},
		{
			"diagnose_listener_config_ok",
			[]string{
				"-config", "./server/test-fixtures/tls_config_ok.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "Start Listeners",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Listener TLS Checks",
							Status: diagnose.OkStatus,
						},
					},
				},
			},
		},
		{
			"diagnose_invalid_https_storage",
			[]string{
				"-config", "./server/test-fixtures/config_bad_https_storage.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "Storage",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:    "Consul TLS Checks",
							Status:  diagnose.ErrorStatus,
							Message: "certificate has expired or is not yet valid",
							Warnings: []string{
								"expired or near expiry",
							},
						},
						{
							Name:   "Consul Storage Direct Server Access",
							Status: diagnose.OkStatus,
						},
					},
				},
			},
		},
		{
			"diagnose_invalid_https_hastorage",
			[]string{
				"-config", "./server/test-fixtures/config_diagnose_hastorage_bad_https.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "Storage",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Consul TLS Checks",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Consul Storage Direct Server Access",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"Consul storage does not connect to local agent, but directly to server.",
							},
						},
					},
				},
				{
					Name:   "HA Storage",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create HA Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Consul HA Storage Direct Server Access",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"Consul storage does not connect to local agent, but directly to server.",
							},
						},
						{
							Name:    "Consul TLS Checks",
							Status:  diagnose.ErrorStatus,
							Message: "certificate has expired or is not yet valid",
							Warnings: []string{
								"expired or near expiry",
							},
						},
					},
				},
				{
					Name:   "Cluster Address Checks",
					Status: diagnose.ErrorStatus,
				},
			},
		},
		{
			"diagnose_seal_transit_tls_check_fail",
			[]string{
				"-config", "./server/test-fixtures/diagnose_seal_transit_tls_check.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "Transit Seal TLS Checks",
					Status: diagnose.WarningStatus,
					Warnings: []string{
						"Found at least one intermediate certificate in the CA certificate file.",
					},
				},
			},
		},
		{
			"diagnose_invalid_https_sr",
			[]string{
				"-config", "./server/test-fixtures/diagnose_bad_https_consul_sr.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "Service Discovery",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:    "Consul Service Discovery TLS Checks",
							Status:  diagnose.ErrorStatus,
							Message: "certificate has expired or is not yet valid",
							Warnings: []string{
								"expired or near expiry",
							},
						},
						{
							Name:   "Consul Service Discovery Direct Server Access Checks",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								diagnose.DirAccessErr,
							},
						},
					},
				},
			},
		},
		{
			"diagnose_direct_storage_access",
			[]string{
				"-config", "./server/test-fixtures/diagnose_ok_storage_direct_access.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "Storage",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Consul TLS Checks",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Consul Storage Direct Server Access",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								diagnose.DirAccessErr,
							},
						},
					},
				},
			},
		},
		{
			"diagnose_raft_no_folder_backend",
			[]string{
				"-config", "./server/test-fixtures/diagnose_raft_no_bolt_folder.hcl",
			},
			[]*diagnose.Result{
				{
					Name:    "Storage",
					Status:  diagnose.ErrorStatus,
					Message: "Diagnose could not initialize storage backend.",
					Children: []*diagnose.Result{
						{
							Name:    "Create Storage Backend",
							Status:  diagnose.ErrorStatus,
							Message: "no such file or directory",
						},
					},
				},
			},
		},
	}

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				client, closer := testVaultServer(t)
				defer closer()

				cmd := testOperatorDiagnoseCommand(t)
				cmd.client = client

				cmd.Run(tc.args)
				result := cmd.diagnose.Finalize(context.Background())

				if err := compareResults(tc.expected, result.Children); err != nil {
					t.Fatalf("Did not find expected test results: %v", err)
				}
			})
		}
	})
}

func compareResults(expected []*diagnose.Result, actual []*diagnose.Result) error {
	for _, exp := range expected {
		found := false
		// Check them all so we don't have to be order specific
		for _, act := range actual {
			fmt.Printf("%+v", act)
			if exp.Name == act.Name {
				found = true
				if err := compareResult(exp, act); err != nil {
					return err
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("could not find expected test result: %s", exp.Name)
		}
	}
	return nil
}

func compareResult(exp *diagnose.Result, act *diagnose.Result) error {
	if exp.Name != act.Name {
		return fmt.Errorf("names mismatch: %s vs %s", exp.Name, act.Name)
	}
	if exp.Status != act.Status {
		if act.Status != diagnose.OkStatus {
			return fmt.Errorf("section %s, status mismatch: %s vs %s, got error %s", exp.Name, exp.Status, act.Status, act.Message)

		}
		return fmt.Errorf("section %s, status mismatch: %s vs %s", exp.Name, exp.Status, act.Status)
	}
	if exp.Message != "" && exp.Message != act.Message && !strings.Contains(act.Message, exp.Message) {
		return fmt.Errorf("section %s, message not found: %s in %s", exp.Name, exp.Message, act.Message)
	}
	if len(exp.Warnings) != len(act.Warnings) {
		return fmt.Errorf("section %s, warning count mismatch: %d vs %d", exp.Name, len(exp.Warnings), len(act.Warnings))
	}
	for j := range exp.Warnings {
		if !strings.Contains(act.Warnings[j], exp.Warnings[j]) {
			return fmt.Errorf("section %s, warning message not found: %s in %s", exp.Name, exp.Warnings[j], act.Warnings[j])
		}
	}
	if len(exp.Children) > len(act.Children) {
		errStrings := []string{}
		for _, c := range act.Children {
			errStrings = append(errStrings, fmt.Sprintf("%+v", c))
		}
		return fmt.Errorf(strings.Join(errStrings, ","))
	}

	if len(exp.Children) > 0 {
		return compareResults(exp.Children, act.Children)
	}

	// Remove raft file if it exists
	os.Remove("./server/test-fixtures/vault.db")
	os.RemoveAll("./server/test-fixtures/raft")

	return nil
}
