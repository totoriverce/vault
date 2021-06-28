package http

import (
	"context"
	"strings"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault"
)

func TestSysMonitorUnknownLogLevel(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{HandlerFunc: Handler})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	request := client.NewRequest("GET", "/v1/sys/monitor")
	request.Params.Add("log_level", "haha")
	_, err := client.RawRequest(request)

	if err == nil {
		t.Fatal("expected to get an error, but didn't")
	} else {
		if !strings.Contains(err.Error(), "Code: 400") {
			t.Fatalf("expected to receive a 400 error, but got %s instead", err)
		}

		if !strings.Contains(err.Error(), "unknown log level") {
			t.Fatalf("expected to receive a message indicating an unknown log level, but got %s instead", err)
		}
	}
}

func TestSysMonitorStreamingLogs(t *testing.T) {
	logger := log.NewInterceptLogger(&log.LoggerOptions{
		Output:     log.DefaultOutput,
		Level:      log.Debug,
		JSONFormat: logging.ParseEnvLogFormat() == logging.JSONFormat,
	})

	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{HandlerFunc: Handler, Logger: logger})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	stopCh := testhelpers.GenerateDebugLogs(t, client)

	debugCount := 0
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	logCh, err := client.Sys().Monitor(ctx, "DEBUG")
	if err != nil {
		t.Fatal(err)
	}

	timeCh := time.After(5 * time.Second)

	for {
		select {
		case log := <-logCh:
			if strings.Contains(log, "[DEBUG]") {
				debugCount++
			}
		case <-timeCh:
			t.Fatal("Failed to get a DEBUG message after 5 seconds")
		}

		// If we've seen multiple lines that match what we want,
		// it's probably safe to assume streaming is working
		if debugCount > 3 {
			stopCh <- struct{}{}
			break
		}
	}

	<-stopCh
}
