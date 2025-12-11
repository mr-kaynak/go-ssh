package components

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StatusBar represents a status bar component
type StatusBar struct {
	*tview.TextView
	message     string
	messageTime time.Time
}

// NewStatusBar creates a new status bar
func NewStatusBar() *StatusBar {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetTextAlign(tview.AlignLeft)
	tv.SetBackgroundColor(tcell.ColorBlack)

	sb := &StatusBar{
		TextView: tv,
	}

	// Show default keybindings
	sb.SetDefaultMessage()

	return sb
}

// SetDefaultMessage sets the default keybinding message
func (s *StatusBar) SetDefaultMessage() {
	s.SetText(" [::b]Enter[::-] view  [::b]c[::-] copy  [::b]n[::-] new  [::b]q[::-] quit  [::b]?[::-] help")
}

// SetMessage displays a temporary message
func (s *StatusBar) SetMessage(message string, color string) {
	s.message = message
	s.messageTime = time.Now()
	s.SetText(fmt.Sprintf(" [%s]%s[-]", color, message))

	// Automatically revert to default message after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		if time.Since(s.messageTime) >= 3*time.Second {
			s.SetDefaultMessage()
		}
	}()
}

// Success displays a success message
func (s *StatusBar) Success(message string) {
	s.SetMessage(message, "green")
}

// Error displays an error message
func (s *StatusBar) Error(message string) {
	s.SetMessage(message, "red")
}

// Info displays an info message
func (s *StatusBar) Info(message string) {
	s.SetMessage(message, "blue")
}
