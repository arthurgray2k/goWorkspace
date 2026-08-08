package multiplexer

import (
	"strings"
	"testing"

	"github.com/arthurgray2k/goWorkspace/internal/config"
	"github.com/arthurgray2k/goWorkspace/internal/exec"
)

func TestGenerateKDLLayout(t *testing.T) {
	tabs := []config.TabConfig{
		{
			Name: "Code",
			Panes: []config.PaneConfig{
				{Name: "shell", Command: ""},
				{Name: "git", Command: "git status", Dir: "./subfolder"},
			},
		},
		{
			Name: "Services",
			Panes: []config.PaneConfig{
				{Name: "server", Command: "go run main.go", Env: map[string]string{"PORT": "8080"}},
			},
		},
	}

	layout := GenerateKDLLayout("goJPeek", tabs, nil, "/tmp/project")

	if !strings.Contains(layout, "tab name=\"Code\"") {
		t.Fatalf("layout missing tab name Code: %s", layout)
	}
	if !strings.Contains(layout, "tab name=\"Services\"") {
		t.Fatalf("layout missing tab name Services: %s", layout)
	}
	if !strings.Contains(layout, "cwd=\"/tmp/project/subfolder\"") {
		t.Fatalf("layout missing subfolder cwd: %s", layout)
	}
	if !strings.Contains(layout, "PORT=8080") {
		t.Fatalf("layout missing env var PORT=8080: %s", layout)
	}
}

func TestLaunchZellij(t *testing.T) {
	t.Setenv("ZELLIJ", "")
	t.Setenv("ZELLIJ_SESSION_NAME", "")

	t.Run("Zellij missing error", func(t *testing.T) {
		runner := exec.NewMockRunner()
		err := Launch(runner, "zellij", "goJPeek", "/tmp", nil, nil, "default")
		if err == nil || !strings.Contains(err.Error(), "Zellij executable not found on PATH") {
			t.Fatalf("expected missing zellij error, got %v", err)
		}
	})

	t.Run("Zellij launch in current terminal", func(t *testing.T) {
		runner := exec.NewMockRunner()
		runner.Available["zellij"] = true

		tabs := []config.TabConfig{
			{Name: "Code", Panes: []config.PaneConfig{{Name: "shell", Command: ""}}},
		}
		err := Launch(runner, "zellij", "goJPeek", "/tmp", tabs, nil, "default")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(runner.RunCalls) != 1 || runner.RunCalls[0].Name != "/usr/bin/zellij" {
			t.Fatalf("expected zellij run call, got %#v", runner.RunCalls)
		}
	})
}
