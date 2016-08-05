// Package certutil contains helper functions that are mostly used
// with the PKI backend but can be generally useful. Functionality
// includes helpers for converting a certificate/private key bundle
// between DER and PEM, printing certificate serial numbers, and more.
//
// Functionality specific to the PKI backend includes some types // and helper methods to make requesting certificates from the // backend easy.
package certutil

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/vault/helper/errutil"
)

// Secret is used to attempt to unmarshal a Vault secret
// JSON response, as a convenience
type Secret struct {
	Data map[string]interface{} `json:"data"`
}

// PrivateKeyType holds a string representation of the type of private key (ec
// or rsa) referenced in CertBundle and ParsedCertBundle. This uses colloquial
// names rather than official names, to eliminate confusion
type PrivateKeyType string

//Well-known PrivateKeyTypes
const (
	UnknownPrivateKey PrivateKeyType = ""
	RSAPrivateKey     PrivateKeyType = "rsa"
	ECPrivateKey      PrivateKeyType = "ec"
)

// TLSUsage controls whether the intended usage of a *tls.Config
// returned from ParsedCertBundle.GetTLSConfig is for server use,
// client use, or both, which affects which values are set
type TLSUsage int

//Well-known TLSUsage types
const (
	TLSUnknown TLSUsage = 0
	TLSServer  TLSUsage = 1 << iota
	TLSClient
)

//BlockType indicates the serialization format of the key
type BlockType string

//Well-known formats
const (
	PKCS1Block BlockType = "RSA PRIVATE KEY"
	PKCS8Block BlockType = "PRIVATE KEY"
	ECBlock    BlockType = "EC PRIVATE KEY"
)

//ParsedPrivateKeyContainer allows common key setting for certs and CSRs
type ParsedPrivateKeyContainer interface {
	SetParsedPrivateKey(crypto.Signer, PrivateKeyType, []byte)
}

// CertBundle contains a key type, a PEM-encoded private key,
// a PEM-encoded certificate, and a string-encoded serial number,
// returned from a successful Issue request
type CertBundle struct {
	PrivateKeyType PrivateKeyType `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	Certificate    string         `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	IssuingCA      string         `json:"issuing_ca" structs:"issuing_ca" mapstructure:"issuing_ca"`
	IssuingCAChain string         `json:"issuing_ca_chain,omitempty" structs:"issuing_ca_chain,omitempty" mapstructure:"issuing_ca_chain"`
	PrivateKey     string         `json:"private_key" structs:"private_key" mapstructure:"private_key"`
	SerialNumber   string         `json:"serial_number" structs:"serial_number" mapstructure:"serial_number"`
}

// ParsedCertBundle contains a key type, a DER-encoded private key,
// and a DER-encoded certificate
type ParsedCertBundle struct {
	PrivateKeyType       PrivateKeyType
	PrivateKeyFormat     BlockType
	PrivateKeyBytes      []byte
	PrivateKey           crypto.Signer
	IssuingCAChainBytes  [][]byte
	IssuingCAChain       []*x509.Certificate
	IssuingCABytes       []byte
	IssuingCA            *x509.Certificate
	CertificateBytes     []byte
	Certificate          *x509.Certificate
	CertificatePath      []*x509.Certificate
	CertificatePathBytes [][]byte
	SerialNumber         *big.Int
}

// CSRBundle contains a key type, a PEM-encoded private key,
// and a PEM-encoded CSR
type CSRBundle struct {
	PrivateKeyType PrivateKeyType `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	CSR            string         `json:"csr" structs:"csr" mapstructure:"csr"`
	PrivateKey     string         `json:"private_key" structs:"private_key" mapstructure:"private_key"`
}

// ParsedCSRBundle contains a key type, a DER-encoded private key,
// and a DER-encoded certificate request
type ParsedCSRBundle struct {
	PrivateKeyType  PrivateKeyType
	PrivateKeyBytes []byte
	PrivateKey      crypto.Signer
	CSRBytes        []byte
	CSR             *x509.CertificateRequest
}

func (c *CertBundle) ToPEMBundle() string {
	var result []string

	if len(c.PrivateKey) > 0 {
		result = append(result, c.PrivateKey)
	}
	if len(c.Certificate) > 0 {
		result = append(result, c.Certificate)
	}
	if len(c.IssuingCA) > 0 {
		result = append(result, c.IssuingCA)
	}
	if len(c.IssuingCAChain) > 0 {
		result = append(result, c.IssuingCAChain)
	}

	return strings.Join(result, "\n")
}

// ToParsedCertBundle converts a string-based certificate bundle
// to a byte-based raw certificate bundle
func (c *CertBundle) ToParsedCertBundle() (*ParsedCertBundle, error) {
	result := &ParsedCertBundle{}
	var err error
	var pemBlock *pem.Block

	if len(c.PrivateKey) > 0 {
		pemBlock, _ = pem.Decode([]byte(c.PrivateKey))
		if pemBlock == nil {
			return nil, errutil.UserError{"Error decoding private key from cert bundle"}
		}

		result.PrivateKeyBytes = pemBlock.Bytes
		result.PrivateKeyFormat = BlockType(strings.TrimSpace(pemBlock.Type))

		switch result.PrivateKeyFormat {
		case ECBlock:
			result.PrivateKeyType, c.PrivateKeyType = ECPrivateKey, ECPrivateKey
		case PKCS1Block:
			c.PrivateKeyType, result.PrivateKeyType = RSAPrivateKey, RSAPrivateKey
		case PKCS8Block:
			t, err := getPKCS8Type(pemBlock.Bytes)
			if err != nil {
				return nil, errutil.UserError{fmt.Sprintf("Error getting key type from pkcs#8: %v", err)}
			}
			result.PrivateKeyType = t
			switch t {
			case ECPrivateKey:
				c.PrivateKeyType = ECPrivateKey
			case RSAPrivateKey:
				c.PrivateKeyType = RSAPrivateKey
			}
		default:
			return nil, errutil.UserError{fmt.Sprintf("Unsupported key block type: %s", pemBlock.Type)}
		}

		result.PrivateKey, err = result.getSigner()
		if err != nil {
			return nil, errutil.UserError{fmt.Sprintf("Error getting signer: %s", err)}
		}
	}

	if len(c.Certificate) > 0 {
		pemBlock, _ = pem.Decode([]byte(c.Certificate))
		if pemBlock == nil {
			return nil, errutil.UserError{"Error decoding certificate from cert bundle"}
		}
		result.CertificateBytes = pemBlock.Bytes
		result.Certificate, err = x509.ParseCertificate(result.CertificateBytes)
		if err != nil {
			return nil, errutil.UserError{"Error encountered parsing certificate bytes from raw bundle"}
		}
		result.CertificatePath = append(result.CertificatePath, result.Certificate)
		result.CertificatePathBytes = append(result.CertificatePathBytes, result.CertificateBytes)
	}

	if len(c.IssuingCA) > 0 {
		pemBlock, _ = pem.Decode([]byte(c.IssuingCA))
		if pemBlock == nil {
			return nil, errutil.UserError{"Error decoding issuing CA from cert bundle"}
		}
		result.IssuingCABytes = pemBlock.Bytes
		result.IssuingCA, err = x509.ParseCertificate(result.IssuingCABytes)
		if err != nil {
			return nil, errutil.UserError{fmt.Sprintf("Error parsing CA certificate: %s", err)}
		}
		result.CertificatePath = append(result.CertificatePath, result.IssuingCA)
		result.CertificatePathBytes = append(result.CertificatePathBytes, result.IssuingCABytes)
	}

	result.SerialNumber = result.Certificate.SerialNumber

	if len(c.IssuingCAChain) > 0 {
		pemBytes := []byte(c.IssuingCAChain)
		for len(pemBytes) > 0 {
			pemBlock, pemBytes = pem.Decode(pemBytes)
			if pemBlock == nil {
				return nil, errutil.UserError{"Error decoding issuing CA chain from cert bundle"}
			}
			result.IssuingCAChainBytes = append(result.IssuingCAChainBytes, pemBlock.Bytes)
			chainCert, err := x509.ParseCertificate(pemBlock.Bytes)
			if err != nil {
				return nil, errutil.UserError{fmt.Sprintf("Error parsing CA chain certificate: %s\n%s", err)}
			}
			result.IssuingCAChain = append(result.IssuingCAChain, chainCert)
		}
		result.CertificatePath = append(result.CertificatePath, result.IssuingCAChain...)
		result.CertificatePathBytes = append(result.CertificatePathBytes, result.IssuingCAChainBytes...)
	}

	// Populate if it isn't there already
	if len(c.SerialNumber) == 0 && len(c.Certificate) > 0 {
		c.SerialNumber = GetHexFormatted(result.Certificate.SerialNumber.Bytes(), ":")
	}

	return result, nil
}

// ToCertBundle converts a byte-based raw DER certificate bundle
// to a PEM-based string certificate bundle
func (p *ParsedCertBundle) ToCertBundle() (*CertBundle, error) {
	result := &CertBundle{}
	block := pem.Block{
		Type: "CERTIFICATE",
	}

	if p.Certificate != nil {
		result.SerialNumber = strings.TrimSpace(GetHexFormatted(p.Certificate.SerialNumber.Bytes(), ":"))
	}

	if p.CertificateBytes != nil && len(p.CertificateBytes) > 0 {
		block.Bytes = p.CertificateBytes
		result.Certificate = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	if p.IssuingCABytes != nil && len(p.IssuingCABytes) > 0 {
		block.Bytes = p.IssuingCABytes
		result.IssuingCA = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	if p.IssuingCAChainBytes != nil && len(p.IssuingCAChainBytes) > 0 {
		var issuingCaChain string
		for _, caCertBytes := range p.IssuingCAChainBytes {
			block.Bytes = caCertBytes
			issuingCaChain = fmt.Sprintf("%s\n%s", issuingCaChain, strings.TrimSpace(string(pem.EncodeToMemory(&block))))
		}
		result.IssuingCAChain = strings.TrimSpace(issuingCaChain)
	}

	if p.PrivateKeyBytes != nil && len(p.PrivateKeyBytes) > 0 {
		block.Type = string(p.PrivateKeyFormat)
		block.Bytes = p.PrivateKeyBytes
		result.PrivateKeyType = p.PrivateKeyType

		//Handle bundle not parsed by us
		if block.Type == "" {
			switch p.PrivateKeyType {
			case ECPrivateKey:
				block.Type = string(ECBlock)
			case RSAPrivateKey:
				block.Type = string(PKCS1Block)
			}
		}

		result.PrivateKey = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	return result, nil
}

// Verify checks if the parsed bundle is valid.  It validates the public
// key of the certificate to the private key and checks the certficate trust
// chain for path issues.
func (p *ParsedCertBundle) Verify() error {
	// If private key exists, check if it matches the public key of cert
	if p.PrivateKey != nil && p.Certificate != nil {
		equal, err := ComparePublicKeys(p.Certificate.PublicKey, p.PrivateKey.Public())
		if err != nil {
			return fmt.Errorf("could not compare public and private keys: %s", err)
		}
		if !equal {
			return fmt.Errorf("Public key of certificate does not match private key")
		}
	}

	if p.IssuingCA != nil {
		if !p.IssuingCA.IsCA {
			return fmt.Errorf("issuing ca is not a certificate authority")
		}
		if !bytes.Equal(p.Certificate.AuthorityKeyId, p.IssuingCA.SubjectKeyId) {
			return fmt.Errorf("issuing certificate authority id does not match the certificate")
		}

		for i, caCert := range p.IssuingCAChain {
			if !caCert.IsCA {
				return fmt.Errorf("certificate %d of certificate chain is not a certificate authority", i+1)
			}

			switch i {
			case 0:
				if !bytes.Equal(p.IssuingCA.AuthorityKeyId, caCert.SubjectKeyId) {
					return fmt.Errorf("certificate %d of certificate chain ca trust path is incorrect", i+1)
				}
			default:
				if !bytes.Equal(p.IssuingCAChain[i-1].AuthorityKeyId, caCert.SubjectKeyId) {
					return fmt.Errorf("certificate %d of certificate chain ca trust path is incorrect", i+1)
				}
			}
		}
	}

	return nil
}

// GetSigner returns a crypto.Signer corresponding to the private key
// contained in this ParsedCertBundle. The Signer contains a Public() function
// for getting the corresponding public. The Signer can also be
// type-converted to private keys
func (p *ParsedCertBundle) getSigner() (crypto.Signer, error) {
	var signer crypto.Signer
	var err error

	if p.PrivateKeyBytes == nil || len(p.PrivateKeyBytes) == 0 {
		return nil, errutil.UserError{"Given parsed cert bundle does not have private key information"}
	}

	switch p.PrivateKeyFormat {
	case ECBlock:
		signer, err = x509.ParseECPrivateKey(p.PrivateKeyBytes)
		if err != nil {
			return nil, errutil.UserError{fmt.Sprintf("Unable to parse CA's private EC key: %s", err)}
		}

	case PKCS1Block:
		signer, err = x509.ParsePKCS1PrivateKey(p.PrivateKeyBytes)
		if err != nil {
			return nil, errutil.UserError{fmt.Sprintf("Unable to parse CA's private RSA key: %s", err)}
		}

	case PKCS8Block:
		if k, err := x509.ParsePKCS8PrivateKey(p.PrivateKeyBytes); err == nil {
			switch k := k.(type) {
			case *rsa.PrivateKey, *ecdsa.PrivateKey:
				return k.(crypto.Signer), nil
			default:
				return nil, errutil.UserError{"Found unknown private key type in pkcs#8 wrapping"}
			}
		}
		return nil, errutil.UserError{fmt.Sprintf("Failed to parse pkcs#8 key: %v", err)}
	default:
		return nil, errutil.UserError{"Unable to determine type of private key; only RSA and EC are supported"}
	}
	return signer, nil
}

// SetParsedPrivateKey sets the private key parameters on the bundle
func (p *ParsedCertBundle) SetParsedPrivateKey(privateKey crypto.Signer, privateKeyType PrivateKeyType, privateKeyBytes []byte) {
	p.PrivateKey = privateKey
	p.PrivateKeyType = privateKeyType
	p.PrivateKeyBytes = privateKeyBytes
}

func getPKCS8Type(bs []byte) (PrivateKeyType, error) {
	k, err := x509.ParsePKCS8PrivateKey(bs)
	if err != nil {
		return UnknownPrivateKey, errutil.UserError{fmt.Sprintf("Failed to parse pkcs#8 key: %v", err)}
	}

	switch k.(type) {
	case *ecdsa.PrivateKey:
		return ECPrivateKey, nil
	case *rsa.PrivateKey:
		return RSAPrivateKey, nil
	default:
		return UnknownPrivateKey, errutil.UserError{"Found unknown private key type in pkcs#8 wrapping"}
	}
}

// ToParsedCSRBundle converts a string-based CSR bundle
// to a byte-based raw CSR bundle
func (c *CSRBundle) ToParsedCSRBundle() (*ParsedCSRBundle, error) {
	result := &ParsedCSRBundle{}
	var err error
	var pemBlock *pem.Block

	if len(c.PrivateKey) > 0 {
		pemBlock, _ = pem.Decode([]byte(c.PrivateKey))
		if pemBlock == nil {
			return nil, errutil.UserError{"Error decoding private key from cert bundle"}
		}
		result.PrivateKeyBytes = pemBlock.Bytes

		switch BlockType(pemBlock.Type) {
		case ECBlock:
			result.PrivateKeyType = ECPrivateKey
		case PKCS1Block:
			result.PrivateKeyType = RSAPrivateKey
		default:
			// Try to figure it out and correct
			if _, err := x509.ParseECPrivateKey(pemBlock.Bytes); err == nil {
				result.PrivateKeyType = ECPrivateKey
				c.PrivateKeyType = "ec"
			} else if _, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes); err == nil {
				result.PrivateKeyType = RSAPrivateKey
				c.PrivateKeyType = "rsa"
			} else {
				return nil, errutil.UserError{fmt.Sprintf("Unknown private key type in bundle: %s", c.PrivateKeyType)}
			}
		}

		result.PrivateKey, err = result.getSigner()
		if err != nil {
			return nil, errutil.UserError{fmt.Sprintf("Error getting signer: %s", err)}
		}
	}

	if len(c.CSR) > 0 {
		pemBlock, _ = pem.Decode([]byte(c.CSR))
		if pemBlock == nil {
			return nil, errutil.UserError{"Error decoding certificate from cert bundle"}
		}
		result.CSRBytes = pemBlock.Bytes
		result.CSR, err = x509.ParseCertificateRequest(result.CSRBytes)
		if err != nil {
			return nil, errutil.UserError{"Error encountered parsing certificate bytes from raw bundle"}
		}
	}

	return result, nil
}

// ToCSRBundle converts a byte-based raw DER certificate bundle
// to a PEM-based string certificate bundle
func (p *ParsedCSRBundle) ToCSRBundle() (*CSRBundle, error) {
	result := &CSRBundle{}
	block := pem.Block{
		Type: "CERTIFICATE REQUEST",
	}

	if p.CSRBytes != nil && len(p.CSRBytes) > 0 {
		block.Bytes = p.CSRBytes
		result.CSR = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	if p.PrivateKeyBytes != nil && len(p.PrivateKeyBytes) > 0 {
		block.Bytes = p.PrivateKeyBytes
		switch p.PrivateKeyType {
		case RSAPrivateKey:
			result.PrivateKeyType = "rsa"
			block.Type = "RSA PRIVATE KEY"
		case ECPrivateKey:
			result.PrivateKeyType = "ec"
			block.Type = "EC PRIVATE KEY"
		default:
			return nil, errutil.InternalError{"Could not determine private key type when creating block"}
		}
		result.PrivateKey = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
	}

	return result, nil
}

// GetSigner returns a crypto.Signer corresponding to the private key
// contained in this ParsedCSRBundle. The Signer contains a Public() function
// for getting the corresponding public. The Signer can also be
// type-converted to private keys
func (p *ParsedCSRBundle) getSigner() (crypto.Signer, error) {
	var signer crypto.Signer
	var err error

	if p.PrivateKeyBytes == nil || len(p.PrivateKeyBytes) == 0 {
		return nil, errutil.UserError{"Given parsed cert bundle does not have private key information"}
	}

	switch p.PrivateKeyType {
	case ECPrivateKey:
		signer, err = x509.ParseECPrivateKey(p.PrivateKeyBytes)
		if err != nil {
			return nil, errutil.UserError{fmt.Sprintf("Unable to parse CA's private EC key: %s", err)}
		}

	case RSAPrivateKey:
		signer, err = x509.ParsePKCS1PrivateKey(p.PrivateKeyBytes)
		if err != nil {
			return nil, errutil.UserError{fmt.Sprintf("Unable to parse CA's private RSA key: %s", err)}
		}

	default:
		return nil, errutil.UserError{"Unable to determine type of private key; only RSA and EC are supported"}
	}
	return signer, nil
}

// SetParsedPrivateKey sets the private key parameters on the bundle
func (p *ParsedCSRBundle) SetParsedPrivateKey(privateKey crypto.Signer, privateKeyType PrivateKeyType, privateKeyBytes []byte) {
	p.PrivateKey = privateKey
	p.PrivateKeyType = privateKeyType
	p.PrivateKeyBytes = privateKeyBytes
}

// GetTLSConfig returns a TLS config generally suitable for client
// authentiation. The returned TLS config can be modified slightly
// to be made suitable for a server requiring client authentication;
// specifically, you should set the value of ClientAuth in the returned
// config to match your needs.
func (p *ParsedCertBundle) GetTLSConfig(usage TLSUsage) (*tls.Config, error) {
	tlsCert := tls.Certificate{
		Certificate: [][]byte{},
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if p.Certificate != nil {
		tlsCert.Leaf = p.Certificate
	}

	if p.PrivateKey != nil {
		tlsCert.PrivateKey = p.PrivateKey
	}

	if p.CertificateBytes != nil && len(p.CertificateBytes) > 0 {
		tlsCert.Certificate = append(tlsCert.Certificate, p.CertificateBytes)
	}

	if p.IssuingCABytes != nil && len(p.IssuingCABytes) > 0 {
		tlsCert.Certificate = append(tlsCert.Certificate, p.IssuingCABytes)

		// Technically we only need one cert, but this doesn't duplicate code
		certBundle, err := p.ToCertBundle()
		if err != nil {
			return nil, fmt.Errorf("Error converting parsed bundle to string bundle when getting TLS config: %s", err)
		}

		caPool := x509.NewCertPool()
		ok := caPool.AppendCertsFromPEM([]byte(certBundle.IssuingCA))
		if !ok {
			return nil, fmt.Errorf("Could not append CA certificate")
		}

		if usage&TLSServer > 0 {
			tlsConfig.ClientCAs = caPool
			tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
		}
		if usage&TLSClient > 0 {
			tlsConfig.RootCAs = caPool
		}
	}

	if tlsCert.Certificate != nil && len(tlsCert.Certificate) > 0 {
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsConfig.BuildNameToCertificate()
	}

	return tlsConfig, nil
}

// IssueData is a structure that is suitable for marshaling into a request;
// either via JSON, or into a map[string]interface{} via the structs package
type IssueData struct {
	TTL        string `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	CommonName string `json:"common_name" structs:"common_name" mapstructure:"common_name"`
	AltNames   string `json:"alt_names" structs:"alt_names" mapstructure:"alt_names"`
	IPSANs     string `json:"ip_sans" structs:"ip_sans" mapstructure:"ip_sans"`
	CSR        string `json:"csr" structs:"csr" mapstructure:"csr"`
}
