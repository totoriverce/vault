package database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

func pathRotateCredentials(b *databaseBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "rotate-root/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of this database connection",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateCredentialsUpdate(),
					ForwardPerformanceSecondary: true,
					ForwardPerformanceStandby:   true,
				},
			},

			HelpSynopsis:    pathCredsCreateReadHelpSyn,
			HelpDescription: pathCredsCreateReadHelpDesc,
		},
		&framework.Path{
			Pattern: "rotate-role/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the static role",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateRoleCredentialsUpdate(),
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},

			HelpSynopsis:    pathCredsCreateReadHelpSyn,
			HelpDescription: pathCredsCreateReadHelpDesc,
		},
	}
}

func (b *databaseBackend) pathRotateCredentialsUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		config, err := b.DatabaseConfig(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}

		db, err := b.GetConnection(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}

		// Take out the backend lock since we are swapping out the connection
		b.Lock()
		defer b.Unlock()

		// Take the write lock on the instance
		db.Lock()
		defer db.Unlock()

		// Generate new credentials
		userName := config.ConnectionDetails["username"].(string)
		oldPassword := config.ConnectionDetails["password"].(string)
		newPassword, err := db.GenerateCredentials(ctx)
		if err != nil {
			return nil, err
		}
		config.ConnectionDetails["password"] = newPassword

		// Write a WAL entry
		walID, err := framework.PutWAL(ctx, req.Storage, rotateRootWALKey, &rotateRootCredentialsWAL{
			ConnectionName: name,
			UserName:       userName,
			OldPassword:    oldPassword,
			NewPassword:    newPassword,
		})
		if err != nil {
			return nil, err
		}

		// Attempt to use SetCredentials for the root credential rotation
		statements := dbplugin.Statements{Rotation: config.RootCredentialsRotateStatements}
		userConfig := dbplugin.StaticUserConfig{
			Username: userName,
			Password: newPassword,
		}
		if _, _, err := db.SetCredentials(ctx, statements, userConfig); err != nil {
			if status.Code(err) == codes.Unimplemented {
				// Fall back to using RotateRootCredentials if unimplemented
				config.ConnectionDetails, err = db.RotateRootCredentials(ctx,
					config.RootCredentialsRotateStatements)
			}
			if err != nil {
				return nil, err
			}
		}

		// Update storage with the new root credentials
		entry, err := logical.StorageEntryJSON(fmt.Sprintf("config/%s", name), config)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		// Delete the WAL entry after successfully rotating root credentials
		if err := framework.DeleteWAL(ctx, req.Storage, walID); err != nil {
			b.Logger().Warn("unable to delete WAL", "error", err, "WAL ID", walID)
		}

		// Close the plugin
		db.closed = true
		if err := db.Database.Close(); err != nil {
			b.Logger().Error("error closing the database plugin connection", "err", err)
		}
		// Even on error, still remove the connection
		delete(b.connections, name)

		return nil, nil
	}
}

func (b *databaseBackend) pathRotateRoleCredentialsUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse("empty role name attribute given"), nil
		}

		role, err := b.StaticRole(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return logical.ErrorResponse("no static role found for role name"), nil
		}

		// In create/update of static accounts, we only care if the operation
		// err'd , and this call does not return credentials
		item, err := b.popFromRotationQueueByKey(name)
		if err != nil {
			item = &queue.Item{
				Key: name,
			}
		}

		resp, err := b.setStaticAccount(ctx, req.Storage, &setStaticAccountInput{
			RoleName: name,
			Role:     role,
		})
		if err != nil {
			b.logger.Warn("unable to rotate credentials in rotate-role", "error", err)
			// Update the priority to re-try this rotation and re-add the item to
			// the queue
			item.Priority = time.Now().Add(10 * time.Second).Unix()

			// Preserve the WALID if it was returned
			if resp.WALID != "" {
				item.Value = resp.WALID
			}
		} else {
			item.Priority = resp.RotationTime.Add(role.StaticAccount.RotationPeriod).Unix()
		}

		// Add their rotation to the queue
		if err := b.pushItem(item); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

const pathRotateCredentialsUpdateHelpSyn = `
Request to rotate the root credentials for a certain database connection.
`

const pathRotateCredentialsUpdateHelpDesc = `
This path attempts to rotate the root credentials for the given database. 
`

const pathRotateRoleCredentialsUpdateHelpSyn = `
Request to rotate the credentials for a static user account.
`
const pathRotateRoleCredentialsUpdateHelpDesc = `
This path attempts to rotate the credentials for the given static user account.
`
