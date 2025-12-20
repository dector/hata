package main

import (
	"context"
	"fmt"
	"log"

	"hata/internal/db"
)

func main() {
	database := db.New()

	ctx := context.Background()

	// Open database (uses default path: data/hata.db)
	if err := database.Open(ctx, ""); err != nil {
		log.Fatalf("failed opening database: %v", err)
	}
	defer database.Close()

	// Run the auto migration tool
	if err := database.RunMigrations(ctx); err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}

	// Get KV repository
	kvRepo := database.Repos().KV()

	// Count all key-value pairs
	count, err := kvRepo.Count(ctx)
	if err != nil {
		log.Fatalf("failed counting sys__kv: %v", err)
	}

	fmt.Printf("Total sys__kv records: %d\n", count)

	// List all keys
	keys, err := kvRepo.ListAllKeys(ctx)
	if err != nil {
		log.Fatalf("failed listing keys: %v", err)
	}

	fmt.Println("sys__kv records:")
	for _, key := range keys {
		value, err := kvRepo.GetByKey(ctx, key)
		if err != nil {
			log.Printf("  %s = <error: %v>\n", key, err)
			continue
		}
		fmt.Printf("  %s = %s\n", key, value)
	}
}
