package roottoken

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/xor"
)

// EncodeToken gets a token and an OTP and encodes the token.
// The OTP must have the same length as the token.
func EncodeToken(token, otp string) (string, error) {
	if len(token) == 0 {
		return "", errors.New("no token provided")
	} else if len(otp) == 0 {
		return "", errors.New("no otp provided")
	}

	// This function performs decoding checks so rather than decode the OTP,
	// just encode the value we're passing in.
	tokenBytes, err := xor.XORBytes([]byte(otp), []byte(token))
	if err != nil {
		return "", fmt.Errorf("xor of root token failed: %w", err)
	}
	return base64.RawStdEncoding.EncodeToString(tokenBytes), nil
}
