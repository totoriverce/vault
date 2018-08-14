// Package base62 provides utilities for working with base62 strings.
// base62 strings will only contain characters: 0-9, a-z, A-Z
package base62

import (
	"math/big"

	uuid "github.com/hashicorp/go-uuid"
)

// Encode converts buf into a base62 string
func Encode(buf []byte) string {
	var encoder big.Int

	encoder.SetBytes(buf)
	return encoder.Text(62)
}

// Decode converts input from base62 to its byte representation
// If the decoding fails, an empty slice is returned.
func Decode(input string) []byte {
	var decoder big.Int

	decoder.SetString(input, 62)
	return decoder.Bytes()
}

// Random generates a random base62-encoded string of the given length.
func Random(length int) (string, error) {
	buf, err := uuid.GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	return Encode(buf)[:length], nil
}
