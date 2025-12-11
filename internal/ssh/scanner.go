package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CommonKeyNames are typical SSH key filenames to look for
var CommonKeyNames = []string{
	"id_rsa",
	"id_ed25519",
	"id_ecdsa",
	"id_dsa",
	"id_ed25519_sk",
	"identity",
}

// Scanner handles scanning for SSH keys
type Scanner struct {
	sshDir string
}

func NewScanner() (*Scanner, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	return &Scanner{sshDir: sshDir}, nil
}

// ScanKeys scans the SSH directory for keys
func (s *Scanner) ScanKeys() ([]*Key, error) {
	// Check if .ssh directory exists
	if _, err := os.Stat(s.sshDir); os.IsNotExist(err) {
		return []*Key{}, nil // Return empty list if .ssh doesn't exist
	}

	entries, err := os.ReadDir(s.sshDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH directory: %w", err)
	}

	keyMap := make(map[string]*Key)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		fullPath := filepath.Join(s.sshDir, name)

		// Skip known non-key files
		if shouldSkipFile(name) {
			continue
		}

		// Check if it's a public key file
		if strings.HasSuffix(name, ".pub") {
			keyName := strings.TrimSuffix(name, ".pub")
			key, exists := keyMap[keyName]
			if !exists {
				key = &Key{
					Name: keyName,
					Path: filepath.Join(s.sshDir, keyName),
				}
				keyMap[keyName] = key
			}
			key.PublicKeyPath = fullPath
			key.HasPublic = true
		} else {
			// Check if it might be a private key
			if isPotentialPrivateKey(name) {
				key, exists := keyMap[name]
				if !exists {
					key = &Key{
						Name: name,
						Path: fullPath,
					}
					keyMap[name] = key
				}
				key.HasPrivate = true
			}
		}
	}

	// Convert map to slice and parse key details
	keys := make([]*Key, 0, len(keyMap))
	parser := NewParser()

	for _, key := range keyMap {
		// Only include keys that have at least a public key
		if key.HasPublic {
			if err := parser.ParseKey(key); err == nil {
				keys = append(keys, key)
			}
		}
	}

	return keys, nil
}

// shouldSkipFile returns true if the file should be skipped
func shouldSkipFile(name string) bool {
	skipPatterns := []string{
		"known_hosts",
		"authorized_keys",
		"config",
		".DS_Store",
	}

	for _, pattern := range skipPatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}
	return false
}

// isPotentialPrivateKey checks if a file might be a private key
func isPotentialPrivateKey(name string) bool {
	// Check against common key names
	for _, keyName := range CommonKeyNames {
		if name == keyName {
			return true
		}
	}

	// Check if it starts with common prefixes and doesn't have common extensions
	prefixes := []string{"id_", "identity", "ssh_", "key_"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) && !strings.Contains(name, ".") {
			return true
		}
	}

	return false
}
