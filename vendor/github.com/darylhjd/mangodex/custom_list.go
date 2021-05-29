package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	CreateCustomListPath        = "list"
	GetCustomListPath           = "list/%s"
	UpdateCustomListPath        = GetCustomListPath
	DeleteCustomListPath        = GetCustomListPath
	AddMangaInListPath          = "manga/%s/list/%s"
	RemoveMangaInListPath       = AddMangaInListPath
	GetPublicCustomListListPath = "user/%s/list"
)

type CustomListList struct {
	Results []CustomListResponse `json:"results"`
	Limit   int                  `json:"limit"`
	Offset  int                  `json:"offset"`
	Total   int                  `json:"total"`
}

type CustomListResponse struct {
	Result        string         `json:"result"`
	Data          CustomList     `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (r *CustomListResponse) GetResult() string {
	return r.Result
}

type CustomList struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes CustomListAttributes `json:"attributes"`
}

type CustomListAttributes struct {
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
	Owner      User   `json:"owner"`
	Version    int    `json:"version"`
}

// CreateCustomList : Create a new custom list.
// https://api.mangadex.org/docs.html#operation/post-list
func (dc *DexClient) CreateCustomList(name, visibility string, manga []string, version int) (*CustomListResponse, error) {
	return dc.CreateCustomListContext(context.Background(), name, visibility, manga, version)
}

// CreateCustomListContext : CreateCustomList with custom context.
func (dc *DexClient) CreateCustomListContext(ctx context.Context, name, visibility string, manga []string, version int) (*CustomListResponse, error) {
	// Create request body.
	req := struct {
		Name       string
		Visibility string
		manga      []string
		Version    int
	}{Name: name, Visibility: visibility, manga: manga, Version: version}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r CustomListResponse
	err = dc.responseOp(ctx, http.MethodPost, CreateCustomListPath, bytes.NewBuffer(rBytes), &r)
	return &r, err
}

// GetCustomList : Get a custom list by ID.
// https://api.mangadex.org/docs.html#operation/get-list-id
func (dc *DexClient) GetCustomList(id string) (*CustomListResponse, error) {
	return dc.GetCustomListContext(context.Background(), id)
}

// GetCustomListContext : GetCustomList with custom context.
func (dc *DexClient) GetCustomListContext(ctx context.Context, id string) (*CustomListResponse, error) {
	var r CustomListResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(GetCustomListPath, id), nil, &r)
	return &r, err
}

// UpdateCustomList : Update a new custom list.
// https://api.mangadex.org/docs.html#operation/put-list-id
func (dc *DexClient) UpdateCustomList(id, name, visibility string, manga []string, version int) (*CustomListResponse, error) {
	return dc.UpdateCustomListContext(context.Background(), id, name, visibility, manga, version)
}

// UpdateCustomListContext : UpdateCustomList with custom context.
func (dc *DexClient) UpdateCustomListContext(ctx context.Context, id, name, visibility string, manga []string, version int) (*CustomListResponse, error) {
	// Create request body.
	req := struct {
		Name       string
		Visibility string
		manga      []string
		Version    int
	}{Name: name, Visibility: visibility, manga: manga, Version: version}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r CustomListResponse
	err = dc.responseOp(ctx, http.MethodPut, fmt.Sprintf(UpdateCustomListPath, id), bytes.NewBuffer(rBytes), &r)
	return &r, err
}

// DeleteCustomList : Delete a custom list by ID.
// https://api.mangadex.org/docs.html#operation/delete-list-id
func (dc *DexClient) DeleteCustomList(id string) (*CustomListResponse, error) {
	return dc.DeleteCustomListContext(context.Background(), id)
}

// DeleteCustomListContext : DeleteCustomList with custom context.
func (dc *DexClient) DeleteCustomListContext(ctx context.Context, id string) (*CustomListResponse, error) {
	var r CustomListResponse
	err := dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(DeleteCustomListPath, id), nil, &r)
	return &r, err
}

// AddMangaInList : Add a Manga to a custom list.
// https://api.mangadex.org/docs.html#operation/post-manga-id-list-listId
func (dc *DexClient) AddMangaInList(mangaId, listId string) error {
	return dc.AddMangaInListContext(context.Background(), mangaId, listId)
}

// AddMangaInListContext : AddMangaInList with custom context.
func (dc *DexClient) AddMangaInListContext(ctx context.Context, mangaId, listId string) error {
	return dc.responseOp(ctx, http.MethodPost, fmt.Sprintf(AddMangaInListPath, mangaId, listId), nil, nil)
}

// RemoveMangaInList : Remove a Manga from a custom list.
// https://api.mangadex.org/docs.html#operation/delete-manga-id-list-listId
func (dc *DexClient) RemoveMangaInList(mangaId, listId string) error {
	return dc.RemoveMangaInListContext(context.Background(), mangaId, listId)
}

// RemoveMangaInListContext : RemoveMangaInList with custom context.
func (dc *DexClient) RemoveMangaInListContext(ctx context.Context, mangaId, listId string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(RemoveMangaInListPath, mangaId, listId), nil, nil)
}

// GetPublicCustomListList : Get a public custom list by ID. Only for public lists.
// https://api.mangadex.org/docs.html#operation/get-user-id-list
func (dc *DexClient) GetPublicCustomListList(id string, limit, offset int) (*CustomListList, error) {
	return dc.GetPublicCustomListListContext(context.Background(), id, limit, offset)
}

// GetPublicCustomListListContext : GetPublicCustomListList with custom context.
func (dc *DexClient) GetPublicCustomListListContext(ctx context.Context, id string, limit, offset int) (*CustomListList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(GetPublicCustomListListPath, id)

	// Set query parameters
	q := u.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	var l CustomListList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}
