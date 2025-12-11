package app

import (
	"github.com/mr-kaynak/go-ssh/internal/tui"
)

// Run starts the SSH key management application
func Run(version, commit, date string) error {
	// Create and run the TUI application
	app := tui.NewApp(version)
	return app.Run()
}
