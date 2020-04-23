package aliasmetadata

/*
	aliasmetadata is a package offering convenience and
	standardization when supporting an `alias_metadata`
	field in a plugin's configuration. This then controls
	what alias metadata is added to an Auth during login.

	To see an example of how to add and use it, check out
	how these structs and fields are used in the AWS auth
	method.

	Or, check out its acceptance test in this package to
	see its integration points.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// Fields is for configuring a back-end's available
// default and additional fields. These are used for
// providing a verbose field description, and for parsing
// user input.
type Fields struct {
	// The field name as it'll be reflected in the user-facing
	// schema.
	FieldName string

	// Default is a list of the default fields that should
	// be included if a user sends "default" in their list
	// of desired fields. These fields should all have a
	// low rate of change because each change can incur a
	// write to storage.
	Default []string

	// AvailableToAdd is a list of fields not included by
	// default, that the user may include.
	AvailableToAdd []string
}

func (f *Fields) all() []string {
	return append(f.Default, f.AvailableToAdd...)
}

// FieldSchema takes the default and additionally available
// fields, and uses them to generate a verbose description
// regarding how to use the "alias_metadata" field.
func FieldSchema(fields *Fields) *framework.FieldSchema {
	return &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Description: description(fields),
		DisplayAttrs: &framework.DisplayAttributes{
			Name:  fields.FieldName,
			Value: "default,field1,field2",
		},
		Default: []string{"default"},
	}
}

func NewHandler(fields *Fields) *Handler {
	return &Handler{
		fields: fields,
	}
}

type Handler struct {
	// aliasMetadata is an explicit list of all the user's configured
	// fields that are being added to alias metadata. It will never
	// include the "default" parameter, and instead includes the actual
	// fields behind "default", if selected. If it has never been set
	// or only the default fields are desired, it is nil.
	aliasMetadata []string

	// fields is a list of the configured default and available
	// fields. It's intentionally not jsonified.
	fields *Fields
}

// AliasMetadata is intended to be used on config reads.
// It gets an explicit list of all the user's configured
// fields that are being added to alias metadata. It will never
// include the "default" parameter, and instead includes the actual
// fields behind "default", if selected.
func (h *Handler) AliasMetadata() []string {
	if h.aliasMetadata == nil {
		return h.fields.Default
	}
	return h.aliasMetadata
}

// ParseAliasMetadata is intended to be used on config create/update.
// It takes a user's selected fields (or lack thereof),
// converts it to a list of explicit fields, and adds it to the Handler
// for later storage.
func (h *Handler) ParseAliasMetadata(data *framework.FieldData) error {
	userProvidedRaw, ok := data.GetOk(h.fields.FieldName)
	if !ok {
		// Nothing further to do here.
		return nil
	}
	userProvided, ok := userProvidedRaw.([]string)
	if !ok {
		return fmt.Errorf("%s is an unexpected type of %T", userProvidedRaw, userProvidedRaw)
	}

	// If the only field the user has chosen was the default field,
	// we don't store anything so we won't have to do a storage
	// migration if the default changes.
	if len(userProvided) == 1 && userProvided[0] == "default" {
		h.aliasMetadata = nil
		return nil
	}

	// uniqueFields protects against weird edge cases like if
	// a user provided "default,field1,field2,default".
	uniqueFields := make(map[string]bool)
	for _, field := range userProvided {
		if field == "default" {
			// Add the field that "default" represents, rather
			// than the explicit field.
			for _, dfltField := range h.fields.Default {
				uniqueFields[dfltField] = true
			}
		} else {
			// Make sure they've sent a supported field so we can
			// error early if not.
			if !strutil.StrListContains(h.fields.all(), field) {
				return fmt.Errorf("%q is not an available field, please select from: %s", field, strings.Join(h.fields.all(), ", "))
			}
			uniqueFields[field] = true
		}
	}
	// Attach the fields we've received so they'll be stored.
	aliasMetadata := make([]string, len(uniqueFields))
	i := 0
	for fieldName := range uniqueFields {
		aliasMetadata[i] = fieldName
		i++
	}
	// Fulfilling the pointer here flags that the user has made
	// a non-default selection.
	h.aliasMetadata = aliasMetadata
	return nil
}

// PopulateDesiredAliasMetadata is intended to be used during login
// just before returning an auth.
// It takes the available alias metadata and,
// if the auth should have it, adds it to the auth's alias metadata.
func (h *Handler) PopulateDesiredAliasMetadata(auth *logical.Auth, available map[string]string) error {
	if auth == nil {
		return errors.New("auth is nil")
	}
	if auth.Alias == nil {
		return errors.New("auth alias is nil")
	}
	if auth.Alias.Name == "" {
		// We need the caller to set the alias name or there will
		// be nothing for these fields to operate upon.
		return errors.New("auth alias name must be set")
	}
	if auth.Alias.Metadata == nil {
		auth.Alias.Metadata = make(map[string]string)
	}
	fieldsToInclude := h.fields.Default
	if h.aliasMetadata != nil {
		fieldsToInclude = h.aliasMetadata
	}
	for availableField, itsValue := range available {
		if strutil.StrListContains(fieldsToInclude, availableField) {
			auth.Alias.Metadata[availableField] = itsValue
		}
	}
	return nil
}

func (h *Handler) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		AliasMetadata []string `json:"alias_metadata"`
	}{
		AliasMetadata: h.aliasMetadata,
	})
}

func (h *Handler) UnmarshalJSON(data []byte) error {
	jsonable := &struct {
		AliasMetadata []string `json:"alias_metadata"`
	}{
		AliasMetadata: h.aliasMetadata,
	}
	if err := json.Unmarshal(data, jsonable); err != nil {
		return err
	}
	h.aliasMetadata = jsonable.AliasMetadata
	return nil
}

func description(fields *Fields) string {
	desc := "The metadata to include on the aliases generated by this plugin."
	if len(fields.Default) > 0 {
		desc += fmt.Sprintf(" When set to 'default', includes: %s.", strings.Join(fields.Default, ", "))
	}
	if len(fields.AvailableToAdd) > 0 {
		desc += fmt.Sprintf(" These fields are available to add: %s.", strings.Join(fields.AvailableToAdd, ", "))
	}
	desc += " Not editing this field means the 'default' fields are included." +
		" Explicitly setting this field to empty overrides the 'default' and means no alias metadata will be included." +
		" Add fields by sending, 'default,field1,field2'." +
		" We advise only including fields that change rarely because each change triggers a storage write."
	return desc
}
