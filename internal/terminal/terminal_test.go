package terminal

import (
	"testing"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
)

func TestBuildLaunchCommand(t *testing.T) {
	t.Setenv("DISPLAY", ":0")

	t.Run("Ghostty terminal detection", func(t *testing.T) {
		runner := exec.NewMockRunner()
		runner.Available["ghostty"] = true

		bin, args, ok := BuildLaunchCommand(runner, "ghostty", "zellij", []string{"attach"})
		if !ok || bin != "/usr/bin/ghostty" {
			t.Fatalf("expected ghostty launch command, got ok=%v, bin=%s", ok, bin)
		}
		if len(args) < 3 || args[0] != "-e" || args[1] != "zellij" {
			t.Fatalf("unexpected ghostty args: %#v", args)
		}
	})

	t.Run("Default terminal detection fallback", func(t *testing.T) {
		runner := exec.NewMockRunner()
		runner.Available["gnome-terminal"] = true

		bin, args, ok := BuildLaunchCommand(runner, "default", "zellij", []string{"attach"})
		if !ok || bin != "/usr/bin/gnome-terminal" {
			t.Fatalf("expected gnome-terminal, got ok=%v, bin=%s", ok, bin)
		}
		if len(args) < 4 || args[0] != "--tab" || args[1] != "--" {
			t.Fatalf("unexpected gnome-terminal args: %#v", args)
		}
	})
}

func TestLaunchTerminal(t *testing.T) {
	t.Setenv("DISPLAY", ":0")

	t.Run("Ghostty terminal launch", func(t *testing.T) {
		runner := exec.NewMockRunner()
		runner.Available["ghostty"] = true

		err := Launch(runner, "ghostty", "zellij", []string{"attach"}, "/path/to/project")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(runner.StartCalls) != 1 || runner.StartCalls[0].Name != "/usr/bin/ghostty" {
			t.Fatalf("expected ghostty start call, got %#v", runner.StartCalls)
		}
	})
}
