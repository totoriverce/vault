package diagnose

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/physical"
)

const (
	success   string = "success"
	secretKey string = "diagnose"
	secretVal string = "diagnoseSecret"

	LatencyWarning    string        = "latency above 100 ms: "
	DirAccessErr      string        = "consul storage does not connect to local agent, but directly to server"
	AddrDNExistErr    string        = "config address does not exist: 127.0.0.1:8500 will be used"
	wrongRWValsPrefix string        = "Storage get and put gave wrong values: "
	latencyThreshold  time.Duration = time.Millisecond * 100
)

func EndToEndLatencyCheckWrite(ctx context.Context, uuid string, b physical.Backend) (time.Duration, error) {
	start := time.Now()
	err := b.Put(context.Background(), &physical.Entry{Key: secretKey, Value: []byte(secretVal)})
	duration := time.Since(start)
	if err != nil {
		return time.Duration(0), err
	}
	if duration > latencyThreshold {
		return duration, nil
	}
	return time.Duration(0), nil
}

func EndToEndLatencyCheckRead(ctx context.Context, uuid string, b physical.Backend) (time.Duration, error) {

	start := time.Now()
	val, err := b.Get(context.Background(), "diagnose")
	duration := time.Since(start)
	if err != nil {
		return time.Duration(0), err
	}
	if val == nil {
		return time.Duration(0), fmt.Errorf("no value found when reading generated data")
	}
	if val.Key != "diagnose" && string(val.Value) != "diagnose" {
		return time.Duration(0), fmt.Errorf(wrongRWValsPrefix+"expecting diagnose, but got %s, %s", val.Key, val.Value)
	}
	if duration > latencyThreshold {
		return duration, nil
	}
	return time.Duration(0), nil
}
func EndToEndLatencyCheckDelete(ctx context.Context, uuid string, b physical.Backend) (time.Duration, error) {

	start := time.Now()
	err := b.Delete(context.Background(), "diagnose")
	duration := time.Since(start)
	if err != nil {
		return time.Duration(0), err
	}
	if duration > latencyThreshold {
		return duration, nil
	}
	return time.Duration(0), nil
}

// ConsulDirectAccess verifies that consul is connecting to local agent,
// versus directly to a remote server. We can only assume that the local address
// is a server, not a client.
func ConsulDirectAccess(config map[string]string) string {
	configAddr, ok := config["address"]
	if !ok {
		return AddrDNExistErr
	}
	if !strings.Contains(configAddr, "localhost") && !strings.Contains(configAddr, "127.0.0.1") {
		return DirAccessErr
	}
	return ""
}
