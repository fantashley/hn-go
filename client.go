package hn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

const defaultBaseURL = "https://hacker-news.firebaseio.com/v0"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *Client) do(ctx context.Context, reqPath string, data interface{}) (int, error) {
	url, err := url.Parse(c.baseURL)
	if err != nil {
		return 0, fmt.Errorf("error parsing base URL: %w", err)
	}
	url.Path = path.Join(url.Path, reqPath)

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	rep, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error during request: %w", err)
	}
	defer rep.Body.Close()

	if err = json.NewDecoder(rep.Body).Decode(data); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}

	return rep.StatusCode, nil
}

func (c *Client) sortedStories(ctx context.Context, sort StorySortBy) ([]Item, error) {
	var (
		stories = make([]uint64, 0, 500)
		path    = string(sort) + "stories.json"
	)

	code, err := c.do(ctx, path, &stories)
	if err != nil {
		return nil, fmt.Errorf("error getting %s stories: %w", sort, err)
	} else if code != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", code)
	}

	return c.GetItems(ctx, stories)
}

func (c *Client) filteredStories(ctx context.Context, filter StoryFilter) ([]Item, error) {
	var (
		stories = make([]uint64, 0, 200)
		path    = string(filter) + "stories.json"
	)

	code, err := c.do(ctx, path, &stories)
	if err != nil {
		return nil, fmt.Errorf("error getting latest stories of type %q: %w", filter, err)
	} else if code != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", code)
	}

	return c.GetItems(ctx, stories)
}
