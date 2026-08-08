package editor

import (
	"strings"
	"testing"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
)

func TestLaunchEditor(t *testing.T) {
	t.Run("VS Code launch success", func(t *testing.T) {
		runner := exec.NewMockRunner()
		runner.Available["code"] = true

		err := Launch(runner, "vscode", "/path/to/project")
		if err != nil {
			t.Fatalf("unexpected error launching code: %v", err)
		}

		if len(runner.StartCalls) != 1 || runner.StartCalls[0].Name != "/usr/bin/code" {
			t.Fatalf("expected 1 start call for /usr/bin/code, got %#v", runner.StartCalls)
		}
	})

	t.Run("VS Code missing error", func(t *testing.T) {
		runner := exec.NewMockRunner()
		err := Launch(runner, "vscode", "/path/to/project")
		if err == nil || !strings.Contains(err.Error(), "not found on PATH") {
			t.Fatalf("expected not found error, got %v", err)
		}
	})

	t.Run("None editor", func(t *testing.T) {
		runner := exec.NewMockRunner()
		err := Launch(runner, "none", "/path/to/project")
		if err != nil {
			t.Fatalf("unexpected error for none editor: %v", err)
		}
		if len(runner.StartCalls) != 0 {
			t.Fatalf("expected no start calls, got %#v", runner.StartCalls)
		}
	})
}
