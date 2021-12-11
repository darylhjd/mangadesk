package mangodex

import (
	"context"
	"net/http"
	"net/url"
)

const (
	MangaListPath = "manga"
)

// MangaService : Provides Manga services provided by the API.
type MangaService service

// MangaList : A response for getting a list of manga.
type MangaList struct {
	Result   string  `json:"result"`
	Response string  `json:"response"`
	Data     []Manga `json:"data"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Total    int     `json:"total"`
}

func (ml *MangaList) GetResult() string {
	return ml.Result
}

// Manga : Struct containing information on a Manga.
type Manga struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Attributes    MangaAttributes `json:"attributes"`
	Relationships []Relationship  `json:"relationships"`
}

// GetTitle : Get title of the Manga.
func (m *Manga) GetTitle(langCode string) string {
	if title := m.Attributes.Title.GetLocalString(langCode); title != "" {
		return title
	}
	return m.Attributes.AltTitles.GetLocalString(langCode)
}

// GetDescription : Get description of the Manga.
func (m *Manga) GetDescription(langCode string) string {
	return m.Attributes.Description.GetLocalString(langCode)
}

// MangaAttributes : Attributes for a Manga.
type MangaAttributes struct {
	Title                  LocalisedStrings `json:"title"`
	AltTitles              LocalisedStrings `json:"altTitles"`
	Description            LocalisedStrings `json:"description"`
	IsLocked               bool             `json:"isLocked"`
	Links                  LocalisedStrings `json:"links"`
	OriginalLanguage       string           `json:"originalLanguage"`
	LastVolume             *string          `json:"lastVolume"`
	LastChapter            *string          `json:"lastChapter"`
	PublicationDemographic *string          `json:"publicationDemographic"`
	Status                 *string          `json:"status"`
	Year                   *int             `json:"year"`
	ContentRating          *string          `json:"contentRating"`
	Tags                   []Tag            `json:"tags"`
	State                  string           `json:"state"`
	Version                int              `json:"version"`
	CreatedAt              string           `json:"createdAt"`
	UpdatedAt              string           `json:"updatedAt"`
}

// GetMangaList : Get a list of Manga.
// https://api.mangadex.org/docs.html#operation/get-search-manga
func (s *MangaService) GetMangaList(params url.Values) (*MangaList, error) {
	return s.GetMangaListContext(context.Background(), params)
}

// GetMangaListContext : GetMangaList with custom context.
func (s *MangaService) GetMangaListContext(ctx context.Context, params url.Values) (*MangaList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = MangaListPath

	// Set query parameters
	u.RawQuery = params.Encode()

	var l MangaList
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}
