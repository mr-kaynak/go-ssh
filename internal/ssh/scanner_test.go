package ssh

import (
	"testing"
)

func TestShouldSkipFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"known_hosts file", "known_hosts", true},
		{"authorized_keys file", "authorized_keys", true},
		{"config file", "config", true},
		{"DS_Store file", ".DS_Store", true},
		{"private key", "id_rsa", false},
		{"public key", "id_rsa.pub", false},
		{"ed25519 key", "id_ed25519", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSkipFile(tt.filename); got != tt.want {
				t.Errorf("shouldSkipFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestIsPotentialPrivateKey(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"id_rsa", "id_rsa", true},
		{"id_ed25519", "id_ed25519", true},
		{"id_ecdsa", "id_ecdsa", true},
		{"identity", "identity", true},
		{"custom key_file", "key_deploy", true},
		{"public key", "id_rsa.pub", false},
		{"config", "config", false},
		{"known_hosts", "known_hosts", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPotentialPrivateKey(tt.filename); got != tt.want {
				t.Errorf("isPotentialPrivateKey(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}
