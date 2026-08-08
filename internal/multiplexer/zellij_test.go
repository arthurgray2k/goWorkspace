package multiplexer

import (
	"strings"
	"testing"

	"github.com/arthurgray2k/goWorkspace/internal/config"
	"github.com/arthurgray2k/goWorkspace/internal/exec"
)

func TestGenerateKDLLayout(t *testing.T) {
	panes := []config.PaneConfig{
		{Name: "shell", Command: ""},
		{Name: "test", Command: ""},
		{Name: "git", Command: "git status"},
	}

	layout := GenerateKDLLayout("goJPeek", panes)

	if !strings.Contains(layout, "tab name=\"goJPeek\"") {
		t.Fatalf("layout missing tab name: %s", layout)
	}
	if !strings.Contains(layout, "pane name=\"shell\"") {
		t.Fatalf("layout missing shell pane: %s", layout)
	}
	if !strings.Contains(layout, "pane name=\"git\"") {
		t.Fatalf("layout missing git pane: %s", layout)
	}
	if !strings.Contains(layout, "git status") {
		t.Fatalf("layout missing git status command: %s", layout)
	}
}

func TestLaunchZellij(t *testing.T) {
	// Unset zellij environment variables for isolation during tests
	t.Setenv("ZELLIJ", "")
	t.Setenv("ZELLIJ_SESSION_NAME", "")

	t.Run("Zellij missing error", func(t *testing.T) {
		runner := exec.NewMockRunner()
		err := Launch(runner, "zellij", "goJPeek", "/tmp", nil, "default")
		if err == nil || !strings.Contains(err.Error(), "Zellij executable not found on PATH") {
			t.Fatalf("expected missing zellij error, got %v", err)
		}
	})

	t.Run("Zellij launch in current terminal", func(t *testing.T) {
		runner := exec.NewMockRunner()
		runner.Available["zellij"] = true

		panes := []config.PaneConfig{{Name: "shell", Command: ""}}
		err := Launch(runner, "zellij", "goJPeek", "/tmp", panes, "default")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(runner.RunCalls) != 1 || runner.RunCalls[0].Name != "/usr/bin/zellij" {
			t.Fatalf("expected zellij run call, got %#v", runner.RunCalls)
		}
	})

	t.Run("Zellij launch in Ghostty window", func(t *testing.T) {
		runner := exec.NewMockRunner()
		runner.Available["zellij"] = true
		runner.Available["ghostty"] = true

		panes := []config.PaneConfig{{Name: "shell", Command: ""}}
		err := Launch(runner, "zellij", "goJPeek", "/tmp", panes, "ghostty")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(runner.StartCalls) != 1 || runner.StartCalls[0].Name != "/usr/bin/ghostty" {
			t.Fatalf("expected ghostty start call, got %#v", runner.StartCalls)
		}
	})

	t.Run("None multiplexer", func(t *testing.T) {
		runner := exec.NewMockRunner()
		err := Launch(runner, "none", "goJPeek", "/tmp", nil, "default")
		if err != nil {
			t.Fatalf("unexpected error for none multiplexer: %v", err)
		}
	})
}
