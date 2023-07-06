// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

const (
	tags = "policymaker,monitoring"
	vhostTopicJSON = `{"vhostOne":{"exchangeOneOne":{"write":".*","read":".*"},"exchangeOneTwo":{"write":".*","read":".*" }}}`
)

func TestBackend_Roles_Static(t *testing.T) {
	cleanup, connectionURI := prepareRabbitMQTestContainer(t)
	defer cleanup()

	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"connection_uri":    connectionURI,
		"username":          "guest",
		"password":          "guest",
		"username_template": "",
	}
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
	}

	tests := []struct {
		name         string
		args         map[string]interface{}
		wantErr      bool
		expectedResp map[string]interface{}
	}{
		{
			name: "invalid static role with no tags or vhost permissisons",
			args: map[string]interface{}{
				"username": "tester",
				"vhost_topics": vhostTopicJSON,
			},
			wantErr: true,
		},
		{
			name: "invalid static role with no username",
			args: map[string]interface{}{
				"tags": tags,
				"vhost_topics": vhostTopicJSON,
			},
			wantErr: true,
		},
		{
			name: "valid static role with tags",
			args: map[string]interface{}{
				"tags": tags,
				"username": "tester",
				"rotation_period": 3,
			},
			wantErr: false,
			expectedResp: map[string]interface{}{
				"tags": tags,
				"username": "tester",
				"rotation_period": 3.0,
			},
		},
		{
			name: "valid static role with revoke on delete",
			args: map[string]interface{}{
				"tags": tags,
				"username": "tester",
				"revoke_user_on_delete": true,
			},
			wantErr: false,
			expectedResp: map[string]interface{}{
				"tags": tags,
				"username": "tester",
				"revoke_user_on_delete": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &logical.Request{
				Operation: logical.CreateOperation,
				Path:      "static-roles/test",
				Storage:   config.StorageView,
				Data:      tt.args,
			}
			// Create the role
			resp, err := b.HandleRequest(context.Background(), req)
			if tt.wantErr {
				assert.True(t, resp.IsError(), "expected error")
				return
			}
			assert.False(t, resp.IsError())
			assert.Nil(t, err)

			// Read the role
			req.Operation = logical.ReadOperation
			resp, err = b.HandleRequest(context.Background(), req)
			assert.False(t, resp.IsError())
			assert.Nil(t, err)
			for k, v := range tt.expectedResp {
				assert.Equal(t, v, resp.Data[k])
			}

			// Delete the role
			req.Operation = logical.DeleteOperation
			resp, err = b.HandleRequest(context.Background(), req)
			assert.False(t, resp.IsError())
			assert.Nil(t, err)
		})
	}
}
