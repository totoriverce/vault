package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/errwrap"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	groupTypeInternal = "internal"
	groupTypeExternal = "external"
)

func groupPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "group$",
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the group. If set, updates the corresponding existing group.",
				},
				"type": {
					Type:        framework.TypeString,
					Description: "Type of the group, 'internal' or 'external'. Defaults to 'internal'",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the group.",
				},
				"metadata": {
					Type: framework.TypeKVPairs,
					Description: `Metadata to be associated with the group.
In CLI, this parameter can be repeated multiple times, and it all gets merged together.
For example:
vault <command> <path> metadata=key1=value1 metadata=key2=value2
					`,
				},
				"policies": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Policies to be tied to the group.",
				},
				"member_group_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Group IDs to be assigned as group members.",
				},
				"member_entity_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Entity IDs to be assigned as group members.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathGroupRegister(),
			},

			HelpSynopsis:    strings.TrimSpace(groupHelp["register"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["register"][1]),
		},
		{
			Pattern: "group/id/" + framework.GenericNameRegex("id"),
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the group.",
				},
				"type": {
					Type:        framework.TypeString,
					Default:     groupTypeInternal,
					Description: "Type of the group, 'internal' or 'external'. Defaults to 'internal'",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the group.",
				},
				"metadata": {
					Type: framework.TypeKVPairs,
					Description: `Metadata to be associated with the group.
In CLI, this parameter can be repeated multiple times, and it all gets merged together.
For example:
vault <command> <path> metadata=key1=value1 metadata=key2=value2
					`,
				},
				"policies": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Policies to be tied to the group.",
				},
				"member_group_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Group IDs to be assigned as group members.",
				},
				"member_entity_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Entity IDs to be assigned as group members.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathGroupIDUpdate(),
				logical.ReadOperation:   i.pathGroupIDRead(),
				logical.DeleteOperation: i.pathGroupIDDelete(),
			},

			HelpSynopsis:    strings.TrimSpace(groupHelp["group-by-id"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-by-id"][1]),
		},
		{
			Pattern: "group/id/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathGroupIDList(),
			},

			HelpSynopsis:    strings.TrimSpace(groupHelp["group-id-list"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-id-list"][1]),
		},
	}
}

func (i *IdentityStore) pathGroupRegister() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		_, ok := d.GetOk("id")
		if ok {
			return i.pathGroupIDUpdate()(ctx, req, d)
		}

		i.groupLock.Lock()
		defer i.groupLock.Unlock()

		return i.handleGroupUpdateCommon(req, d, nil)
	}
}

func (i *IdentityStore) pathGroupIDUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupID := d.Get("id").(string)
		if groupID == "" {
			return logical.ErrorResponse("empty group ID"), nil
		}

		i.groupLock.Lock()
		defer i.groupLock.Unlock()

		group, err := i.MemDBGroupByID(groupID, true)
		if err != nil {
			return nil, err
		}
		if group == nil {
			return logical.ErrorResponse("invalid group ID"), nil
		}

		return i.handleGroupUpdateCommon(req, d, group)
	}
}

func (i *IdentityStore) handleGroupUpdateCommon(req *logical.Request, d *framework.FieldData, group *identity.Group) (*logical.Response, error) {
	var err error
	var newGroup bool
	if group == nil {
		group = &identity.Group{}
		newGroup = true
	}

	// Update the policies if supplied
	policiesRaw, ok := d.GetOk("policies")
	if ok {
		group.Policies = policiesRaw.([]string)
	}

	groupTypeRaw, ok := d.GetOk("type")
	if ok {
		groupType := groupTypeRaw.(string)
		if group.Type != "" && groupType != group.Type {
			return logical.ErrorResponse(fmt.Sprintf("group type cannot be changed")), nil
		}

		group.Type = groupType
	}

	// If group type is not set, default to internal type
	if group.Type == "" {
		group.Type = groupTypeInternal
	}

	if group.Type != groupTypeInternal && group.Type != groupTypeExternal {
		return logical.ErrorResponse(fmt.Sprintf("invalid group type %q", group.Type)), nil
	}

	// Get the name
	groupName := d.Get("name").(string)
	if groupName != "" {
		// Check if there is a group already existing for the given name
		groupByName, err := i.MemDBGroupByName(groupName, false)
		if err != nil {
			return nil, err
		}

		// If this is a new group and if there already exists a group by this
		// name, error out. If the name of an existing group is about to be
		// modified into something which is already tied to a different group,
		// error out.
		switch {
		case (newGroup && groupByName != nil), (groupByName != nil && group.ID != "" && groupByName.ID != group.ID):
			return logical.ErrorResponse("group name is already in use"), nil
		}
		group.Name = groupName
	}

	metadata, ok, err := d.GetOkErr("metadata")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to parse metadata: %v", err)), nil
	}
	if ok {
		group.Metadata = metadata.(map[string]string)
	}

	memberEntityIDsRaw, ok := d.GetOk("member_entity_ids")
	if ok {
		if group.Type == groupTypeExternal {
			return logical.ErrorResponse("member entities can't be set manually for external groups"), nil
		}
		group.MemberEntityIDs = memberEntityIDsRaw.([]string)
		if len(group.MemberEntityIDs) > 512 {
			return logical.ErrorResponse("member entity IDs exceeding the limit of 512"), nil
		}
	}

	memberGroupIDsRaw, ok := d.GetOk("member_group_ids")
	var memberGroupIDs []string
	if ok {
		if group.Type == groupTypeExternal {
			return logical.ErrorResponse("member groups can't be set for external groups"), nil
		}
		memberGroupIDs = memberGroupIDsRaw.([]string)
	}

	err = i.sanitizeAndUpsertGroup(group, memberGroupIDs)
	if err != nil {
		return nil, err
	}

	respData := map[string]interface{}{
		"id":   group.ID,
		"name": group.Name,
	}
	return &logical.Response{
		Data: respData,
	}, nil
}

func (i *IdentityStore) pathGroupIDRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupID := d.Get("id").(string)
		if groupID == "" {
			return logical.ErrorResponse("empty group id"), nil
		}

		group, err := i.MemDBGroupByID(groupID, false)
		if err != nil {
			return nil, err
		}

		return i.handleGroupReadCommon(group)
	}
}

func (i *IdentityStore) handleGroupReadCommon(group *identity.Group) (*logical.Response, error) {
	if group == nil {
		return nil, nil
	}

	respData := map[string]interface{}{}
	respData["id"] = group.ID
	respData["name"] = group.Name
	respData["policies"] = group.Policies
	respData["member_entity_ids"] = group.MemberEntityIDs
	respData["parent_group_ids"] = group.ParentGroupIDs
	respData["metadata"] = group.Metadata
	respData["creation_time"] = ptypes.TimestampString(group.CreationTime)
	respData["last_update_time"] = ptypes.TimestampString(group.LastUpdateTime)
	respData["modify_index"] = group.ModifyIndex
	respData["type"] = group.Type

	aliasMap := map[string]interface{}{}
	if group.Alias != nil {
		aliasMap["id"] = group.Alias.ID
		aliasMap["canonical_id"] = group.Alias.CanonicalID
		aliasMap["mount_type"] = group.Alias.MountType
		aliasMap["mount_accessor"] = group.Alias.MountAccessor
		aliasMap["mount_path"] = group.Alias.MountPath
		aliasMap["metadata"] = group.Alias.Metadata
		aliasMap["name"] = group.Alias.Name
		aliasMap["merged_from_canonical_ids"] = group.Alias.MergedFromCanonicalIDs
		aliasMap["creation_time"] = ptypes.TimestampString(group.Alias.CreationTime)
		aliasMap["last_update_time"] = ptypes.TimestampString(group.Alias.LastUpdateTime)
	}

	respData["alias"] = aliasMap

	var memberGroupIDs []string
	memberGroups, err := i.MemDBGroupsByParentGroupID(group.ID, false)
	if err != nil {
		return nil, err
	}
	for _, memberGroup := range memberGroups {
		memberGroupIDs = append(memberGroupIDs, memberGroup.ID)
	}

	respData["member_group_ids"] = memberGroupIDs

	return &logical.Response{
		Data: respData,
	}, nil
}

func (i *IdentityStore) pathGroupIDDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupID := d.Get("id").(string)
		if groupID == "" {
			return logical.ErrorResponse("empty group ID"), nil
		}

		if groupID == "" {
			return nil, fmt.Errorf("missing group ID")
		}

		// Acquire the lock to modify the group storage entry
		i.groupLock.Lock()
		defer i.groupLock.Unlock()

		// Create a MemDB transaction to delete group
		txn := i.db.Txn(true)
		defer txn.Abort()

		group, err := i.MemDBGroupByIDInTxn(txn, groupID, false)
		if err != nil {
			return nil, err
		}

		// If there is no group for the ID, do nothing
		if group == nil {
			return nil, nil
		}

		// Delete group alias from memdb
		if group.Type == groupTypeExternal && group.Alias != nil {
			err = i.MemDBDeleteAliasByIDInTxn(txn, group.Alias.ID, true)
			if err != nil {
				return nil, err
			}
		}

		// Delete the group using the same transaction
		err = i.MemDBDeleteGroupByIDInTxn(txn, group.ID)
		if err != nil {
			return nil, err
		}

		// Delete the group from storage
		err = i.groupPacker.DeleteItem(group.ID)
		if err != nil {
			return nil, err
		}

		// Committing the transaction *after* successfully deleting group
		txn.Commit()

		return nil, nil
	}
}

// pathGroupIDList lists the IDs of all the groups in the identity store
func (i *IdentityStore) pathGroupIDList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		ws := memdb.NewWatchSet()

		txn := i.db.Txn(false)

		iter, err := txn.Get(groupsTable, "id")
		if err != nil {
			return nil, errwrap.Wrapf("failed to fetch iterator for group in memdb: {{err}}", err)
		}

		ws.Add(iter.WatchCh())

		var groupIDs []string
		groupInfo := map[string]interface{}{}

		type mountInfo struct {
			MountType string
			MountPath string
		}
		mountAccessorMap := map[string]mountInfo{}

		for {
			raw := iter.Next()
			if raw == nil {
				break
			}
			group := raw.(*identity.Group)
			groupIDs = append(groupIDs, group.ID)
			groupInfoEntry := map[string]interface{}{
				"name":                group.Name,
				"num_member_entities": len(group.MemberEntityIDs),
				"num_parent_groups":   len(group.ParentGroupIDs),
			}
			if group.Alias != nil {
				entry := map[string]interface{}{
					"id":             group.Alias.ID,
					"name":           group.Alias.Name,
					"mount_accessor": group.Alias.MountAccessor,
				}

				mi, ok := mountAccessorMap[group.Alias.MountAccessor]
				if ok {
					entry["mount_type"] = mi.MountType
					entry["mount_path"] = mi.MountPath
				} else {
					mi = mountInfo{}
					if mountValidationResp := i.core.router.validateMountByAccessor(group.Alias.MountAccessor); mountValidationResp != nil {
						mi.MountType = mountValidationResp.MountType
						mi.MountPath = mountValidationResp.MountPath
						entry["mount_type"] = mi.MountType
						entry["mount_path"] = mi.MountPath
					}
					mountAccessorMap[group.Alias.MountAccessor] = mi
				}

				groupInfoEntry["alias"] = entry
			}
			groupInfo[group.ID] = groupInfoEntry
		}

		return logical.ListResponseWithInfo(groupIDs, groupInfo), nil
	}
}

var groupHelp = map[string][2]string{
	"register": {
		"Create a new group.",
		"",
	},
	"group-by-id": {
		"Update or delete an existing group using its ID.",
		"",
	},
	"group-id-list": {
		"List all the group IDs.",
		"",
	},
}
