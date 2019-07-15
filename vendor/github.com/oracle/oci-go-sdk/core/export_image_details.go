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
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// ExportImageDetails The destination details for the image export.
// Set `destinationType` to `objectStorageTuple`
// and use ExportImageViaObjectStorageTupleDetails
// when specifying the namespace, bucket name, and object name.
// Set `destinationType` to `objectStorageUri` and
// use ExportImageViaObjectStorageUriDetails
// when specifying the Object Storage URL.
type ExportImageDetails interface {
}

type exportimagedetails struct {
	JsonData        []byte
	DestinationType string `json:"destinationType"`
}

// UnmarshalJSON unmarshals json
func (m *exportimagedetails) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalerexportimagedetails exportimagedetails
	s := struct {
		Model Unmarshalerexportimagedetails
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.DestinationType = s.Model.DestinationType

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *exportimagedetails) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.DestinationType {
	case "objectStorageUri":
		mm := ExportImageViaObjectStorageUriDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "objectStorageTuple":
		mm := ExportImageViaObjectStorageTupleDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

func (m exportimagedetails) String() string {
	return common.PointerString(m)
}
