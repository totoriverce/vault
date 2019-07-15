// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// LaunchDbSystemDetails Used for creating a new DB system. Does not use backups or an existing database for the creation of the initial database.
type LaunchDbSystemDetails struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment the DB system  belongs in.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The availability domain where the DB system is located.
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the subnet the DB system is associated with.
	// **Subnet Restrictions:**
	// - For bare metal DB systems and for single node virtual machine DB systems, do not use a subnet that overlaps with 192.168.16.16/28.
	// - For Exadata and virtual machine 2-node RAC DB systems, do not use a subnet that overlaps with 192.168.128.0/20.
	// These subnets are used by the Oracle Clusterware private interconnect on the database instance.
	// Specifying an overlapping subnet will cause the private interconnect to malfunction.
	// This restriction applies to both the client subnet and the backup subnet.
	SubnetId *string `mandatory:"true" json:"subnetId"`

	// The shape of the DB system. The shape determines resources allocated to the DB system.
	// - For virtual machine shapes, the number of CPU cores and memory
	// - For bare metal and Exadata shapes, the number of CPU cores, memory, and storage
	// To get a list of shapes, use the ListDbSystemShapes operation.
	Shape *string `mandatory:"true" json:"shape"`

	// The public key portion of the key pair to use for SSH access to the DB system. Multiple public keys can be provided. The length of the combined keys cannot exceed 40,000 characters.
	SshPublicKeys []string `mandatory:"true" json:"sshPublicKeys"`

	// The hostname for the DB system. The hostname must begin with an alphabetic character, and
	// can contain alphanumeric characters and hyphens (-). The maximum length of the hostname is 16 characters for bare metal and virtual machine DB systems, and 12 characters for Exadata DB systems.
	// The maximum length of the combined hostname and domain is 63 characters.
	// **Note:** The hostname must be unique within the subnet. If it is not unique,
	// the DB system will fail to provision.
	Hostname *string `mandatory:"true" json:"hostname"`

	// The number of CPU cores to enable for a bare metal or Exadata DB system. The valid values depend on the specified shape:
	// - BM.DenseIO1.36 - Specify a multiple of 2, from 2 to 36.
	// - BM.DenseIO2.52 - Specify a multiple of 2, from 2 to 52.
	// - Exadata.Quarter1.84 - Specify a multiple of 2, from 22 to 84.
	// - Exadata.Half1.168 - Specify a multiple of 4, from 44 to 168.
	// - Exadata.Full1.336 - Specify a multiple of 8, from 88 to 336.
	// - Exadata.Quarter2.92 - Specify a multiple of 2, from 0 to 92.
	// - Exadata.Half2.184 - Specify a multiple of 4, from 0 to 184.
	// - Exadata.Full2.368 - Specify a multiple of 8, from 0 to 368.
	// This parameter is not used for virtual machine DB systems because virtual machine DB systems have a set number of cores for each shape.
	// For information about the number of cores for a virtual machine DB system shape, see Virtual Machine DB Systems (https://docs.cloud.oracle.com/Content/Database/Concepts/overview.htm#virtualmachine)
	CpuCoreCount *int `mandatory:"true" json:"cpuCoreCount"`

	DbHome *CreateDbHomeDetails `mandatory:"true" json:"dbHome"`

	// A Fault Domain is a grouping of hardware and infrastructure within an availability domain.
	// Fault Domains let you distribute your instances so that they are not on the same physical
	// hardware within a single availability domain. A hardware failure or maintenance
	// that affects one Fault Domain does not affect DB systems in other Fault Domains.
	// If you do not specify the Fault Domain, the system selects one for you. To change the Fault
	// Domain for a DB system, terminate it and launch a new DB system in the preferred Fault Domain.
	// If the node count is greater than 1, you can specify which Fault Domains these nodes will be distributed into.
	// The system assigns your nodes automatically to the Fault Domains you specify so that
	// no Fault Domain contains more than one node.
	// To get a list of Fault Domains, use the
	// ListFaultDomains operation in the
	// Identity and Access Management Service API.
	// Example: `FAULT-DOMAIN-1`
	FaultDomains []string `mandatory:"false" json:"faultDomains"`

	// The user-friendly name for the DB system. The name does not have to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the backup network subnet the DB system is associated with. Applicable only to Exadata DB systems.
	// **Subnet Restrictions:** See the subnet restrictions information for **subnetId**.
	BackupSubnetId *string `mandatory:"false" json:"backupSubnetId"`

	// The list of Network Security Group OCIDs (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) associated with this DB system.
	// A maximum of 5 allowed.
	NsgIds []string `mandatory:"false" json:"nsgIds"`

	// The list of Network Security Group OCIDs (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) associated with the backup network of this DB system.
	// Applicable only to Exadata DB systems.
	// A maximum of 5 allowed.
	BackupNetworkNsgIds []string `mandatory:"false" json:"backupNetworkNsgIds"`

	// The time zone to use for the DB system. For details, see DB System Time Zones (https://docs.cloud.oracle.com/Content/Database/References/timezones.htm).
	TimeZone *string `mandatory:"false" json:"timeZone"`

	// If true, Sparse Diskgroup is configured for Exadata dbsystem. If False, Sparse diskgroup is not configured.
	SparseDiskgroup *bool `mandatory:"false" json:"sparseDiskgroup"`

	// A domain name used for the DB system. If the Oracle-provided Internet and VCN
	// Resolver is enabled for the specified subnet, the domain name for the subnet is used
	// (do not provide one). Otherwise, provide a valid DNS domain name. Hyphens (-) are not permitted.
	Domain *string `mandatory:"false" json:"domain"`

	// The cluster name for Exadata and 2-node RAC virtual machine DB systems. The cluster name must begin with an an alphabetic character, and may contain hyphens (-). Underscores (_) are not permitted. The cluster name can be no longer than 11 characters and is not case sensitive.
	ClusterName *string `mandatory:"false" json:"clusterName"`

	// The percentage assigned to DATA storage (user data and database files).
	// The remaining percentage is assigned to RECO storage (database redo logs, archive logs, and recovery manager backups).
	// Specify 80 or 40. The default is 80 percent assigned to DATA storage. Not applicable for virtual machine DB systems.
	DataStoragePercentage *int `mandatory:"false" json:"dataStoragePercentage"`

	// Size (in GB) of the initial data volume that will be created and attached to a virtual machine DB system. You can scale up storage after provisioning, as needed. Note that the total storage size attached will be more than the amount you specify to allow for REDO/RECO space and software volume.
	InitialDataStorageSizeInGB *int `mandatory:"false" json:"initialDataStorageSizeInGB"`

	// The number of nodes to launch for a 2-node RAC virtual machine DB system.
	NodeCount *int `mandatory:"false" json:"nodeCount"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The Oracle Database Edition that applies to all the databases on the DB system.
	// Exadata DB systems and 2-node RAC DB systems require ENTERPRISE_EDITION_EXTREME_PERFORMANCE.
	DatabaseEdition LaunchDbSystemDetailsDatabaseEditionEnum `mandatory:"true" json:"databaseEdition"`

	// The type of redundancy configured for the DB system.
	// Normal is 2-way redundancy, recommended for test and development systems.
	// High is 3-way redundancy, recommended for production systems.
	DiskRedundancy LaunchDbSystemDetailsDiskRedundancyEnum `mandatory:"false" json:"diskRedundancy,omitempty"`

	// The Oracle license model that applies to all the databases on the DB system. The default is LICENSE_INCLUDED.
	LicenseModel LaunchDbSystemDetailsLicenseModelEnum `mandatory:"false" json:"licenseModel,omitempty"`
}

//GetCompartmentId returns CompartmentId
func (m LaunchDbSystemDetails) GetCompartmentId() *string {
	return m.CompartmentId
}

//GetFaultDomains returns FaultDomains
func (m LaunchDbSystemDetails) GetFaultDomains() []string {
	return m.FaultDomains
}

//GetDisplayName returns DisplayName
func (m LaunchDbSystemDetails) GetDisplayName() *string {
	return m.DisplayName
}

//GetAvailabilityDomain returns AvailabilityDomain
func (m LaunchDbSystemDetails) GetAvailabilityDomain() *string {
	return m.AvailabilityDomain
}

//GetSubnetId returns SubnetId
func (m LaunchDbSystemDetails) GetSubnetId() *string {
	return m.SubnetId
}

//GetBackupSubnetId returns BackupSubnetId
func (m LaunchDbSystemDetails) GetBackupSubnetId() *string {
	return m.BackupSubnetId
}

//GetNsgIds returns NsgIds
func (m LaunchDbSystemDetails) GetNsgIds() []string {
	return m.NsgIds
}

//GetBackupNetworkNsgIds returns BackupNetworkNsgIds
func (m LaunchDbSystemDetails) GetBackupNetworkNsgIds() []string {
	return m.BackupNetworkNsgIds
}

//GetShape returns Shape
func (m LaunchDbSystemDetails) GetShape() *string {
	return m.Shape
}

//GetTimeZone returns TimeZone
func (m LaunchDbSystemDetails) GetTimeZone() *string {
	return m.TimeZone
}

//GetSparseDiskgroup returns SparseDiskgroup
func (m LaunchDbSystemDetails) GetSparseDiskgroup() *bool {
	return m.SparseDiskgroup
}

//GetSshPublicKeys returns SshPublicKeys
func (m LaunchDbSystemDetails) GetSshPublicKeys() []string {
	return m.SshPublicKeys
}

//GetHostname returns Hostname
func (m LaunchDbSystemDetails) GetHostname() *string {
	return m.Hostname
}

//GetDomain returns Domain
func (m LaunchDbSystemDetails) GetDomain() *string {
	return m.Domain
}

//GetCpuCoreCount returns CpuCoreCount
func (m LaunchDbSystemDetails) GetCpuCoreCount() *int {
	return m.CpuCoreCount
}

//GetClusterName returns ClusterName
func (m LaunchDbSystemDetails) GetClusterName() *string {
	return m.ClusterName
}

//GetDataStoragePercentage returns DataStoragePercentage
func (m LaunchDbSystemDetails) GetDataStoragePercentage() *int {
	return m.DataStoragePercentage
}

//GetInitialDataStorageSizeInGB returns InitialDataStorageSizeInGB
func (m LaunchDbSystemDetails) GetInitialDataStorageSizeInGB() *int {
	return m.InitialDataStorageSizeInGB
}

//GetNodeCount returns NodeCount
func (m LaunchDbSystemDetails) GetNodeCount() *int {
	return m.NodeCount
}

//GetFreeformTags returns FreeformTags
func (m LaunchDbSystemDetails) GetFreeformTags() map[string]string {
	return m.FreeformTags
}

//GetDefinedTags returns DefinedTags
func (m LaunchDbSystemDetails) GetDefinedTags() map[string]map[string]interface{} {
	return m.DefinedTags
}

func (m LaunchDbSystemDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m LaunchDbSystemDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeLaunchDbSystemDetails LaunchDbSystemDetails
	s := struct {
		DiscriminatorParam string `json:"source"`
		MarshalTypeLaunchDbSystemDetails
	}{
		"NONE",
		(MarshalTypeLaunchDbSystemDetails)(m),
	}

	return json.Marshal(&s)
}

// LaunchDbSystemDetailsDatabaseEditionEnum Enum with underlying type: string
type LaunchDbSystemDetailsDatabaseEditionEnum string

// Set of constants representing the allowable values for LaunchDbSystemDetailsDatabaseEditionEnum
const (
	LaunchDbSystemDetailsDatabaseEditionStandardEdition                     LaunchDbSystemDetailsDatabaseEditionEnum = "STANDARD_EDITION"
	LaunchDbSystemDetailsDatabaseEditionEnterpriseEdition                   LaunchDbSystemDetailsDatabaseEditionEnum = "ENTERPRISE_EDITION"
	LaunchDbSystemDetailsDatabaseEditionEnterpriseEditionHighPerformance    LaunchDbSystemDetailsDatabaseEditionEnum = "ENTERPRISE_EDITION_HIGH_PERFORMANCE"
	LaunchDbSystemDetailsDatabaseEditionEnterpriseEditionExtremePerformance LaunchDbSystemDetailsDatabaseEditionEnum = "ENTERPRISE_EDITION_EXTREME_PERFORMANCE"
)

var mappingLaunchDbSystemDetailsDatabaseEdition = map[string]LaunchDbSystemDetailsDatabaseEditionEnum{
	"STANDARD_EDITION":                       LaunchDbSystemDetailsDatabaseEditionStandardEdition,
	"ENTERPRISE_EDITION":                     LaunchDbSystemDetailsDatabaseEditionEnterpriseEdition,
	"ENTERPRISE_EDITION_HIGH_PERFORMANCE":    LaunchDbSystemDetailsDatabaseEditionEnterpriseEditionHighPerformance,
	"ENTERPRISE_EDITION_EXTREME_PERFORMANCE": LaunchDbSystemDetailsDatabaseEditionEnterpriseEditionExtremePerformance,
}

// GetLaunchDbSystemDetailsDatabaseEditionEnumValues Enumerates the set of values for LaunchDbSystemDetailsDatabaseEditionEnum
func GetLaunchDbSystemDetailsDatabaseEditionEnumValues() []LaunchDbSystemDetailsDatabaseEditionEnum {
	values := make([]LaunchDbSystemDetailsDatabaseEditionEnum, 0)
	for _, v := range mappingLaunchDbSystemDetailsDatabaseEdition {
		values = append(values, v)
	}
	return values
}

// LaunchDbSystemDetailsDiskRedundancyEnum Enum with underlying type: string
type LaunchDbSystemDetailsDiskRedundancyEnum string

// Set of constants representing the allowable values for LaunchDbSystemDetailsDiskRedundancyEnum
const (
	LaunchDbSystemDetailsDiskRedundancyHigh   LaunchDbSystemDetailsDiskRedundancyEnum = "HIGH"
	LaunchDbSystemDetailsDiskRedundancyNormal LaunchDbSystemDetailsDiskRedundancyEnum = "NORMAL"
)

var mappingLaunchDbSystemDetailsDiskRedundancy = map[string]LaunchDbSystemDetailsDiskRedundancyEnum{
	"HIGH":   LaunchDbSystemDetailsDiskRedundancyHigh,
	"NORMAL": LaunchDbSystemDetailsDiskRedundancyNormal,
}

// GetLaunchDbSystemDetailsDiskRedundancyEnumValues Enumerates the set of values for LaunchDbSystemDetailsDiskRedundancyEnum
func GetLaunchDbSystemDetailsDiskRedundancyEnumValues() []LaunchDbSystemDetailsDiskRedundancyEnum {
	values := make([]LaunchDbSystemDetailsDiskRedundancyEnum, 0)
	for _, v := range mappingLaunchDbSystemDetailsDiskRedundancy {
		values = append(values, v)
	}
	return values
}

// LaunchDbSystemDetailsLicenseModelEnum Enum with underlying type: string
type LaunchDbSystemDetailsLicenseModelEnum string

// Set of constants representing the allowable values for LaunchDbSystemDetailsLicenseModelEnum
const (
	LaunchDbSystemDetailsLicenseModelLicenseIncluded     LaunchDbSystemDetailsLicenseModelEnum = "LICENSE_INCLUDED"
	LaunchDbSystemDetailsLicenseModelBringYourOwnLicense LaunchDbSystemDetailsLicenseModelEnum = "BRING_YOUR_OWN_LICENSE"
)

var mappingLaunchDbSystemDetailsLicenseModel = map[string]LaunchDbSystemDetailsLicenseModelEnum{
	"LICENSE_INCLUDED":       LaunchDbSystemDetailsLicenseModelLicenseIncluded,
	"BRING_YOUR_OWN_LICENSE": LaunchDbSystemDetailsLicenseModelBringYourOwnLicense,
}

// GetLaunchDbSystemDetailsLicenseModelEnumValues Enumerates the set of values for LaunchDbSystemDetailsLicenseModelEnum
func GetLaunchDbSystemDetailsLicenseModelEnumValues() []LaunchDbSystemDetailsLicenseModelEnum {
	values := make([]LaunchDbSystemDetailsLicenseModelEnum, 0)
	for _, v := range mappingLaunchDbSystemDetailsLicenseModel {
		values = append(values, v)
	}
	return values
}
