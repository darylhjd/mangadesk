package mangodex

import (
	"context"
	"net/http"
	"net/url"
)

const (
	MangaListPath     = "manga"
	ScanGroupListPath = "group"
	ChapterListPath   = "chapter"
	CoverArtListPath  = "cover"
	AuthorListPath    = "author"
)

// MangaList : Get a list of manga.
// https://api.mangadex.org/docs.html#operation/get-search-manga
func (dc *DexClient) MangaList(params url.Values) (*MangaList, error) {
	return dc.MangaListContext(context.Background(), params)
}

// MangaListContext : MangaList with custom context.
func (dc *DexClient) MangaListContext(ctx context.Context, params url.Values) (*MangaList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = MangaListPath

	// Set query parameters
	u.RawQuery = params.Encode()

	var l MangaList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// ScanGroupList : Return a list of scanlation groups.
// https://api.mangadex.org/docs.html#operation/get-search-group
func (dc *DexClient) ScanGroupList(params url.Values) (*ScanGroupList, error) {
	return dc.ScanGroupListContext(context.Background(), params)
}

// ScanGroupListContext : ScanGroupList with custom context.
func (dc *DexClient) ScanGroupListContext(ctx context.Context, params url.Values) (*ScanGroupList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = ScanGroupListPath

	// Set query parameters
	u.RawQuery = params.Encode()

	var l ScanGroupList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// ChapterList : Get a list of chapters.
// https://api.mangadex.org/docs.html#operation/get-chapter
func (dc *DexClient) ChapterList(params url.Values) (*ChapterList, error) {
	return dc.ChapterListContext(context.Background(), params)
}

// ChapterListContext : ChapterList with custom context.
func (dc *DexClient) ChapterListContext(ctx context.Context, params url.Values) (*ChapterList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = ChapterListPath

	// Set query parameters
	u.RawQuery = params.Encode()

	var l ChapterList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// CoverArtList : Get a list of manga covers.
// https://api.mangadex.org/docs.html#operation/get-cover
func (dc *DexClient) CoverArtList(params url.Values) (*CoverArtList, error) {
	return dc.CoverArtListContext(context.Background(), params)
}

// CoverArtListContext : CoverArtList with custom context.
func (dc *DexClient) CoverArtListContext(ctx context.Context, params url.Values) (*CoverArtList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = CoverArtListPath

	// Set query parameters
	u.RawQuery = params.Encode()

	var l CoverArtList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// AuthorList : Return a list of authors.
// https://api.mangadex.org/docs.html#operation/get-author
func (dc *DexClient) AuthorList(params url.Values) (*AuthorList, error) {
	return dc.AuthorListContext(context.Background(), params)
}

// AuthorListContext : AuthorList with custom context.
func (dc *DexClient) AuthorListContext(ctx context.Context, params url.Values) (*AuthorList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = AuthorListPath

	// Set query parameters
	u.RawQuery = params.Encode()

	var l AuthorList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}
