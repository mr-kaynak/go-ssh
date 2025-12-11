package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HelpView displays help information
type HelpView struct {
	*tview.TextView
	onClose func()
}

// NewHelpView creates a new help view
func NewHelpView() *HelpView {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetScrollable(true)

	hv := &HelpView{
		TextView: tv,
	}

	hv.setupContent()
	hv.setupInputCapture()

	return hv
}

// setupContent sets up the help content
func (hv *HelpView) setupContent() {
	content := `[::b]go-ssh - SSH Key Management[::-]

A minimalist tool for managing SSH keys with an interactive terminal UI.

[::b]Keyboard Shortcuts[::-]

  [::b]Navigation[::-]
  ↑/↓ or j/k      Navigate through keys
  Enter           View key details
  Esc or q        Go back / Quit
  ?               Show this help

  [::b]Actions[::-]
  c or y          Copy public key to clipboard
  n               Create new SSH key

  [::b]Detail View[::-]
  ↑/↓ or j/k      Scroll through details
  c or y          Copy public key
  q or b          Back to list

[::b]Features[::-]

  • View all SSH keys in ~/.ssh directory
  • Display key metadata (type, fingerprint, comment)
  • Copy public keys to clipboard
  • View detailed key information
  • Create new SSH keys interactively
  • Secure and read-only by default

[::b]Security Notes[::-]

  • Private keys are never displayed in the UI
  • Only public key content is shown and copied
  • File permissions are checked for security
  • Encrypted private keys are detected

[gray]Press any key to close this help...[-]`

	hv.SetText(content)
}

// setupInputCapture sets up keyboard shortcuts
func (hv *HelpView) setupInputCapture() {
	hv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if hv.onClose != nil {
			hv.onClose()
		}
		return nil
	})
}

// OnClose sets the callback for when help is closed
func (hv *HelpView) OnClose(handler func()) {
	hv.onClose = handler
}

// GetTitle returns the view title
func (hv *HelpView) GetTitle() string {
	return " Help "
}
