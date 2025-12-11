package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mr-kaynak/go-ssh/internal/ssh"
	"github.com/rivo/tview"
)

// CreateView displays a form for creating new SSH keys
type CreateView struct {
	*tview.Form
	onCreate func(name, comment, passphrase string, keyType ssh.KeyType)
	onCancel func()
}

func NewCreateView() *CreateView {
	form := tview.NewForm()
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetButtonBackgroundColor(tcell.ColorSteelBlue)
	form.SetButtonTextColor(tcell.ColorWhite)
	form.SetFieldBackgroundColor(tcell.ColorBlack)
	form.SetFieldTextColor(tcell.ColorWhite)
	form.SetLabelColor(tcell.ColorGray)

	cv := &CreateView{
		Form: form,
	}

	cv.setupForm()

	return cv
}

// setupForm sets up the form fields
func (cv *CreateView) setupForm() {
	// Key name field
	cv.AddInputField("Key Name", "id_ed25519_new", 40, nil, nil)

	// Key type dropdown
	keyTypes := []string{"ED25519 (recommended)", "RSA 4096", "ECDSA 521"}
	cv.AddDropDown("Key Type", keyTypes, 0, nil)

	// Comment field (optional)
	cv.AddInputField("Comment (optional)", "", 40, nil, nil)

	// Passphrase field (optional)
	cv.AddPasswordField("Passphrase (optional)", "", 40, '*', nil)

	// Passphrase confirmation
	cv.AddPasswordField("Confirm Passphrase", "", 40, '*', nil)

	// Buttons
	cv.AddButton("Create", func() {
		cv.handleCreate()
	})

	cv.AddButton("Cancel", func() {
		if cv.onCancel != nil {
			cv.onCancel()
		}
	})

	// Set up input capture for Escape key
	cv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			if cv.onCancel != nil {
				cv.onCancel()
			}
			return nil
		}
		return event
	})
}

// handleCreate handles the create button click
func (cv *CreateView) handleCreate() {
	if cv.onCreate == nil {
		return
	}

	// Get form values
	name := cv.GetFormItemByLabel("Key Name").(*tview.InputField).GetText()
	comment := cv.GetFormItemByLabel("Comment (optional)").(*tview.InputField).GetText()
	passphrase := cv.GetFormItemByLabel("Passphrase (optional)").(*tview.InputField).GetText()
	passphraseConfirm := cv.GetFormItemByLabel("Confirm Passphrase").(*tview.InputField).GetText()

	// Validate passphrase match
	if passphrase != passphraseConfirm {
		// TODO: Show error message
		return
	}

	// Get key type
	keyTypeIndex, _ := cv.GetFormItemByLabel("Key Type").(*tview.DropDown).GetCurrentOption()
	var keyType ssh.KeyType
	switch keyTypeIndex {
	case 0:
		keyType = ssh.KeyTypeED25519
	case 1:
		keyType = ssh.KeyTypeRSA
	case 2:
		keyType = ssh.KeyTypeECDSA
	default:
		keyType = ssh.KeyTypeED25519
	}

	cv.onCreate(name, comment, passphrase, keyType)
}

// OnCreate sets the callback for when create is confirmed
func (cv *CreateView) OnCreate(handler func(name, comment, passphrase string, keyType ssh.KeyType)) {
	cv.onCreate = handler
}

// OnCancel sets the callback for when creation is cancelled
func (cv *CreateView) OnCancel(handler func()) {
	cv.onCancel = handler
}

// GetTitle returns the view title
func (cv *CreateView) GetTitle() string {
	return " Create New SSH Key "
}

// Reset resets the form to default values
func (cv *CreateView) Reset() {
	cv.GetFormItemByLabel("Key Name").(*tview.InputField).SetText("id_ed25519_new")
	cv.GetFormItemByLabel("Key Type").(*tview.DropDown).SetCurrentOption(0)
	cv.GetFormItemByLabel("Comment (optional)").(*tview.InputField).SetText("")
	cv.GetFormItemByLabel("Passphrase (optional)").(*tview.InputField).SetText("")
	cv.GetFormItemByLabel("Confirm Passphrase").(*tview.InputField).SetText("")
}
