package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	GetMDHomeURLPath = "at-home/server/%s"
	MDHomeReportURL  = "https://api.mangadex.network/report"
)

type MDHomeClient struct {
	Client  http.Client
	BaseURL string
	Quality string
	Hash    string
}

type ReportPayload struct {
	URL      string
	Success  bool
	Cached   bool
	Bytes    int
	Duration int64
}

// NewMDHomeClient : Get MangaDex@Home client for a chapter.
// https://api.mangadex.org/docs.html#operation/get-at-home-server-chapterId
func (dc *DexClient) NewMDHomeClient(chapId, quality, hash string, forcePort443 bool) (*MDHomeClient, error) {
	return dc.NewMDHomeClientContext(context.Background(), chapId, quality, hash, forcePort443)
}

// NewMDHomeClientContext : NewMDHomeClient with custom context.
func (dc *DexClient) NewMDHomeClientContext(ctx context.Context, chapId, quality, hash string, forcePort443 bool) (*MDHomeClient, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(GetMDHomeURLPath, chapId)

	// Set query parameters
	q := u.Query()
	q.Add("forcePort443", strconv.FormatBool(forcePort443))
	u.RawQuery = q.Encode()

	var r map[string]string
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r)
	if err != nil {
		return nil, err
	}

	return &MDHomeClient{
		Client:  http.Client{},
		BaseURL: r["baseUrl"],
		Quality: quality,
		Hash:    hash,
	}, nil
}

// Report : Report success of getting chapter page data.
func (c *MDHomeClient) Report(r ReportPayload) (*http.Response, error) {
	return c.ReportContext(context.Background(), r)
}

// ReportContext : Report with custom context.
func (c *MDHomeClient) ReportContext(ctx context.Context, r ReportPayload) (*http.Response, error) {
	rBytes, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, MDHomeReportURL, bytes.NewBuffer(rBytes))
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// GetChapterPage : Return page data for a chapter with the filename of that page.
func (c *MDHomeClient) GetChapterPage(filename string) ([]byte, error) {
	return c.GetChapterPageWithContext(context.Background(), filename)
}

// GetChapterPageWithContext : GetChapterPage with custom context.
func (c *MDHomeClient) GetChapterPageWithContext(ctx context.Context, filename string) ([]byte, error) {
	path := strings.Join([]string{c.BaseURL, c.Quality, c.Hash, filename}, "/")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	r := ReportPayload{
		URL:      path,
		Success:  true,
		Cached:   false,
		Bytes:    0,
		Duration: 0,
	}

	start := time.Now()
	resp, err := c.Client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		var errM string
		if err != nil {
			errM = err.Error()
		} else {
			errM = fmt.Sprintf("%d status code", resp.StatusCode)
		}
		r.Success = false
		r.Duration = time.Since(start).Milliseconds()
		_, _ = c.ReportContext(ctx, r) // Make report
		return nil, fmt.Errorf("unable to get chapter data: %s", errM)
	}

	b, err := ioutil.ReadAll(resp.Body)
	r.Duration = time.Since(start).Milliseconds()
	r.Bytes = len(b)
	r.Cached = resp.Header.Get("X-Cache") == "HIT"
	_, _ = c.ReportContext(ctx, r) // Make report
	return b, err
}
