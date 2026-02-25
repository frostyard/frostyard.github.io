package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	// Create a temp directory for the test project
	tmpDir := t.TempDir()

	// Set up minimal content structure
	contentDir := filepath.Join(tmpDir, "content")

	// content/docs/_index.md
	docsDir := filepath.Join(contentDir, "docs")
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "_index.md"), []byte(`---
title: "Docs"
---
`), 0o644); err != nil {
		t.Fatal(err)
	}

	// content/docs/intro.md
	introContent := `---
title: "Introduction"
weight: 1
---
# Introduction

Hello world.
`
	if err := os.WriteFile(filepath.Join(docsDir, "intro.md"), []byte(introContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Set up output directory
	outputDir := filepath.Join(tmpDir, "dist")

	// Run build
	cfg := Config{
		ContentDir: contentDir,
		StaticDir:  filepath.Join(tmpDir, "static"), // does not exist, should be skipped
		OutputDir:  outputDir,
		Root:       tmpDir,
	}

	if err := Build(cfg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify dist/docs/intro/index.html exists
	introHTML := filepath.Join(outputDir, "docs", "intro", "index.html")
	data, err := os.ReadFile(introHTML)
	if err != nil {
		t.Fatalf("Expected %s to exist: %v", introHTML, err)
	}

	html := string(data)

	// Verify the HTML contains expected content
	if !strings.Contains(html, "Introduction") {
		t.Errorf("Expected HTML to contain 'Introduction', got:\n%s", html)
	}
	if !strings.Contains(html, "Hello world.") {
		t.Errorf("Expected HTML to contain 'Hello world.', got:\n%s", html)
	}
}
