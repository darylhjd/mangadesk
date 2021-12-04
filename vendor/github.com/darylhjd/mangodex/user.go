package mangodex

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GetUserFollowedMangaListPath = "user/follows/manga"
)

// UserService : Provides User services provided by the API.
type UserService service

// GetUserFollowedMangaList : Return list of followed Manga.
// https://api.mangadex.org/docs.html#operation/get-user-follows-manga
func (s *UserService) GetUserFollowedMangaList(limit, offset int, includes []string) (*MangaList, error) {
	return s.GetUserFollowedMangaListContext(context.Background(), limit, offset, includes)
}

// GetUserFollowedMangaListContext : GetUserFollowedMangaListPath with custom context.
func (s *UserService) GetUserFollowedMangaListContext(ctx context.Context, limit, offset int, includes []string) (*MangaList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = GetUserFollowedMangaListPath

	// Set required query parameters
	q := u.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	for _, i := range includes {
		q.Add("includes[]", i)
	}
	u.RawQuery = q.Encode()

	var l MangaList
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}
