package app

import (
	"github.com/mr-kaynak/go-ssh/internal/tui"
	"github.com/mr-kaynak/go-ssh/internal/update"
)

// Run starts the SSH key management application
func Run(version, commit, date string) error {
	// Check for updates in background (non-blocking)
	go update.CheckAndNotify(version)

	// Create and run the TUI application
	app := tui.NewApp(version)
	return app.Run()
}
