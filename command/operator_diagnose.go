package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"

	"github.com/hashicorp/vault/helper/constants"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	cserver "github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/internalshared/reloadutil"
	physconsul "github.com/hashicorp/vault/physical/consul"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/version"
	sr "github.com/hashicorp/vault/serviceregistration"
	srconsul "github.com/hashicorp/vault/serviceregistration/consul"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/diagnose"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

const OperatorDiagnoseEnableEnv = "VAULT_DIAGNOSE"

const CoreUninitializedErr = "diagnose cannot attempt this step because core could not be initialized"
const BackendUninitializedErr = "diagnose cannot attempt this step because backend could not be initialized"
const CoreConfigUninitializedErr = "diagnose cannot attempt this step because core config could not be set"

var (
	_ cli.Command             = (*OperatorDiagnoseCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorDiagnoseCommand)(nil)
)

type OperatorDiagnoseCommand struct {
	*BaseCommand
	diagnose *diagnose.Session

	flagDebug    bool
	flagSkips    []string
	flagConfigs  []string
	cleanupGuard sync.Once

	reloadFuncsLock      *sync.RWMutex
	reloadFuncs          *map[string][]reloadutil.ReloadFunc
	ServiceRegistrations map[string]sr.Factory
	startedCh            chan struct{} // for tests
	reloadedCh           chan struct{} // for tests
	skipEndEnd           bool          // for tests
}

func (c *OperatorDiagnoseCommand) Synopsis() string {
	return "Troubleshoot problems starting Vault"
}

func (c *OperatorDiagnoseCommand) Help() string {
	helpText := `
Usage: vault operator diagnose 

  This command troubleshoots Vault startup issues, such as TLS configuration or
  auto-unseal. It should be run using the same environment variables and configuration
  files as the "vault server" command, so that startup problems can be accurately
  reproduced.

  Start diagnose with a configuration file:
    
     $ vault operator diagnose -config=/etc/vault/config.hcl

  Perform a diagnostic check while Vault is still running:

     $ vault operator diagnose -config=/etc/vault/config.hcl -skip=listener

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *OperatorDiagnoseCommand) Flags() *FlagSets {
	set := NewFlagSets(c.UI)
	f := set.NewFlagSet("Command Options")

	f.StringSliceVar(&StringSliceVar{
		Name:   "config",
		Target: &c.flagConfigs,
		Completion: complete.PredictOr(
			complete.PredictFiles("*.hcl"),
			complete.PredictFiles("*.json"),
			complete.PredictDirs("*"),
		),
		Usage: "Path to a Vault configuration file or directory of configuration " +
			"files. This flag can be specified multiple times to load multiple " +
			"configurations. If the path is a directory, all files which end in " +
			".hcl or .json are loaded.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   "skip",
		Target: &c.flagSkips,
		Usage:  "Skip the health checks named as arguments. May be 'listener', 'storage', or 'autounseal'.",
	})

	f.BoolVar(&BoolVar{
		Name:    "debug",
		Target:  &c.flagDebug,
		Default: false,
		Usage:   "Dump all information collected by Diagnose.",
	})

	f.StringVar(&StringVar{
		Name:   "format",
		Target: &c.flagFormat,
		Usage:  "The output format",
	})
	return set
}

func (c *OperatorDiagnoseCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *OperatorDiagnoseCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

const (
	status_unknown = "[      ] "
	status_ok      = "\u001b[32m[  ok  ]\u001b[0m "
	status_failed  = "\u001b[31m[failed]\u001b[0m "
	status_warn    = "\u001b[33m[ warn ]\u001b[0m "
	same_line      = "\u001b[F"
)

func (c *OperatorDiagnoseCommand) Run(args []string) int {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 3
	}
	return c.RunWithParsedFlags()
}

func (c *OperatorDiagnoseCommand) RunWithParsedFlags() int {

	if len(c.flagConfigs) == 0 {
		c.UI.Error("Must specify a configuration file using -config.")
		return 3
	}

	if c.diagnose == nil {
		if c.flagFormat == "json" {
			c.diagnose = diagnose.New(&ioutils.NopWriter{})
		} else {
			c.UI.Output(version.GetVersion().FullVersionNumber(true))
			c.diagnose = diagnose.New(os.Stdout)
		}
	}
	ctx := diagnose.Context(context.Background(), c.diagnose)
	c.diagnose.SetSkipList(c.flagSkips)
	err := c.offlineDiagnostics(ctx)

	results := c.diagnose.Finalize(ctx)
	if c.flagFormat == "json" {
		resultsJS, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error marshalling results: %v", err)
			return 4
		}
		c.UI.Output(string(resultsJS))
	} else {
		c.UI.Output("\nResults:")
		w, _, err := term.GetSize(0)
		if err == nil {
			results.Write(os.Stdout, w)
		} else {
			results.Write(os.Stdout, 0)
		}
	}

	if err != nil {
		return 4
	}
	// Use a different return code
	switch results.Status {
	case diagnose.WarningStatus:
		return 2
	case diagnose.ErrorStatus:
		return 1
	}
	return 0
}

func (c *OperatorDiagnoseCommand) offlineDiagnostics(ctx context.Context) error {
	rloadFuncs := make(map[string][]reloadutil.ReloadFunc)
	server := &ServerCommand{
		// TODO: set up a different one?
		// In particular, a UI instance that won't output?
		BaseCommand: c.BaseCommand,

		// TODO: refactor to a common place?
		AuditBackends:        auditBackends,
		CredentialBackends:   credentialBackends,
		LogicalBackends:      logicalBackends,
		PhysicalBackends:     physicalBackends,
		ServiceRegistrations: serviceRegistrations,

		// TODO: other ServerCommand options?

		logger:          log.NewInterceptLogger(nil),
		allLoggers:      []log.Logger{},
		reloadFuncs:     &rloadFuncs,
		reloadFuncsLock: new(sync.RWMutex),
	}

	ctx, span := diagnose.StartSpan(ctx, "initialization")
	defer span.End()

	// OS Specific checks
	diagnose.OSChecks(ctx)

	var config *cserver.Config

	diagnose.Test(ctx, "Parse configuration", func(ctx context.Context) (err error) {
		server.flagConfigs = c.flagConfigs
		var configErrors []configutil.ConfigError
		config, configErrors, err = server.parseConfig()
		if err != nil {
			return err
		}
		for _, ce := range configErrors {
			diagnose.Warn(ctx, ce.String())
		}
		return nil
	})
	if config == nil {
		return fmt.Errorf("No vault server configuration found.")
	}

	var metricSink *metricsutil.ClusterMetricSink
	var metricsHelper *metricsutil.MetricsHelper

	var backend *physical.Backend
	diagnose.Test(ctx, "storage", func(ctx context.Context) error {

		// Ensure that there is a storage stanza
		if config.Storage == nil {
			return fmt.Errorf("no storage stanza found in config")
		}

		diagnose.Test(ctx, "create-storage-backend", func(ctx context.Context) error {
			b, err := server.setupStorage(config)
			if err != nil {
				return err
			}
			if b == nil {
				return fmt.Errorf("storage backend was initialized to nil")
			}
			backend = &b
			return nil
		})

		if backend == nil {
			diagnose.Fail(ctx, "Diagnose could not initialize storage backend.")
			span.End()
			return fmt.Errorf("Diagnose could not initialize storage backend.")
		}

		// Check for raft quorum status
		if config.Storage.Type == storageTypeRaft {
			path := os.Getenv(raft.EnvVaultRaftPath)
			if path == "" {
				path, ok := config.Storage.Config["path"]
				if !ok {
					diagnose.SpotError(ctx, "raft folder permission checks", fmt.Errorf("storage file path is required"))
				}
				diagnose.RaftFileChecks(ctx, path)
			}
			if backend != nil {
				diagnose.RaftStorageQuorum(ctx, (*backend).(*raft.RaftBackend))
			} else {
				diagnose.SpotError(ctx, "raft quorum", fmt.Errorf("could not determine quorum status without initialized backend"))
			}
		}

		// Consul storage checks
		if config.Storage != nil && config.Storage.Type == storageTypeConsul {
			diagnose.Test(ctx, "test-storage-tls-consul", func(ctx context.Context) error {
				err := physconsul.SetupSecureTLS(ctx, api.DefaultConfig(), config.Storage.Config, server.logger, true)
				if err != nil {
					return err
				}
				return nil
			})

			diagnose.Test(ctx, "test-consul-direct-access-storage", func(ctx context.Context) error {
				dirAccess := diagnose.ConsulDirectAccess(config.Storage.Config)
				if dirAccess != "" {
					diagnose.Warn(ctx, dirAccess)
				}
				return nil
			})
		}

		// Attempt to use storage backend
		if !c.skipEndEnd && config.Storage.Type != storageTypeRaft {
			diagnose.Test(ctx, "test-access-storage", diagnose.WithTimeout(30*time.Second, func(ctx context.Context) error {
				maxDurationCrudOperation := "write"
				maxDuration := time.Duration(0)
				uuidSuffix, err := uuid.GenerateUUID()
				if err != nil {
					return err
				}
				uuid := "diagnose/latency/" + uuidSuffix
				dur, err := diagnose.EndToEndLatencyCheckWrite(ctx, uuid, *backend)
				if err != nil {
					return err
				}
				maxDuration = dur
				dur, err = diagnose.EndToEndLatencyCheckRead(ctx, uuid, *backend)
				if err != nil {
					return err
				}
				if dur > maxDuration {
					maxDuration = dur
					maxDurationCrudOperation = "read"
				}
				dur, err = diagnose.EndToEndLatencyCheckDelete(ctx, uuid, *backend)
				if err != nil {
					return err
				}
				if dur > maxDuration {
					maxDuration = dur
					maxDurationCrudOperation = "delete"
				}

				if maxDuration > time.Duration(0) {
					diagnose.Warn(ctx, diagnose.LatencyWarning+fmt.Sprintf("duration: %s, ", maxDuration)+fmt.Sprintf("operation: %s", maxDurationCrudOperation))
				}
				return nil
			}))
		}
		return nil
	})

	// Return from top-level span when backend is nil
	if backend == nil {
		return fmt.Errorf("Diagnose could not initialize storage backend.")
	}

	var configSR sr.ServiceRegistration
	diagnose.Test(ctx, "service-discovery", func(ctx context.Context) error {
		if config.ServiceRegistration == nil || config.ServiceRegistration.Config == nil {
			diagnose.Skipped(ctx, "no service registration configured")
			return nil
		}
		srConfig := config.ServiceRegistration.Config

		diagnose.Test(ctx, "test-serviceregistration-tls-consul", func(ctx context.Context) error {
			// SetupSecureTLS for service discovery uses the same cert and key to set up physical
			// storage. See the consul package in physical for details.
			err := srconsul.SetupSecureTLS(ctx, api.DefaultConfig(), srConfig, server.logger, true)
			if err != nil {
				return err
			}
			return nil
		})

		if config.ServiceRegistration != nil && config.ServiceRegistration.Type == "consul" {
			diagnose.Test(ctx, "test-consul-direct-access-service-discovery", func(ctx context.Context) error {
				dirAccess := diagnose.ConsulDirectAccess(config.ServiceRegistration.Config)
				if dirAccess != "" {
					diagnose.Warn(ctx, dirAccess)
				}
				return nil
			})
		}
		return nil
	})

	diagnose.Test(ctx, "seal-transit-tls-checks", func(ctx context.Context) error {
		var checkSealTransit bool
		for _, seal := range config.Seals {
			if seal.Type == "transit" {
				checkSealTransit = true

				tlsSkipVerify, _ := seal.Config["tls_skip_verify"]
				if tlsSkipVerify == "true" {
					diagnose.Warn(ctx, "TLS verification is Skipped! Using this option is highly discouraged and decreases the security of data transmissions to and from the Vault server.")
					return nil
				}

				// Checking tls_client_cert and tls_client_key
				tlsClientCert, ok := seal.Config["tls_client_cert"]
				if !ok {
					diagnose.Warn(ctx, "Missing tls_client_cert in the config")
					return nil
				}
				tlsClientKey, ok := seal.Config["tls_client_key"]
				if !ok {
					diagnose.Warn(ctx, "Missing tls_client_key in the config")
					return nil
				}
				_, err := diagnose.TLSFileChecks(tlsClientCert, tlsClientKey)
				if err != nil {
					return fmt.Errorf("TLS file check failed for tls_client_cert and tls_client_key with the following error: %w", err)
				}

				// checking tls_ca_cert
				tlsCACert, ok := seal.Config["tls_ca_cert"]
				if !ok {
					diagnose.Warn(ctx, "Mising tls_ca_cert in the config")
					return nil
				}
				_, err = diagnose.TLSCAFileCheck(tlsCACert)
				if err != nil {
					return fmt.Errorf("TLS file check failed for tls_ca_cert with the following error: %w", err)
				}
			}
		}
		if !checkSealTransit {
			diagnose.Skipped(ctx, "No transit seal found!")
		}
		return nil
	})

	sealcontext, sealspan := diagnose.StartSpan(ctx, "create-seal")
	var seals []vault.Seal
	var sealConfigError error

	barrierSeal, barrierWrapper, unwrapSeal, seals, sealConfigError, err := setSeal(server, config, make([]string, 0), make(map[string]string))

	// Check error here
	if err != nil {
		diagnose.Fail(sealcontext, err.Error())
		goto SEALFAIL
	}
	if sealConfigError != nil {
		diagnose.Fail(sealcontext, "seal could not be configured: seals may already be initialized")
		goto SEALFAIL
	}

	if seals != nil {
		for _, seal := range seals {
			// Ensure that the seal finalizer is called, even if using verify-only
			defer func(seal *vault.Seal) {
				sealType := (*seal).BarrierType()
				finalizeSealContext, finalizeSealSpan := diagnose.StartSpan(ctx, "finalize-seal-"+sealType)
				err = (*seal).Finalize(finalizeSealContext)
				if err != nil {
					diagnose.Fail(finalizeSealContext, "error finalizing seal")
					finalizeSealSpan.End()
				}
				finalizeSealSpan.End()
			}(&seal)
		}
	}

	if barrierSeal == nil {
		diagnose.Fail(sealcontext, "could not create barrier seal! Most likely proper Seal configuration information was not set, but no error was generated")
	}

SEALFAIL:
	sealspan.End()
	var coreConfig vault.CoreConfig
	diagnose.Test(ctx, "setup-core", func(ctx context.Context) error {
		var secureRandomReader io.Reader
		// prepare a secure random reader for core
		secureRandomReader, err = configutil.CreateSecureRandomReaderFunc(config.SharedConfig, barrierWrapper)
		if err != nil {
			return diagnose.SpotError(ctx, "init-randreader", err)
		}
		diagnose.SpotOk(ctx, "init-randreader", "")
		coreConfig = createCoreConfig(server, config, *backend, configSR, barrierSeal, unwrapSeal, metricsHelper, metricSink, secureRandomReader)
		return nil
	})

	var disableClustering bool
	diagnose.Test(ctx, "setup-ha-storage", func(ctx context.Context) error {
		diagnose.Test(ctx, "create-ha-storage-backend", func(ctx context.Context) error {
			// Initialize the separate HA storage backend, if it exists
			disableClustering, err = initHaBackend(server, config, &coreConfig, *backend)
			if err != nil {
				return err
			}
			return nil
		})

		diagnose.Test(ctx, "test-consul-direct-access-storage", func(ctx context.Context) error {
			if config.HAStorage == nil {
				diagnose.Skipped(ctx, "no HA storage configured")
			} else {
				dirAccess := diagnose.ConsulDirectAccess(config.HAStorage.Config)
				if dirAccess != "" {
					diagnose.Warn(ctx, dirAccess)
				}
			}
			return nil
		})
		if config.HAStorage != nil && config.HAStorage.Type == storageTypeConsul {
			diagnose.Test(ctx, "test-ha-storage-tls-consul", func(ctx context.Context) error {
				err = physconsul.SetupSecureTLS(ctx, api.DefaultConfig(), config.HAStorage.Config, server.logger, true)
				if err != nil {
					return err
				}
				return nil
			})
		}
		return nil
	})

	// Determine the redirect address from environment variables
	err = determineRedirectAddr(server, &coreConfig, config)
	if err != nil {
		return diagnose.SpotError(ctx, "determine-redirect", err)
	}
	diagnose.SpotOk(ctx, "determine-redirect", "")

	err = findClusterAddress(server, &coreConfig, config, disableClustering)
	if err != nil {
		return diagnose.SpotError(ctx, "find-cluster-addr", err)
	}
	diagnose.SpotOk(ctx, "find-cluster-addr", "")

	var vaultCore *vault.Core

	// Run all the checks that are utilized when initializing a core object
	// without actually calling core.Init. These are in the init-core section
	// as they are runtime checks.
	diagnose.Test(ctx, "init-core", func(ctx context.Context) error {
		var newCoreError error
		if coreConfig.RawConfig == nil {
			return fmt.Errorf(CoreConfigUninitializedErr)
		}
		core, newCoreError := vault.CreateCore(&coreConfig)
		if newCoreError != nil {
			if vault.IsFatalError(newCoreError) {
				return fmt.Errorf("Error initializing core: %s", newCoreError)
			}
			diagnose.Warn(ctx, wrapAtLength(
				"WARNING! A non-fatal error occurred during initialization. Please "+
					"check the logs for more information."))
		} else {
			vaultCore = core
		}
		return nil
	})

	if vaultCore == nil {
		return fmt.Errorf("Diagnose could not initialize the vault core from the vault server configuration.")
	}

	licenseCtx, licenseSpan := diagnose.StartSpan(ctx, "autoloaded license")
	// If we are not in enterprise, return from the check
	if !constants.IsEnterprise {
		diagnose.Skipped(licenseCtx, "License check will not run on OSS Vault.")
	} else {
		// Load License from environment variables. These take precedence over the
		// configured license.
		if envLicensePath := os.Getenv(EnvVaultLicensePath); envLicensePath != "" {
			coreConfig.LicensePath = envLicensePath
		}
		if envLicense := os.Getenv(EnvVaultLicense); envLicense != "" {
			coreConfig.License = envLicense
		}
		vault.DiagnoseCheckLicense(licenseCtx, vaultCore, coreConfig)
	}
	licenseSpan.End()

	var lns []listenerutil.Listener
	diagnose.Test(ctx, "init-listeners", func(ctx context.Context) error {
		disableClustering := config.HAStorage != nil && config.HAStorage.DisableClustering
		infoKeys := make([]string, 0, 10)
		info := make(map[string]string)
		var listeners []listenerutil.Listener
		var status int

		diagnose.ListenerChecks(ctx, config.Listeners)

		diagnose.Test(ctx, "create-listeners", func(ctx context.Context) error {
			status, listeners, _, err = server.InitListeners(config, disableClustering, &infoKeys, &info)
			if status != 0 {
				return err
			}
			return nil
		})

		lns = listeners

		// Make sure we close all listeners from this point on
		listenerCloseFunc := func() {
			for _, ln := range lns {
				ln.Listener.Close()
			}
		}

		c.cleanupGuard.Do(listenerCloseFunc)

		return nil
	})

	// TODO: Diagnose logging configuration

	// The unseal diagnose check will simply attempt to use the barrier to encrypt and
	// decrypt a mock value. It will not call runUnseal.
	diagnose.Test(ctx, "unseal", diagnose.WithTimeout(30*time.Second, func(ctx context.Context) error {
		if barrierWrapper == nil {
			return fmt.Errorf("Diagnose could not create a barrier seal object")
		}
		barrierUUID, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Diagnose could not create unique UUID for unsealing")
		}
		barrierEncValue := "diagnose-" + barrierUUID
		ciphertext, err := barrierWrapper.Encrypt(ctx, []byte(barrierEncValue), nil)
		if err != nil {
			return fmt.Errorf("Error encrypting with seal barrier: %w", err)
		}
		plaintext, err := barrierWrapper.Decrypt(ctx, ciphertext, nil)
		if err != nil {
			return fmt.Errorf("Error decrypting with seal barrier: %w", err)

		}
		if string(plaintext) != barrierEncValue {
			return fmt.Errorf("barrier returned incorrect decrypted value for mock data")
		}
		return nil
	}))

	// The following block contains static checks that are run during the
	// startHttpServers portion of server run. In other words, they are static
	// checks during resource creation.
	diagnose.Test(ctx, "start-servers", func(ctx context.Context) error {
		for _, ln := range lns {
			if ln.Config == nil {
				return fmt.Errorf("Found nil listener config after parsing")
			}
		}
		return nil
	})
	return nil
}
