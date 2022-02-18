package api

import (
	"context"
	"errors"

	"github.com/mitchellh/mapstructure"
)

func (c *Sys) RekeyStatus() (*RekeyStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyStatusWithContext(ctx)
}

func (c *Sys) RekeyStatusWithContext(ctx context.Context) (*RekeyStatusResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/rekey/init")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRecoveryKeyStatus() (*RekeyStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyRecoveryKeyStatusWithContext(ctx)
}

func (c *Sys) RekeyRecoveryKeyStatusWithContext(ctx context.Context) (*RekeyStatusResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/rekey-recovery-key/init")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyVerificationStatus() (*RekeyVerificationStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyVerificationStatusWithContext(ctx)
}

func (c *Sys) RekeyVerificationStatusWithContext(ctx context.Context) (*RekeyVerificationStatusResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/rekey/verify")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyVerificationStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRecoveryKeyVerificationStatus() (*RekeyVerificationStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyRecoveryKeyVerificationStatusWithContext(ctx)
}

func (c *Sys) RekeyRecoveryKeyVerificationStatusWithContext(ctx context.Context) (*RekeyVerificationStatusResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/rekey-recovery-key/verify")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyVerificationStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyInit(config *RekeyInitRequest) (*RekeyStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyInitWithContext(ctx, config)
}

func (c *Sys) RekeyInitWithContext(ctx context.Context, config *RekeyInitRequest) (*RekeyStatusResponse, error) {
	r := c.c.NewRequest("PUT", "/v1/sys/rekey/init")
	if err := r.SetJSONBody(config); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRecoveryKeyInit(config *RekeyInitRequest) (*RekeyStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyRecoveryKeyInitWithContext(ctx, config)
}

func (c *Sys) RekeyRecoveryKeyInitWithContext(ctx context.Context, config *RekeyInitRequest) (*RekeyStatusResponse, error) {
	r := c.c.NewRequest("PUT", "/v1/sys/rekey-recovery-key/init")
	if err := r.SetJSONBody(config); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyCancel() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyCancelWithContext(ctx)
}

func (c *Sys) RekeyCancelWithContext(ctx context.Context) error {
	r := c.c.NewRequest("DELETE", "/v1/sys/rekey/init")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RekeyRecoveryKeyCancel() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyRecoveryKeyCancelWithContext(ctx)
}

func (c *Sys) RekeyRecoveryKeyCancelWithContext(ctx context.Context) error {
	r := c.c.NewRequest("DELETE", "/v1/sys/rekey-recovery-key/init")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RekeyVerificationCancel() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyVerificationCancelWithContext(ctx)
}

func (c *Sys) RekeyVerificationCancelWithContext(ctx context.Context) error {
	r := c.c.NewRequest("DELETE", "/v1/sys/rekey/verify")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RekeyRecoveryKeyVerificationCancel() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyRecoveryKeyVerificationCancelWithContext(ctx)
}

func (c *Sys) RekeyRecoveryKeyVerificationCancelWithContext(ctx context.Context) error {
	r := c.c.NewRequest("DELETE", "/v1/sys/rekey-recovery-key/verify")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RekeyUpdate(shard, nonce string) (*RekeyUpdateResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyUpdateWithContext(ctx, shard, nonce)
}

func (c *Sys) RekeyUpdateWithContext(ctx context.Context, shard, nonce string) (*RekeyUpdateResponse, error) {
	body := map[string]interface{}{
		"key":   shard,
		"nonce": nonce,
	}

	r := c.c.NewRequest("PUT", "/v1/sys/rekey/update")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyUpdateResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRecoveryKeyUpdate(shard, nonce string) (*RekeyUpdateResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyRecoveryKeyUpdateWithContext(ctx, shard, nonce)
}

func (c *Sys) RekeyRecoveryKeyUpdateWithContext(ctx context.Context, shard, nonce string) (*RekeyUpdateResponse, error) {
	body := map[string]interface{}{
		"key":   shard,
		"nonce": nonce,
	}

	r := c.c.NewRequest("PUT", "/v1/sys/rekey-recovery-key/update")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyUpdateResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRetrieveBackup() (*RekeyRetrieveResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyRetrieveBackupWithContext(ctx)
}

func (c *Sys) RekeyRetrieveBackupWithContext(ctx context.Context) (*RekeyRetrieveResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/rekey/backup")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result RekeyRetrieveResponse
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (c *Sys) RekeyRetrieveRecoveryBackup() (*RekeyRetrieveResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyRetrieveRecoveryBackupWithContext(ctx)
}

func (c *Sys) RekeyRetrieveRecoveryBackupWithContext(ctx context.Context) (*RekeyRetrieveResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/rekey/recovery-key-backup")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result RekeyRetrieveResponse
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (c *Sys) RekeyDeleteBackup() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyDeleteBackupWithContext(ctx)
}

func (c *Sys) RekeyDeleteBackupWithContext(ctx context.Context) error {
	r := c.c.NewRequest("DELETE", "/v1/sys/rekey/backup")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}

	return err
}

func (c *Sys) RekeyDeleteRecoveryBackup() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyDeleteRecoveryBackupWithContext(ctx)
}

func (c *Sys) RekeyDeleteRecoveryBackupWithContext(ctx context.Context) error {
	r := c.c.NewRequest("DELETE", "/v1/sys/rekey/recovery-key-backup")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}

	return err
}

func (c *Sys) RekeyVerificationUpdate(shard, nonce string) (*RekeyVerificationUpdateResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyVerificationUpdateWithContext(ctx, shard, nonce)
}

func (c *Sys) RekeyVerificationUpdateWithContext(ctx context.Context, shard, nonce string) (*RekeyVerificationUpdateResponse, error) {
	body := map[string]interface{}{
		"key":   shard,
		"nonce": nonce,
	}

	r := c.c.NewRequest("PUT", "/v1/sys/rekey/verify")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyVerificationUpdateResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRecoveryKeyVerificationUpdate(shard, nonce string) (*RekeyVerificationUpdateResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RekeyRecoveryKeyVerificationUpdateWithContext(ctx, shard, nonce)
}

func (c *Sys) RekeyRecoveryKeyVerificationUpdateWithContext(ctx context.Context, shard, nonce string) (*RekeyVerificationUpdateResponse, error) {
	body := map[string]interface{}{
		"key":   shard,
		"nonce": nonce,
	}

	r := c.c.NewRequest("PUT", "/v1/sys/rekey-recovery-key/verify")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyVerificationUpdateResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

type RekeyInitRequest struct {
	SecretShares        int      `json:"secret_shares"`
	SecretThreshold     int      `json:"secret_threshold"`
	StoredShares        int      `json:"stored_shares"`
	PGPKeys             []string `json:"pgp_keys"`
	Backup              bool
	RequireVerification bool `json:"require_verification"`
}

type RekeyStatusResponse struct {
	Nonce                string   `json:"nonce"`
	Started              bool     `json:"started"`
	T                    int      `json:"t"`
	N                    int      `json:"n"`
	Progress             int      `json:"progress"`
	Required             int      `json:"required"`
	PGPFingerprints      []string `json:"pgp_fingerprints"`
	Backup               bool     `json:"backup"`
	VerificationRequired bool     `json:"verification_required"`
	VerificationNonce    string   `json:"verification_nonce"`
}

type RekeyUpdateResponse struct {
	Nonce                string   `json:"nonce"`
	Complete             bool     `json:"complete"`
	Keys                 []string `json:"keys"`
	KeysB64              []string `json:"keys_base64"`
	PGPFingerprints      []string `json:"pgp_fingerprints"`
	Backup               bool     `json:"backup"`
	VerificationRequired bool     `json:"verification_required"`
	VerificationNonce    string   `json:"verification_nonce,omitempty"`
}

type RekeyRetrieveResponse struct {
	Nonce   string              `json:"nonce" mapstructure:"nonce"`
	Keys    map[string][]string `json:"keys" mapstructure:"keys"`
	KeysB64 map[string][]string `json:"keys_base64" mapstructure:"keys_base64"`
}

type RekeyVerificationStatusResponse struct {
	Nonce    string `json:"nonce"`
	Started  bool   `json:"started"`
	T        int    `json:"t"`
	N        int    `json:"n"`
	Progress int    `json:"progress"`
}

type RekeyVerificationUpdateResponse struct {
	Nonce    string `json:"nonce"`
	Complete bool   `json:"complete"`
}
