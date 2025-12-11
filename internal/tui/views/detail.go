package views

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mr-kaynak/go-ssh/internal/ssh"
	"github.com/rivo/tview"
)

// DetailView displays detailed information about an SSH key
type DetailView struct {
	*tview.TextView
	key    *ssh.Key
	onBack func()
	onCopy func(*ssh.Key)
}

func NewDetailView() *DetailView {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetScrollable(true)
	tv.SetWordWrap(true)

	dv := &DetailView{
		TextView: tv,
	}

	dv.setupInputCapture()

	return dv
}

// setupInputCapture sets up keyboard shortcuts
func (dv *DetailView) setupInputCapture() {
	dv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if dv.onBack != nil {
				dv.onBack()
			}
			return nil
		}

		switch event.Rune() {
		case 'q', 'b':
			if dv.onBack != nil {
				dv.onBack()
			}
			return nil
		case 'c', 'y':
			if dv.onCopy != nil && dv.key != nil {
				dv.onCopy(dv.key)
			}
			return nil
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		}

		return event
	})
}

// SetKey sets the key to display
func (dv *DetailView) SetKey(key *ssh.Key) {
	dv.key = key
	dv.updateDisplay()
}

// updateDisplay refreshes the detail display
func (dv *DetailView) updateDisplay() {
	if dv.key == nil {
		dv.SetText("No key selected")
		return
	}

	var b strings.Builder

	// Title
	b.WriteString(fmt.Sprintf("[::b]%s[::-]\n\n", dv.key.Name))

	// Key Type
	typeColor := "blue"
	if dv.key.Type == ssh.KeyTypeED25519 {
		typeColor = "green"
	}
	b.WriteString(fmt.Sprintf("[gray]Type:[-]         [%s]%s[-]\n", typeColor, dv.key.Type))

	// Fingerprints
	b.WriteString(fmt.Sprintf("[gray]Fingerprint:[-]  %s\n", dv.key.Fingerprint))
	if dv.key.FingerprintMD5 != "" {
		b.WriteString(fmt.Sprintf("[gray]MD5:[-]          %s\n", dv.key.FingerprintMD5))
	}

	b.WriteString("\n")

	// Comment
	if dv.key.Comment != "" {
		b.WriteString(fmt.Sprintf("[gray]Comment:[-]      %s\n\n", dv.key.Comment))
	}

	// Paths
	b.WriteString(fmt.Sprintf("[gray]Private Key:[-]  %s\n", dv.key.Path))
	if dv.key.HasPublic {
		b.WriteString(fmt.Sprintf("[gray]Public Key:[-]   %s\n", dv.key.PublicKeyPath))
	}

	b.WriteString("\n")

	// Status
	if dv.key.HasPrivate {
		if dv.key.IsEncrypted {
			b.WriteString("[yellow]Status:[-]       Private key is encrypted\n")
		} else {
			b.WriteString("[green]Status:[-]       Private key exists\n")
		}
	} else {
		b.WriteString("[gray]Status:[-]       Public key only\n")
	}

	b.WriteString("\n")

	// Modified time
	if !dv.key.Modified.IsZero() {
		b.WriteString(fmt.Sprintf("[gray]Modified:[-]     %s\n", dv.key.Modified.Format("2006-01-02 15:04:05")))
	}

	b.WriteString("\n[gray]---[-]\n\n")

	// Public key content (truncated)
	if dv.key.PublicKey != "" {
		publicKey := strings.TrimSpace(dv.key.PublicKey)
		if len(publicKey) > 200 {
			publicKey = publicKey[:197] + "..."
		}
		b.WriteString(fmt.Sprintf("[gray]Public Key Content:[-]\n[dim]%s[-]\n", publicKey))
	}

	dv.SetText(b.String())
}

// OnBack sets the callback for when back is requested
func (dv *DetailView) OnBack(handler func()) {
	dv.onBack = handler
}

// OnCopy sets the callback for when copy is requested
func (dv *DetailView) OnCopy(handler func(*ssh.Key)) {
	dv.onCopy = handler
}

// GetTitle returns the view title
func (dv *DetailView) GetTitle() string {
	if dv.key != nil {
		return fmt.Sprintf(" %s - Details ", dv.key.Name)
	}
	return " Key Details "
}
