// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/logical"
)

// version defines the version of audit events.
const version = "v0.1"

// Audit subtypes.
const (
	RequestType  subtype = "AuditRequest"
	ResponseType subtype = "AuditResponse"
)

// Audit formats.
const (
	JSONFormat  format = "json"
	JSONxFormat format = "jsonx"
)

// Check AuditEvent implements the timeProvider at compile time.
var _ timeProvider = (*AuditEvent)(nil)

// AuditEvent is the audit event.
type AuditEvent struct {
	ID        string            `json:"id"`
	Version   string            `json:"version"`
	Subtype   subtype           `json:"subtype"` // the subtype of the audit event.
	Timestamp time.Time         `json:"timestamp"`
	Data      *logical.LogInput `json:"data"`
}

// format defines types of format audit events support.
type format string

// subtype defines the type of audit event.
type subtype string

// NewEvent should be used to create an audit event. The subtype field is needed
// for audit events. It will generate an ID if no ID is supplied. Supported
// options: WithID, WithNow.
func NewEvent(s subtype, opt ...Option) (*AuditEvent, error) {
	// Get the default options
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("error applying options: %w", err)
	}

	if opts.withID == "" {
		var err error

		opts.withID, err = event.NewID(string(event.AuditType))
		if err != nil {
			return nil, fmt.Errorf("error creating ID for event: %w", err)
		}
	}

	audit := &AuditEvent{
		ID:        opts.withID,
		Timestamp: opts.withNow,
		Version:   version,
		Subtype:   s,
	}

	if err := audit.validate(); err != nil {
		return nil, err
	}
	return audit, nil
}

// validate attempts to ensure the audit event in its present state is valid.
func (a *AuditEvent) validate() error {
	if a == nil {
		return fmt.Errorf("event is nil: %w", event.ErrInvalidParameter)
	}

	if a.ID == "" {
		return fmt.Errorf("missing ID: %w", event.ErrInvalidParameter)
	}

	if a.Version != version {
		return fmt.Errorf("event version unsupported: %w", event.ErrInvalidParameter)
	}

	if a.Timestamp.IsZero() {
		return fmt.Errorf("event timestamp cannot be the zero time instant: %w", event.ErrInvalidParameter)
	}

	err := a.Subtype.validate()
	if err != nil {
		return err
	}

	return nil
}

// validate ensures that subtype is one of the set of allowed event subtypes.
func (t subtype) validate() error {
	switch t {
	case RequestType, ResponseType:
		return nil
	default:
		return fmt.Errorf("invalid event subtype %q: %w", t, event.ErrInvalidParameter)
	}
}

// validate ensures that format is one of the set of allowed event formats.
func (f format) validate() error {
	switch f {
	case JSONFormat, JSONxFormat:
		return nil
	default:
		return fmt.Errorf("invalid format %q: %w", f, event.ErrInvalidParameter)
	}
}

// String returns the string version of a format.
func (f format) String() string {
	return string(f)
}

// MetricTag returns a tag corresponding to this subtype to include in metrics.
// If a tag cannot be found the value is returned 'as-is' in string format.
func (t subtype) MetricTag() string {
	switch t {
	case RequestType:
		return "log_request"
	case ResponseType:
		return "log_response"
	}

	return t.String()
}

// String returns the subtype as a human-readable string.
func (t subtype) String() string {
	switch t {
	case RequestType:
		return "request"
	case ResponseType:
		return "response"
	}

	return string(t)
}

// formattedTime returns the UTC time the AuditEvent was created in the RFC3339Nano
// format (which removes trailing zeros from the seconds field).
func (a *AuditEvent) formattedTime() string {
	return a.Timestamp.UTC().Format(time.RFC3339Nano)
}
