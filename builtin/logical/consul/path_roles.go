package consul

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},
	}
}

func pathRoles() *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},

			"policy": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Policy document, base64 encoded. Required
for 'client' tokens. Required for Consul pre-1.4`,
			},

			"policies": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `List of policies attached to the token. Required
for Consul 1.4 or above`,
			},

			"token_type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "client",
				Description: `Which type of token to create: 'client'
or 'management'. If a 'management' token,
the "policy" parameter is not required.
Defaults to 'client'.`,
			},

			"lease": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: "Lease time of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   pathRolesRead,
			logical.UpdateOperation: pathRolesWrite,
			logical.DeleteOperation: pathRolesDelete,
		},
	}
}

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "policy/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, "policy/"+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	if result.TokenType == "" {
		result.TokenType = "client"
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"lease":      int64(result.Lease.Seconds()),
			"token_type": result.TokenType,
		},
	}
	if result.Policy != "" {
		resp.Data["policy"] = base64.StdEncoding.EncodeToString([]byte(result.Policy))
	}
	if len(result.Policies) > 0 {
		resp.Data["policies"] = result.Policies
	}
	return resp, nil
}

func pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	tokenType := d.Get("token_type").(string)
	policy := d.Get("policy").(string)
	policies := d.Get("policies").([]string)
	if len(policies) == 0 {
		switch tokenType {
		case "client":
		case "management":
		default:
			return logical.ErrorResponse(
				"token_type must be \"client\" or \"management\""), nil
		}
	}

	if policy != "" && len(policies) > 0 {
		return logical.ErrorResponse(
			"Use either a policy document, or a list of policies, depending on your Consul version"), nil
	}

	name := d.Get("name").(string)

	var policyRaw []byte
	var err error
	if len(policies) == 0 {
		if tokenType != "management" {
			if policy == "" {
				return logical.ErrorResponse(
					"policy cannot be empty when not using management tokens"), nil
			}
		}
	}
	policyRaw, err = base64.StdEncoding.DecodeString(d.Get("policy").(string))
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error decoding policy base64: %s", err)), nil
	}

	var lease time.Duration
	leaseParamRaw, ok := d.GetOk("lease")
	if ok {
		lease = time.Second * time.Duration(leaseParamRaw.(int))
	}

	entry, err := logical.StorageEntryJSON("policy/"+name, roleConfig{
		Policy:    string(policyRaw),
		Policies:  []string(policies),
		Lease:     lease,
		TokenType: tokenType,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if err := req.Storage.Delete(ctx, "policy/"+name); err != nil {
		return nil, err
	}
	return nil, nil
}

type roleConfig struct {
	Policy    string        `json:"policy"`
	Policies  []string      `json:"policies"`
	Lease     time.Duration `json:"lease"`
	TokenType string        `json:"token_type"`
}
