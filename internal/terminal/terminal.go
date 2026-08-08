package terminal

import (
	"fmt"
	"strings"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
)

// Launch executes a command inside the specified terminal emulator.
func Launch(runner exec.Runner, termType string, cmd string, args []string, projectDir string) error {
	normalized := strings.ToLower(strings.TrimSpace(termType))
	if normalized == "" || normalized == "none" || normalized == "default" || normalized == "system" {
		// No separate terminal window requested, run command directly in current context
		if cmd == "" {
			return nil
		}
		return runner.Run(projectDir, cmd, args...)
	}

	var binary string
	var termArgs []string

	switch normalized {
	case "ghostty":
		binary = "ghostty"
		if cmd != "" {
			termArgs = append([]string{"-e", cmd}, args...)
		}
	default:
		binary = normalized
		if cmd != "" {
			termArgs = append([]string{"-e", cmd}, args...)
		}
	}

	path, err := runner.LookPath(binary)
	if err != nil {
		return fmt.Errorf("terminal emulator executable %q not found on PATH. Please install %s or use --terminal default", binary, termType)
	}

	err = runner.Start(projectDir, path, termArgs...)
	if err != nil {
		return fmt.Errorf("failed to launch terminal %s (%s): %w", termType, path, err)
	}

	return nil
}
