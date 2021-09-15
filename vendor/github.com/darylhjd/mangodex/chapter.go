package mangodex

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const (
	DeleteChapterPath     = "chapter/%s"
	GetChapterPath        = DeleteChapterPath
	UpdateChapterPath     = DeleteChapterPath
	MarkChapterReadPath   = "chapter/%s/read"
	MarkChapterUnreadPath = MarkChapterReadPath
)

type ChapterList struct {
	Data   []Chapter `json:"data"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
	Total  int       `json:"total"`
}

type ChapterResponse struct {
	Result        string         `json:"result"`
	Data          Chapter        `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (r *ChapterResponse) GetResult() string {
	return r.Result
}

type Chapter struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Attributes    ChapterAttributes `json:"attributes"`
	Relationships []Relationship    `json:"relationships"`
}

type ChapterAttributes struct {
	Title              string   `json:"title"`
	Volume             *string  `json:"volume"`
	Chapter            *string  `json:"chapter"`
	TranslatedLanguage string   `json:"translatedLanguage"`
	Hash               string   `json:"hash"`
	Data               []string `json:"data"`
	DataSaver          []string `json:"dataSaver"`
	Uploader           string   `json:"uploader"`
	Version            int      `json:"version"`
	CreatedAt          string   `json:"createdAt"`
	UpdatedAt          string   `json:"updatedAt"`
	PublishAt          string   `json:"publishAt"`
}

// DeleteChapter : Remove a chapter by ID.
// https://api.mangadex.org/docs.html#operation/delete-chapter-id
func (dc *DexClient) DeleteChapter(id string) error {
	return dc.DeleteChapterContext(context.Background(), id)
}

// DeleteChapterContext : DeleteChapter with custom context.
func (dc *DexClient) DeleteChapterContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(DeleteChapterPath, id), nil, nil)
}

// GetChapter : Get a chapter by ID.
// https://api.mangadex.org/docs.html#operation/get-chapter-id
func (dc *DexClient) GetChapter(id string) (*ChapterResponse, error) {
	return dc.GetChapterContext(context.Background(), id)
}

// GetChapterContext : GetChapter with custom context.
func (dc *DexClient) GetChapterContext(ctx context.Context, id string) (*ChapterResponse, error) {
	var r ChapterResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(GetChapterPath, id), nil, &r)
	return &r, err
}

// UpdateChapter : Update a chapter by ID
// https://api.mangadex.org/docs.html#operation/put-chapter-id
func (dc *DexClient) UpdateChapter(id string, upChapter io.Reader) (*ChapterResponse, error) {
	return dc.UpdateChapterContext(context.Background(), id, upChapter)
}

// UpdateChapterContext : UpdateChapter with custom context.
func (dc *DexClient) UpdateChapterContext(ctx context.Context, id string, upChapter io.Reader) (*ChapterResponse, error) {
	var r ChapterResponse
	err := dc.responseOp(ctx, http.MethodPut, fmt.Sprintf(UpdateChapterPath, id), upChapter, &r)
	return &r, err
}

// MarkChapterRead : Mark chapter as read.
// https://api.mangadex.org/docs.html#operation/chapter-id-read
func (dc *DexClient) MarkChapterRead(id string) error {
	return dc.MarkChapterReadContext(context.Background(), id)
}

// MarkChapterReadContext : MarkChapterRead with custom context.
func (dc *DexClient) MarkChapterReadContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodPost, fmt.Sprintf(MarkChapterReadPath, id), nil, nil)
}

// MarkChapterUnread : Mark chapter as unread.
// https://api.mangadex.org/docs.html#operation/chapter-id-unread
func (dc *DexClient) MarkChapterUnread(id string) error {
	return dc.MarkChapterUnreadContext(context.Background(), id)
}

// MarkChapterUnreadContext : MarkChapterUnread with custom context.
func (dc *DexClient) MarkChapterUnreadContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(MarkChapterUnreadPath, id), nil, nil)
}
