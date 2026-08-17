package main

import (
	"fmt"
	"os"
)

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
