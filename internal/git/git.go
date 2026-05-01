package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsRepo checks if the current directory is a git repository
func IsRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// GetRepoName returns the name of the current directory
func GetRepoName() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return filepath.Base(dir)
}

// GetFiles returns a list of files tracked by git in the current repository
func GetFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return nil, nil
	}
	return strings.Split(trimmed, "\n"), nil
}
