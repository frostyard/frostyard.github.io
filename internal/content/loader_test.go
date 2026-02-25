package content

import (
	"os"
	"path/filepath"
	"testing"
)

// helper to write a file at the given path under dir, creating parent dirs as needed.
func writeFile(t *testing.T, dir, relPath, content string) {
	t.Helper()
	full := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", filepath.Dir(full), err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", full, err)
	}
}

func TestLoadContent(t *testing.T) {
	tmp := t.TempDir()

	// content/docs/_index.md
	writeFile(t, tmp, "content/docs/_index.md", `---
title: "Documentation"
---

Welcome to the docs.
`)

	// content/docs/getting-started/_index.md
	writeFile(t, tmp, "content/docs/getting-started/_index.md", `---
title: "Getting Started"
weight: 1
---

Getting started guide.
`)

	// content/docs/getting-started/install.md
	writeFile(t, tmp, "content/docs/getting-started/install.md", `---
title: "Installation"
weight: 1
---

How to install.
`)

	// content/docs/getting-started/quickstart.md
	writeFile(t, tmp, "content/docs/getting-started/quickstart.md", `---
title: "Quickstart"
weight: 2
---

Quick start guide.
`)

	// content/docs/tools/_index.md
	writeFile(t, tmp, "content/docs/tools/_index.md", `---
title: "Tools"
weight: 2
---

Tools section.
`)

	// content/docs/tools/nbc/_index.md
	writeFile(t, tmp, "content/docs/tools/nbc/_index.md", `---
title: "nbc"
weight: 1
---

NBC compiler.
`)

	// content/docs/tools/nbc/usage.md
	writeFile(t, tmp, "content/docs/tools/nbc/usage.md", `---
title: "Usage"
weight: 1
---

How to use nbc.
`)

	// content/blog/posts/hello.md
	writeFile(t, tmp, "content/blog/posts/hello.md", `---
title: "Hello World"
date: "2026-01-15"
author: "bjk"
tags:
  - general
---

Hello world blog post.
`)

	contentDir := filepath.Join(tmp, "content")
	site, err := LoadContent(contentDir)
	if err != nil {
		t.Fatalf("LoadContent returned error: %v", err)
	}

	// Verify site.Pages is non-empty
	if len(site.Pages) == 0 {
		t.Error("site.Pages is empty, want non-empty")
	}

	// Verify site.Posts has 1 post titled "Hello World"
	if len(site.Posts) != 1 {
		t.Fatalf("len(site.Posts) = %d, want 1", len(site.Posts))
	}
	if site.Posts[0].Title != "Hello World" {
		t.Errorf("site.Posts[0].Title = %q, want %q", site.Posts[0].Title, "Hello World")
	}

	// Verify site.Sections is non-empty
	if len(site.Sections) == 0 {
		t.Error("site.Sections is empty, want non-empty")
	}

	// Verify blog post metadata
	post := site.Posts[0]
	if post.Author != "bjk" {
		t.Errorf("post.Author = %q, want %q", post.Author, "bjk")
	}
	if len(post.Tags) != 1 || post.Tags[0] != "general" {
		t.Errorf("post.Tags = %v, want [general]", post.Tags)
	}

	// Verify the blog post is not counted as an index page
	if post.IsIndex {
		t.Error("blog post IsIndex = true, want false")
	}

	// Verify section tree structure: find the "Documentation" root section
	var docsSection *Section
	for _, s := range site.Sections {
		if s.Title == "Documentation" {
			docsSection = s
			break
		}
	}
	if docsSection == nil {
		t.Fatal("could not find 'Documentation' section in site.Sections")
	}

	// Documentation should have two subsections: Getting Started (weight 1) and Tools (weight 2)
	if len(docsSection.Subsections) != 2 {
		t.Fatalf("docs subsections count = %d, want 2", len(docsSection.Subsections))
	}

	// Verify sorting by weight: Getting Started first, Tools second
	if docsSection.Subsections[0].Title != "Getting Started" {
		t.Errorf("first subsection = %q, want %q", docsSection.Subsections[0].Title, "Getting Started")
	}
	if docsSection.Subsections[1].Title != "Tools" {
		t.Errorf("second subsection = %q, want %q", docsSection.Subsections[1].Title, "Tools")
	}

	// Getting Started should have 2 pages: Installation (weight 1) then Quickstart (weight 2)
	gs := docsSection.Subsections[0]
	if len(gs.Pages) != 2 {
		t.Fatalf("Getting Started pages count = %d, want 2", len(gs.Pages))
	}
	if gs.Pages[0].Title != "Installation" {
		t.Errorf("first page = %q, want %q", gs.Pages[0].Title, "Installation")
	}
	if gs.Pages[1].Title != "Quickstart" {
		t.Errorf("second page = %q, want %q", gs.Pages[1].Title, "Quickstart")
	}

	// Tools should have one subsection: nbc
	tools := docsSection.Subsections[1]
	if len(tools.Subsections) != 1 {
		t.Fatalf("Tools subsections count = %d, want 1", len(tools.Subsections))
	}
	if tools.Subsections[0].Title != "nbc" {
		t.Errorf("Tools subsection = %q, want %q", tools.Subsections[0].Title, "nbc")
	}

	// nbc should have 1 page: Usage
	nbc := tools.Subsections[0]
	if len(nbc.Pages) != 1 {
		t.Fatalf("nbc pages count = %d, want 1", len(nbc.Pages))
	}
	if nbc.Pages[0].Title != "Usage" {
		t.Errorf("nbc page = %q, want %q", nbc.Pages[0].Title, "Usage")
	}
}

func TestLoadContentSkipsDrafts(t *testing.T) {
	tmp := t.TempDir()

	// visible.md — not a draft
	writeFile(t, tmp, "content/visible.md", `---
title: "Visible Page"
draft: false
---

This page should appear.
`)

	// hidden.md — is a draft
	writeFile(t, tmp, "content/hidden.md", `---
title: "Hidden Page"
draft: true
---

This page should NOT appear.
`)

	contentDir := filepath.Join(tmp, "content")
	site, err := LoadContent(contentDir)
	if err != nil {
		t.Fatalf("LoadContent returned error: %v", err)
	}

	// Should have exactly 1 page (the visible one)
	if len(site.Pages) != 1 {
		t.Fatalf("len(site.Pages) = %d, want 1", len(site.Pages))
	}
	if site.Pages[0].Title != "Visible Page" {
		t.Errorf("site.Pages[0].Title = %q, want %q", site.Pages[0].Title, "Visible Page")
	}

	// Make sure the hidden page is not present
	for _, p := range site.Pages {
		if p.Title == "Hidden Page" {
			t.Error("found 'Hidden Page' in site.Pages, but it should be skipped as a draft")
		}
	}
}
