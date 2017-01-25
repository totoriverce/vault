package transit

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// BatchRewrapItemRequest represents an item in the batch rewrap
// request
type BatchRewrapItemRequest struct {
	// Context for key derivation. This is required for derived keys.
	Context string `json:"context" structs:"context" mapstructure:"context"`

	// DecodedContext, for internal use, which is the base64 decoded version of
	// the Context field
	DecodedContext []byte `json:"decoded_context" structs:"decoded_context" mapstructure:"decoded_context"`

	// Ciphertext which needs rewrap
	Ciphertext string `json:"ciphertext" structs:"ciphertext" mapstructure:"ciphertext"`

	// Nonce to be used when v1 convergent encryption is used
	Nonce string `json:"nonce" structs:"nonce" mapstructure:"nonce"`

	// DecodedNonce, for internal use, which is the base64 decoded version of
	// the Nonce field
	DecodedNonce []byte `json:"decoded_nonce" structs:"decoded_nonce" mapstructure:"decoded_nonce"`
}

// BatchRewrapItemResponse represents an item in the batch rewrap
// response
type BatchRewrapItemResponse struct {
	// Ciphertext is a rewrapped version of the ciphertext in the corresponding
	// batch request item
	Ciphertext string `json:"ciphertext" structs:"ciphertext" mapstructure:"ciphertext"`

	// Error, if set represents a failure encountered while encrypting rewrapping a
	// corresponding batch request item
	Error string `json:"error" structs:"error" mapstructure:"error"`
}

func (b *backend) pathRewrap() *framework.Path {
	return &framework.Path{
		Pattern: "rewrap/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"ciphertext": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Ciphertext value to rewrap",
			},

			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Base64 encoded context for key derivation. Required for derived keys.",
			},

			"nonce": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Nonce for when convergent encryption is used",
			},

			"batch": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Base64 encoded list of items to be rewrapped in a single batch. When this
parameter is set, if the parameters 'ciphertext', 'context' and 'nonce' are
also set, they will be ignored. JSON format for the input goes like this:
[
  {
    "context": "context1",
    "ciphertext": "vault:v1:/DupSiSbX/ATkGmKAmhqD0tvukByrx6gmps7dVI="
  },
  {
    "context": "context2",
    "ciphertext": "vault:v1:XjsPWPjqPrBi1N2Ms2s1QM798YyFWnO4TR4lsFA="
  },
  ...
]`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRewrapWrite,
		},

		HelpSynopsis:    pathRewrapHelpSyn,
		HelpDescription: pathRewrapHelpDesc,
	}
}

func (b *backend) pathRewrapWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	batchInputRaw := d.Get("batch").(string)
	var batchInputItems []BatchRewrapItemRequest
	var batchInput []byte
	var err error
	if len(batchInputRaw) != 0 {
		batchInput, err = base64.StdEncoding.DecodeString(batchInputRaw)
		if err != nil {
			return logical.ErrorResponse("failed to base64-decode batch input"), logical.ErrInvalidRequest
		}

		if err := jsonutil.DecodeJSON([]byte(batchInput), &batchInputItems); err != nil {
			return nil, fmt.Errorf("invalid input: %v", err)
		}
	} else {
		ciphertext := d.Get("ciphertext").(string)
		if len(ciphertext) == 0 {
			return logical.ErrorResponse("missing ciphertext to decrypt"), logical.ErrInvalidRequest
		}

		batchInputItems = make([]BatchRewrapItemRequest, 1)
		batchInputItems[0] = BatchRewrapItemRequest{
			Ciphertext: ciphertext,
			Context:    d.Get("context").(string),
			Nonce:      d.Get("nonce").(string),
		}
	}

	batchResponseItems := make([]BatchRewrapItemResponse, len(batchInputItems))
	if len(batchInputItems) == 0 {
		return logical.ErrorResponse("missing input to process"), logical.ErrInvalidRequest
	}

	contextSet := batchInputItems[0].Context != ""

	for i, item := range batchInputItems {
		if (item.Context == "" && contextSet) || (item.Context != "" && !contextSet) {
			return logical.ErrorResponse("context should be set either in all the request blocks or in none"), logical.ErrInvalidRequest
		}

		if item.Ciphertext == "" {
			batchResponseItems[i] = BatchRewrapItemResponse{
				Error: "missing ciphertext to decrypt",
			}
			continue
		}

		if len(item.Context) != 0 {
			batchInputItems[i].DecodedContext, err = base64.StdEncoding.DecodeString(item.Context)
			if err != nil {
				batchResponseItems[i] = BatchRewrapItemResponse{
					Error: "failed to base64-decode context",
				}
				continue
			}
		}

		if len(item.Nonce) != 0 {
			batchInputItems[i].DecodedNonce, err = base64.StdEncoding.DecodeString(item.Nonce)
			if err != nil {
				batchResponseItems[i] = BatchRewrapItemResponse{
					Error: "failed to base64-decode nonce",
				}
			}
			continue
		}
	}

	// Get the policy
	p, lock, err := b.lm.GetPolicyShared(req.Storage, d.Get("name").(string))
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("policy not found"), logical.ErrInvalidRequest
	}

	for i, item := range batchInputItems {
		if batchResponseItems[i].Error != "" {
			continue
		}

		plaintext, err := p.Decrypt(item.DecodedContext, item.DecodedNonce, item.Ciphertext)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
			case errutil.InternalError:
				return nil, err
			default:
				return nil, err
			}
		}

		if plaintext == "" {
			batchResponseItems[i] = BatchRewrapItemResponse{
				Error: "empty plaintext returned during rewrap",
			}
			continue
		}

		ciphertext, err := p.Encrypt(item.DecodedContext, item.DecodedNonce, plaintext)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
			case errutil.InternalError:
				return nil, err
			default:
				return nil, err
			}
		}

		if ciphertext == "" {
			batchResponseItems[i] = BatchRewrapItemResponse{
				Error: "empty ciphertext returned",
			}
			continue
		}

		batchResponseItems[i] = BatchRewrapItemResponse{
			Ciphertext: ciphertext,
		}
	}

	resp := &logical.Response{}
	if len(batchInputRaw) != 0 {
		batchResponseJSON, err := jsonutil.EncodeJSON(batchResponseItems)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON encode batch response")
		}
		resp.Data = map[string]interface{}{
			"data": string(batchResponseJSON),
		}
	} else {
		if batchResponseItems[0].Error != "" {
			return nil, fmt.Errorf(batchResponseItems[0].Error)
		}
		resp.Data = map[string]interface{}{
			"ciphertext": batchResponseItems[0].Ciphertext,
		}
	}

	return resp, nil
}

const pathRewrapHelpSyn = `Rewrap ciphertext`

const pathRewrapHelpDesc = `
After key rotation, this function can be used to rewrap the given ciphertext or
a batch of given ciphertext blocks with the latest version of the named key.
If the given ciphertext is already using the latest version of the key, this
function is a no-op.
`
