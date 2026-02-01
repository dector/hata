package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"

	"hata/internal/api"
	"hata/internal/db"
	"hata/internal/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/term"
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

func startServer(database db.DB, ctx context.Context) {
	// Setup chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Create API handlers
	authHandler := api.NewAuthHandler(database.Repos())

	// Add routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Hata!"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api/latest", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
		})
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
	// Note: database remains open for the server to use

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

// CLI command handlers

func handleManageCommand() {
	if len(os.Args) < 3 {
		printManageUsage()
		os.Exit(1)
	}

	switch os.Args[2] {
	case "user":
		handleUserCommand()
	default:
		fmt.Printf("Unknown manage command: %s\n", os.Args[2])
		printManageUsage()
		os.Exit(1)
	}
}

func handleUserCommand() {
	if len(os.Args) < 4 {
		printUserUsage()
		os.Exit(1)
	}

	switch os.Args[3] {
	case "create":
		handleUserCreate()
	default:
		fmt.Printf("Unknown user command: %s\n", os.Args[3])
		printUserUsage()
		os.Exit(1)
	}
}

func handleUserCreate() {
	// Create a new FlagSet for the create command
	createFlags := flag.NewFlagSet("create", flag.ExitOnError)
	email := createFlags.String("email", "", "User's email address (used for login)")
	password := createFlags.String("password", "", "User's password")
	name := createFlags.String("name", "", "User's display name")

	createFlags.Usage = func() {
		printUserCreateUsage()
	}

	// Parse flags from remaining arguments
	createFlags.Parse(os.Args[4:])

	// Validate required flags
	if *email == "" || *name == "" {
		fmt.Println("Error: --email and --name flags are required")
		printUserCreateUsage()
		os.Exit(1)
	}

	// If password not provided, prompt for it interactively
	if *password == "" {
		var passwordStr string

		// Check if stdin is a terminal
		if term.IsTerminal(int(syscall.Stdin)) {
			// Interactive terminal - use secure password input
			fmt.Print("Password: ")
			passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
			fmt.Println() // Print newline after password input
			if err != nil {
				log.Fatalf("Failed to read password: %v", err)
			}
			passwordStr = string(passwordBytes)
		} else {
			// Non-interactive (pipe/redirect) - read from stdin
			fmt.Fprint(os.Stderr, "Password: ")
			fmt.Scanln(&passwordStr)
		}

		password = &passwordStr

		if *password == "" {
			fmt.Println("Error: password cannot be empty")
			os.Exit(1)
		}
	}

	// Initialize database
	database := db.New()
	ctx := context.Background()

	if err := database.Open(ctx, ""); err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	if err := database.RunMigrations(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Hash password
	passwordHash, err := util.HashPassword(*password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create user (email is stored in username field)
	user, err := database.Repos().User().Create(ctx, *email, passwordHash, *name)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("✓ User created successfully!\n")
	fmt.Printf("  ID: %d\n", user.ID)
	fmt.Printf("  Email: %s\n", user.Username)
	fmt.Printf("  Display Name: %s\n", user.DisplayName)
}

// Usage functions

func printUsage() {
	fmt.Println("Hata Server")
	fmt.Println("\nUsage:")
	fmt.Println("  hata                Start the server (default)")
	fmt.Println("  hata serve          Start the server")
	fmt.Println("  hata manage ...     Run management commands")
	fmt.Println("\nRun 'hata manage' for more information on management commands")
}

func printManageUsage() {
	fmt.Println("Hata Management Commands")
	fmt.Println("\nUsage:")
	fmt.Println("  hata manage user ...    User management commands")
	fmt.Println("\nRun 'hata manage user' for user-specific commands")
}

func printUserUsage() {
	fmt.Println("Hata User Management")
	fmt.Println("\nUsage:")
	fmt.Println("  hata manage user create --email <email> --name <name> [--password <password>]")
	fmt.Println("\nExamples:")
	fmt.Println("  hata manage user create --email alice@example.com --password secret123 --name \"Alice Smith\"")
	fmt.Println("  hata manage user create --email alice@example.com --name \"Alice Smith\"  # password will be prompted")
}

func printUserCreateUsage() {
	fmt.Println("Create a new user")
	fmt.Println("\nUsage:")
	fmt.Println("  hata manage user create --email <email> --name <name> [--password <password>]")
	fmt.Println("\nFlags:")
	fmt.Println("  --email      User's email address (required, used for login)")
	fmt.Println("  --name       User's display name (required)")
	fmt.Println("  --password   User's password (optional, will be prompted if not provided)")
	fmt.Println("\nExamples:")
	fmt.Println("  hata manage user create --email alice@example.com --password secret123 --name \"Alice Smith\"")
	fmt.Println("  hata manage user create --email bob@example.com --name \"Bob Jones\"  # interactive password")
	fmt.Println("  hata manage user create --name \"Charlie\" --email charlie@example.com  # flags in any order")
}
