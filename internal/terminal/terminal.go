package terminal

import (
	"os"
	"strings"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
)

// TerminalSpec describes how to launch a specific terminal emulator.
type TerminalSpec struct {
	Name string
	Exec func(cmd string, args []string) (string, []string)
}

var Candidates = []TerminalSpec{
	{
		Name: "ghostty",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"-e", cmd}, args...)
			return "ghostty", fullArgs
		},
	},
	{
		Name: "foot",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{cmd}, args...)
			return "foot", fullArgs
		},
	},
	{
		Name: "qterminal",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"-e", cmd}, args...)
			return "qterminal", fullArgs
		},
	},
	{
		Name: "lxterminal",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"-e", cmd}, args...)
			return "lxterminal", fullArgs
		},
	},
	{
		Name: "mate-terminal",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"--tab", "-e", cmd}, args...)
			return "mate-terminal", fullArgs
		},
	},
	{
		Name: "gnome-terminal",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"--tab", "--", cmd}, args...)
			return "gnome-terminal", fullArgs
		},
	},
	{
		Name: "konsole",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"--new-tab", "-e", cmd}, args...)
			return "konsole", fullArgs
		},
	},
	{
		Name: "xfce4-terminal",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"--tab", "-e", cmd}, args...)
			return "xfce4-terminal", fullArgs
		},
	},
	{
		Name: "tilix",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"-e", cmd}, args...)
			return "tilix", fullArgs
		},
	},
	{
		Name: "terminator",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"-e", cmd}, args...)
			return "terminator", fullArgs
		},
	},
	{
		Name: "kitty",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"-e", cmd}, args...)
			return "kitty", fullArgs
		},
	},
	{
		Name: "alacritty",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"-e", cmd}, args...)
			return "alacritty", fullArgs
		},
	},
	{
		Name: "terminology",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"-e", cmd}, args...)
			return "terminology", fullArgs
		},
	},
	{
		Name: "urxvt",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"-e", cmd}, args...)
			return "urxvt", fullArgs
		},
	},
	{
		Name: "x-terminal-emulator",
		Exec: func(cmd string, args []string) (string, []string) {
			fullArgs := append([]string{"-e", cmd}, args...)
			return "x-terminal-emulator", fullArgs
		},
	},
}

// BuildLaunchCommand attempts to find an available terminal emulator on PATH to launch targetCmd in a new tab/window.
func BuildLaunchCommand(runner exec.Runner, termType string, targetCmd string, targetArgs []string) (string, []string, bool) {
	norm := strings.ToLower(strings.TrimSpace(termType))
	if norm == "none" || norm == "same" || norm == "current" {
		return "", nil, false
	}

	// Do not attempt GUI terminal launch if headless environment without DISPLAY/WAYLAND_DISPLAY
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		return "", nil, false
	}

	// 1. If explicit terminal requested
	if norm != "" && norm != "default" && norm != "system" {
		for _, spec := range Candidates {
			if strings.EqualFold(spec.Name, norm) {
				path, err := runner.LookPath(spec.Name)
				if err == nil {
					_, args := spec.Exec(targetCmd, targetArgs)
					return path, args, true
				}
			}
		}
		// Custom binary fallback
		path, err := runner.LookPath(norm)
		if err == nil {
			fullArgs := append([]string{"-e", targetCmd}, targetArgs...)
			return path, fullArgs, true
		}
	}

	// 2. Default/auto-detect mode: search candidates in order of preference
	for _, spec := range Candidates {
		path, err := runner.LookPath(spec.Name)
		if err == nil {
			_, args := spec.Exec(targetCmd, targetArgs)
			return path, args, true
		}
	}

	return "", nil, false
}

// Launch executes a command inside a terminal window/tab or in current session if no GUI terminal is found.
func Launch(runner exec.Runner, termType string, cmd string, args []string, projectDir string) error {
	if cmd == "" {
		return nil
	}

	termBin, termArgs, ok := BuildLaunchCommand(runner, termType, cmd, args)
	if ok {
		return runner.Start(projectDir, termBin, termArgs...)
	}

	// Fallback to running directly in current session
	return runner.Run(projectDir, cmd, args...)
}
