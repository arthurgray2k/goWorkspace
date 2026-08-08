package detect

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	TemplateGo      = "go"
	TemplateDotnet  = "dotnet"
	TemplateNode    = "node"
	TemplateRust    = "rust"
	TemplateGeneric = "generic"
)

// DetectProjectType examines files in directory to deterministically identify the project type.
func DetectProjectType(dir string) string {
	if dir == "" {
		return TemplateGeneric
	}

	// 1. Go project detection: go.mod
	if fileExists(filepath.Join(dir, "go.mod")) {
		return TemplateGo
	}

	// 2. Node project detection: package.json
	if fileExists(filepath.Join(dir, "package.json")) {
		return TemplateNode
	}

	// 3. Rust project detection: Cargo.toml
	if fileExists(filepath.Join(dir, "Cargo.toml")) {
		return TemplateRust
	}

	// 4. .NET project detection: *.sln or *.csproj
	if hasPattern(dir, "*.sln") || hasPattern(dir, "*.csproj") {
		return TemplateDotnet
	}

	return TemplateGeneric
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func hasPattern(dir string, pattern string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matched, err := filepath.Match(pattern, entry.Name())
		if err == nil && matched {
			return true
		}
		if strings.HasSuffix(pattern, "csproj") && strings.HasSuffix(entry.Name(), ".csproj") {
			return true
		}
		if strings.HasSuffix(pattern, "sln") && strings.HasSuffix(entry.Name(), ".sln") {
			return true
		}
	}
	return false
}
