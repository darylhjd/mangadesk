package mangodex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	BaseAPI  = "https://api.mangadex.org"
	PingPath = "ping"
)

type DexClient struct {
	client       http.Client
	header       http.Header
	logger       *log.Logger
	RefreshToken string
	isLoggedIn   bool
}

// NewDexClient : New anonymous client. To login as an authenticated user, use DexClient.Login.
func NewDexClient() *DexClient {
	// Create client
	client := http.Client{}

	// Create header
	header := http.Header{}
	header.Add("Accept", "application/json") // Set default accepted encoding

	// Create default logger for the client
	logger := log.New(os.Stderr, "mango", log.LstdFlags|log.Lshortfile)

	return &DexClient{
		client:     client,
		header:     header,
		logger:     logger,
		isLoggedIn: false,
	}
}

// Ping : Ping the API server.
func (dc *DexClient) Ping(ctx context.Context) error {
	u, _ := url.Parse(BaseAPI)
	u.Path = PingPath

	var res string
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	switch {
	case err != nil:
		return err
	case res != "pong":
		return errors.New("unexpected response for ping")
	default:
		return nil
	}
}

// Request : Sends a request to the MangaDex API.
func (dc *DexClient) Request(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	// Set header for HTTP Authentication
	req.Header = dc.header

	// Send request
	resp, err := dc.client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("non-2xx status code -> %d", resp.StatusCode)
	}

	return resp, nil
}

// RequestAndDecode : Convenience wrapper to also decode response to required data type
func (dc *DexClient) RequestAndDecode(ctx context.Context, method, url string, body io.Reader, s interface{}) (*http.Response, error) {
	resp, err := dc.Request(ctx, method, url, body)
	if err != nil {
		return resp, err
	} else if resp.StatusCode != 200 {
		return resp, nil
	}

	err = json.NewDecoder(resp.Body).Decode(s)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			dc.logger.Println("could not close response body.")
		}
	}(resp.Body)
	return resp, err
}
