package ssh

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

// Parser handles parsing SSH keys
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// ParseKey parses an SSH key and extracts metadata
func (p *Parser) ParseKey(key *Key) error {
	if key.HasPublic && key.PublicKeyPath != "" {
		// Read public key content
		content, err := os.ReadFile(key.PublicKeyPath)
		if err != nil {
			return fmt.Errorf("failed to read public key: %w", err)
		}

		key.PublicKey = string(content)

		// Parse the public key
		publicKey, comment, _, _, err := ssh.ParseAuthorizedKey(content)
		if err != nil {
			return fmt.Errorf("failed to parse public key: %w", err)
		}

		// Set comment
		key.Comment = comment

		// Determine key type
		key.Type = p.detectKeyType(publicKey.Type())

		// Calculate fingerprints
		key.Fingerprint = p.calculateFingerprint(publicKey)
		key.FingerprintMD5 = p.calculateFingerprintMD5(publicKey)

		// Get file info
		if info, err := os.Stat(key.PublicKeyPath); err == nil {
			key.Modified = info.ModTime()
		}
	}

	// Check if private key is encrypted (if it exists)
	if key.HasPrivate && key.Path != "" {
		key.IsEncrypted = p.isPrivateKeyEncrypted(key.Path)

		// Get file info from private key if public doesn't exist
		if !key.HasPublic {
			if info, err := os.Stat(key.Path); err == nil {
				key.Modified = info.ModTime()
			}
		}
	}

	return nil
}

// detectKeyType detects the type of SSH key from the algorithm string
func (p *Parser) detectKeyType(algorithm string) KeyType {
	switch {
	case strings.Contains(algorithm, "rsa"):
		return KeyTypeRSA
	case strings.Contains(algorithm, "ed25519"):
		return KeyTypeED25519
	case strings.Contains(algorithm, "ecdsa"):
		return KeyTypeECDSA
	case strings.Contains(algorithm, "dsa"), strings.Contains(algorithm, "dss"):
		return KeyTypeDSA
	default:
		return KeyTypeUnknown
	}
}

// calculateFingerprint calculates the SHA256 fingerprint of a public key
func (p *Parser) calculateFingerprint(publicKey ssh.PublicKey) string {
	hash := sha256.Sum256(publicKey.Marshal())
	encoded := base64.RawStdEncoding.EncodeToString(hash[:])
	return fmt.Sprintf("SHA256:%s", encoded)
}

// calculateFingerprintMD5 calculates the MD5 fingerprint of a public key (legacy)
func (p *Parser) calculateFingerprintMD5(publicKey ssh.PublicKey) string {
	hash := md5.Sum(publicKey.Marshal())
	parts := make([]string, len(hash))
	for i, b := range hash {
		parts[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(parts, ":")
}

// isPrivateKeyEncrypted checks if a private key file is encrypted
func (p *Parser) isPrivateKeyEncrypted(path string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	contentStr := string(content)
	// Check for common encryption headers
	encryptionMarkers := []string{
		"ENCRYPTED",
		"Proc-Type: 4,ENCRYPTED",
		"DEK-Info:",
	}

	for _, marker := range encryptionMarkers {
		if strings.Contains(contentStr, marker) {
			return true
		}
	}

	return false
}
