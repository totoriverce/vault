package quotas

import (
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"go.uber.org/goleak"
)

type clientResult struct {
	atomicNumAllow *atomic.Int32
	atomicNumFail  *atomic.Int32
}

func TestNewRateLimitQuota(t *testing.T) {
	testCases := []struct {
		name      string
		rlq       *RateLimitQuota
		expectErr bool
	}{
		{"valid rate", NewRateLimitQuota("test-rate-limiter", "qa", "/foo/bar", 16.7, time.Second, 0), false},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink())
			require.Equal(t, tc.expectErr, err != nil, err)
			if err == nil {
				require.Nil(t, tc.rlq.close())
			}
		})
	}
}

func TestRateLimitQuota_Close(t *testing.T) {
	rlq := NewRateLimitQuota("test-rate-limiter", "qa", "/foo/bar", 16.7, time.Second, time.Minute)
	require.NoError(t, rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink()))
	require.NoError(t, rlq.close())

	time.Sleep(time.Second) // allow enough time for purgeClientsLoop to receive on closeCh
	require.False(t, rlq.getPurgeBlocked(), "expected blocked client purging to be disabled after explicit close")
}

func TestRateLimitQuota_Allow(t *testing.T) {
	rlq := &RateLimitQuota{
		Name:          "test-rate-limiter",
		Type:          TypeRateLimit,
		NamespacePath: "qa",
		MountPath:     "/foo/bar",
		Rate:          16.7,

		// override values to lower durations for testing purposes
		purgeInterval: 10 * time.Second,
		staleAge:      10 * time.Second,
	}

	require.NoError(t, rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink()))
	defer rlq.close()

	var wg sync.WaitGroup

	reqFunc := func(addr string, atomicNumAllow, atomicNumFail *atomic.Int32) {
		defer wg.Done()

		resp, err := rlq.allow(&Request{ClientAddress: addr})
		if err != nil {
			return
		}

		if resp.Allowed {
			atomicNumAllow.Add(1)
		} else {
			atomicNumFail.Add(1)
		}
	}

	results := make(map[string]*clientResult)

	start := time.Now()
	end := start.Add(5 * time.Second)

	for time.Now().Before(end) {
		for i := 0; i < 5; i++ {
			wg.Add(1)

			addr := fmt.Sprintf("127.0.0.%d", i)
			cr, ok := results[addr]
			if !ok {
				results[addr] = &clientResult{atomicNumAllow: atomic.NewInt32(0), atomicNumFail: atomic.NewInt32(0)}
				cr = results[addr]
			}

			go reqFunc(addr, cr.atomicNumAllow, cr.atomicNumFail)

			time.Sleep(2 * time.Millisecond)
		}
	}

	wg.Wait()

	// evaluate the ideal RPS as (ceil(RPS) + (RPS * totalSeconds))
	elapsed := time.Since(start)
	ideal := math.Ceil(rlq.Rate) + (rlq.Rate * float64(elapsed) / float64(time.Second))

	for addr, cr := range results {
		numAllow := cr.atomicNumAllow.Load()
		numFail := cr.atomicNumFail.Load()

		// ensure there were some failed requests for the namespace
		require.NotZerof(t, numFail, "expected some requests to fail; addr: %s, numSuccess: %d, numFail: %d, elapsed: %d", addr, numAllow, numFail, elapsed)

		// ensure that we should never get more requests than allowed for the namespace
		want := int32(ideal + 1)
		require.Falsef(t, numAllow > want, "too many successful requests; addr: %s, want: %d, numSuccess: %d, numFail: %d, elapsed: %d", addr, want, numAllow, numFail, elapsed)
	}
}

func TestRateLimitQuota_Allow_WithBlock(t *testing.T) {
	rlq := &RateLimitQuota{
		Name:          "test-rate-limiter",
		Type:          TypeRateLimit,
		NamespacePath: "qa",
		MountPath:     "/foo/bar",
		Rate:          16.7,
		BlockInterval: 10 * time.Second,

		// override values to lower durations for testing purposes
		purgeInterval: 10 * time.Second,
		staleAge:      10 * time.Second,
	}

	require.NoError(t, rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink()))
	defer rlq.close()
	require.True(t, rlq.getPurgeBlocked())

	var wg sync.WaitGroup

	reqFunc := func(addr string, atomicNumAllow, atomicNumFail *atomic.Int32) {
		defer wg.Done()

		resp, err := rlq.allow(&Request{ClientAddress: addr})
		if err != nil {
			return
		}

		if resp.Allowed {
			atomicNumAllow.Add(1)
		} else {
			atomicNumFail.Add(1)
		}
	}

	results := make(map[string]*clientResult)

	start := time.Now()
	end := start.Add(5 * time.Second)

	for time.Now().Before(end) {
		for i := 0; i < 5; i++ {
			wg.Add(1)

			addr := fmt.Sprintf("127.0.0.%d", i)
			cr, ok := results[addr]
			if !ok {
				results[addr] = &clientResult{atomicNumAllow: atomic.NewInt32(0), atomicNumFail: atomic.NewInt32(0)}
				cr = results[addr]
			}

			go reqFunc(addr, cr.atomicNumAllow, cr.atomicNumFail)

			time.Sleep(2 * time.Millisecond)
		}
	}

	wg.Wait()

	for _, cr := range results {
		numAllow := cr.atomicNumAllow.Load()
		numFail := cr.atomicNumFail.Load()

		// Since blocking is enabled, each client should only have 'rate' successful
		// requests, whereas all subsequent requests fail.
		require.Equal(t, int32(17), numAllow)
		require.NotZero(t, numFail)
	}

	func() {
		timeout := time.After(rlq.purgeInterval * 2)
		ticker := time.Tick(time.Second)
		for {
			select {
			case <-timeout:
				require.Failf(t, "timeout exceeded waiting for blocked clients to be purged", "num blocked: %d", rlq.numBlockedClients())

			case <-ticker:
				if rlq.numBlockedClients() == 0 {
					return
				}
			}
		}
	}()
}

func TestRateLimitQuota_Update(t *testing.T) {
	defer goleak.VerifyNone(t)
	qm, err := NewManager(logging.NewVaultLogger(log.Trace), nil, metricsutil.BlackholeSink())
	require.NoError(t, err)

	quota := NewRateLimitQuota("quota1", "", "", 10, time.Second, 0)
	require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), quota, true))
	require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), quota, true))

	require.Nil(t, quota.close())
}
