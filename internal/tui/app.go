package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mr-kaynak/go-ssh/internal/clipboard"
	"github.com/mr-kaynak/go-ssh/internal/ssh"
	"github.com/mr-kaynak/go-ssh/internal/tui/components"
	"github.com/mr-kaynak/go-ssh/internal/tui/views"
	"github.com/rivo/tview"
)

// App represents the TUI application
type App struct {
	app       *tview.Application
	pages     *tview.Pages
	layout    *tview.Flex
	statusBar *components.StatusBar
	version   string

	// Views
	listView   *views.ListView
	detailView *views.DetailView
	createView *views.CreateView
	helpView   *views.HelpView

	// Services
	scanner   *ssh.Scanner
	generator *ssh.Generator

	// State
	keys []*ssh.Key
}

// NewApp creates a new TUI application
func NewApp(version string) *App {
	scanner, err := ssh.NewScanner()
	if err != nil {
		panic(fmt.Sprintf("Failed to create scanner: %v", err))
	}

	generator, err := ssh.NewGenerator()
	if err != nil {
		panic(fmt.Sprintf("Failed to create generator: %v", err))
	}

	return &App{
		app:       tview.NewApplication(),
		version:   version,
		scanner:   scanner,
		generator: generator,
	}
}

// Run starts the TUI application
func (a *App) Run() error {
	// Initialize components
	a.initComponents()

	// Load SSH keys
	if err := a.loadKeys(); err != nil {
		return fmt.Errorf("failed to load SSH keys: %w", err)
	}

	// Set up the layout
	a.setupLayout()

	// Show the list view by default
	a.showListView()

	// Run the application
	return a.app.Run()
}

// initComponents initializes all UI components
func (a *App) initComponents() {
	a.statusBar = components.NewStatusBar()
	a.listView = views.NewListView()
	a.detailView = views.NewDetailView()
	a.createView = views.NewCreateView()
	a.helpView = views.NewHelpView()

	// Set up list view callbacks
	a.listView.OnSelect(a.handleKeySelect)
	a.listView.OnCopy(a.handleKeyCopy)
	a.listView.OnNew(a.showCreateView)
	a.listView.OnHelp(a.showHelpView)
	a.listView.OnQuit(a.handleQuit)

	// Set up detail view callbacks
	a.detailView.OnBack(a.showListView)
	a.detailView.OnCopy(a.handleKeyCopy)

	// Set up create view callbacks
	a.createView.OnCreate(a.handleKeyCreate)
	a.createView.OnCancel(a.showListView)

	// Set up help view callback
	a.helpView.OnClose(a.showListView)
}

// setupLayout sets up the main layout
func (a *App) setupLayout() {
	a.pages = tview.NewPages()

	// Create main layout with status bar at the bottom
	a.layout = tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Add pages and status bar to layout
	a.layout.AddItem(a.pages, 0, 1, true)
	a.layout.AddItem(a.statusBar, 1, 0, false)

	// Add all pages
	a.addPage("list", a.listView, a.listView.GetTitle())
	a.addPage("detail", a.detailView, a.detailView.GetTitle())
	a.addPage("create", a.createView, a.createView.GetTitle())
	a.addPage("help", a.helpView, a.helpView.GetTitle())

	a.app.SetRoot(a.layout, true)
}

// addPage adds a page with a border and title
func (a *App) addPage(name string, primitive tview.Primitive, title string) {
	bordered := tview.NewFrame(primitive).
		SetBorders(0, 0, 1, 0, 0, 0).
		AddText(title, true, tview.AlignLeft, tcell.ColorWhite)

	a.pages.AddPage(name, bordered, true, false)
}

// loadKeys loads SSH keys from the file system
func (a *App) loadKeys() error {
	keys, err := a.scanner.ScanKeys()
	if err != nil {
		return err
	}

	a.keys = keys
	a.listView.SetKeys(keys)

	return nil
}

// showListView shows the list view
func (a *App) showListView() {
	// Reload keys to show any newly created keys
	if err := a.loadKeys(); err != nil {
		a.statusBar.Error(fmt.Sprintf("Failed to load keys: %v", err))
	}

	a.pages.SwitchToPage("list")
	a.app.SetFocus(a.listView)
	a.statusBar.SetDefaultMessage()
}

// showDetailView shows the detail view
func (a *App) showDetailView(key *ssh.Key) {
	a.detailView.SetKey(key)

	// Update the page title
	a.pages.RemovePage("detail")
	a.addPage("detail", a.detailView, a.detailView.GetTitle())

	a.pages.SwitchToPage("detail")
	a.app.SetFocus(a.detailView)

	// Update status bar
	a.statusBar.SetText(" [::b]c[::-] copy  [::b]q[::-] back")
}

// showCreateView shows the create view
func (a *App) showCreateView() {
	a.createView.Reset()
	a.pages.SwitchToPage("create")
	a.app.SetFocus(a.createView)

	// Update status bar
	a.statusBar.SetText(" [::b]Tab[::-] next field  [::b]Enter[::-] submit  [::b]Esc[::-] cancel")
}

// showHelpView shows the help view
func (a *App) showHelpView() {
	a.pages.SwitchToPage("help")
	a.app.SetFocus(a.helpView)

	// Update status bar
	a.statusBar.SetText(" Press any key to close")
}

// handleKeySelect handles key selection
func (a *App) handleKeySelect(key *ssh.Key) {
	a.showDetailView(key)
}

// handleKeyCopy handles copying a key to clipboard
func (a *App) handleKeyCopy(key *ssh.Key) {
	if key == nil {
		a.statusBar.Error("No key selected")
		return
	}

	if !key.HasPublic {
		a.statusBar.Error("No public key available")
		return
	}

	if err := clipboard.Copy(key.PublicKey); err != nil {
		a.statusBar.Error(fmt.Sprintf("Failed to copy: %v", err))
		return
	}

	a.statusBar.Success(fmt.Sprintf("Copied %s to clipboard", key.Name))
}

// handleKeyCreate handles creating a new key
func (a *App) handleKeyCreate(name, comment, passphrase string, keyType ssh.KeyType) {
	// Validate name
	if name == "" {
		a.statusBar.Error("Key name is required")
		return
	}

	// Check if key already exists
	if a.generator.KeyExists(name) {
		a.statusBar.Error(fmt.Sprintf("Key '%s' already exists", name))
		return
	}

	// Set default bits for key types
	bits := 0
	switch keyType {
	case ssh.KeyTypeRSA:
		bits = 4096
	case ssh.KeyTypeECDSA:
		bits = 521
	}

	// Create generation options
	opts := ssh.GeneratorOptions{
		Name:       name,
		Type:       keyType,
		Bits:       bits,
		Comment:    comment,
		Passphrase: passphrase,
	}

	// Generate the key
	if err := a.generator.Generate(opts); err != nil {
		a.statusBar.Error(fmt.Sprintf("Failed to create key: %v", err))
		return
	}

	a.statusBar.Success(fmt.Sprintf("Created key: %s", name))
	a.showListView()
}

// handleQuit quits the application
func (a *App) handleQuit() {
	a.app.Stop()
}
