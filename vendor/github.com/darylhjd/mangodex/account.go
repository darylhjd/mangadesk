package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	CreateAccountPath          = "account/create"
	ActivateAccountPath        = "account/activate/%s"
	ResendActivationCodePath   = "account/activate/resend"
	RecoverAccountPath         = "account/recover"
	CompleteAccountRecoverPath = "account/recover/%s"
)

// CreateAccount : Create a new account.
// https://api.mangadex.org/docs.html#operation/post-account-create
func (dc *DexClient) CreateAccount(user, pass, email string) (*UserResponse, error) {
	return dc.CreateAccountContext(context.Background(), user, pass, email)
}

// CreateAccountContext : CreateAccount with custom context.
func (dc *DexClient) CreateAccountContext(ctx context.Context, user, pass, email string) (*UserResponse, error) {
	// Create request body
	req := map[string]string{
		"username": user, "password": pass, "email": email,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r UserResponse
	err = dc.responseOp(ctx, http.MethodPost, CreateAccountPath, bytes.NewBuffer(rBytes), &r)
	return &r, err
}

// ActivateAccount : Activate account.
// https://api.mangadex.org/docs.html#operation/get-account-activate-code
func (dc *DexClient) ActivateAccount(code string) error {
	return dc.ActivateAccountContext(context.Background(), code)
}

// ActivateAccountContext : ActivateAccount with custom context.
func (dc *DexClient) ActivateAccountContext(ctx context.Context, code string) error {
	return dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(ActivateAccountPath, code), nil, nil)
}

// ResendActivationCode : Resend activation code.
// https://api.mangadex.org/docs.html#operation/post-account-activate-resend
func (dc *DexClient) ResendActivationCode(email string) error {
	return dc.ResendActivationCodeContext(context.Background(), email)
}

// ResendActivationCodeContext : ResendActivationCode with custom context.
func (dc *DexClient) ResendActivationCodeContext(ctx context.Context, email string) error {
	// Create request body
	req := map[string]string{
		"email": email,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	return dc.responseOp(ctx, http.MethodPost, ResendActivationCodePath, bytes.NewBuffer(rBytes), nil)
}

// RecoverAccount : Recover an account.
// https://api.mangadex.org/docs.html#operation/post-account-recover
func (dc *DexClient) RecoverAccount(email string) error {
	return dc.RecoverAccountContext(context.Background(), email)
}

// RecoverAccountContext : RecoverAccount with custom context.
func (dc *DexClient) RecoverAccountContext(ctx context.Context, email string) error {
	// Create request body.
	req := map[string]string{
		"email": email,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	return dc.responseOp(ctx, http.MethodPost, RecoverAccountPath, bytes.NewBuffer(rBytes), nil)
}

// CompleteAccountRecover : Complete account recovery.
// https://api.mangadex.org/docs.html#operation/post-account-recover-code
func (dc *DexClient) CompleteAccountRecover(code, newp string) error {
	return dc.CompleteAccountRecoverContext(context.Background(), code, newp)
}

// CompleteAccountRecoverContext : CompleteAccountRecover with custom context.
func (dc *DexClient) CompleteAccountRecoverContext(ctx context.Context, code, newp string) error {
	// Create request body
	req := map[string]string{
		"newPassword": newp,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	return dc.responseOp(ctx, http.MethodPost, fmt.Sprintf(CompleteAccountRecoverPath, code), bytes.NewBuffer(rBytes), nil)
}
