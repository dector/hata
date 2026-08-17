package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"hata/internal/db"

	"github.com/sqids/sqids-go"
)

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
