package views

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mr-kaynak/go-ssh/internal/ssh"
	"github.com/rivo/tview"
)

// ListView displays a list of SSH keys
type ListView struct {
	*tview.Table
	keys         []*ssh.Key
	onSelect     func(*ssh.Key)
	onCopy       func(*ssh.Key)
	onNew        func()
	onHelp       func()
	onQuit       func()
	filteredKeys []*ssh.Key
	filterText   string
}

func NewListView() *ListView {
	table := tview.NewTable()
	table.SetBorders(false)
	table.SetSelectable(true, false)
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(tcell.ColorSteelBlue).
		Foreground(tcell.ColorWhite))

	lv := &ListView{
		Table: table,
		keys:  []*ssh.Key{},
	}

	lv.setupHeader()
	lv.setupInputCapture()

	return lv
}

// setupHeader sets up the table header
func (lv *ListView) setupHeader() {
	headers := []string{"Name", "Type", "Fingerprint", "Comment"}
	for col, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorBlack).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false).
			SetExpansion(1)
		lv.SetCell(0, col, cell)
	}
}

// setupInputCapture sets up keyboard shortcuts
func (lv *ListView) setupInputCapture() {
	lv.SetInputCapture(lv.handleKeyEvent)
}

// handleKeyEvent handles keyboard events
func (lv *ListView) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	// Handle rune-based keys
	if event.Key() == tcell.KeyRune {
		if result := lv.handleRuneKey(event.Rune()); result != nil {
			return result
		}
		return nil
	}

	// Handle special keys
	if result := lv.handleSpecialKey(event.Key()); result != nil {
		return result
	}

	// Return original event if not handled
	return event
}

// handleRuneKey handles character key presses
func (lv *ListView) handleRuneKey(r rune) *tcell.EventKey {
	switch r {
	case 'q':
		if lv.onQuit != nil {
			lv.onQuit()
		}
	case 'c', 'y':
		lv.handleCopyKey()
	case 'n':
		if lv.onNew != nil {
			lv.onNew()
		}
	case '?':
		if lv.onHelp != nil {
			lv.onHelp()
		}
	case 'j':
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	case 'k':
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	case '/':
		// TODO: Implement search/filter
	default:
		return nil // Signal not handled
	}
	return tcell.NewEventKey(tcell.KeyRune, 0, tcell.ModNone) // Signal handled
}

// handleSpecialKey handles special key presses
func (lv *ListView) handleSpecialKey(key tcell.Key) *tcell.EventKey {
	if key == tcell.KeyEnter {
		if lv.onSelect != nil {
			selectedKey := lv.getSelectedKey()
			if selectedKey != nil {
				lv.onSelect(selectedKey)
			}
		}
		return tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone) // Signal handled
	}
	return nil // Signal not handled
}

// handleCopyKey handles copy action
func (lv *ListView) handleCopyKey() {
	if lv.onCopy != nil {
		key := lv.getSelectedKey()
		if key != nil {
			lv.onCopy(key)
		}
	}
}

// SetKeys updates the list with new keys
func (lv *ListView) SetKeys(keys []*ssh.Key) {
	lv.keys = keys
	lv.filteredKeys = keys
	lv.updateDisplay()
}

// updateDisplay refreshes the table display
func (lv *ListView) updateDisplay() {
	// Clear existing rows (keep header)
	for row := lv.GetRowCount() - 1; row > 0; row-- {
		lv.RemoveRow(row)
	}

	// Add keys to table
	for i, key := range lv.filteredKeys {
		row := i + 1

		// Truncate long fingerprints for display
		fingerprint := key.Fingerprint
		if len(fingerprint) > 40 {
			fingerprint = fingerprint[:37] + "..."
		}

		// Truncate comment if too long
		comment := key.Comment
		if len(comment) > 30 {
			comment = comment[:27] + "..."
		}

		// Name
		lv.SetCell(row, 0, tview.NewTableCell(key.Name).
			SetTextColor(tcell.ColorWhite).
			SetExpansion(1))

		// Type
		typeColor := tcell.ColorSteelBlue
		if key.Type == ssh.KeyTypeED25519 {
			typeColor = tcell.ColorGreen
		}
		lv.SetCell(row, 1, tview.NewTableCell(string(key.Type)).
			SetTextColor(typeColor).
			SetExpansion(1))

		// Fingerprint
		lv.SetCell(row, 2, tview.NewTableCell(fingerprint).
			SetTextColor(tcell.ColorGray).
			SetExpansion(2))

		// Comment
		lv.SetCell(row, 3, tview.NewTableCell(comment).
			SetTextColor(tcell.ColorDimGray).
			SetExpansion(1))
	}

	// Select first item if available
	if len(lv.filteredKeys) > 0 {
		lv.Select(1, 0)
	}
}

// getSelectedKey returns the currently selected key
func (lv *ListView) getSelectedKey() *ssh.Key {
	row, _ := lv.GetSelection()
	if row <= 0 || row > len(lv.filteredKeys) {
		return nil
	}
	return lv.filteredKeys[row-1]
}

// OnSelect sets the callback for when a key is selected
func (lv *ListView) OnSelect(handler func(*ssh.Key)) {
	lv.onSelect = handler
}

// OnCopy sets the callback for when copy is requested
func (lv *ListView) OnCopy(handler func(*ssh.Key)) {
	lv.onCopy = handler
}

// OnNew sets the callback for when new key is requested
func (lv *ListView) OnNew(handler func()) {
	lv.onNew = handler
}

// OnHelp sets the callback for when help is requested
func (lv *ListView) OnHelp(handler func()) {
	lv.onHelp = handler
}

// OnQuit sets the callback for when quit is requested
func (lv *ListView) OnQuit(handler func()) {
	lv.onQuit = handler
}

// GetTitle returns the view title
func (lv *ListView) GetTitle() string {
	count := len(lv.filteredKeys)
	total := len(lv.keys)

	if lv.filterText != "" {
		return fmt.Sprintf(" SSH Keys (%d/%d) - Filter: %s ", count, total, lv.filterText)
	}
	return fmt.Sprintf(" SSH Keys (%d) ", count)
}

// Filter filters the keys by the given text
func (lv *ListView) Filter(text string) {
	lv.filterText = text
	if text == "" {
		lv.filteredKeys = lv.keys
	} else {
		lv.filteredKeys = []*ssh.Key{}
		text = strings.ToLower(text)
		for _, key := range lv.keys {
			if strings.Contains(strings.ToLower(key.Name), text) ||
				strings.Contains(strings.ToLower(key.Comment), text) ||
				strings.Contains(strings.ToLower(string(key.Type)), text) {
				lv.filteredKeys = append(lv.filteredKeys, key)
			}
		}
	}
	lv.updateDisplay()
}
