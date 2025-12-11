#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}go-ssh Installation Script${NC}"
echo ""

# Detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

# Detect shell
SHELL_RC=""
if [ -n "$ZSH_VERSION" ]; then
    SHELL_RC="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_RC="$HOME/.bashrc"
else
    SHELL_RC="$HOME/.profile"
fi

echo -e "Detected OS: ${YELLOW}$OS${NC}"
echo -e "Detected Architecture: ${YELLOW}$ARCH${NC}"
echo -e "Shell config: ${YELLOW}$SHELL_RC${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go first:"
    echo "  - macOS: brew install go"
    echo "  - Ubuntu/Debian: sudo apt install golang"
    echo "  - Fedora: sudo dnf install golang"
    echo "  - Or visit: https://golang.org/doc/install"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "Go version: ${GREEN}$GO_VERSION${NC}"
echo ""

# Install go-ssh
echo -e "${YELLOW}Installing go-ssh...${NC}"
go install github.com/mr-kaynak/go-ssh/cmd/go-ssh@latest

# Get GOPATH
GOPATH=$(go env GOPATH)
if [ -z "$GOPATH" ]; then
    GOPATH="$HOME/go"
fi
GOBIN="$GOPATH/bin"

echo -e "${GREEN}go-ssh installed to: $GOBIN/go-ssh${NC}"
echo ""

# Check if GOBIN is in PATH
if [[ ":$PATH:" == *":$GOBIN:"* ]]; then
    echo -e "${GREEN}$GOBIN is already in your PATH${NC}"
else
    echo -e "${YELLOW}Adding $GOBIN to PATH in $SHELL_RC${NC}"

    # Add to shell config
    if ! grep -q "export PATH=.*$GOBIN" "$SHELL_RC" 2>/dev/null; then
        echo "" >> "$SHELL_RC"
        echo "# Added by go-ssh installer" >> "$SHELL_RC"
        echo "export PATH=\"\$PATH:$GOBIN\"" >> "$SHELL_RC"
        echo -e "${GREEN}Added to $SHELL_RC${NC}"
        echo -e "${YELLOW}Please run: source $SHELL_RC${NC}"
    else
        echo -e "${GREEN}PATH already configured in $SHELL_RC${NC}"
    fi
fi

echo ""
echo -e "${GREEN}Installation complete!${NC}"
echo ""
echo "To start using go-ssh:"
echo -e "  1. ${YELLOW}source $SHELL_RC${NC} (or restart your terminal)"
echo -e "  2. Run: ${GREEN}go-ssh${NC}"
echo ""
echo "For more information: https://github.com/mr-kaynak/go-ssh"
