// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/logical"
)

// ProcessManual will attempt to create an (audit) event with the specified data
// and manually iterate over the supplied nodes calling Process on each.
// Order of IDs in the NodeID slice determines the order they are processed.
// (Audit) Event will be of RequestType (as opposed to ResponseType).
// The last node must be a sink node (eventlogger.NodeTypeSink).
func ProcessManual(ctx context.Context, data *logical.LogInput, ids []eventlogger.NodeID, nodes map[eventlogger.NodeID]eventlogger.Node) error {
	// Create an audit event.
	a, err := NewEvent(RequestType)
	if err != nil {
		return err
	}

	// Insert the data into the audit event.
	a.Data = data

	// Create an eventlogger event with the audit event as the payload.
	e := &eventlogger.Event{
		Type:      eventlogger.EventType(event.AuditType.String()),
		CreatedAt: time.Now(),
		Formatted: make(map[string][]byte),
		Payload:   a,
	}

	var lastSeen eventlogger.NodeType

	// Process nodes data order, updating the event with the result.
	// This means we *should* do:
	// 1. formatter
	// 2. sink
	for _, id := range ids {
		node, ok := nodes[id]
		if !ok {
			return fmt.Errorf("node not found: %v", id)
		}
		e, err = node.Process(ctx, e)
		if err != nil {
			return err
		}

		// Track the last node we have processed, as we should end with a sink.
		lastSeen = node.Type()
	}

	if lastSeen != eventlogger.NodeTypeSink {
		return fmt.Errorf("last node must be a sink: %v", lastSeen)
	}

	return nil
}
