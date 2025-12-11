# go-ssh

A minimalist SSH key management tool with an interactive terminal UI.

[![Release](https://img.shields.io/github/v/release/mr-kaynak/go-ssh)](https://github.com/mr-kaynak/go-ssh/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/mr-kaynak/go-ssh/ci.yml?branch=main)](https://github.com/mr-kaynak/go-ssh/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/mr-kaynak/go-ssh)](https://goreportcard.com/report/github.com/mr-kaynak/go-ssh)
[![License](https://img.shields.io/github/license/mr-kaynak/go-ssh)](https://github.com/mr-kaynak/go-ssh/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mr-kaynak/go-ssh)](https://github.com/mr-kaynak/go-ssh/blob/main/go.mod)

## Features

- **Interactive TUI** - Clean, minimalist terminal interface built with tview
- **View SSH Keys** - List all SSH keys in your `~/.ssh` directory
- **Key Details** - View detailed information including fingerprints, type, and metadata
- **Copy to Clipboard** - Quickly copy public keys with a single keypress
- **Generate Keys** - Create new SSH keys interactively with various options
- **Secure** - Read-only by default, private keys never displayed
- **Cross-platform** - Works on macOS, Linux, and Windows

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap mr-kaynak/tap
brew install go-ssh
```

### Windows

#### Option 1: Download Binary (Recommended)
1. Download the latest Windows binary from [GitHub Releases](https://github.com/mr-kaynak/go-ssh/releases)
2. Extract the ZIP file
3. Add the binary to your PATH or run directly from the extracted folder

#### Option 2: Build from Source
```powershell
git clone https://github.com/mr-kaynak/go-ssh.git
cd go-ssh
go build -o go-ssh.exe ./cmd/go-ssh
```

#### Option 3: Using Go Install
```powershell
go install github.com/mr-kaynak/go-ssh/cmd/go-ssh@latest
```

### Linux/macOS (Alternative Methods)

#### Download Binary
Download the latest release from [GitHub Releases](https://github.com/mr-kaynak/go-ssh/releases)

```bash
# Linux (amd64)
curl -L https://github.com/mr-kaynak/go-ssh/releases/latest/download/go-ssh_linux_amd64.tar.gz | tar xz
sudo mv go-ssh /usr/local/bin/

# macOS (arm64)
curl -L https://github.com/mr-kaynak/go-ssh/releases/latest/download/go-ssh_darwin_arm64.tar.gz | tar xz
sudo mv go-ssh /usr/local/bin/
```

#### Build from Source
```bash
git clone https://github.com/mr-kaynak/go-ssh.git
cd go-ssh
make install
```

## Usage

Simply run:

```bash
go-ssh
```

### Keyboard Shortcuts

#### Main View
- `↑/↓` or `j/k` - Navigate through keys
- `Enter` - View key details
- `c` or `y` - Copy public key to clipboard
- `n` - Create new SSH key
- `?` - Show help
- `q` - Quit

#### Detail View
- `↑/↓` or `j/k` - Scroll through details
- `c` or `y` - Copy public key
- `q` or `b` - Back to list
- `Esc` - Back to list

#### Create View
- `Tab` - Next field
- `Shift+Tab` - Previous field
- `Enter` - Submit form
- `Esc` - Cancel

## Screenshots

### Main List View
```
┌─ SSH Keys (3) ──────────────────────────────────┐
│ Name           Type      Fingerprint             │
├─────────────────────────────────────────────────┤
│ id_ed25519     ED25519   SHA256:abc123...        │
│ id_rsa         RSA       SHA256:def456...        │
│ github_deploy  ED25519   SHA256:ghi789...        │
└─────────────────────────────────────────────────┘
 [Enter] view  [c] copy  [n] new  [q] quit  [?] help
```

### Key Details
```
┌─ id_ed25519 - Details ──────────────────────────┐
│ Type:         ED25519                            │
│ Fingerprint:  SHA256:abc123def456...             │
│ MD5:          12:34:56:78:90:ab:cd:ef            │
│                                                  │
│ Comment:      user@hostname                      │
│                                                  │
│ Private Key:  /Users/user/.ssh/id_ed25519       │
│ Public Key:   /Users/user/.ssh/id_ed25519.pub   │
│                                                  │
│ Status:       Private key is encrypted           │
│ Modified:     2025-12-10 15:30:45                │
└─────────────────────────────────────────────────┘
 [c] copy  [q] back
```

## Key Types Supported

- **ED25519** (Recommended) - Modern, fast, and secure
- **RSA** - Traditional, widely supported (2048, 3072, 4096 bits)
- **ECDSA** - Elliptic curve (256, 384, 521 bits)
- **DSA** - Legacy support only (not recommended)

## Security

- Private keys are **never displayed** in the UI
- Only public key content is shown and copied
- File permissions are checked for security warnings
- Encrypted private keys are detected and indicated
- All key generation uses `ssh-keygen` for maximum security
- Read-only operations by default

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for convenience)

### Build

```bash
make build
```

### Run in Development

```bash
make dev
```

### Run Tests

```bash
make test
```

### Create Release

```bash
# Tag a version
git tag v0.1.0
git push origin v0.1.0

# GoReleaser will automatically build and publish
```

## Project Structure

```
go-ssh/
├── cmd/go-ssh/         # Main entry point
├── internal/
│   ├── app/            # Application orchestrator
│   ├── ssh/            # SSH key operations
│   ├── tui/            # Terminal UI components
│   │   ├── views/      # UI views (list, detail, create)
│   │   └── components/ # Reusable UI components
│   ├── clipboard/      # Clipboard operations
│   └── config/         # Configuration
├── Makefile            # Build automation
└── .goreleaser.yaml    # Release configuration
```

## Dependencies

- [tview](https://github.com/rivo/tview) - Terminal UI framework
- [tcell](https://github.com/gdamore/tcell) - Terminal handling
- [golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) - SSH key parsing
- [clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [tview](https://github.com/rivo/tview) by rivo
- Inspired by the need for a simple, secure SSH key management tool

## Support

If you encounter any issues or have questions:

- Open an issue on [GitHub Issues](https://github.com/mr-kaynak/go-ssh/issues)
- Check existing issues for solutions

---

Made with ❤️ by [mr-kaynak](https://github.com/mr-kaynak)