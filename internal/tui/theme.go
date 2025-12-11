package tui

import (
	"github.com/gdamore/tcell/v2"
)

// Theme defines the color scheme for the application
type Theme struct {
	Background      tcell.Color
	Foreground      tcell.Color
	AccentPrimary   tcell.Color
	AccentSecondary tcell.Color
	Success         tcell.Color
	Warning         tcell.Color
	Error           tcell.Color
	Border          tcell.Color
	Title           tcell.Color
	StatusBar       tcell.Color
}

// DefaultTheme returns the minimalist default theme
func DefaultTheme() *Theme {
	return &Theme{
		Background:      tcell.ColorBlack,
		Foreground:      tcell.ColorWhite,
		AccentPrimary:   tcell.ColorSteelBlue,
		AccentSecondary: tcell.ColorDarkCyan,
		Success:         tcell.ColorGreen,
		Warning:         tcell.ColorYellow,
		Error:           tcell.ColorRed,
		Border:          tcell.ColorDimGray,
		Title:           tcell.ColorWhite,
		StatusBar:       tcell.ColorDimGray,
	}
}
