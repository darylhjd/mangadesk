package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	CreateMangaPath                = MangaListPath
	GetMangaVolumesAndChaptersPath = "manga/%s/aggregate"
	ViewMangaPath                  = "manga/%s"
	UpdateMangaPath                = ViewMangaPath
	DeleteMangaPath                = ViewMangaPath
	UnfollowMangaPath              = "manga/%s/follow"
	FollowMangaPath                = UnfollowMangaPath
	MangaFeedPath                  = "manga/%s/feed"
	MangaReadMarkersPath           = "manga/%s/read"
	GetRandomMangaPath             = "manga/random"
	TagListPath                    = "manga/tag"
	GetMangaReadingStatusPath      = "manga/%s/status"
	UpdateMangaReadingStatusPath   = GetMangaReadingStatusPath
)

type MangaList struct {
	Data   []Manga `json:"data"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
	Total  int     `json:"total"`
}

type MangaResponse struct {
	Result        string         `json:"result"`
	Data          Manga          `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (mr *MangaResponse) GetResult() string {
	return mr.Result
}

type Manga struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Attributes    MangaAttributes `json:"attributes"`
	Relationships []Relationship  `json:"relationships"`
}

type Relationship struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type MangaAttributes struct {
	Title                  LocalisedStrings   `json:"title"`
	AltTitles              []LocalisedStrings `json:"altTitles"`
	Description            LocalisedStrings   `json:"description"`
	IsLocked               bool               `json:"isLocked"`
	OriginalLanguage       string             `json:"originalLanguage"`
	LastVolume             *string            `json:"lastVolume"`
	LastChapter            *string            `json:"lastChapter"`
	PublicationDemographic *string            `json:"publicationDemographic"`
	Status                 *string            `json:"status"`
	Year                   *int               `json:"year"`
	ContentRating          *string            `json:"contentRating"`
	Tags                   []Tag              `json:"tags"`
	Version                int                `json:"version"`
	CreatedAt              string             `json:"createdAt"`
	UpdatedAt              string             `json:"updatedAt"`
}

type LocalisedStrings map[string]string

type ChapterReadMarkersResponse struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

func (rmr *ChapterReadMarkersResponse) GetResult() string {
	return rmr.Result
}

type TagResponse struct {
	Result        string         `json:"result"`
	Data          Tag            `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (tg *TagResponse) GetResult() string {
	return tg.Result
}

type Tag struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Attributes TagAttributes `json:"attributes"`
}

type TagAttributes struct {
	Name    LocalisedStrings `json:"name"`
	Group   string           `json:"group"`
	Version int              `json:"version"`
}

type MangaVolChapsResponse struct {
	Result  string               `json:"result"`
	Volumes map[string]VolumeAgg `json:"volumes"`
}

func (r *MangaVolChapsResponse) GetResult() string {
	return r.Result
}

type VolumeAgg struct {
	Volume   string                `json:"volume"`
	Count    int                   `json:"count"`
	Chapters map[string]ChapterAgg `json:"chapters"`
}

type ChapterAgg struct {
	Chapter string `json:"chapter"`
	Count   int    `json:"count"`
}

type AllMangaReadingStatusResponse struct {
	Result string            `json:"result"`
	Status map[string]string `json:"statuses"`
}

func (s *AllMangaReadingStatusResponse) GetResult() string {
	return s.Result
}

type MangaReadingStatusResponse struct {
	Result string `json:"result"`
	Status string `json:"status"`
}

func (r *MangaReadingStatusResponse) GetResult() string {
	return r.Result
}

// CreateManga : Create a new manga.
// https://api.mangadex.org/docs.html#operation/post-manga
func (dc *DexClient) CreateManga(newManga io.Reader) (*MangaResponse, error) {
	return dc.CreateMangaContext(context.Background(), newManga)
}

// CreateMangaContext : CreateManga with custom context.
func (dc *DexClient) CreateMangaContext(ctx context.Context, newManga io.Reader) (*MangaResponse, error) {
	var mr MangaResponse
	err := dc.responseOp(ctx, http.MethodPost, CreateMangaPath, newManga, &mr)
	return &mr, err
}

// GetMangaVolumesAndChapters : Get volume and chapters aggregate for a manga.
// https://api.mangadex.org/docs.html#tag/Manga/paths/~1manga~1{id}~1aggregate/get
func (dc *DexClient) GetMangaVolumesAndChapters(id string, ls []string) (*MangaVolChapsResponse, error) {
	return dc.GetMangaVolumesAndChaptersContext(context.Background(), id, ls)
}

// GetMangaVolumesAndChaptersContext : GetMangaVolumesAndChapters with custom context.
func (dc *DexClient) GetMangaVolumesAndChaptersContext(ctx context.Context, id string, ls []string) (*MangaVolChapsResponse, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(GetMangaVolumesAndChaptersPath, id)

	// Set query parameters
	q := u.Query()
	for _, l := range ls {
		q.Add("translatedLanguage", l)
	}
	u.RawQuery = q.Encode()

	var r MangaVolChapsResponse
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r)
	return &r, err
}

// ViewManga : View a manga by ID.
// https://api.mangadex.org/docs.html#operation/get-manga-id
func (dc *DexClient) ViewManga(id string) (*MangaResponse, error) {
	return dc.ViewMangaContext(context.Background(), id)
}

// ViewMangaContext : ViewManga with custom context.
func (dc *DexClient) ViewMangaContext(ctx context.Context, id string) (*MangaResponse, error) {
	var mr MangaResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(ViewMangaPath, id), nil, &mr)
	return &mr, err
}

// UpdateManga : Update a Manga.
// https://api.mangadex.org/docs.html#operation/put-manga-id
func (dc *DexClient) UpdateManga(id string, upManga io.Reader) (*MangaResponse, error) {
	return dc.UpdateMangaContext(context.Background(), id, upManga)
}

// UpdateMangaContext : UpdateManga with custom context.
func (dc *DexClient) UpdateMangaContext(ctx context.Context, id string, upManga io.Reader) (*MangaResponse, error) {
	var mr MangaResponse
	err := dc.responseOp(ctx, http.MethodPut, fmt.Sprintf(UpdateMangaPath, id), upManga, &mr)
	return &mr, err
}

// DeleteManga : Delete a Manga through ID.
// https://api.mangadex.org/docs.html#operation/delete-manga-id
func (dc *DexClient) DeleteManga(id string) error {
	return dc.DeleteMangaContext(context.Background(), id)
}

// DeleteMangaContext : DeleteManga with custom context.
func (dc *DexClient) DeleteMangaContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(DeleteMangaPath, id), nil, nil)
}

// UnfollowManga : Unfollow a Manga by ID.
// https://api.mangadex.org/docs.html#operation/delete-manga-id-follow
func (dc *DexClient) UnfollowManga(id string) error {
	return dc.UnfollowMangaContext(context.Background(), id)
}

// UnfollowMangaContext : UnfollowManga with custom context.
func (dc *DexClient) UnfollowMangaContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(UnfollowMangaPath, id), nil, nil)
}

// FollowManga : Follow a Manga by ID.
// https://api.mangadex.org/docs.html#operation/post-manga-id-follow
func (dc *DexClient) FollowManga(id string) error {
	return dc.FollowMangaContext(context.Background(), id)
}

// FollowMangaContext : FollowManga with custom context.
func (dc *DexClient) FollowMangaContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodPost, fmt.Sprintf(FollowMangaPath, id), nil, nil)
}

// MangaFeed : Get Manga feed by ID.
// https://api.mangadex.org/docs.html#operation/get-manga-id-feed
func (dc *DexClient) MangaFeed(id string, params url.Values) (*ChapterList, error) {
	return dc.MangaFeedContext(context.Background(), id, params)
}

// MangaFeedContext : MangaFeed with custom context.
func (dc *DexClient) MangaFeedContext(ctx context.Context, id string, params url.Values) (*ChapterList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaFeedPath, id)

	// Set request parameters
	u.RawQuery = params.Encode()

	var l ChapterList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// MangaReadMarkers : Get list of Chapter IDs that are marked as read for a specified manga ID.
// https://api.mangadex.org/docs.html#operation/get-manga-chapter-readmarkers
func (dc *DexClient) MangaReadMarkers(id string) (*ChapterReadMarkersResponse, error) {
	return dc.MangaReadMarkersContext(context.Background(), id)
}

// MangaReadMarkersContext : MangaReadMarkers with custom context.
func (dc *DexClient) MangaReadMarkersContext(ctx context.Context, id string) (*ChapterReadMarkersResponse, error) {
	var rmr ChapterReadMarkersResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(MangaReadMarkersPath, id), nil, &rmr)
	return &rmr, err
}

// GetRandomManga : Return a random Manga.
// https://api.mangadex.org/docs.html#operation/get-manga-random
func (dc *DexClient) GetRandomManga() (*MangaResponse, error) {
	return dc.GetRandomMangaContext(context.Background())
}

// GetRandomMangaContext : GetRandomManga with custom context.
func (dc *DexClient) GetRandomMangaContext(ctx context.Context) (*MangaResponse, error) {
	var mr MangaResponse
	err := dc.responseOp(ctx, http.MethodGet, GetRandomMangaPath, nil, &mr)
	return &mr, err
}

// TagList : Get tag list.
// https://api.mangadex.org/docs.html#operation/get-manga-tag
func (dc *DexClient) TagList() (*[]TagResponse, error) {
	return dc.TagListContext(context.Background())
}

// TagListContext : TagList with custom context.
func (dc *DexClient) TagListContext(ctx context.Context) (*[]TagResponse, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = TagListPath

	var tg []TagResponse
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &tg)
	return &tg, err
}

// GetMangaReadingStatus : Get reading status for a manga.
// https://api.mangadex.org/docs.html#operation/get-manga-id-status
func (dc *DexClient) GetMangaReadingStatus(id string) (*MangaReadingStatusResponse, error) {
	return dc.GetMangaReadingStatusContext(context.Background(), id)
}

// GetMangaReadingStatusContext : GetMangaReadingStatus with custom context.
func (dc *DexClient) GetMangaReadingStatusContext(ctx context.Context, id string) (*MangaReadingStatusResponse, error) {
	var r MangaReadingStatusResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(GetMangaReadingStatusPath, id), nil, &r)
	return &r, err
}

// UpdateMangaReadingStatus : Update reading status for a manga.
func (dc *DexClient) UpdateMangaReadingStatus(id string, status ReadStatus) error {
	return dc.UpdateMangaReadingStatusContext(context.Background(), id, status)
}

// UpdateMangaReadingStatusContext : UpdateMangaReadingStatus with custom context.
func (dc *DexClient) UpdateMangaReadingStatusContext(ctx context.Context, id string, status ReadStatus) error {
	// Create required request body.
	req := map[string]ReadStatus{
		"status": status,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	return dc.responseOp(ctx, http.MethodPost,
		fmt.Sprintf(UpdateMangaReadingStatusPath, id), bytes.NewBuffer(rBytes),
		nil)
}
