package main

import (
	"context"
	"log"

	"github.com/fantashley/hn-go"
)

func main() {
	client := hn.NewClient("", nil)
	ctx := context.Background()

	userID := "fantashley"
	user, err := client.GetUser(ctx, userID)
	if err != nil {
		log.Panicf("Error getting user %q: %v", userID, err)
	}

	log.Printf("%+v", user)

	var items []hn.Item

	for _, submitted := range user.Submitted {
		item, err := client.GetItem(ctx, submitted)
		if err != nil {
			log.Panicf("Error getting item %d: %v", submitted, err)
		}

		items = append(items, item)
	}

	log.Printf("%+v", items)
}
