package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	CreateAuthorPath = AuthorListPath
	GetAuthorPath    = "author/%s"
	UpdateAuthorPath = GetAuthorPath
	DeleteAuthorPath = GetAuthorPath
)

type AuthorList struct {
	Results []AuthorResponse `json:"results"`
	Limit   int              `json:"limit"`
	Offset  int              `json:"offset"`
	Total   int              `json:"total"`
}

type AuthorResponse struct {
	Result        string         `json:"result"`
	Data          Author         `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (ar *AuthorResponse) GetResult() string {
	return ar.Result
}

type Author struct {
	ID         string           `json:"string"`
	Type       string           `json:"type"`
	Attributes AuthorAttributes `json:"attributes"`
}

type AuthorAttributes struct {
	Name      string             `json:"name"`
	ImageURL  string             `json:"imageUrl"`
	Biography []LocalisedStrings `json:"biography"`
	Version   int                `json:"version"`
	CreatedAt string             `json:"createdAt"`
	UpdatedAt string             `json:"updatedAt"`
}

// CreateAuthor : Create Author.
// https://api.mangadex.org/docs.html#operation/post-author
func (dc *DexClient) CreateAuthor(name string, version int) (*AuthorResponse, error) {
	return dc.CreateAuthorContext(context.Background(), name, version)
}

// CreateAuthorContext : CreateAuthor with custom context.
func (dc *DexClient) CreateAuthorContext(ctx context.Context, name string, version int) (*AuthorResponse, error) {
	// Create required request body.
	req := struct {
		Name    string
		Version int
	}{Name: name, Version: version}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r AuthorResponse
	err = dc.responseOp(ctx, http.MethodPost, CreateAuthorPath, bytes.NewBuffer(rBytes), &r)
	return &r, err
}

// GetAuthor : Get Author.
// https://api.mangadex.org/docs.html#operation/get-author-id
func (dc *DexClient) GetAuthor(id string) (*AuthorResponse, error) {
	return dc.GetAuthorContext(context.Background(), id)
}

// GetAuthorContext : GetAuthor with custom context.
func (dc *DexClient) GetAuthorContext(ctx context.Context, id string) (*AuthorResponse, error) {
	var r AuthorResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(GetAuthorPath, id), nil, &r)
	return &r, err
}

// UpdateAuthor : Update Author.
// https://api.mangadex.org/docs.html#operation/put-author-id
func (dc *DexClient) UpdateAuthor(id, name string, version int) (*AuthorResponse, error) {
	return dc.UpdateAuthorContext(context.Background(), id, name, version)
}

// UpdateAuthorContext : UpdateAuthor with custom context.
func (dc *DexClient) UpdateAuthorContext(ctx context.Context, id, name string, version int) (*AuthorResponse, error) {
	// Create required request body.
	req := struct {
		Name    string
		Version int
	}{Name: name, Version: version}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r AuthorResponse
	err = dc.responseOp(ctx, http.MethodPut, fmt.Sprintf(UpdateAuthorPath, id), bytes.NewBuffer(rBytes), &r)
	return &r, err
}

// DeleteAuthor : Delete Author.
func (dc *DexClient) DeleteAuthor(id string) error {
	return dc.DeleteAuthorContext(context.Background(), id)
}

// DeleteAuthorContext : DeleteAuthor with custom context.
func (dc *DexClient) DeleteAuthorContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(DeleteAuthorPath, id), nil, nil)
}
