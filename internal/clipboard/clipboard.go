package clipboard

import (
	"fmt"

	"github.com/atotto/clipboard"
)

// Copy copies text to the system clipboard
func Copy(text string) error {
	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}
	return nil
}

// Paste gets text from the system clipboard
func Paste() (string, error) {
	text, err := clipboard.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read from clipboard: %w", err)
	}
	return text, nil
}
