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
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("go-ssh version %s (commit: %s, built: %s)\n", version, commit, date)
			os.Exit(0)
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		}
	}

	if err := app.Run(version, commit, date); err != nil {
		log.Fatal(err)
	}
}

func printHelp() {
	fmt.Printf(`go-ssh - Minimalist SSH key management with interactive TUI

Usage:
  go-ssh [flags]

Flags:
  -h, --help       Show this help message
  -v, --version    Show version information

Interactive Commands (once running):
  j/k, ↑/↓        Navigate through keys
  Enter           View key details
  c, y            Copy public key to clipboard
  n               Create new SSH key
  ?               Show help
  q               Quit

Examples:
  go-ssh                 # Launch interactive TUI
  go-ssh --version       # Show version
  go-ssh --help          # Show this help

For more information: https://github.com/mr-kaynak/go-ssh
`)
}
