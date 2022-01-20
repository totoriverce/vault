package dbplugin

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// Serve is called from within a plugin and wraps the provided
// Database implementation in a databasePluginRPCServer object and starts a
// RPC server.
func Serve(db Database) {
	plugin.Serve(ServeConfig(db))
}

func ServeConfig(db Database) *plugin.ServeConfig {
	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// pluginSets is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		5: {
			"database": &GRPCDatabasePlugin{
				Impl: db,
			},
		},
	}

	conf := &plugin.ServeConfig{
		HandshakeConfig:  HandshakeConfig,
		VersionedPlugins: pluginSets,
		GRPCServer:       plugin.DefaultGRPCServer,
	}

	return conf
}

func ServeMultiplex(factory func() (Database, error)) {
	plugin.Serve(ServeConfigMultiplex(factory))
}

func ServeConfigMultiplex(factory func() (Database, error)) *plugin.ServeConfig {
	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// pluginSets is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		6: {
			"database": &GRPCDatabasePlugin{
				FactoryFunc: factory,
			},
		},
	}

	conf := &plugin.ServeConfig{
		HandshakeConfig:  HandshakeConfig,
		VersionedPlugins: pluginSets,
		GRPCServer:       plugin.DefaultGRPCServer,
	}

	return conf
}
