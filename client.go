package hn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	apiPathItem = "item"
	apiPathUser = "user"
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

func (c *Client) GetUser(ctx context.Context, id string) (User, error) {
	var (
		user User
		path = path.Join(apiPathUser, id+".json")
	)

	code, err := c.do(ctx, path, &user)
	if err != nil {
		return User{}, fmt.Errorf("error getting user %q: %w", id, err)
	} else if code != http.StatusOK {
		return User{}, fmt.Errorf("unexpected status code %d", code)
	}

	return user, nil
}

func (c *Client) GetItem(ctx context.Context, id uint64) (Item, error) {
	var (
		item Item
		path = path.Join(apiPathItem, strconv.FormatUint(id, 10)+".json")
	)

	code, err := c.do(ctx, path, &item)
	if err != nil {
		return Item{}, fmt.Errorf("error getting item %d: %w", id, err)
	} else if code != http.StatusOK {
		return Item{}, fmt.Errorf("unexpected status code %d", code)
	}

	return item, nil
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
