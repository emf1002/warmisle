package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  reset-password <username>  Reset user password to default")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "reset-password":
		fmt.Println("reset-password: not yet implemented")
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
