package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	LastCheck   time.Time `json:"last_check"`
	LastVersion string    `json:"last_version"`
}

func GetStateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ".go-ssh")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "update.json"), nil
}

func LoadState() (*State, error) {
	path, err := GetStateFilePath()
	if err != nil {
		return &State{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{}, nil
		}
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return &State{}, nil
	}

	return &state, nil
}

func (s *State) Save() error {
	path, err := GetStateFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (s *State) ShouldCheck() bool {
	return time.Since(s.LastCheck) > 24*time.Hour
}
