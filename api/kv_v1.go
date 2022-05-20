package api

import (
	"context"
	"fmt"
)

//// KV v1 methods ////

// Get returns a secret from the KV v1 secrets engine.
func (kv *kvv1) Get(ctx context.Context, secretPath string) (*KVSecret, error) {
	pathToRead := fmt.Sprintf("%s/%s", kv.mountPath, secretPath)

	secret, err := kv.c.Logical().ReadWithContext(ctx, pathToRead)
	if err != nil {
		return nil, fmt.Errorf("error encountered while reading secret at %s: %w", pathToRead, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at %s", pathToRead)
	}

	return &KVSecret{
		Data:     secret.Data,
		Metadata: nil,
		Raw:      secret,
	}, nil
}

// Put inserts a key-value secret (e.g. {"password": "Hashi123"}) into the
// KV v1 secrets engine.
//
// If the secret already exists, it will be overwritten.
func (kv *kvv1) Put(ctx context.Context, secretPath string, data map[string]interface{}) error {
	pathToWriteTo := fmt.Sprintf("%s/%s", kv.mountPath, secretPath)

	_, err := kv.c.Logical().WriteWithContext(ctx, pathToWriteTo, data)
	if err != nil {
		return fmt.Errorf("error writing secret to %s: %w", pathToWriteTo, err)
	}

	return nil
}

// Delete deletes a secret from the KV v1 secrets engine.
func (kv *kvv1) Delete(ctx context.Context, secretPath string) error {
	pathToDelete := fmt.Sprintf("%s/%s", kv.mountPath, secretPath)

	_, err := kv.c.Logical().DeleteWithContext(ctx, pathToDelete)
	if err != nil {
		return fmt.Errorf("error deleting secret at %s: %w", pathToDelete, err)
	}

	return nil
}
