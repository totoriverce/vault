// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"github.com/oracle/oci-go-sdk/common"
)

// CreateInstancePoolDetails The data to create an instance pool.
type CreateInstancePoolDetails struct {

	// The OCID of the compartment containing the instance pool
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the instance configuration associated with the instance pool.
	InstanceConfigurationId *string `mandatory:"true" json:"instanceConfigurationId"`

	// The placement configurations for the instance pool. Provide one placement configuration for
	// each availability domain.
	// To use the instance pool with a regional subnet, provide a placement configuration for
	// each availability domain, and include the regional subnet in each placement
	// configuration.
	PlacementConfigurations []CreateInstancePoolPlacementConfigurationDetails `mandatory:"true" json:"placementConfigurations"`

	// The number of instances that should be in the instance pool.
	Size *int `mandatory:"true" json:"size"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The user-friendly name.  Does not have to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The load balancers to attach to the instance pool.
	LoadBalancers []AttachLoadBalancerDetails `mandatory:"false" json:"loadBalancers"`
}

func (m CreateInstancePoolDetails) String() string {
	return common.PointerString(m)
}
