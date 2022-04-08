package pki

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/sdk/logical"
)

func isKeyDefaultSet(ctx context.Context, s logical.Storage) (bool, error) {
	config, err := getKeysConfig(ctx, s)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(config.DefaultKeyId.String()) != "", nil
}

func isIssuerDefaultSet(ctx context.Context, s logical.Storage) (bool, error) {
	config, err := getIssuersConfig(ctx, s)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(config.DefaultIssuerId.String()) != "", nil
}

func updateDefaultKeyId(ctx context.Context, s logical.Storage, id keyId) error {
	config, err := getKeysConfig(ctx, s)
	if err != nil {
		return err
	}

	if config.DefaultKeyId != id {
		return setKeysConfig(ctx, s, &keyConfig{
			DefaultKeyId: id,
		})
	}

	return nil
}

func updateDefaultIssuerId(ctx context.Context, s logical.Storage, id issuerId) error {
	config, err := getIssuersConfig(ctx, s)
	if err != nil {
		return err
	}

	if config.DefaultIssuerId != id {
		return setIssuersConfig(ctx, s, &issuerConfig{
			DefaultIssuerId: id,
		})
	}

	return nil
}
