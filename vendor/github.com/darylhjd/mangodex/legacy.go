package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

const (
	LegacyIDMappingPath = "legacy/mapping"
)

type LegacyMappingResponse struct {
	Result        string         `json:"result"`
	Data          MappingID      `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (r *LegacyMappingResponse) GetResult() string {
	return r.Result
}

type MappingID struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes MappingIDAttributes `json:"attributes"`
}

type MappingIDAttributes struct {
	Type     string `json:"type"`
	LegacyID int    `json:"legacyId"`
	NewID    string `json:"newId"`
}

// LegacyIDMapping : Map Legacy IDs.
// https://api.mangadex.org/docs.html#operation/post-legacy-mapping
func (dc *DexClient) LegacyIDMapping(typ string, ids []string) (*LegacyMappingResponse, error) {
	return dc.LegacyIDMappingContext(context.Background(), typ, ids)
}

// LegacyIDMappingContext : LegacyIDMapping with custom context.
func (dc *DexClient) LegacyIDMappingContext(ctx context.Context, typ string, ids []string) (*LegacyMappingResponse, error) {
	// Create request body.
	req := struct {
		Type string
		IDs  []string
	}{Type: typ, IDs: ids}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r LegacyMappingResponse
	err = dc.responseOp(ctx, http.MethodPost, LegacyIDMappingPath, bytes.NewBuffer(rBytes), &r)
	return &r, err
}
