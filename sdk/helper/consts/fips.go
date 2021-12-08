//go:build !fips_140_3
// +build !fips_140_3

package consts

// IsFIPS returns true if Vault is operating in a FIPS-140-{2,3} mode.
func IsFIPS() bool {
	return false
}
