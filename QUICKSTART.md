# Quick Start Guide

Get up and running with go-ssh in 5 minutes!

## Installation

### macOS/Linux: Homebrew (Recommended)

```bash
brew tap mr-kaynak/tap
brew install go-ssh
```

### Windows: Download Binary

1. Go to [GitHub Releases](https://github.com/mr-kaynak/go-ssh/releases)
2. Download `go-ssh_windows_amd64.zip` (or `arm64` for ARM)
3. Extract the ZIP file
4. Run `go-ssh.exe`

Optional: Add to PATH for easy access from anywhere

### Alternative: Build from Source (All Platforms)

```bash
git clone https://github.com/mr-kaynak/go-ssh.git
cd go-ssh
make install  # macOS/Linux
# OR
go build -o go-ssh.exe ./cmd/go-ssh  # Windows
```

## First Run

Simply run:
```bash
go-ssh
```

The application will:
1. Scan your `~/.ssh` directory
2. Display all found SSH keys in an interactive list
3. Show key metadata (type, fingerprint, comment)

## Basic Usage

### Navigation
- Use `↑`/`↓` arrow keys or `j`/`k` to move through the list
- Press `Enter` to view detailed information about a key

### Copy Public Key
1. Select a key from the list
2. Press `c` or `y` to copy its public key to clipboard
3. Paste it wherever needed (GitHub, servers, etc.)

### Create New SSH Key
1. Press `n` to open the creation form
2. Fill in:
   - **Key Name**: e.g., `id_ed25519_github`
   - **Key Type**: Choose ED25519 (recommended), RSA 4096, or ECDSA 521
   - **Comment** (optional): e.g., `your.email@example.com`
   - **Passphrase** (optional): For additional security
3. Press `Enter` on "Create" button
4. Your new key will appear in the list

### Get Help
Press `?` at any time to view keyboard shortcuts and help information.

## Common Tasks

### Add SSH Key to GitHub
1. Run `go-ssh`
2. Navigate to your GitHub key (e.g., `id_ed25519`)
3. Press `c` to copy the public key
4. Go to GitHub → Settings → SSH Keys → New SSH Key
5. Paste the key

### View Key Fingerprint
1. Select a key from the list
2. Press `Enter` to view details
3. You'll see both SHA256 and MD5 fingerprints

### Check if Key is Encrypted
1. Select a key and press `Enter`
2. Look for the "Status" field
3. It will show if the private key is encrypted

## Tips

- **ED25519 is recommended** for new keys (faster, more secure, shorter)
- **Always use a passphrase** for keys that access important resources
- Keys are shown **read-only** - go-ssh never modifies existing keys
- The tool only displays **public key content** for security

## Keyboard Reference

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate |
| `Enter` | View details |
| `c` or `y` | Copy public key |
| `n` | Create new key |
| `q` | Quit / Back |
| `?` | Help |
| `Esc` | Back to list |

## Next Steps

- Read the full [README](README.md) for detailed documentation
- Check [CONTRIBUTING.md](CONTRIBUTING.md) to contribute
- Report issues on [GitHub Issues](https://github.com/mr-kaynak/go-ssh/issues)

Enjoy managing your SSH keys! 🚀