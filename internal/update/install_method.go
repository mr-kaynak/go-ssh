package update

import (
	"os"
	"path/filepath"
	"strings"
)

type InstallMethod int

const (
	InstallMethodUnknown InstallMethod = iota
	InstallMethodHomebrew
	InstallMethodGoInstall
	InstallMethodSystem
	InstallMethodBinary
)

func (m InstallMethod) String() string {
	switch m {
	case InstallMethodHomebrew:
		return "Homebrew"
	case InstallMethodGoInstall:
		return "go install"
	case InstallMethodSystem:
		return "system package manager"
	case InstallMethodBinary:
		return "binary"
	default:
		return "unknown"
	}
}

func DetectInstallMethod() InstallMethod {
	exe, err := os.Executable()
	if err != nil {
		return InstallMethodUnknown
	}

	exePath, err := filepath.EvalSymlinks(exe)
	if err != nil {
		exePath = exe
	}

	if strings.Contains(exePath, "/Cellar/") ||
		strings.Contains(exePath, "/opt/homebrew/") ||
		strings.Contains(exePath, "linuxbrew") {
		return InstallMethodHomebrew
	}

	if gopath := os.Getenv("GOPATH"); gopath != "" {
		gopathBin := filepath.Join(gopath, "bin")
		if strings.HasPrefix(exePath, gopathBin) {
			return InstallMethodGoInstall
		}
	}

	home, err := os.UserHomeDir()
	if err == nil {
		defaultGoPath := filepath.Join(home, "go", "bin")
		if strings.HasPrefix(exePath, defaultGoPath) {
			return InstallMethodGoInstall
		}
	}

	if strings.HasPrefix(exePath, "/usr/bin/") ||
		strings.HasPrefix(exePath, "/usr/local/bin/") {
		return InstallMethodSystem
	}

	return InstallMethodBinary
}

func GetUpdateInstructions(method InstallMethod, currentVer, latestVer string) string {
	switch method {
	case InstallMethodHomebrew:
		return "go-ssh was installed via Homebrew.\nTo update, run: brew upgrade go-ssh\n\n" +
			"Current version: " + currentVer + "\n" +
			"Latest version:  " + latestVer
	case InstallMethodGoInstall:
		return "go-ssh was installed via 'go install'.\nTo update, run: go install github.com/mr-kaynak/go-ssh/cmd/go-ssh@latest\n\n" +
			"Current version: " + currentVer + "\n" +
			"Latest version:  " + latestVer
	case InstallMethodSystem:
		return "go-ssh appears to be installed via a system package manager.\nPlease use your package manager to update.\n\n" +
			"Current version: " + currentVer + "\n" +
			"Latest version:  " + latestVer
	default:
		return ""
	}
}
