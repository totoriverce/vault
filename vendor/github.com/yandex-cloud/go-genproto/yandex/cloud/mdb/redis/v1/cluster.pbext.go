// Code generated by protoc-gen-goext. DO NOT EDIT.

package redis

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	config "github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/redis/v1/config"
	timeofday "google.golang.org/genproto/googleapis/type/timeofday"
)

func (m *Cluster) SetId(v string) {
	m.Id = v
}

func (m *Cluster) SetFolderId(v string) {
	m.FolderId = v
}

func (m *Cluster) SetCreatedAt(v *timestamp.Timestamp) {
	m.CreatedAt = v
}

func (m *Cluster) SetName(v string) {
	m.Name = v
}

func (m *Cluster) SetDescription(v string) {
	m.Description = v
}

func (m *Cluster) SetLabels(v map[string]string) {
	m.Labels = v
}

func (m *Cluster) SetEnvironment(v Cluster_Environment) {
	m.Environment = v
}

func (m *Cluster) SetMonitoring(v []*Monitoring) {
	m.Monitoring = v
}

func (m *Cluster) SetConfig(v *ClusterConfig) {
	m.Config = v
}

func (m *Cluster) SetNetworkId(v string) {
	m.NetworkId = v
}

func (m *Cluster) SetHealth(v Cluster_Health) {
	m.Health = v
}

func (m *Cluster) SetStatus(v Cluster_Status) {
	m.Status = v
}

func (m *Cluster) SetSharded(v bool) {
	m.Sharded = v
}

func (m *Cluster) SetMaintenanceWindow(v *MaintenanceWindow) {
	m.MaintenanceWindow = v
}

func (m *Cluster) SetPlannedOperation(v *MaintenanceOperation) {
	m.PlannedOperation = v
}

func (m *Monitoring) SetName(v string) {
	m.Name = v
}

func (m *Monitoring) SetDescription(v string) {
	m.Description = v
}

func (m *Monitoring) SetLink(v string) {
	m.Link = v
}

type ClusterConfig_RedisConfig = isClusterConfig_RedisConfig

func (m *ClusterConfig) SetRedisConfig(v ClusterConfig_RedisConfig) {
	m.RedisConfig = v
}

func (m *ClusterConfig) SetVersion(v string) {
	m.Version = v
}

func (m *ClusterConfig) SetRedisConfig_5_0(v *config.RedisConfigSet5_0) {
	m.RedisConfig = &ClusterConfig_RedisConfig_5_0{
		RedisConfig_5_0: v,
	}
}

func (m *ClusterConfig) SetRedisConfig_6_0(v *config.RedisConfigSet6_0) {
	m.RedisConfig = &ClusterConfig_RedisConfig_6_0{
		RedisConfig_6_0: v,
	}
}

func (m *ClusterConfig) SetResources(v *Resources) {
	m.Resources = v
}

func (m *ClusterConfig) SetBackupWindowStart(v *timeofday.TimeOfDay) {
	m.BackupWindowStart = v
}

func (m *ClusterConfig) SetAccess(v *Access) {
	m.Access = v
}

func (m *Shard) SetName(v string) {
	m.Name = v
}

func (m *Shard) SetClusterId(v string) {
	m.ClusterId = v
}

func (m *Host) SetName(v string) {
	m.Name = v
}

func (m *Host) SetClusterId(v string) {
	m.ClusterId = v
}

func (m *Host) SetZoneId(v string) {
	m.ZoneId = v
}

func (m *Host) SetSubnetId(v string) {
	m.SubnetId = v
}

func (m *Host) SetResources(v *Resources) {
	m.Resources = v
}

func (m *Host) SetRole(v Host_Role) {
	m.Role = v
}

func (m *Host) SetHealth(v Host_Health) {
	m.Health = v
}

func (m *Host) SetServices(v []*Service) {
	m.Services = v
}

func (m *Host) SetShardName(v string) {
	m.ShardName = v
}

func (m *Service) SetType(v Service_Type) {
	m.Type = v
}

func (m *Service) SetHealth(v Service_Health) {
	m.Health = v
}

func (m *Resources) SetResourcePresetId(v string) {
	m.ResourcePresetId = v
}

func (m *Resources) SetDiskSize(v int64) {
	m.DiskSize = v
}

func (m *Access) SetDataLens(v bool) {
	m.DataLens = v
}
