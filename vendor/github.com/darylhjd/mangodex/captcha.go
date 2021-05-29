package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

const (
	SolveCaptchaPath = "captcha/solve"
)

// SolveCaptcha : Solve a captcha
// https://api.mangadex.org/docs.html#operation/post-captcha-solve
func (dc *DexClient) SolveCaptcha(cc string) error {
	return dc.SolveCaptchaContext(context.Background(), cc)
}

// SolveCaptchaContext : SolveCaptcha with custom context.
func (dc *DexClient) SolveCaptchaContext(ctx context.Context, cc string) error {
	// Create request body
	req := map[string]string{
		"captchaChallenge": cc,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	return dc.responseOp(ctx, http.MethodPost, SolveCaptchaPath, bytes.NewBuffer(rBytes), nil)
}
