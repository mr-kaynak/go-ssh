package ssh

import (
	"time"
)

// KeyType represents the type of SSH key
type KeyType string

const (
	KeyTypeRSA     KeyType = "RSA"
	KeyTypeED25519 KeyType = "ED25519"
	KeyTypeECDSA   KeyType = "ECDSA"
	KeyTypeDSA     KeyType = "DSA"
	KeyTypeUnknown KeyType = "Unknown"
)

// Key represents an SSH key with its metadata
type Key struct {
	Name           string    // Base filename (e.g., "id_rsa")
	Path           string    // Full path to the private key
	PublicKeyPath  string    // Full path to the public key
	Type           KeyType   // Key type (RSA, ED25519, etc.)
	Fingerprint    string    // SHA256 fingerprint
	FingerprintMD5 string    // MD5 fingerprint (legacy)
	Comment        string    // Key comment if present
	BitLength      int       // Bit length for RSA/DSA keys
	Created        time.Time // Creation time
	Modified       time.Time // Modification time
	HasPrivate     bool      // Whether private key exists
	HasPublic      bool      // Whether public key exists
	IsEncrypted    bool      // Whether private key is encrypted
	PublicKey      string    // Public key content
}
