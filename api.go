package hn

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"sync"

	"github.com/hashicorp/go-multierror"
)

const (
	apiPathItem = "item"
	apiPathUser = "user"
)

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

func (c *Client) GetUsers(ctx context.Context, ids []string) ([]User, error) {
	var (
		reqErr  error
		errMut  sync.Mutex
		users   = []User{}
		userMap = make(map[string]User, len(ids))
		mapMut  sync.Mutex
		wg      sync.WaitGroup
	)

	userRequest := func(id string) {
		defer wg.Done()
		user, err := c.GetUser(ctx, id)
		if err != nil {
			errMut.Lock()
			defer errMut.Unlock()
			reqErr = multierror.Append(reqErr, fmt.Errorf("error getting user %q: %w", id, err))
			return
		}

		mapMut.Lock()
		defer mapMut.Unlock()
		userMap[id] = user
	}

	for _, id := range ids {
		wg.Add(1)
		go userRequest(id)
	}
	wg.Wait()

	if len(userMap) == 0 && reqErr != nil {
		return nil, fmt.Errorf("error getting all users: %w", reqErr)
	} else if reqErr != nil {
		log.Printf("Errors encountered getting one or more users: %v", reqErr)
	}

	for _, id := range ids {
		if user, ok := userMap[id]; ok {
			users = append(users, user)
		}
	}

	return users, nil
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

func (c *Client) GetItems(ctx context.Context, ids []uint64) ([]Item, error) {
	var (
		reqErr  error
		errMut  sync.Mutex
		items   = []Item{}
		itemMap = make(map[uint64]Item, len(ids))
		mapMut  sync.Mutex
		wg      sync.WaitGroup
	)

	itemRequest := func(id uint64) {
		defer wg.Done()
		item, err := c.GetItem(ctx, id)
		if err != nil {
			errMut.Lock()
			defer errMut.Unlock()
			reqErr = multierror.Append(reqErr, fmt.Errorf("error getting item %d: %w", id, err))
			return
		}

		mapMut.Lock()
		defer mapMut.Unlock()
		itemMap[id] = item
	}

	for _, id := range ids {
		wg.Add(1)
		go itemRequest(id)
	}
	wg.Wait()

	if len(itemMap) == 0 && reqErr != nil {
		return nil, fmt.Errorf("error getting all items: %w", reqErr)
	} else if reqErr != nil {
		log.Printf("Errors encountered getting one or more items: %v", reqErr)
	}

	for _, id := range ids {
		if item, ok := itemMap[id]; ok {
			items = append(items, item)
		}
	}

	return items, nil
}

func (c *Client) Changes(ctx context.Context) ([]Item, []User, error) {
	var changes struct {
		Items    []uint64 `json:"items"`
		Profiles []string `json:"profiles"`
	}
	path := "updates.json"

	code, err := c.do(ctx, path, &changes)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting changes: %w", err)
	} else if code != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code %d", code)
	}

	items, err := c.GetItems(ctx, changes.Items)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting items: %w", err)
	}

	users, err := c.GetUsers(ctx, changes.Profiles)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting users: %w", err)
	}

	return items, users, nil
}

func (c *Client) MaxItemID(ctx context.Context) (uint64, error) {
	var (
		itemID uint64
		path   = "maxitem.json"
	)

	code, err := c.do(ctx, path, &itemID)
	if err != nil {
		return 0, fmt.Errorf("error getting max item ID: %w", err)
	} else if code != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code %d", code)
	}

	return itemID, nil
}

func (c *Client) NewStories(ctx context.Context) ([]Item, error) {
	return c.sortedStories(ctx, StorySortByNew)
}

func (c *Client) TopStories(ctx context.Context) ([]Item, error) {
	return c.sortedStories(ctx, StorySortByTop)
}

func (c *Client) BestStories(ctx context.Context) ([]Item, error) {
	return c.sortedStories(ctx, StorySortByBest)
}

func (c *Client) AskStories(ctx context.Context) ([]Item, error) {
	return c.filteredStories(ctx, StoryFilterAsk)
}

func (c *Client) ShowStories(ctx context.Context) ([]Item, error) {
	return c.filteredStories(ctx, StoryFilterShow)
}

func (c *Client) JobStories(ctx context.Context) ([]Item, error) {
	return c.filteredStories(ctx, StoryFilterJob)
}
