package testhelpers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	raftlib "github.com/hashicorp/raft"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/go-testing-interface"
)

type GenerateRootKind int

const (
	GenerateRootRegular GenerateRootKind = iota
	GenerateRootDR
	GenerateRecovery
)

// Generates a root token on the target cluster.
func GenerateRoot(t testing.T, cluster *vault.TestCluster, kind GenerateRootKind) string {
	t.Helper()
	token, err := GenerateRootWithError(t, cluster, kind)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func GenerateRootWithError(t testing.T, cluster *vault.TestCluster, kind GenerateRootKind) (string, error) {
	t.Helper()
	// If recovery keys supported, use those to perform root token generation instead
	var keys [][]byte
	if cluster.Cores[0].SealAccess().RecoveryKeySupported() {
		keys = cluster.RecoveryKeys
	} else {
		keys = cluster.BarrierKeys
	}
	client := cluster.Cores[0].Client

	var err error
	var status *api.GenerateRootStatusResponse
	switch kind {
	case GenerateRootRegular:
		status, err = client.Sys().GenerateRootInit("", "")
	case GenerateRootDR:
		status, err = client.Sys().GenerateDROperationTokenInit("", "")
	case GenerateRecovery:
		status, err = client.Sys().GenerateRecoveryOperationTokenInit("", "")
	}
	if err != nil {
		return "", err
	}

	if status.Required > len(keys) {
		return "", fmt.Errorf("need more keys than have, need %d have %d", status.Required, len(keys))
	}

	otp := status.OTP

	for i, key := range keys {
		if i >= status.Required {
			break
		}

		strKey := base64.StdEncoding.EncodeToString(key)
		switch kind {
		case GenerateRootRegular:
			status, err = client.Sys().GenerateRootUpdate(strKey, status.Nonce)
		case GenerateRootDR:
			status, err = client.Sys().GenerateDROperationTokenUpdate(strKey, status.Nonce)
		case GenerateRecovery:
			status, err = client.Sys().GenerateRecoveryOperationTokenUpdate(strKey, status.Nonce)
		}
		if err != nil {
			return "", err
		}
	}
	if !status.Complete {
		return "", errors.New("generate root operation did not end successfully")
	}

	tokenBytes, err := base64.RawStdEncoding.DecodeString(status.EncodedToken)
	if err != nil {
		return "", err
	}
	tokenBytes, err = xor.XORBytes(tokenBytes, []byte(otp))
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}

// RandomWithPrefix is used to generate a unique name with a prefix, for
// randomizing names in acceptance tests
func RandomWithPrefix(name string) string {
	return fmt.Sprintf("%s-%d", name, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

func EnsureCoresSealed(t testing.T, c *vault.TestCluster) {
	t.Helper()
	for _, core := range c.Cores {
		EnsureCoreSealed(t, core)
	}
}

func EnsureCoreSealed(t testing.T, core *vault.TestClusterCore) {
	t.Helper()
	core.Seal(t)
	timeout := time.Now().Add(60 * time.Second)
	for {
		if time.Now().After(timeout) {
			t.Fatal("timeout waiting for core to seal")
		}
		if core.Core.Sealed() {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func EnsureCoresUnsealed(t testing.T, c *vault.TestCluster) {
	t.Helper()
	for i, core := range c.Cores {
		err := AttemptUnsealCore(c, core)
		if err != nil {
			t.Fatalf("failed to unseal core %d: %v", i, err)
		}
	}
}

func EnsureCoreUnsealed(t testing.T, c *vault.TestCluster, core *vault.TestClusterCore) {
	t.Helper()
	err := AttemptUnsealCore(c, core)
	if err != nil {
		t.Fatalf("failed to unseal core: %v", err)
	}
}

func AttemptUnsealCores(c *vault.TestCluster) error {
	for i, core := range c.Cores {
		err := AttemptUnsealCore(c, core)
		if err != nil {
			return fmt.Errorf("failed to unseal core %d: %v", i, err)
		}
	}
	return nil
}

func AttemptUnsealCore(c *vault.TestCluster, core *vault.TestClusterCore) error {
	if !core.Sealed() {
		return nil
	}

	core.SealAccess().ClearCaches(context.Background())
	if err := core.UnsealWithStoredKeys(context.Background()); err != nil {
		return err
	}

	client := core.Client
	client.Sys().ResetUnsealProcess()
	for j := 0; j < len(c.BarrierKeys); j++ {
		statusResp, err := client.Sys().Unseal(base64.StdEncoding.EncodeToString(c.BarrierKeys[j]))
		if err != nil {
			// Sometimes when we get here it's already unsealed on its own
			// and then this fails for DR secondaries so check again
			if core.Sealed() {
				return err
			} else {
				return nil
			}
		}
		if statusResp == nil {
			return fmt.Errorf("nil status response during unseal")
		}
		if !statusResp.Sealed {
			break
		}
	}
	if core.Sealed() {
		return fmt.Errorf("core is still sealed")
	}
	return nil
}

func EnsureStableActiveNode(t testing.T, cluster *vault.TestCluster) {
	t.Helper()
	deriveStableActiveCore(t, cluster)
}

func DeriveStableActiveCore(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	t.Helper()
	return deriveStableActiveCore(t, cluster)
}

func deriveStableActiveCore(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	t.Helper()
	activeCore := DeriveActiveCore(t, cluster)
	minDuration := time.NewTimer(3 * time.Second)

	for i := 0; i < 30; i++ {
		leaderResp, err := activeCore.Client.Sys().Leader()
		if err != nil {
			t.Fatal(err)
		}
		if !leaderResp.IsSelf {
			minDuration.Reset(3 * time.Second)
		}
		time.Sleep(200 * time.Millisecond)
	}

	select {
	case <-minDuration.C:
	default:
		if stopped := minDuration.Stop(); stopped {
			t.Fatal("unstable active node")
		}
		// Drain the value
		<-minDuration.C
	}

	return activeCore
}

func DeriveActiveCore(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	t.Helper()
	for i := 0; i < 20; i++ {
		for _, core := range cluster.Cores {
			leaderResp, err := core.Client.Sys().Leader()
			if err != nil {
				t.Fatal(err)
			}
			if leaderResp.IsSelf {
				return core
			}
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatal("could not derive the active core")
	return nil
}

func DeriveStandbyCores(t testing.T, cluster *vault.TestCluster) []*vault.TestClusterCore {
	t.Helper()
	cores := make([]*vault.TestClusterCore, 0, 2)
	for _, core := range cluster.Cores {
		leaderResp, err := core.Client.Sys().Leader()
		if err != nil {
			t.Fatal(err)
		}
		if !leaderResp.IsSelf {
			cores = append(cores, core)
		}
	}

	return cores
}

func WaitForNCoresUnsealed(t testing.T, cluster *vault.TestCluster, n int) {
	t.Helper()
	for i := 0; i < 30; i++ {
		unsealed := 0
		for _, core := range cluster.Cores {
			if !core.Core.Sealed() {
				unsealed++
			}
		}

		if unsealed >= n {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("%d cores were not unsealed", n)
}

func SealCores(t testing.T, cluster *vault.TestCluster) {
	t.Helper()
	for _, core := range cluster.Cores {
		if err := core.Shutdown(); err != nil {
			t.Fatal(err)
		}
		timeout := time.Now().Add(3 * time.Second)
		for {
			if time.Now().After(timeout) {
				t.Fatal("timeout waiting for core to seal")
			}
			if core.Sealed() {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func WaitForNCoresSealed(t testing.T, cluster *vault.TestCluster, n int) {
	t.Helper()
	for i := 0; i < 60; i++ {
		sealed := 0
		for _, core := range cluster.Cores {
			if core.Core.Sealed() {
				sealed++
			}
		}

		if sealed >= n {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("%d cores were not sealed", n)
}

func WaitForActiveNode(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	t.Helper()
	for i := 0; i < 30; i++ {
		for _, core := range cluster.Cores {
			if standby, _ := core.Core.Standby(); !standby {
				return core
			}
		}

		time.Sleep(time.Second)
	}

	t.Fatalf("node did not become active")
	return nil
}

func WaitForStandbyNode(t testing.T, core *vault.TestClusterCore) {
	t.Helper()
	for i := 0; i < 30; i++ {
		if isLeader, _, clusterAddr, _ := core.Core.Leader(); isLeader != true && clusterAddr != "" {
			return
		}

		time.Sleep(time.Second)
	}

	t.Fatalf("node did not become standby")
}

func RekeyCluster(t testing.T, cluster *vault.TestCluster, recovery bool) [][]byte {
	t.Helper()
	cluster.Logger.Info("rekeying cluster", "recovery", recovery)
	client := cluster.Cores[0].Client

	initFunc := client.Sys().RekeyInit
	if recovery {
		initFunc = client.Sys().RekeyRecoveryKeyInit
	}
	init, err := initFunc(&api.RekeyInitRequest{
		SecretShares:    5,
		SecretThreshold: 3,
	})
	if err != nil {
		t.Fatal(err)
	}

	var statusResp *api.RekeyUpdateResponse
	keys := cluster.BarrierKeys
	if cluster.Cores[0].Core.SealAccess().RecoveryKeySupported() {
		keys = cluster.RecoveryKeys
	}

	updateFunc := client.Sys().RekeyUpdate
	if recovery {
		updateFunc = client.Sys().RekeyRecoveryKeyUpdate
	}
	for j := 0; j < len(keys); j++ {
		statusResp, err = updateFunc(base64.StdEncoding.EncodeToString(keys[j]), init.Nonce)
		if err != nil {
			t.Fatal(err)
		}
		if statusResp == nil {
			t.Fatal("nil status response during unseal")
		}
		if statusResp.Complete {
			break
		}
	}
	cluster.Logger.Info("cluster rekeyed", "recovery", recovery)

	if cluster.Cores[0].Core.SealAccess().RecoveryKeySupported() && !recovery {
		return nil
	}
	if len(statusResp.KeysB64) != 5 {
		t.Fatal("wrong number of keys")
	}

	newKeys := make([][]byte, 5)
	for i, key := range statusResp.KeysB64 {
		newKeys[i], err = base64.StdEncoding.DecodeString(key)
		if err != nil {
			t.Fatal(err)
		}
	}
	return newKeys
}

// TestRaftServerAddressProvider is a ServerAddressProvider that uses the
// ClusterAddr() of each node to provide raft addresses.
//
// Note that TestRaftServerAddressProvider should only be used in cases where
// cores that are part of a raft configuration have already had
// startClusterListener() called (via either unsealing or raft joining).
type TestRaftServerAddressProvider struct {
	Cluster *vault.TestCluster
}

func (p *TestRaftServerAddressProvider) ServerAddr(id raftlib.ServerID) (raftlib.ServerAddress, error) {
	for _, core := range p.Cluster.Cores {
		if core.NodeID == string(id) {
			parsed, err := url.Parse(core.ClusterAddr())
			if err != nil {
				return "", err
			}

			return raftlib.ServerAddress(parsed.Host), nil
		}
	}

	return "", errors.New("could not find cluster addr")
}

func RaftClusterJoinNodes(t testing.T, cluster *vault.TestCluster) {
	addressProvider := &TestRaftServerAddressProvider{Cluster: cluster}

	atomic.StoreUint32(&vault.TestingUpdateClusterAddr, 1)

	leader := cluster.Cores[0]

	// Seal the leader so we can install an address provider
	{
		EnsureCoreSealed(t, leader)
		leader.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		cluster.UnsealCore(t, leader)
		vault.TestWaitActive(t, leader.Core)
	}

	leaderInfos := []*raft.LeaderJoinInfo{
		{
			LeaderAPIAddr: leader.Client.Address(),
			TLSConfig:     leader.TLSConfig,
		},
	}

	// Join followers
	for i := 1; i < len(cluster.Cores); i++ {
		core := cluster.Cores[i]
		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderInfos, false)
		if err != nil {
			t.Fatal(err)
		}

		cluster.UnsealCore(t, core)
	}

	WaitForNCoresUnsealed(t, cluster, len(cluster.Cores))
}

// HardcodedServerAddressProvider is a ServerAddressProvider that uses
// a hardcoded map of raft node addresses.
//
// It is useful in cases where the raft configuration is known ahead of time,
// but some of the cores have not yet had startClusterListener() called (via
// either unsealing or raft joining), and thus do not yet have a ClusterAddr()
// assigned.
type HardcodedServerAddressProvider struct {
	Entries map[raftlib.ServerID]raftlib.ServerAddress
}

func (p *HardcodedServerAddressProvider) ServerAddr(id raftlib.ServerID) (raftlib.ServerAddress, error) {
	if addr, ok := p.Entries[id]; ok {
		return addr, nil
	}
	return "", errors.New("could not find cluster addr")
}

// NewHardcodedServerAddressProvider is a convenience function that makes a
// ServerAddressProvider from a given cluster address base port.
func NewHardcodedServerAddressProvider(numCores, baseClusterPort int) raftlib.ServerAddressProvider {
	entries := make(map[raftlib.ServerID]raftlib.ServerAddress)

	for i := 0; i < numCores; i++ {
		id := fmt.Sprintf("core-%d", i)
		addr := fmt.Sprintf("127.0.0.1:%d", baseClusterPort+i)
		entries[raftlib.ServerID(id)] = raftlib.ServerAddress(addr)
	}

	return &HardcodedServerAddressProvider{
		entries,
	}
}

// VerifyRaftConfiguration checks that we have a valid raft configuration, i.e.
// the correct number of servers, having the correct NodeIDs, and exactly one
// leader.
func VerifyRaftConfiguration(core *vault.TestClusterCore, numCores int) error {
	backend := core.UnderlyingRawStorage.(*raft.RaftBackend)
	ctx := namespace.RootContext(context.Background())
	config, err := backend.GetConfiguration(ctx)
	if err != nil {
		return err
	}

	servers := config.Servers
	if len(servers) != numCores {
		return fmt.Errorf("Found %d servers, not %d", len(servers), numCores)
	}

	leaders := 0
	for i, s := range servers {
		if s.NodeID != fmt.Sprintf("core-%d", i) {
			return fmt.Errorf("Found unexpected node ID %q", s.NodeID)
		}
		if s.Leader {
			leaders++
		}
	}

	if leaders != 1 {
		return fmt.Errorf("Found %d leaders", leaders)
	}

	return nil
}

func RaftAppliedIndex(core *vault.TestClusterCore) uint64 {
	return core.UnderlyingRawStorage.(*raft.RaftBackend).AppliedIndex()
}

func WaitForRaftApply(t testing.T, core *vault.TestClusterCore, index uint64) {
	t.Helper()

	backend := core.UnderlyingRawStorage.(*raft.RaftBackend)
	for i := 0; i < 30; i++ {
		if backend.AppliedIndex() >= index {
			return
		}

		time.Sleep(time.Second)
	}

	t.Fatalf("node did not apply index")
}

// AwaitLeader waits for one of the cluster's nodes to become leader.
func AwaitLeader(t testing.T, cluster *vault.TestCluster) (int, error) {
	timeout := time.Now().Add(30 * time.Second)
	for {
		if time.Now().After(timeout) {
			break
		}

		for i, core := range cluster.Cores {
			if core.Core.Sealed() {
				continue
			}

			isLeader, _, _, _ := core.Leader()
			if isLeader {
				return i, nil
			}
		}

		time.Sleep(time.Second)
	}

	return 0, fmt.Errorf("timeout waiting leader")
}

func GenerateDebugLogs(t testing.T, client *api.Client) chan struct{} {
	t.Helper()

	stopCh := make(chan struct{})
	ticker := time.NewTicker(time.Second)
	var err error

	go func() {
		for {
			select {
			case <-stopCh:
				ticker.Stop()
				stopCh <- struct{}{}
				return
			case <-ticker.C:
				err = client.Sys().Mount("foo", &api.MountInput{
					Type: "kv",
					Options: map[string]string{
						"version": "1",
					},
				})
				if err != nil {
					t.Fatal(err)
				}

				err = client.Sys().Unmount("foo")
				if err != nil {
					t.Fatal(err)
				}
			}
		}
	}()

	return stopCh
}

func VerifyRaftPeers(t testing.T, client *api.Client, expected map[string]bool) {
	t.Helper()

	resp, err := client.Logical().Read("sys/storage/raft/configuration")
	if err != nil {
		t.Fatalf("error reading raft config: %v", err)
	}

	if resp == nil || resp.Data == nil {
		t.Fatal("missing response data")
	}

	config, ok := resp.Data["config"].(map[string]interface{})
	if !ok {
		t.Fatal("missing config in response data")
	}

	servers, ok := config["servers"].([]interface{})
	if !ok {
		t.Fatal("missing servers in response data config")
	}

	// Iterate through the servers and remove the node found in the response
	// from the expected collection
	for _, s := range servers {
		server := s.(map[string]interface{})
		delete(expected, server["node_id"].(string))
	}

	// If the collection is non-empty, it means that the peer was not found in
	// the response.
	if len(expected) != 0 {
		t.Fatalf("failed to read configuration successfully, expected peers no found in configuration list: %v", expected)
	}
}

func TestMetricSinkProvider(gaugeInterval time.Duration) func(string) (*metricsutil.ClusterMetricSink, *metricsutil.MetricsHelper) {
	return func(clusterName string) (*metricsutil.ClusterMetricSink, *metricsutil.MetricsHelper) {
		inm := metrics.NewInmemSink(1000000*time.Hour, 2000000*time.Hour)
		clusterSink := metricsutil.NewClusterMetricSink(clusterName, inm)
		clusterSink.GaugeInterval = gaugeInterval
		return clusterSink, metricsutil.NewMetricsHelper(inm, false)
	}
}

func SysMetricsReq(client *api.Client, cluster *vault.TestCluster, unauth bool) (*SysMetricsJSON, error) {
	r := client.NewRequest("GET", "/v1/sys/metrics")
	if !unauth {
		r.Headers.Set("X-Vault-Token", cluster.RootToken)
	}
	var data SysMetricsJSON
	resp, err := client.RawRequestWithContext(context.Background(), r)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Response.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, errors.New("failed to unmarshal:" + err.Error())
	}
	return &data, nil
}

type SysMetricsJSON struct {
	Gauges []GaugeJSON `json:"Gauges"`
}

type GaugeJSON struct {
	Name   string                 `json:"Name"`
	Value  int                    `json:"Value"`
	Labels map[string]interface{} `json:"Labels"`
}
