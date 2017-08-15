package plugin

import (
	"fmt"
	"net/rpc"
	"sync"

	"github.com/hashicorp/vault/logical"
	bplugin "github.com/hashicorp/vault/logical/plugin"
)

// Factory returns a configured plugin logical.Backend.
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	_, ok := conf.Config["plugin_name"]
	if !ok {
		return nil, fmt.Errorf("plugin_name not provided")
	}
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}

	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend returns an instance of the backend, either as a plugin if external
// or as a concrete implementation if builtin, casted as logical.Backend.
func Backend(conf *logical.BackendConfig) (logical.Backend, error) {
	var b backend
	name := conf.Config["plugin_name"]
	sys := conf.System

	raw, err := bplugin.NewBackend(name, sys, conf.Logger)
	if err != nil {
		return nil, err
	}
	b.Backend = raw
	b.config = conf

	return &b, nil
}

// backend is a thin wrapper around plugin.BackendPluginClient
type backend struct {
	logical.Backend
	sync.Mutex

	config *logical.BackendConfig
}

func (b *backend) reloadBackend() error {
	b.Lock()
	defer b.Unlock()

	pluginName := b.config.Config["plugin_name"]
	b.Logger().Trace("plugin: reloading plugin backend", "plugin", pluginName)

	nb, err := bplugin.NewBackend(pluginName, b.config.System)
	if err != nil {
		return err
	}
	err = nb.Setup(b.config)
	if err != nil {
		return err
	}
	b.Backend = nb

	return nil
}

// HandleRequest is a thin wrapper implementation of HandleRequest that includes automatic plugin reload.
func (b *backend) HandleRequest(req *logical.Request) (*logical.Response, error) {
	resp, err := b.Backend.HandleRequest(req)
	if err == rpc.ErrShutdown {
		// Reload plugin if it's an rpc.ErrShutdown
		err := b.reloadBackend()
		if err != nil {
			return nil, err
		}

		// Try handle request once more
		return b.Backend.HandleRequest(req)
	}
	return resp, err
}

// HandleExistenceCheck is a thin wrapper implementation of HandleRequest that includes automatic plugin reload.
func (b *backend) HandleExistenceCheck(req *logical.Request) (bool, bool, error) {
	checkFound, exists, err := b.Backend.HandleExistenceCheck(req)
	if err == rpc.ErrShutdown {
		// Reload plugin if it's an rpc.ErrShutdown
		err := b.reloadBackend()
		if err != nil {
			return false, false, err
		}

		// Try handle request once more
		return b.Backend.HandleExistenceCheck(req)
	}
	return checkFound, exists, err
}
