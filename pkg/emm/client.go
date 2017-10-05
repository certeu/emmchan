package emm

import (
	"net/http"
	"time"
)

const (
	userAgent = "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.3; Trident/7.0; .NET4.0E; .NET4.0C)"
	timeout   = 30 * time.Second
)

type Client struct {
	httpClient *http.Client
	UserAgent  string
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	httpClient.Timeout = timeout

	c := &Client{
		httpClient: httpClient,
		UserAgent:  userAgent,
	}

	return c
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return c.httpClient.Do(req)
}