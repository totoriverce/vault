package okta

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/sstarcher/go-okta"
	"time"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `config`,
		Fields: map[string]*framework.FieldSchema{
			"organization": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Okta organization to authenticate against",
			},
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Okta admin API token",
			},
			"base_url": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The API endpoint to use. Useful if you
are using Okta development accounts.`,
			},
			"ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Duration after which authentication will be expired`,
			},
			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Maximum duration after which authentication will be expired`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.CreateOperation: b.pathConfigWrite,
			logical.UpdateOperation: b.pathConfigWrite,
		},

		ExistenceCheck: b.pathConfigExistenceCheck,

		HelpSynopsis: pathConfigHelp,
	}
}

// Config returns the configuration for this backend.
func (b *backend) Config(s logical.Storage) (*ConfigEntry, error) {
	entry, err := s.Get("config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result ConfigEntry
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (b *backend) pathConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	cfg, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"Org":     cfg.Org,
			"BaseURL": cfg.BaseURL,
			"TTL":     cfg.TTL,
			"MaxTTL":  cfg.MaxTTL,
		},
	}

	return resp, nil
}

func (b *backend) pathConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	org := d.Get("organization").(string)
	cfg, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	// Due to the existence check, entry will only be nil if it's a create
	// operation, so just create a new one
	if cfg == nil {
		cfg = &ConfigEntry{
			Org: org,
		}
	}

	token, ok := d.GetOk("token")
	if ok {
		cfg.Token = token.(string)
	} else if req.Operation == logical.CreateOperation {
		cfg.Token = d.Get("token").(string)
	}

	baseURL, ok := d.GetOk("base_url")
	if ok {
		baseURLString := baseURL.(string)
		if len(baseURLString) != 0 {
			_, err = url.Parse(baseURLString)
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("Error parsing given base_url: %s", err)), nil
			}
			cfg.BaseURL = baseURLString
		}
	} else if req.Operation == logical.CreateOperation {
		cfg.BaseURL = d.Get("base_url").(string)
	}

	var ttl time.Duration
	ttlRaw, ok := d.GetOk("ttl")
	if !ok || len(ttlRaw.(string)) == 0 {
		ttl = 0
	} else {
		ttl, err = time.ParseDuration(ttlRaw.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid 'ttl':%s", err)), nil
		}
	}

	var maxTTL time.Duration
	maxTTLRaw, ok := d.GetOk("max_ttl")
	if !ok || len(maxTTLRaw.(string)) == 0 {
		maxTTL = 0
	} else {
		maxTTL, err = time.ParseDuration(maxTTLRaw.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid 'max_ttl':%s", err)), nil
		}
	}

	cfg.TTL = ttl
	cfg.MaxTTL = maxTTL

	jsonCfg, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(jsonCfg); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigExistenceCheck(
	req *logical.Request, d *framework.FieldData) (bool, error) {
	cfg, err := b.Config(req.Storage)
	if err != nil {
		return false, err
	}

	return cfg != nil, nil
}

// OktaClient creates a basic okta client connection
func (c *ConfigEntry) OktaClient() *okta.Client {
	client := okta.NewClient(c.Org)
	if c.BaseURL != "" {
		client.Url = c.BaseURL
	}

	if c.Token != "" {
		client.ApiToken = c.Token
	}

	return client
}

// ConfigEntry for Okta
type ConfigEntry struct {
	Org     string        `json:"organization"`
	Token   string        `json:"token"`
	BaseURL string        `json:"base_url"`
	TTL     time.Duration `json:"ttl"`
	MaxTTL  time.Duration `json:"max_ttl"`
}

const pathConfigHelp = `
This endpoint allows you to configure the Okta and its
configuration options.

The Okta organization are the characters at the front of the URL for Okta.
Example https://ORG.okta.com
`
