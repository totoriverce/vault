// Code generated by sdkgen. DO NOT EDIT.

package mongodb

import (
	"context"

	"google.golang.org/grpc"
)

// MongoDB provides access to "mongodb" component of Yandex.Cloud
type MongoDB struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

// NewMongoDB creates instance of MongoDB
func NewMongoDB(g func(ctx context.Context) (*grpc.ClientConn, error)) *MongoDB {
	return &MongoDB{g}
}

// Backup gets BackupService client
func (m *MongoDB) Backup() *BackupServiceClient {
	return &BackupServiceClient{getConn: m.getConn}
}

// Cluster gets ClusterService client
func (m *MongoDB) Cluster() *ClusterServiceClient {
	return &ClusterServiceClient{getConn: m.getConn}
}

// Database gets DatabaseService client
func (m *MongoDB) Database() *DatabaseServiceClient {
	return &DatabaseServiceClient{getConn: m.getConn}
}

// ResourcePreset gets ResourcePresetService client
func (m *MongoDB) ResourcePreset() *ResourcePresetServiceClient {
	return &ResourcePresetServiceClient{getConn: m.getConn}
}

// User gets UserService client
func (m *MongoDB) User() *UserServiceClient {
	return &UserServiceClient{getConn: m.getConn}
}
