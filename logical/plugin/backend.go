package plugin

import (
	"net/rpc"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/pb"
)

// BackendPlugin is the plugin.Plugin implementation
type BackendPlugin struct {
	Factory      func(*logical.BackendConfig) (logical.Backend, error)
	metadataMode bool
}

// Server gets called when on plugin.Serve()
func (b *BackendPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &backendPluginServer{factory: b.Factory, broker: broker}, nil
}

// Client gets called on plugin.NewClient()
func (b BackendPlugin) Client(broker *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &backendPluginClient{client: c, broker: broker, metadataMode: b.metadataMode}, nil
}

func (b BackendPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterBackendServer(s, &backendGRPCPluginServer{
		broker:  broker,
		factory: b.Factory,
	})
	return nil
}

func (p *BackendPlugin) GRPCClient(broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &backendGRPCPluginClient{
		client: pb.NewBackendClient(c),
		broker: broker,
	}, nil
}
