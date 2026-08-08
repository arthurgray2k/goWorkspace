package editor

import (
	"fmt"
	"strings"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
)

// Launch executes the configured editor for a project directory.
func Launch(runner exec.Runner, editorType string, projectDir string) error {
	normalized := strings.ToLower(strings.TrimSpace(editorType))
	if normalized == "" || normalized == "none" {
		return nil
	}

	var binary string
	switch normalized {
	case "vscode", "code":
		binary = "code"
	case "zed":
		binary = "zed"
	default:
		// Attempt to use specified tool name as binary name
		binary = normalized
	}

	path, err := runner.LookPath(binary)
	if err != nil {
		return fmt.Errorf("editor executable %q not found on PATH. Please install %s or use --no-editor", binary, editorType)
	}

	// Launch editor as detached background process
	err = runner.Start(projectDir, path, ".")
	if err != nil {
		return fmt.Errorf("failed to launch editor %s (%s): %w", editorType, path, err)
	}

	return nil
}
