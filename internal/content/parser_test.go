package content

import (
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	input := []byte(`---
title: "Install NBC"
description: "How to install the NBC compiler"
section: tools
weight: 10
draft: true
---

This is the body of the page.

Some **bold** text and a [link](https://example.com).
`)

	page, err := ParsePage(input, "content/docs/tools/nbc/install.md")
	if err != nil {
		t.Fatalf("ParsePage returned error: %v", err)
	}

	if page.Title != "Install NBC" {
		t.Errorf("Title = %q, want %q", page.Title, "Install NBC")
	}
	if page.Description != "How to install the NBC compiler" {
		t.Errorf("Description = %q, want %q", page.Description, "How to install the NBC compiler")
	}
	if page.Section != "tools" {
		t.Errorf("Section = %q, want %q", page.Section, "tools")
	}
	if page.Weight != 10 {
		t.Errorf("Weight = %d, want %d", page.Weight, 10)
	}
	if !page.Draft {
		t.Error("Draft = false, want true")
	}
	if page.Content == "" {
		t.Error("Content is empty, want non-empty HTML")
	}
}

func TestParseIndexPage(t *testing.T) {
	input := []byte(`---
title: "Images"
description: "Image gallery section"
icon: "image"
---

Welcome to the images section.
`)

	page, err := ParsePage(input, "content/docs/images/_index.md")
	if err != nil {
		t.Fatalf("ParsePage returned error: %v", err)
	}

	if !page.IsIndex {
		t.Error("IsIndex = false, want true")
	}
	if page.Title != "Images" {
		t.Errorf("Title = %q, want %q", page.Title, "Images")
	}
	if page.Icon != "image" {
		t.Errorf("Icon = %q, want %q", page.Icon, "image")
	}
}

func TestParseBlogPost(t *testing.T) {
	input := []byte(`---
title: "My First Post"
date: "2025-03-15"
author: "bjk"
tags:
  - go
  - web
---

This is a blog post about Go and web development.
`)

	page, err := ParsePage(input, "content/blog/posts/my-post.md")
	if err != nil {
		t.Fatalf("ParsePage returned error: %v", err)
	}

	if page.Date != "2025-03-15" {
		t.Errorf("Date = %q, want %q", page.Date, "2025-03-15")
	}
	if page.Author != "bjk" {
		t.Errorf("Author = %q, want %q", page.Author, "bjk")
	}
	if len(page.Tags) != 2 {
		t.Errorf("len(Tags) = %d, want 2", len(page.Tags))
	}
	if page.ParsedDate.IsZero() {
		t.Error("ParsedDate is zero, want non-zero")
	}
	if page.ParsedDate.Year() != 2025 || page.ParsedDate.Month() != 3 || page.ParsedDate.Day() != 15 {
		t.Errorf("ParsedDate = %v, want 2025-03-15", page.ParsedDate)
	}
}

func TestParseHeadings(t *testing.T) {
	input := []byte(`---
title: "Headings Test"
---

# First Heading

Some text.

## Second Heading

More text.

### Third Heading

Even more text.
`)

	page, err := ParsePage(input, "content/docs/headings.md")
	if err != nil {
		t.Fatalf("ParsePage returned error: %v", err)
	}

	if len(page.Headings) < 3 {
		t.Fatalf("len(Headings) = %d, want at least 3", len(page.Headings))
	}

	tests := []struct {
		level int
		text  string
	}{
		{1, "First Heading"},
		{2, "Second Heading"},
		{3, "Third Heading"},
	}

	for i, tt := range tests {
		if page.Headings[i].Level != tt.level {
			t.Errorf("Headings[%d].Level = %d, want %d", i, page.Headings[i].Level, tt.level)
		}
		if page.Headings[i].Text != tt.text {
			t.Errorf("Headings[%d].Text = %q, want %q", i, page.Headings[i].Text, tt.text)
		}
		if page.Headings[i].ID == "" {
			t.Errorf("Headings[%d].ID is empty, want non-empty", i)
		}
	}
}

func TestParsePagePath(t *testing.T) {
	tests := []struct {
		sourcePath string
		wantPath   string
	}{
		{"content/docs/tools/nbc/install.md", "/docs/tools/nbc/install/"},
		{"content/docs/images/_index.md", "/docs/images/"},
		{"content/blog/posts/my-post.md", "/blog/posts/my-post/"},
		{"content/docs/faq.md", "/docs/faq/"},
	}

	for _, tt := range tests {
		t.Run(tt.sourcePath, func(t *testing.T) {
			got := computePath(tt.sourcePath)
			if got != tt.wantPath {
				t.Errorf("computePath(%q) = %q, want %q", tt.sourcePath, got, tt.wantPath)
			}
		})
	}
}
