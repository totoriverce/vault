// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMetricLabelerAuditSink_Label ensures we always get the right label based
// on the input value of the error.
func TestMetricLabelerAuditSink_Label(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err      error
		expected string
	}{
		"nil": {
			err:      nil,
			expected: "vault.audit.sink.success",
		},
		"error": {
			err:      errors.New("I am an error"),
			expected: "vault.audit.sink.failure",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := &MetricLabelerAuditSink{}
			result := m.Label(nil, tc.err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestMetricLabelerAuditFallback_Label ensures we always get the right label based
// on the input value of the error for fallback devices.
func TestMetricLabelerAuditFallback_Label(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err      error
		expected string
	}{
		"nil": {
			err:      nil,
			expected: "vault.audit.fallback.success",
		},
		"error": {
			err:      errors.New("I am an error"),
			expected: "vault.audit.sink.failure",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := &MetricLabelerAuditFallback{}
			result := m.Label(nil, tc.err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
