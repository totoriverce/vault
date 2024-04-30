//go:build !enterprise

package pki

import (
	"context"
	"crypto/x509"
	"errors"
	"math/big"

	"github.com/hashicorp/vault/sdk/logical"
)

var ErrMetadataIsEntOnly = errors.New("certificate metadata is only supported on Vault Enterprise")

func storeMetadata(ctx context.Context, storage logical.Storage, certificate *x509.Certificate, metadata interface{}) error {
	return ErrMetadataIsEntOnly
}

func GetCertificateMetadata(ctx context.Context, storage logical.Storage, serialNumber *big.Int) (*CertificateMetadata, error) {
	return nil, ErrMetadataIsEntOnly
}
