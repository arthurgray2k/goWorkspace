package terminal

import (
	"strings"
	"testing"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
)

func TestLaunchTerminal(t *testing.T) {
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

	t.Run("Ghostty missing error", func(t *testing.T) {
		runner := exec.NewMockRunner()
		err := Launch(runner, "ghostty", "", nil, "/path/to/project")
		if err == nil || !strings.Contains(err.Error(), "not found on PATH") {
			t.Fatalf("expected missing error, got %v", err)
		}
	})

	t.Run("Default terminal context", func(t *testing.T) {
		runner := exec.NewMockRunner()
		err := Launch(runner, "default", "", nil, "/path/to/project")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(runner.StartCalls) != 0 {
			t.Fatalf("expected no start calls, got %#v", runner.StartCalls)
		}
	})
}
