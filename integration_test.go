//go:build integration
// +build integration

package hn_test

import (
	"context"
	"testing"

	"github.com/fantashley/hn-go"
)

const testUser = "fantashley"

func TestAPI(t *testing.T) {
	ctx := context.Background()
	client := hn.NewClient("", nil)
	testUserItems(ctx, t, client)
	testStories(ctx, t, client)
}

func testUserItems(ctx context.Context, t *testing.T, c *hn.Client) {
	user, err := c.GetUser(ctx, testUser)
	if err != nil {
		t.Fatalf("Error getting user: %v", err)
	}

	items, err := c.GetItems(ctx, user.Submitted)
	if err != nil {
		t.Errorf("Error getting items for user %q: %v", testUser, err)
	}

	t.Logf("User %q has %d items", user.ID, len(items))

	items, users, err := c.Changes(ctx)
	if err != nil {
		t.Fatalf("Error getting changed items and users: %v", err)
	}

	t.Logf("%d items and %d users changed", len(items), len(users))
}

func testStories(ctx context.Context, t *testing.T, c *hn.Client) {
	tests := []struct {
		name string
		fn   func(context.Context) ([]hn.Item, error)
	}{
		{
			name: "Top Stories",
			fn:   c.TopStories,
		},
		{
			name: "Best Stories",
			fn:   c.BestStories,
		},
		{
			name: "New Stories",
			fn:   c.NewStories,
		},
		{
			name: "Ask Stories",
			fn:   c.AskStories,
		},
		{
			name: "Job Stories",
			fn:   c.JobStories,
		},
		{
			name: "Show Stories",
			fn:   c.ShowStories,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stories, err := test.fn(ctx)
			if err != nil {
				t.Fatalf("Error getting stories: %v", err)
			}

			t.Logf("%d stories found", len(stories))
		})
	}
}
