package multiplexer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arthurgray2k/goWorkspace/internal/config"
	"github.com/arthurgray2k/goWorkspace/internal/exec"
	"github.com/arthurgray2k/goWorkspace/internal/terminal"
)

// IsInsideZellij checks if the current execution environment is inside a Zellij session.
func IsInsideZellij() bool {
	return os.Getenv("ZELLIJ") != "" || os.Getenv("ZELLIJ_SESSION_NAME") != ""
}

// GenerateKDLLayout creates a valid Zellij KDL layout string for the given workspace panes.
func GenerateKDLLayout(sessionName string, panes []config.PaneConfig) string {
	var sb strings.Builder

	sb.WriteString("layout {\n")
	sb.WriteString("    default_tab_template {\n")
	sb.WriteString("        pane size=1 borderless=true {\n")
	sb.WriteString("            plugin location=\"zellij:tab-bar\"\n")
	sb.WriteString("        }\n")
	sb.WriteString("        children\n")
	sb.WriteString("        pane size=1 borderless=true {\n")
	sb.WriteString("            plugin location=\"zellij:status-bar\"\n")
	sb.WriteString("        }\n")
	sb.WriteString("    }\n")
	sb.WriteString(fmt.Sprintf("    tab name=\"%s\" {\n", escapeKDLString(sessionName)))

	userShell := os.Getenv("SHELL")
	if userShell == "" {
		userShell = "/bin/bash"
	}

	if len(panes) == 0 {
		sb.WriteString("        pane name=\"shell\"\n")
	} else {
		for _, pane := range panes {
			paneName := pane.Name
			if paneName == "" {
				paneName = "pane"
			}
			cmdStr := strings.TrimSpace(pane.Command)
			if cmdStr == "" {
				sb.WriteString(fmt.Sprintf("        pane name=\"%s\"\n", escapeKDLString(paneName)))
			} else {
				// Wrap initial command in user's shell so output displays and drops into an interactive shell
				shellCmd := fmt.Sprintf("%s; exec %s", cmdStr, userShell)
				sb.WriteString(fmt.Sprintf("        pane name=\"%s\" command=\"%s\" {\n", escapeKDLString(paneName), escapeKDLString(userShell)))
				sb.WriteString(fmt.Sprintf("            args \"-c\" \"%s\"\n", escapeKDLString(shellCmd)))
				sb.WriteString("        }\n")
			}
		}
	}

	sb.WriteString("    }\n")
	sb.WriteString("}\n")

	return sb.String()
}

func escapeKDLString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// Launch handles creation or attachment to a Zellij session using the generated layout.
func Launch(runner exec.Runner, muxType string, sessionName string, projectDir string, panes []config.PaneConfig, termType string) error {
	normalizedMux := strings.ToLower(strings.TrimSpace(muxType))
	if normalizedMux == "" || normalizedMux == "none" {
		return nil
	}

	if IsInsideZellij() {
		fmt.Printf("Notice: Already running inside Zellij session (%s). Skipping nested session launch.\n", os.Getenv("ZELLIJ_SESSION_NAME"))
		return nil
	}

	path, err := runner.LookPath("zellij")
	if err != nil {
		return fmt.Errorf("Zellij executable not found on PATH. Please install Zellij (https://zellij.dev) or run with --no-multiplexer")
	}

	// Sanitize session name
	cleanSessionName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, sessionName)

	if cleanSessionName == "" {
		cleanSessionName = "goworkspace"
	}

	// Create temp layout file
	kdlContent := GenerateKDLLayout(cleanSessionName, panes)
	tmpDir := os.TempDir()
	layoutPath := filepath.Join(tmpDir, fmt.Sprintf("gws-zellij-%s.kdl", cleanSessionName))

	err = os.WriteFile(layoutPath, []byte(kdlContent), 0600)
	if err != nil {
		return fmt.Errorf("failed to create temporary Zellij layout file: %w", err)
	}

	zellijArgs := []string{"--layout", layoutPath, "attach", "--create", cleanSessionName}

	// Launch inside a new window/tab of Ghostty or default system terminal emulator if available
	termBin, termArgs, ok := terminal.BuildLaunchCommand(runner, termType, path, zellijArgs)
	if ok {
		return runner.Start(projectDir, termBin, termArgs...)
	}

	// Fallback: Run Zellij directly in current terminal session
	return runner.Run(projectDir, path, zellijArgs...)
}
