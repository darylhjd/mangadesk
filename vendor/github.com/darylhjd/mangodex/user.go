package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GetUserPath                         = "user/%s"
	GetLoggedUserPath                   = "user/me"
	GetUserFollowedScanGroupListPath    = "user/follows/group"
	GetUserFollowedUsersListPath        = "user/follows/user"
	GetUserFollowedMangaListPath        = "user/follows/manga"
	GetUserFollowedMangaChapterFeedPath = "user/follows/manga/feed"
	GetUserMangaReadingStatusPath       = "manga/status"
	GetUserCustomListListPath           = "user/list"
)

type UserList struct {
	Results []UserResponse `json:"results"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
	Total   int            `json:"total"`
}

type UserResponse struct {
	Result        string         `json:"result"`
	Data          User           `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (r *UserResponse) GetResult() string {
	return r.Result
}

type User struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes UserAttributes `json:"attributes"`
}

type UserAttributes struct {
	Username string `json:"username"`
	Version  int    `json:"version"`
}

// GetUser : Return a UserResponse.
// https://api.mangadex.org/docs.html#operation/get-user-id
func (dc *DexClient) GetUser(id string) (*UserResponse, error) {
	return dc.GetUserContext(context.Background(), id)
}

// GetUserContext : GetUser with custom context.
func (dc *DexClient) GetUserContext(ctx context.Context, id string) (*UserResponse, error) {
	var r UserResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(GetUserPath, id), nil, &r)
	return &r, err
}

// GetLoggedUser : Return logged UserResponse.
// https://api.mangadex.org/docs.html#operation/get-user-follows-group
func (dc *DexClient) GetLoggedUser() (*UserResponse, error) {
	return dc.GetLoggedUserContext(context.Background())
}

// GetLoggedUserContext : GetLoggedUser with custom context.
func (dc *DexClient) GetLoggedUserContext(ctx context.Context) (*UserResponse, error) {
	var r UserResponse
	err := dc.responseOp(ctx, http.MethodGet, GetLoggedUserPath, nil, &r)
	return &r, err
}

// GetUserFollowedScanGroupList : Return list of followed ScanGroup.
// https://api.mangadex.org/docs.html#operation/get-user-follows-group
func (dc *DexClient) GetUserFollowedScanGroupList(limit, offset int) (*ScanGroupList, error) {
	return dc.GetUserFollowedScanGroupListContext(context.Background(), limit, offset)
}

// GetUserFollowedScanGroupListContext : GetUserFollowedScanGroupList with custom context.
func (dc *DexClient) GetUserFollowedScanGroupListContext(ctx context.Context, limit, offset int) (*ScanGroupList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = GetUserFollowedScanGroupListPath

	// Set query parameters
	q := u.Query()
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	var l ScanGroupList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// GetUserFollowedUsersList : Return list of followed User.
// https://api.mangadex.org/docs.html#operation/get-user-follows-user
func (dc *DexClient) GetUserFollowedUsersList(limit, offset int) (*UserList, error) {
	return dc.GetUserFollowedUsersListContext(context.Background(), limit, offset)
}

// GetUserFollowedUsersListContext : GetUserFollowedUsersList with custom context.
func (dc *DexClient) GetUserFollowedUsersListContext(ctx context.Context, limit, offset int) (*UserList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = GetUserFollowedUsersListPath

	// Set query parameters
	q := u.Query()
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	var l UserList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// GetUserFollowedMangaList : Return list of followed Manga.
// https://api.mangadex.org/docs.html#operation/get-user-follows-manga
func (dc *DexClient) GetUserFollowedMangaList(limit, offset int) (*MangaList, error) {
	return dc.GetUserFollowedMangaListContext(context.Background(), limit, offset)
}

// GetUserFollowedMangaListContext : GetUserFollowedMangaListPath with custom context.
func (dc *DexClient) GetUserFollowedMangaListContext(ctx context.Context, limit, offset int) (*MangaList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = GetUserFollowedMangaListPath

	// Set required query parameters
	q := u.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	var l MangaList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// From other operations

// GetUserFollowedMangaChapterFeed : Return Chapter feed.
// https://api.mangadex.org/docs.html#operation/get-user-follows-manga-feed
func (dc *DexClient) GetUserFollowedMangaChapterFeed(params url.Values) (*ChapterList, error) {
	return dc.GetUserFollowedMangaChapterFeedContext(context.Background(), params)
}

// GetUserFollowedMangaChapterFeedContext : GetUserFollowedMangaChapterFeedPath with custom context.
func (dc *DexClient) GetUserFollowedMangaChapterFeedContext(ctx context.Context, params url.Values) (*ChapterList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = GetUserFollowedMangaChapterFeedPath

	// Set query parameters
	u.RawQuery = params.Encode()

	var l ChapterList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// GetAllUserMangaReadingStatus : Get reading status for all manga for logged user.
// https://api.mangadex.org/docs.html#operation/get-manga-status
func (dc *DexClient) GetAllUserMangaReadingStatus() (*AllMangaReadingStatusResponse, error) {
	return dc.GetAllUserMangaReadingStatusContext(context.Background())
}

// GetAllUserMangaReadingStatusContext : GetAllUserMangaReadingStatus with custom context.
func (dc *DexClient) GetAllUserMangaReadingStatusContext(ctx context.Context) (*AllMangaReadingStatusResponse, error) {
	var r AllMangaReadingStatusResponse
	err := dc.responseOp(ctx, http.MethodGet, GetUserMangaReadingStatusPath, nil, &r)
	return &r, err
}

// GetUserCustomListList : Get list of custom lists.
// https://api.mangadex.org/docs.html#operation/get-user-list
func (dc *DexClient) GetUserCustomListList(limit, offset int) (*CustomListList, error) {
	return dc.GetUserCustomListListContext(context.Background(), limit, offset)
}

// GetUserCustomListListContext : GetUserCustomListList with custom context.
func (dc *DexClient) GetUserCustomListListContext(ctx context.Context, limit, offset int) (*CustomListList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = GetUserCustomListListPath

	// Set query parameters
	q := u.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	var l CustomListList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}
