// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Identity and Access Management Service API
//
// APIs for managing users, groups, compartments, and policies.
//

package identity

import (
	"github.com/oracle/oci-go-sdk/common"
)

// CreateUserDetails The representation of CreateUserDetails
type CreateUserDetails struct {

	// The OCID of the tenancy containing the user.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The name you assign to the user during creation. This is the user's login for the Console.
	// The name must be unique across all users in the tenancy and cannot be changed.
	Name *string `mandatory:"true" json:"name"`

	// The description you assign to the user during creation. Does not have to be unique, and it's changeable.
	Description *string `mandatory:"true" json:"description"`

	// The email you assign to the user. Has to be unique across the tenancy.
	Email *string `mandatory:"false" json:"email"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m CreateUserDetails) String() string {
	return common.PointerString(m)
}
