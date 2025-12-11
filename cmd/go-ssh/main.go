package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mr-kaynak/go-ssh/internal/app"
	"github.com/mr-kaynak/go-ssh/internal/update"
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
		case "--update":
			handleUpdate(version)
			os.Exit(0)
		case "--check-update":
			checkUpdate(version)
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
  -h, --help           Show this help message
  -v, --version        Show version information
      --update         Update go-ssh to the latest version
      --check-update   Check if a new version is available

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

func handleUpdate(currentVersion string) {
	method := update.DetectInstallMethod()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	checker := update.NewChecker(currentVersion)
	release, found, err := checker.CheckForUpdate(ctx)

	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	if !found {
		fmt.Printf("You're already on the latest version: %s\n", currentVersion)
		return
	}

	if method != update.InstallMethodBinary {
		fmt.Println(update.GetUpdateInstructions(method, currentVersion, release.Version()))
		return
	}

	fmt.Printf("Update available: %s → %s\n\n", currentVersion, release.Version())
	fmt.Printf("Update now? [Y/n]: ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "" && response != "y" && response != "yes" {
		fmt.Println("Update cancelled")
		return
	}

	fmt.Println("Downloading update...")

	if err := checker.PerformUpdate(ctx, release); err != nil {
		fmt.Printf("Update failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Updated successfully to %s\n", release.Version())
	fmt.Println("Restart go-ssh to use the new version")
}

func checkUpdate(currentVersion string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	checker := update.NewChecker(currentVersion)
	release, found, err := checker.CheckForUpdate(ctx)

	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	if !found {
		fmt.Printf("You're on the latest version: %s\n", currentVersion)
		return
	}

	fmt.Printf("Update available: %s → %s\n", currentVersion, release.Version())
	fmt.Printf("\nRun 'go-ssh --update' to update\n")
}
