package vault

import (
	"fmt"

	memdb "github.com/hashicorp/go-memdb"
)

func identityStoreSchema() *memdb.DBSchema {
	iStoreSchema := &memdb.DBSchema{
		Tables: make(map[string]*memdb.TableSchema),
	}

	schemas := []func() *memdb.TableSchema{
		entityTableSchema,
		aliasesTableSchema,
		groupTableSchema,
	}

	for _, schemaFunc := range schemas {
		schema := schemaFunc()
		if _, ok := iStoreSchema.Tables[schema.Name]; ok {
			panic(fmt.Sprintf("duplicate table name: %s", schema.Name))
		}
		iStoreSchema.Tables[schema.Name] = schema
	}

	return iStoreSchema
}

func aliasesTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "aliases",
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ID",
				},
			},
			"entity_id": &memdb.IndexSchema{
				Name:   "entity_id",
				Unique: false,
				Indexer: &memdb.StringFieldIndex{
					Field: "EntityID",
				},
			},
			"mount_type": &memdb.IndexSchema{
				Name:   "mount_type",
				Unique: false,
				Indexer: &memdb.StringFieldIndex{
					Field: "MountType",
				},
			},
			"factors": &memdb.IndexSchema{
				Name:   "factors",
				Unique: true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "MountAccessor",
						},
						&memdb.StringFieldIndex{
							Field: "Name",
						},
					},
				},
			},
			"metadata": &memdb.IndexSchema{
				Name:         "metadata",
				Unique:       false,
				AllowMissing: true,
				Indexer: &memdb.StringMapFieldIndex{
					Field: "Metadata",
				},
			},
		},
	}
}

func entityTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "entities",
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ID",
				},
			},
			"name": &memdb.IndexSchema{
				Name:   "name",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Name",
				},
			},
			"metadata": &memdb.IndexSchema{
				Name:         "metadata",
				Unique:       false,
				AllowMissing: true,
				Indexer: &memdb.StringMapFieldIndex{
					Field: "Metadata",
				},
			},
			"merged_entity_ids": &memdb.IndexSchema{
				Name:         "merged_entity_ids",
				Unique:       true,
				AllowMissing: true,
				Indexer: &memdb.StringSliceFieldIndex{
					Field: "MergedEntityIDs",
				},
			},
			"bucket_key_hash": &memdb.IndexSchema{
				Name:         "bucket_key_hash",
				Unique:       false,
				AllowMissing: false,
				Indexer: &memdb.StringFieldIndex{
					Field: "BucketKeyHash",
				},
			},
		},
	}
}

func groupTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "groups",
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ID",
				},
			},
			"name": {
				Name:   "name",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Name",
				},
			},
			"member_entity_ids": {
				Name:         "member_entity_ids",
				Unique:       false,
				AllowMissing: true,
				Indexer: &memdb.StringSliceFieldIndex{
					Field: "MemberEntityIDs",
				},
			},
			"parent_group_ids": {
				Name:         "parent_group_ids",
				Unique:       false,
				AllowMissing: true,
				Indexer: &memdb.StringSliceFieldIndex{
					Field: "ParentGroupIDs",
				},
			},
			"policies": {
				Name:         "policies",
				Unique:       false,
				AllowMissing: true,
				Indexer: &memdb.StringSliceFieldIndex{
					Field: "Policies",
				},
			},
			"bucket_key_hash": &memdb.IndexSchema{
				Name:         "bucket_key_hash",
				Unique:       false,
				AllowMissing: false,
				Indexer: &memdb.StringFieldIndex{
					Field: "BucketKeyHash",
				},
			},
		},
	}
}
