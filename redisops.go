package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// Updates the Redis data store based on the fetched data
func updateRedisDataStore(ctx context.Context, data []Item, client *redis.Client) {
	// Update Redis data store based on the fetched data
	for _, item := range data {
		key := fmt.Sprintf("%s:%s", item.Type, item.ID)
		err := client.HSet(ctx, key, "Data", item.Data).Err()
		if err != nil {
			log.Printf("Error updating Redis data store for key %s: %v\n", key, err)
		}
	}
}
