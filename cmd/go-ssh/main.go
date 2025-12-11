package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mr-kaynak/go-ssh/internal/app"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Handle version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("go-ssh version %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	// Run the application
	if err := app.Run(version, commit, date); err != nil {
		log.Fatal(err)
	}
}
