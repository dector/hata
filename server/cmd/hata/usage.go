package main

import "fmt"

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
