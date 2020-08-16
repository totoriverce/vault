// Code generated by protoc-gen-goext. DO NOT EDIT.

package mysql

import (
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
)

func (m *User) SetName(v string) {
	m.Name = v
}

func (m *User) SetClusterId(v string) {
	m.ClusterId = v
}

func (m *User) SetPermissions(v []*Permission) {
	m.Permissions = v
}

func (m *User) SetGlobalPermissions(v []GlobalPermission) {
	m.GlobalPermissions = v
}

func (m *User) SetConnectionLimits(v *ConnectionLimits) {
	m.ConnectionLimits = v
}

func (m *User) SetAuthenticationPlugin(v AuthPlugin) {
	m.AuthenticationPlugin = v
}

func (m *Permission) SetDatabaseName(v string) {
	m.DatabaseName = v
}

func (m *Permission) SetRoles(v []Permission_Privilege) {
	m.Roles = v
}

func (m *ConnectionLimits) SetMaxQuestionsPerHour(v *wrappers.Int64Value) {
	m.MaxQuestionsPerHour = v
}

func (m *ConnectionLimits) SetMaxUpdatesPerHour(v *wrappers.Int64Value) {
	m.MaxUpdatesPerHour = v
}

func (m *ConnectionLimits) SetMaxConnectionsPerHour(v *wrappers.Int64Value) {
	m.MaxConnectionsPerHour = v
}

func (m *ConnectionLimits) SetMaxUserConnections(v *wrappers.Int64Value) {
	m.MaxUserConnections = v
}

func (m *UserSpec) SetName(v string) {
	m.Name = v
}

func (m *UserSpec) SetPassword(v string) {
	m.Password = v
}

func (m *UserSpec) SetPermissions(v []*Permission) {
	m.Permissions = v
}

func (m *UserSpec) SetGlobalPermissions(v []GlobalPermission) {
	m.GlobalPermissions = v
}

func (m *UserSpec) SetConnectionLimits(v *ConnectionLimits) {
	m.ConnectionLimits = v
}

func (m *UserSpec) SetAuthenticationPlugin(v AuthPlugin) {
	m.AuthenticationPlugin = v
}
