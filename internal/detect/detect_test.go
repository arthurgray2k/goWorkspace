package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectProjectType(t *testing.T) {
	t.Run("Go Project", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module foo"), 0644)
		if got := DetectProjectType(dir); got != TemplateGo {
			t.Fatalf("expected %q, got %q", TemplateGo, got)
		}
	})

	t.Run("Node Project", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0644)
		if got := DetectProjectType(dir); got != TemplateNode {
			t.Fatalf("expected %q, got %q", TemplateNode, got)
		}
	})

	t.Run("Rust Project", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte("[package]"), 0644)
		if got := DetectProjectType(dir); got != TemplateRust {
			t.Fatalf("expected %q, got %q", TemplateRust, got)
		}
	})

	t.Run(".NET Project", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "MyApp.csproj"), []byte("<Project></Project>"), 0644)
		if got := DetectProjectType(dir); got != TemplateDotnet {
			t.Fatalf("expected %q, got %q", TemplateDotnet, got)
		}
	})

	t.Run("Generic Fallback", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Hello"), 0644)
		if got := DetectProjectType(dir); got != TemplateGeneric {
			t.Fatalf("expected %q, got %q", TemplateGeneric, got)
		}
	})
}
