package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"hata/internal/db"
	"hata/internal/util"

	"golang.org/x/term"
)

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
