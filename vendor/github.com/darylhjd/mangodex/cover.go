package mangodex

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

const (
	UploadMangaCoverPath = "manga/%s/cover"
	GetMangaCoverPath    = "cover/%s"
	EditMangaCoverPath   = "manga/%s/cover/%s"
	DeleteMangaCoverPath = EditMangaCoverPath
)

type CoverArtList struct {
	Results []CoverResponse `json:"results"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	Total   int             `json:"total"`
}

type CoverResponse struct {
	Result        string         `json:"result"`
	Data          Cover          `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (r *CoverResponse) GetResult() string {
	return r.Result
}

type Cover struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes CoverAttributes `json:"attributes"`
}

type CoverAttributes struct {
	Volume      *string `json:"volume"`
	FileName    string  `json:"fileName"`
	Description *string `json:"description"`
	Version     int     `json:"version"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// UploadMangaCover : Upload a manga cover.
// https://api.mangadex.org/docs.html#operation/upload-cover
func (dc *DexClient) UploadMangaCover(id, filename string) (*CoverResponse, error) {
	return dc.UploadMangaCoverContext(context.Background(), id, filename)
}

// UploadMangaCoverContext : UploadMangaCover with custom context.
func (dc *DexClient) UploadMangaCoverContext(ctx context.Context, id, fPath string) (*CoverResponse, error) {
	// Get the file information
	pic, err := ioutil.ReadFile(fPath)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	part, err := writer.CreateFormFile("file", filepath.Base(fPath))
	if err != nil {
		return nil, err
	}
	_, err = part.Write(pic)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Remember to change the content type
	dc.header.Set("Content-Type", writer.FormDataContentType())
	defer dc.header.Del("Content-Type")

	var r CoverResponse
	err = dc.responseOp(ctx, http.MethodPost, fmt.Sprintf(UploadMangaCoverPath, id), &b, &r)
	return &r, err
}

// GetMangaCover : Get a manga cover by ID.
// https://api.mangadex.org/docs.html#operation/get-cover
func (dc *DexClient) GetMangaCover(id string) (*CoverResponse, error) {
	return dc.GetMangaCoverContext(context.Background(), id)
}

// GetMangaCoverContext : GetMangaCover with custom context.
func (dc *DexClient) GetMangaCoverContext(ctx context.Context, id string) (*CoverResponse, error) {
	var r CoverResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(GetMangaCoverPath, id), nil, &r)
	return &r, err
}

// EditMangaCover : Edit a manga cover
// https://api.mangadex.org/docs.html#operation/edit-cover
func (dc *DexClient) EditMangaCover(mangaId, coverId string, r io.Reader) (*CoverResponse, error) {
	return dc.EditMangaCoverContext(context.Background(), mangaId, coverId, r)
}

// EditMangaCoverContext : EditMangaCover with custom context.
func (dc *DexClient) EditMangaCoverContext(ctx context.Context, mangaId, coverId string, r io.Reader) (*CoverResponse, error) {
	var res CoverResponse
	err := dc.responseOp(ctx, http.MethodPost, fmt.Sprintf(EditMangaCoverPath, mangaId, coverId), r, &res)
	return &res, err
}

// DeleteMangaCover : Delete a manga cover.
// https://api.mangadex.org/docs.html#operation/delete-cover
func (dc *DexClient) DeleteMangaCover(mangaId, coverId string) error {
	return dc.DeleteMangaCoverContext(context.Background(), mangaId, coverId)
}

// DeleteMangaCoverContext : DeleteMangaCover with custom context.
func (dc *DexClient) DeleteMangaCoverContext(ctx context.Context, mangaId, coverId string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(DeleteMangaCoverPath, mangaId, coverId), nil, nil)
}
