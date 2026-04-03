package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	"hata/internal/api"
	"hata/internal/db"
	"hata/internal/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sqids/sqids-go"
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
	serverHandler := api.NewServerHandler()
	houseHandler := api.NewHouseHandler(database.Repos())
	deviceHandler := api.NewDeviceHandler(database.Repos())

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
		r.Get("/house", houseHandler.List)
		r.Get("/house/{houseId}/device", deviceHandler.ListByHouse)
		r.Get("/device", deviceHandler.ListByUser)
		r.Get("/ping", serverHandler.Ping)
	})

	host := "http://localhost"
	port := "4501"

	portAccess := port
	if os.Getenv("AIR_PROXY") == "1" {
		portAccess = "4500"
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
	case "house":
		handleHouseCommand()
	default:
		fmt.Printf("Unknown manage command: %s\n", os.Args[2])
		printManageUsage()
		os.Exit(1)
	}
}

func handleHouseCommand() {
	if len(os.Args) < 4 {
		printHouseUsage()
		os.Exit(1)
	}

	switch os.Args[3] {
	case "create":
		handleHouseCreate()
	case "assign":
		handleHouseAssign()
	default:
		fmt.Printf("Unknown house command: %s\n", os.Args[3])
		printHouseUsage()
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
	case "list":
		handleUserList()
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

func handleUserList() {
	// Create a new FlagSet for the list command
	listFlags := flag.NewFlagSet("list", flag.ExitOnError)
	limit := listFlags.Int("limit", 10, "Maximum number of users to display")
	offset := listFlags.Int("offset", 0, "Number of users to skip")

	listFlags.Usage = func() {
		printUserListUsage()
	}

	// Parse flags from remaining arguments
	listFlags.Parse(os.Args[4:])

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

	// List users
	users, total, err := database.Repos().User().List(ctx, *limit, *offset)
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}

	// Display users
	if len(users) == 0 {
		fmt.Println("No users found")
		fmt.Printf("\nUsers: 0/%d\n", total)
		return
	}

	fmt.Printf("%-5s %-30s %-30s\n", "ID", "Display Name", "Username")
	fmt.Println("─────────────────────────────────────────────────────────────────────")
	for _, user := range users {
		fmt.Printf("%-5d %-30s %-30s\n", user.ID, user.DisplayName, user.Username)
	}
	fmt.Printf("\nUsers: %d/%d\n", len(users), total)
}

func handleHouseCreate() {
	createFlags := flag.NewFlagSet("create", flag.ExitOnError)
	name := createFlags.String("name", "", "House display name")

	createFlags.Usage = func() {
		printHouseCreateUsage()
	}

	createFlags.Parse(os.Args[4:])

	if strings.TrimSpace(*name) == "" {
		fmt.Println("Error: --name flag is required")
		printHouseCreateUsage()
		os.Exit(1)
	}

	houseID, err := generateHouseID()
	if err != nil {
		log.Fatalf("Failed to generate house ID: %v", err)
	}

	database := db.New()
	ctx := context.Background()

	if err := database.Open(ctx, ""); err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	if err := database.RunMigrations(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	house, err := database.Repos().House().Create(ctx, houseID, strings.TrimSpace(*name))
	if err != nil {
		log.Fatalf("Failed to create house: %v", err)
	}

	fmt.Printf("✓ House created successfully!\n")
	fmt.Printf("  ID: %s\n", house.ID)
	fmt.Printf("  Display Name: %s\n", house.DisplayName)
}

func handleHouseAssign() {
	if len(os.Args) < 7 {
		fmt.Println("Error: missing arguments")
		printHouseAssignUsage()
		os.Exit(1)
	}

	houseID := os.Args[4]
	username := os.Args[5]
	role, ok := normalizeHouseRole(os.Args[6])
	if !ok {
		fmt.Printf("Error: invalid role %q\n", os.Args[6])
		printHouseAssignUsage()
		os.Exit(1)
	}

	database := db.New()
	ctx := context.Background()

	if err := database.Open(ctx, ""); err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	if err := database.RunMigrations(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	house, err := database.Repos().House().GetByID(ctx, houseID)
	if err != nil {
		log.Fatalf("Failed to fetch house: %v", err)
	}
	if house == nil {
		log.Fatalf("House %q not found", houseID)
	}

	user, err := database.Repos().User().GetByUsername(ctx, username)
	if err != nil {
		log.Fatalf("Failed to fetch user: %v", err)
	}
	if user == nil {
		log.Fatalf("User %q not found", username)
	}

	assignment, err := database.Repos().HouseRole().Assign(ctx, houseID, user.ID, role)
	if err != nil {
		log.Fatalf("Failed to assign house role: %v", err)
	}

	fmt.Printf("✓ House role assigned successfully!\n")
	fmt.Printf("  House ID: %s\n", assignment.HouseID)
	fmt.Printf("  User: %s\n", username)
	fmt.Printf("  Role: %s\n", assignment.Role)
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
	fmt.Println("  hata manage house ...   House management commands")
	fmt.Println("\nRun 'hata manage user' or 'hata manage house' for more information")
}

func printUserUsage() {
	fmt.Println("Hata User Management")
	fmt.Println("\nUsage:")
	fmt.Println("  hata manage user create --email <email> --name <name> [--password <password>]")
	fmt.Println("  hata manage user list [--limit <n>] [--offset <n>]")
	fmt.Println("\nExamples:")
	fmt.Println("  hata manage user create --email alice@example.com --password secret123 --name \"Alice Smith\"")
	fmt.Println("  hata manage user create --email alice@example.com --name \"Alice Smith\"  # password will be prompted")
	fmt.Println("  hata manage user list")
	fmt.Println("  hata manage user list --limit 20 --offset 10")
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

func printUserListUsage() {
	fmt.Println("List users with pagination")
	fmt.Println("\nUsage:")
	fmt.Println("  hata manage user list [--limit <n>] [--offset <n>]")
	fmt.Println("\nFlags:")
	fmt.Println("  --limit      Maximum number of users to display (default: 10)")
	fmt.Println("  --offset     Number of users to skip (default: 0)")
	fmt.Println("\nExamples:")
	fmt.Println("  hata manage user list")
	fmt.Println("  hata manage user list --limit 20")
	fmt.Println("  hata manage user list --limit 20 --offset 10")
}

func printHouseUsage() {
	fmt.Println("Hata House Management")
	fmt.Println("\nUsage:")
	fmt.Println("  hata manage house create --name <display name>")
	fmt.Println("  hata manage house assign <house-id> <username> <role>")
	fmt.Println("\nRoles:")
	fmt.Println("  admin | owner | habitant | guest")
	fmt.Println("\nExamples:")
	fmt.Println("  hata manage house create --name \"Main Home\"")
	fmt.Println("  hata manage house assign H7k1a user@example.com admin")
}

func printHouseCreateUsage() {
	fmt.Println("Create a new house")
	fmt.Println("\nUsage:")
	fmt.Println("  hata manage house create --name <display name>")
	fmt.Println("\nFlags:")
	fmt.Println("  --name     House display name (required)")
	fmt.Println("\nExamples:")
	fmt.Println("  hata manage house create --name \"Main Home\"")
}

func printHouseAssignUsage() {
	fmt.Println("Assign a user to a house with a role")
	fmt.Println("\nUsage:")
	fmt.Println("  hata manage house assign <house-id> <username> <role>")
	fmt.Println("\nRoles:")
	fmt.Println("  admin | owner | habitant | guest")
	fmt.Println("\nExamples:")
	fmt.Println("  hata manage house assign H7k1a user@example.com owner")
}

func normalizeHouseRole(role string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "admin", "owner", "habitant", "guest":
		return strings.ToLower(strings.TrimSpace(role)), true
	default:
		return "", false
	}
}

func generateHouseID() (string, error) {
	s, err := sqids.New()
	if err != nil {
		return "", fmt.Errorf("failed to create sqids instance: %w", err)
	}

	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	n := binary.BigEndian.Uint64(buf[:])
	if n == 0 {
		n = 1
	}

	id, err := s.Encode([]uint64{n})
	if err != nil {
		return "", fmt.Errorf("failed to encode house id: %w", err)
	}

	return id, nil
}
