// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

// NOTE: this file was copied from
// https://github.com/hashicorp/vault/blob/main/sdk/helper/consts/plugin_runtime_types.go
// Any changes made should be made to both files at the same time.

import "fmt"

var PluginRuntimeTypes = []PluginRuntimeType{
	PluginRuntimeTypeUnknown,
	PluginRuntimeTypeContainer,
}

type PluginRuntimeType uint32

// This is a list of PluginRuntimeTypes used by Vault.
const (
	PluginRuntimeTypeUnknown PluginRuntimeType = iota
	PluginRuntimeTypeContainer
)

func (r PluginRuntimeType) String() string {
	switch r {
	case PluginRuntimeTypeContainer:
		return "container"
	case PluginRuntimeTypeUnknown:
		fallthrough
	default:
		return "unsupported"
	}
}

func ParsePluginRuntimeType(PluginRuntimeType string) (PluginRuntimeType, error) {
	switch PluginRuntimeType {
	case "container":
		return PluginRuntimeTypeContainer, nil
	default:
		return PluginRuntimeTypeUnknown, fmt.Errorf("%q is not a supported plugin runtime type", PluginRuntimeType)
	}
}
