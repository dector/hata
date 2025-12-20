package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"hata/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	database := db.New()

	ctx := context.Background()

	testDb(database, ctx)

	startServer(database, ctx)
}

func startServer(database db.DB, ctx context.Context) {
	// Setup chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Add routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Hata!"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	host := "http://localhost"
	port := "8080"

	portAccess := port
	if os.Getenv("AIR") == "1" {
		portAccess = "3000"
	}

	// Start server
	log.Printf("Starting server on %s:%s\n", host, portAccess)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("failed starting server: %v", err)
	}
}

func testDb(database db.DB, ctx context.Context) {
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
