package ssh

import (
	"testing"
)

func TestDetectKeyType(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name      string
		algorithm string
		want      KeyType
	}{
		{"RSA key", "ssh-rsa", KeyTypeRSA},
		{"ED25519 key", "ssh-ed25519", KeyTypeED25519},
		{"ECDSA key", "ecdsa-sha2-nistp256", KeyTypeECDSA},
		{"DSA key", "ssh-dss", KeyTypeDSA},
		{"Unknown key", "unknown-algo", KeyTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.detectKeyType(tt.algorithm); got != tt.want {
				t.Errorf("detectKeyType(%q) = %v, want %v", tt.algorithm, got, tt.want)
			}
		})
	}
}
