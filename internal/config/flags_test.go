package config

import (
	"testing"
)

func TestParsePaneFlag(t *testing.T) {
	t.Run("Name and command", func(t *testing.T) {
		p := ParsePaneFlag("logs:tail -f app.log")
		if p.Name != "logs" || p.Command != "tail -f app.log" {
			t.Fatalf("unexpected pane: %#v", p)
		}
	})

	t.Run("Name only", func(t *testing.T) {
		p := ParsePaneFlag("shell")
		if p.Name != "shell" || p.Command != "" {
			t.Fatalf("unexpected pane: %#v", p)
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		p := ParsePaneFlag("")
		if p.Name != "" || p.Command != "" {
			t.Fatalf("unexpected pane: %#v", p)
		}
	})
}
