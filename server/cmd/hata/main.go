package main

import (
	"context"
	"fmt"
	"os"

	"hata/internal/db"
)

func main() {
	// Check for CLI commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "manage":
			handleManageCommand()
			return
		case "serve":
			// Start server explicitly
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			printUsage()
			os.Exit(1)
		}
	}

	// Default behavior: start server
	database := db.New()

	ctx := context.Background()

	testDb(database, ctx)

	startServer(database, ctx)
}
