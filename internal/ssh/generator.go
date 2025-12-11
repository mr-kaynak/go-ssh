package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GeneratorOptions contains options for SSH key generation
type GeneratorOptions struct {
	Name       string  // Key name (e.g., "id_ed25519")
	Type       KeyType // Key type
	Bits       int     // Bit length (for RSA)
	Comment    string  // Optional comment
	Passphrase string  // Optional passphrase
}

// Generator handles SSH key generation
type Generator struct {
	sshDir string
}

// NewGenerator creates a new SSH key generator
func NewGenerator() (*Generator, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	sshDir := filepath.Join(homeDir, ".ssh")

	// Ensure .ssh directory exists
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	return &Generator{sshDir: sshDir}, nil
}

// Generate generates a new SSH key
func (g *Generator) Generate(opts GeneratorOptions) error {
	// Validate options
	if err := g.validateOptions(opts); err != nil {
		return err
	}

	// Build the key path
	keyPath := filepath.Join(g.sshDir, opts.Name)

	// Check if key already exists
	if _, err := os.Stat(keyPath); err == nil {
		return fmt.Errorf("key already exists: %s", keyPath)
	}

	// Build ssh-keygen command
	args := g.buildSSHKeygenArgs(keyPath, opts)

	// Execute ssh-keygen
	cmd := exec.Command("ssh-keygen", args...)

	// If no passphrase, provide empty input
	if opts.Passphrase == "" {
		cmd.Stdin = strings.NewReader("\n\n")
	} else {
		cmd.Stdin = strings.NewReader(opts.Passphrase + "\n" + opts.Passphrase + "\n")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate key: %w\nOutput: %s", err, string(output))
	}

	// Ensure proper permissions
	if err := os.Chmod(keyPath, 0600); err != nil {
		return fmt.Errorf("failed to set key permissions: %w", err)
	}

	return nil
}

// validateOptions validates the generation options
func (g *Generator) validateOptions(opts GeneratorOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("key name is required")
	}

	// Validate key name (no spaces, no slashes, reasonable length)
	if strings.ContainsAny(opts.Name, " /\\") {
		return fmt.Errorf("key name cannot contain spaces or slashes")
	}

	if len(opts.Name) > 255 {
		return fmt.Errorf("key name too long")
	}

	// Validate key type
	switch opts.Type {
	case KeyTypeRSA:
		if opts.Bits == 0 {
			opts.Bits = 4096 // Default to 4096 for RSA
		}
		if opts.Bits < 2048 {
			return fmt.Errorf("RSA key size must be at least 2048 bits")
		}
	case KeyTypeED25519:
		// ED25519 has fixed key size
	case KeyTypeECDSA:
		if opts.Bits == 0 {
			opts.Bits = 521 // Default to 521 for ECDSA
		}
		if opts.Bits != 256 && opts.Bits != 384 && opts.Bits != 521 {
			return fmt.Errorf("ECDSA key size must be 256, 384, or 521 bits")
		}
	default:
		return fmt.Errorf("unsupported key type: %s", opts.Type)
	}

	return nil
}

// buildSSHKeygenArgs builds the arguments for ssh-keygen command
func (g *Generator) buildSSHKeygenArgs(keyPath string, opts GeneratorOptions) []string {
	args := []string{
		"-f", keyPath,
		"-N", "", // Passphrase will be provided via stdin
	}

	// Add key type
	switch opts.Type {
	case KeyTypeRSA:
		args = append(args, "-t", "rsa")
		args = append(args, "-b", fmt.Sprintf("%d", opts.Bits))
	case KeyTypeED25519:
		args = append(args, "-t", "ed25519")
	case KeyTypeECDSA:
		args = append(args, "-t", "ecdsa")
		args = append(args, "-b", fmt.Sprintf("%d", opts.Bits))
	}

	// Add comment if provided
	if opts.Comment != "" {
		args = append(args, "-C", opts.Comment)
	}

	return args
}

// KeyExists checks if a key with the given name already exists
func (g *Generator) KeyExists(name string) bool {
	keyPath := filepath.Join(g.sshDir, name)
	_, err := os.Stat(keyPath)
	return err == nil
}
