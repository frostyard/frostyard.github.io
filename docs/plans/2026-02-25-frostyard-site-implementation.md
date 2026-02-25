# Frostyard Site Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a custom Go static site generator that renders markdown content with YAML frontmatter through Templ templates, styled with Tailwind CSS, producing a frost-themed documentation site for the Frostyard ecosystem.

**Architecture:** A Go CLI (`cmd/frostyard/main.go`) orchestrates the build pipeline: glob markdown files from `content/`, parse frontmatter+markdown with goldmark/yaml.v3, build a section tree for navigation, render each page through Templ component templates, run Tailwind CLI for CSS, copy static assets, generate sitemap/RSS, and run Pagefind for search indexing. A dev server with file watching and live reload supports local development.

**Tech Stack:** Go 1.26, Templ v0.3.x, Tailwind CSS v4 (standalone CLI), goldmark + GFM extensions, Chroma (syntax highlighting), yaml.v3, fsnotify, Pagefind

---

## Task 1: Clean Up Old Site and Initialize Go Module

Remove the MkDocs/submodule infrastructure and set up the Go project skeleton.

**Files:**
- Delete: `.gitmodules`, `mkdocs.yml`, `requirements.txt`, `.github/workflows/update-submodules.yml`
- Delete directories: `site/`, `.venv/`, `snow/`, `nbc/`, `chairlift/`, `debian-bootc-core/`, `debian-bootc-gnome/`, `snowfield/`, `cayo/`, `snowdrift/`, `first-setup/`
- Delete: `docs/index.md`, `docs/atomic.md`, `docs/blog/` (keep `docs/plans/`)
- Create: `go.mod`, `go.sum`
- Create: `.gitignore` (updated)

**Step 1: Save nbc docs content for later migration**

Before deleting submodules, copy the nbc docs we'll port later:

```bash
mkdir -p /tmp/frostyard-nbc-docs
cp -r nbc/docs/* /tmp/frostyard-nbc-docs/
```

**Step 2: Deinitialize and remove all git submodules**

```bash
git submodule deinit --all -f
git rm --cached debian-bootc-core debian-bootc-gnome snow snowfield cayo snowdrift chairlift first-setup nbc
rm -rf debian-bootc-core debian-bootc-gnome snow snowfield cayo snowdrift chairlift first-setup nbc
rm -f .gitmodules
```

**Step 3: Remove old MkDocs files and built output**

```bash
rm -f mkdocs.yml requirements.txt
rm -rf site/ .venv/
rm -rf .github/workflows/update-submodules.yml
rm -f docs/index.md docs/atomic.md
rm -rf docs/blog/
```

**Step 4: Initialize Go module**

```bash
go mod init github.com/frostyard/site
```

**Step 5: Create updated .gitignore**

Write `.gitignore`:

```
dist/
.venv/
*.exe
.DS_Store
node_modules/
/tailwindcss
/pagefind
```

**Step 6: Create project directory skeleton**

```bash
mkdir -p cmd/frostyard
mkdir -p internal/build internal/content internal/render internal/server
mkdir -p templates/layouts templates/components templates/pages
mkdir -p content/docs/getting-started content/docs/images content/docs/tools
mkdir -p content/blog/posts
mkdir -p static/images static/fonts
```

**Step 7: Commit**

```bash
git add -A
git commit -m "chore: remove MkDocs and submodules, init Go project skeleton"
```

---

## Task 2: Content Parser — Frontmatter + Markdown

Build the core content parsing library: read a markdown file, extract YAML frontmatter, parse markdown to HTML.

**Files:**
- Create: `internal/content/content.go`
- Create: `internal/content/parser.go`
- Create: `internal/content/parser_test.go`

**Step 1: Write the content types**

Create `internal/content/content.go`:

```go
package content

import (
	"html/template"
	"time"
)

// Page represents a parsed markdown page.
type Page struct {
	// Frontmatter fields
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Section     string   `yaml:"section"`
	Weight      int      `yaml:"weight"`
	Draft       bool     `yaml:"draft"`
	Icon        string   `yaml:"icon"`
	Date        string   `yaml:"date"`
	Author      string   `yaml:"author"`
	Tags        []string `yaml:"tags"`

	// Computed fields
	Content    template.HTML // Rendered HTML from markdown
	Path       string        // URL path (e.g., "/docs/tools/nbc/")
	SourcePath string        // Filesystem path to the .md file
	Slug       string        // URL-friendly name derived from filename
	IsIndex    bool          // True if this is an _index.md file
	ParsedDate time.Time     // Parsed from Date string
	Headings   []Heading     // Extracted headings for TOC
}

// Heading represents a heading extracted from markdown for TOC generation.
type Heading struct {
	Level int
	ID    string
	Text  string
}

// Section represents a group of pages in the navigation tree.
type Section struct {
	Title       string
	Description string
	Icon        string
	Path        string
	Pages       []*Page
	Subsections []*Section
	IndexPage   *Page // The _index.md page for this section
	Weight      int
}

// Site holds all parsed content for the site.
type Site struct {
	Pages    []*Page
	Sections []*Section
	Posts    []*Page // Blog posts, sorted by date descending
}
```

**Step 2: Write the failing test for the parser**

Create `internal/content/parser_test.go`:

```go
package content

import (
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	input := `---
title: "Test Page"
description: "A test page"
section: tools/nbc
weight: 10
draft: false
---

# Hello World

This is a test.
`

	page, err := ParsePage([]byte(input), "content/docs/tools/nbc/test.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if page.Title != "Test Page" {
		t.Errorf("expected title 'Test Page', got %q", page.Title)
	}
	if page.Description != "A test page" {
		t.Errorf("expected description 'A test page', got %q", page.Description)
	}
	if page.Weight != 10 {
		t.Errorf("expected weight 10, got %d", page.Weight)
	}
	if page.Draft {
		t.Error("expected draft to be false")
	}
	if page.Content == "" {
		t.Error("expected non-empty content")
	}
	if page.IsIndex {
		t.Error("expected IsIndex to be false for non-index file")
	}
}

func TestParseIndexPage(t *testing.T) {
	input := `---
title: "Images"
description: "Frostyard bootc container images"
icon: "server"
---

Overview content here.
`

	page, err := ParsePage([]byte(input), "content/docs/images/_index.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !page.IsIndex {
		t.Error("expected IsIndex to be true for _index.md")
	}
	if page.Title != "Images" {
		t.Errorf("expected title 'Images', got %q", page.Title)
	}
	if page.Icon != "server" {
		t.Errorf("expected icon 'server', got %q", page.Icon)
	}
}

func TestParseBlogPost(t *testing.T) {
	input := `---
title: "Introducing Snow 2.0"
date: "2026-02-20"
author: "bjk"
tags: ["snow", "release"]
description: "What's new in Snow 2.0"
---

Blog content here.
`

	page, err := ParsePage([]byte(input), "content/blog/posts/introducing-snow-2.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if page.Title != "Introducing Snow 2.0" {
		t.Errorf("expected title 'Introducing Snow 2.0', got %q", page.Title)
	}
	if page.Author != "bjk" {
		t.Errorf("expected author 'bjk', got %q", page.Author)
	}
	if len(page.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(page.Tags))
	}
	if page.ParsedDate.IsZero() {
		t.Error("expected parsed date to be non-zero")
	}
}

func TestParseHeadings(t *testing.T) {
	input := `---
title: "Test"
---

# Main Title

## Section One

### Subsection

## Section Two
`

	page, err := ParsePage([]byte(input), "content/docs/test.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(page.Headings) < 3 {
		t.Fatalf("expected at least 3 headings, got %d", len(page.Headings))
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
		input := "---\ntitle: \"Test\"\n---\nContent."
		page, err := ParsePage([]byte(input), tt.sourcePath)
		if err != nil {
			t.Fatalf("error parsing %s: %v", tt.sourcePath, err)
		}
		if page.Path != tt.wantPath {
			t.Errorf("sourcePath=%s: expected path %q, got %q", tt.sourcePath, tt.wantPath, page.Path)
		}
	}
}
```

**Step 3: Run tests to verify they fail**

```bash
cd /home/bjk/projects/frostyard/site && go test ./internal/content/ -v
```

Expected: FAIL — `ParsePage` not defined.

**Step 4: Implement the parser**

Create `internal/content/parser.go`:

```go
package content

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

// ParsePage parses a markdown file with YAML frontmatter.
// sourcePath is relative to the project root (e.g., "content/docs/tools/nbc/install.md").
func ParsePage(data []byte, sourcePath string) (*Page, error) {
	frontmatter, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("parsing frontmatter in %s: %w", sourcePath, err)
	}

	page := &Page{
		SourcePath: sourcePath,
	}

	if err := yaml.Unmarshal(frontmatter, page); err != nil {
		return nil, fmt.Errorf("unmarshaling frontmatter in %s: %w", sourcePath, err)
	}

	// Parse date if present
	if page.Date != "" {
		t, err := time.Parse("2006-01-02", page.Date)
		if err == nil {
			page.ParsedDate = t
		}
	}

	// Compute path and slug
	page.IsIndex = filepath.Base(sourcePath) == "_index.md"
	page.Path = computePath(sourcePath)
	page.Slug = computeSlug(sourcePath)

	// Parse markdown to HTML
	html, headings, err := renderMarkdown(body)
	if err != nil {
		return nil, fmt.Errorf("rendering markdown in %s: %w", sourcePath, err)
	}
	page.Content = template.HTML(html)
	page.Headings = headings

	return page, nil
}

// splitFrontmatter separates YAML frontmatter from markdown body.
func splitFrontmatter(data []byte) (frontmatter, body []byte, err error) {
	const delimiter = "---"
	s := string(data)

	// Must start with ---
	if !strings.HasPrefix(strings.TrimSpace(s), delimiter) {
		return nil, data, nil
	}

	s = strings.TrimSpace(s)
	// Find the closing ---
	rest := s[len(delimiter)+1:] // skip first --- and newline
	idx := strings.Index(rest, "\n"+delimiter)
	if idx == -1 {
		return nil, data, fmt.Errorf("unclosed frontmatter")
	}

	fm := rest[:idx]
	bodyStart := idx + len("\n"+delimiter)
	if bodyStart < len(rest) {
		// Skip the newline after closing ---
		bodyStr := rest[bodyStart:]
		if len(bodyStr) > 0 && bodyStr[0] == '\n' {
			bodyStr = bodyStr[1:]
		}
		body = []byte(bodyStr)
	}

	return []byte(fm), body, nil
}

// computePath derives the URL path from the source file path.
func computePath(sourcePath string) string {
	// Strip "content/" prefix
	p := strings.TrimPrefix(sourcePath, "content/")
	// Remove .md extension
	p = strings.TrimSuffix(p, ".md")
	// For _index.md, use the directory path
	if strings.HasSuffix(p, "/_index") {
		p = strings.TrimSuffix(p, "/_index")
	}
	// Ensure leading and trailing slashes
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}
	return p
}

// computeSlug derives a URL-friendly name from the filename.
func computeSlug(sourcePath string) string {
	base := filepath.Base(sourcePath)
	base = strings.TrimSuffix(base, ".md")
	if base == "_index" {
		// Use parent directory name
		return filepath.Base(filepath.Dir(sourcePath))
	}
	return base
}

// renderMarkdown converts markdown bytes to HTML and extracts headings.
func renderMarkdown(source []byte) (string, []Heading, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	// Extract headings from AST
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)
	var headings []Heading
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if h, ok := n.(*ast.Heading); ok {
			id, _ := n.AttributeString("id")
			idStr := ""
			if id != nil {
				idStr = string(id.([]byte))
			}
			var textBuf bytes.Buffer
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if t, ok := c.(*ast.Text); ok {
					textBuf.Write(t.Segment.Value(source))
				}
			}
			headings = append(headings, Heading{
				Level: h.Level,
				ID:    idStr,
				Text:  textBuf.String(),
			})
		}
		return ast.WalkContinue, nil
	})

	// Render to HTML
	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, source, doc); err != nil {
		return "", nil, err
	}

	return buf.String(), headings, nil
}
```

**Step 5: Install dependencies and run tests**

```bash
cd /home/bjk/projects/frostyard/site && go mod tidy && go test ./internal/content/ -v
```

Expected: ALL PASS

**Step 6: Commit**

```bash
git add internal/content/ go.mod go.sum
git commit -m "feat: add content parser with frontmatter and markdown rendering"
```

---

## Task 3: Content Loader — Glob Files and Build Section Tree

Read all markdown files from `content/` and organize them into a navigable section tree.

**Files:**
- Create: `internal/content/loader.go`
- Create: `internal/content/loader_test.go`

**Step 1: Write the failing test**

Create `internal/content/loader_test.go`:

```go
package content

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadContent(t *testing.T) {
	// Create a temp content directory
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")

	// Create test structure
	files := map[string]string{
		"content/docs/_index.md": `---
title: "Documentation"
description: "Frostyard docs"
---
`,
		"content/docs/getting-started/_index.md": `---
title: "Getting Started"
weight: 1
---
`,
		"content/docs/getting-started/install.md": `---
title: "Installation"
weight: 1
---
Install steps here.
`,
		"content/docs/getting-started/quickstart.md": `---
title: "Quickstart"
weight: 2
---
Quick start guide.
`,
		"content/docs/tools/_index.md": `---
title: "Tools"
weight: 2
---
`,
		"content/docs/tools/nbc/_index.md": `---
title: "nbc"
weight: 1
---
`,
		"content/docs/tools/nbc/usage.md": `---
title: "Usage"
weight: 1
---
Usage info.
`,
		"content/blog/posts/hello.md": `---
title: "Hello World"
date: "2026-01-15"
author: "bjk"
tags: ["general"]
---
First post.
`,
	}

	for path, body := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	site, err := LoadContent(contentDir)
	if err != nil {
		t.Fatalf("LoadContent failed: %v", err)
	}

	if len(site.Pages) == 0 {
		t.Fatal("expected pages, got none")
	}

	if len(site.Posts) != 1 {
		t.Fatalf("expected 1 blog post, got %d", len(site.Posts))
	}

	if site.Posts[0].Title != "Hello World" {
		t.Errorf("expected post title 'Hello World', got %q", site.Posts[0].Title)
	}

	// Check section tree
	if len(site.Sections) == 0 {
		t.Fatal("expected sections, got none")
	}
}

func TestLoadContentSkipsDrafts(t *testing.T) {
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")

	files := map[string]string{
		"content/docs/visible.md": `---
title: "Visible"
---
Content.
`,
		"content/docs/hidden.md": `---
title: "Hidden"
draft: true
---
Content.
`,
	}

	for path, body := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	site, err := LoadContent(contentDir)
	if err != nil {
		t.Fatalf("LoadContent failed: %v", err)
	}

	for _, p := range site.Pages {
		if p.Title == "Hidden" {
			t.Error("draft page should not be in site.Pages")
		}
	}
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./internal/content/ -v -run TestLoad
```

Expected: FAIL — `LoadContent` not defined.

**Step 3: Implement the loader**

Create `internal/content/loader.go`:

```go
package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LoadContent reads all markdown files from contentDir and returns a Site.
func LoadContent(contentDir string) (*Site, error) {
	var allPages []*Page

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		// Make sourcePath relative to parent of contentDir
		relPath, err := filepath.Rel(filepath.Dir(contentDir), path)
		if err != nil {
			return err
		}

		page, err := ParsePage(data, relPath)
		if err != nil {
			return err
		}

		allPages = append(allPages, page)
		return nil
	})
	if err != nil {
		return nil, err
	}

	site := &Site{}

	// Separate drafts, blog posts, and regular pages
	for _, p := range allPages {
		if p.Draft {
			continue
		}
		if strings.HasPrefix(p.Path, "/blog/posts/") && !p.IsIndex {
			site.Posts = append(site.Posts, p)
		}
		site.Pages = append(site.Pages, p)
	}

	// Sort posts by date descending
	sort.Slice(site.Posts, func(i, j int) bool {
		return site.Posts[i].ParsedDate.After(site.Posts[j].ParsedDate)
	})

	// Build section tree
	site.Sections = buildSectionTree(site.Pages)

	return site, nil
}

// buildSectionTree organizes pages into a hierarchical section structure.
func buildSectionTree(pages []*Page) []*Section {
	sectionMap := make(map[string]*Section)

	// First pass: create sections from _index.md pages
	for _, p := range pages {
		if !p.IsIndex {
			continue
		}
		sectionMap[p.Path] = &Section{
			Title:       p.Title,
			Description: p.Description,
			Icon:        p.Icon,
			Path:        p.Path,
			IndexPage:   p,
			Weight:      p.Weight,
		}
	}

	// Second pass: assign non-index pages to their parent section
	for _, p := range pages {
		if p.IsIndex {
			continue
		}
		parentPath := filepath.Dir(strings.TrimSuffix(p.Path, "/")) + "/"
		if sec, ok := sectionMap[parentPath]; ok {
			sec.Pages = append(sec.Pages, p)
		}
	}

	// Sort pages within each section by weight
	for _, sec := range sectionMap {
		sort.Slice(sec.Pages, func(i, j int) bool {
			return sec.Pages[i].Weight < sec.Pages[j].Weight
		})
	}

	// Build hierarchy: assign subsections to parent sections
	for path, sec := range sectionMap {
		parentPath := filepath.Dir(strings.TrimSuffix(path, "/")) + "/"
		if parent, ok := sectionMap[parentPath]; ok && parentPath != path {
			parent.Subsections = append(parent.Subsections, sec)
		}
	}

	// Sort subsections by weight
	for _, sec := range sectionMap {
		sort.Slice(sec.Subsections, func(i, j int) bool {
			return sec.Subsections[i].Weight < sec.Subsections[j].Weight
		})
	}

	// Return top-level sections (those whose parent isn't in sectionMap)
	var roots []*Section
	for path, sec := range sectionMap {
		parentPath := filepath.Dir(strings.TrimSuffix(path, "/")) + "/"
		if _, ok := sectionMap[parentPath]; !ok {
			roots = append(roots, sec)
		}
	}

	sort.Slice(roots, func(i, j int) bool {
		return roots[i].Weight < roots[j].Weight
	})

	return roots
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/content/ -v
```

Expected: ALL PASS

**Step 5: Commit**

```bash
git add internal/content/loader.go internal/content/loader_test.go
git commit -m "feat: add content loader with section tree builder"
```

---

## Task 4: Templ Templates — Base Layout and Components

Create the Templ template files for the site shell, navigation, sidebar, footer, and page layouts.

**Files:**
- Create: `templates/layouts/base.templ`
- Create: `templates/layouts/docs.templ`
- Create: `templates/layouts/blog.templ`
- Create: `templates/layouts/landing.templ`
- Create: `templates/components/nav.templ`
- Create: `templates/components/sidebar.templ`
- Create: `templates/components/footer.templ`
- Create: `templates/components/toc.templ`

**Step 1: Create the base layout**

Create `templates/layouts/base.templ`:

```templ
package layouts

import "github.com/frostyard/site/templates/components"

type PageMeta struct {
	Title       string
	Description string
	Path        string
	SiteName    string
}

templ Base(meta PageMeta) {
	<!DOCTYPE html>
	<html lang="en" class="dark">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>
				if meta.Title != "" {
					{ meta.Title } — { meta.SiteName }
				} else {
					{ meta.SiteName }
				}
			</title>
			if meta.Description != "" {
				<meta name="description" content={ meta.Description }/>
			}
			<link rel="stylesheet" href="/css/style.css"/>
			<link href="/pagefind/pagefind-ui.css" rel="stylesheet"/>
			<script src="/pagefind/pagefind-ui.js"></script>
		</head>
		<body class="bg-slate-900 text-slate-100 min-h-screen flex flex-col">
			<!-- Frost gradient line -->
			<div class="h-0.5 bg-gradient-to-r from-sky-400 via-blue-400 to-sky-500"></div>
			@components.Nav(meta.SiteName, meta.Path)
			<main class="flex-1">
				{ children... }
			</main>
			@components.Footer()
			<script>
				// Dark mode toggle
				function toggleDarkMode() {
					document.documentElement.classList.toggle('dark');
					localStorage.setItem('theme',
						document.documentElement.classList.contains('dark') ? 'dark' : 'light'
					);
				}
				// Apply saved theme
				if (localStorage.getItem('theme') === 'light') {
					document.documentElement.classList.remove('dark');
				}
			</script>
		</body>
	</html>
}
```

**Step 2: Create the nav component**

Create `templates/components/nav.templ`:

```templ
package components

type NavLink struct {
	Label  string
	Path   string
}

var mainNav = []NavLink{
	{Label: "Docs", Path: "/docs/"},
	{Label: "Blog", Path: "/blog/"},
	{Label: "Downloads", Path: "/downloads/"},
	{Label: "Community", Path: "/community/"},
}

templ Nav(siteName string, currentPath string) {
	<nav class="sticky top-0 z-50 bg-slate-900/95 backdrop-blur border-b border-slate-800">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex items-center justify-between h-14">
				<div class="flex items-center gap-8">
					<a href="/" class="text-lg font-semibold text-sky-400 hover:text-sky-300">
						{ siteName }
					</a>
					<div class="hidden md:flex items-center gap-1">
						for _, link := range mainNav {
							<a
								href={ templ.SafeURL(link.Path) }
								class={ "px-3 py-1.5 rounded-md text-sm transition-colors",
									templ.KV("text-sky-400 bg-slate-800", isActive(currentPath, link.Path)),
									templ.KV("text-slate-300 hover:text-slate-100 hover:bg-slate-800/50", !isActive(currentPath, link.Path)) }
							>
								{ link.Label }
							</a>
						}
					</div>
				</div>
				<div class="flex items-center gap-3">
					<div id="search" class="hidden sm:block"></div>
					<button
						onclick="toggleDarkMode()"
						class="p-2 rounded-md text-slate-400 hover:text-slate-200 hover:bg-slate-800"
						aria-label="Toggle dark mode"
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
						</svg>
					</button>
					<!-- Mobile hamburger -->
					<button
						class="md:hidden p-2 rounded-md text-slate-400 hover:text-slate-200"
						onclick="document.getElementById('mobile-menu').classList.toggle('hidden')"
						aria-label="Menu"
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
						</svg>
					</button>
				</div>
			</div>
		</div>
		<!-- Mobile menu -->
		<div id="mobile-menu" class="hidden md:hidden border-t border-slate-800 px-4 py-2">
			for _, link := range mainNav {
				<a href={ templ.SafeURL(link.Path) } class="block px-3 py-2 text-sm text-slate-300 hover:text-slate-100 rounded-md">
					{ link.Label }
				</a>
			}
		</div>
	</nav>
	<script>
		window.addEventListener('DOMContentLoaded', function() {
			if (typeof PagefindUI !== 'undefined') {
				new PagefindUI({
					element: "#search",
					showSubResults: true,
					showImages: false,
				});
			}
		});
	</script>
}

func isActive(currentPath, linkPath string) bool {
	if linkPath == "/" {
		return currentPath == "/"
	}
	return len(currentPath) >= len(linkPath) && currentPath[:len(linkPath)] == linkPath
}
```

**Step 3: Create the sidebar component**

Create `templates/components/sidebar.templ`:

```templ
package components

type SidebarSection struct {
	Title       string
	Path        string
	Pages       []SidebarLink
	Subsections []SidebarSection
}

type SidebarLink struct {
	Title  string
	Path   string
}

templ Sidebar(sections []SidebarSection, currentPath string) {
	<aside class="hidden lg:block w-64 shrink-0">
		<nav class="sticky top-16 overflow-y-auto max-h-[calc(100vh-4rem)] py-6 pr-4">
			for _, section := range sections {
				@sidebarSection(section, currentPath)
			}
		</nav>
	</aside>
}

templ sidebarSection(section SidebarSection, currentPath string) {
	<div class="mb-4">
		<a
			href={ templ.SafeURL(section.Path) }
			class="block text-sm font-semibold text-slate-300 mb-1 hover:text-sky-400"
		>
			{ section.Title }
		</a>
		<ul class="ml-2 border-l border-slate-700 space-y-0.5">
			for _, page := range section.Pages {
				<li>
					<a
						href={ templ.SafeURL(page.Path) }
						class={ "block pl-3 py-1 text-sm border-l -ml-px transition-colors",
							templ.KV("border-sky-400 text-sky-400", currentPath == page.Path),
							templ.KV("border-transparent text-slate-400 hover:text-slate-200 hover:border-slate-500", currentPath != page.Path) }
					>
						{ page.Title }
					</a>
				</li>
			}
			for _, sub := range section.Subsections {
				<li class="pt-2">
					@sidebarSection(sub, currentPath)
				</li>
			}
		</ul>
	</div>
}
```

**Step 4: Create the TOC component**

Create `templates/components/toc.templ`:

```templ
package components

type TOCHeading struct {
	Level int
	ID    string
	Text  string
}

templ TOC(headings []TOCHeading) {
	if len(headings) > 1 {
		<aside class="hidden xl:block w-56 shrink-0">
			<nav class="sticky top-16 overflow-y-auto max-h-[calc(100vh-4rem)] py-6 pl-4">
				<h4 class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">On this page</h4>
				<ul class="space-y-1 text-sm">
					for _, h := range headings {
						if h.Level == 2 || h.Level == 3 {
							<li class={ templ.KV("ml-3", h.Level == 3) }>
								<a
									href={ templ.SafeURL("#" + h.ID) }
									class="block py-0.5 text-slate-400 hover:text-slate-200 transition-colors"
								>
									{ h.Text }
								</a>
							</li>
						}
					}
				</ul>
			</nav>
		</aside>
	}
}
```

**Step 5: Create the footer component**

Create `templates/components/footer.templ`:

```templ
package components

templ Footer() {
	<footer class="border-t border-slate-800 mt-16">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-slate-500">
				<p>Frostyard — Secure, reproducible Linux images</p>
				<div class="flex items-center gap-4">
					<a href="https://github.com/frostyard" class="hover:text-slate-300 transition-colors">GitHub</a>
					<a href="/community/" class="hover:text-slate-300 transition-colors">Community</a>
				</div>
			</div>
		</div>
	</footer>
}
```

**Step 6: Create the docs page layout**

Create `templates/layouts/docs.templ`:

```templ
package layouts

import (
	"github.com/frostyard/site/templates/components"
)

templ Docs(meta PageMeta, sidebar []components.SidebarSection, toc []components.TOCHeading) {
	@Base(meta) {
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex gap-0">
			@components.Sidebar(sidebar, meta.Path)
			<article class="flex-1 min-w-0 py-8 px-4 lg:px-8">
				<div class="prose prose-invert prose-slate max-w-none prose-headings:text-slate-100 prose-a:text-sky-400 prose-code:bg-slate-800 prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded prose-pre:bg-slate-800/50 prose-pre:border prose-pre:border-slate-700">
					{ children... }
				</div>
			</article>
			@components.TOC(toc)
		</div>
	}
}
```

**Step 7: Create the blog layout**

Create `templates/layouts/blog.templ`:

```templ
package layouts

templ Blog(meta PageMeta) {
	@Base(meta) {
		<div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<article class="prose prose-invert prose-slate max-w-none prose-headings:text-slate-100 prose-a:text-sky-400 prose-code:bg-slate-800 prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded prose-pre:bg-slate-800/50 prose-pre:border prose-pre:border-slate-700">
				{ children... }
			</article>
		</div>
	}
}
```

**Step 8: Create the landing page layout**

Create `templates/layouts/landing.templ`:

```templ
package layouts

templ Landing(meta PageMeta) {
	@Base(meta) {
		{ children... }
	}
}
```

**Step 9: Generate Templ code and verify it compiles**

```bash
cd /home/bjk/projects/frostyard/site && templ generate && go build ./...
```

Expected: Clean compile, no errors.

**Step 10: Commit**

```bash
git add templates/ go.mod go.sum
git commit -m "feat: add Templ templates for base layout, nav, sidebar, TOC, footer, and page layouts"
```

---

## Task 5: Build Pipeline — Orchestrate Content to HTML

Wire up the build command: load content, render through templates, write HTML files.

**Files:**
- Create: `internal/build/build.go`
- Create: `internal/build/build_test.go`
- Create: `internal/render/render.go`

**Step 1: Create the renderer**

Create `internal/render/render.go` — converts content.Page + Templ templates → HTML strings:

```go
package render

import (
	"bytes"
	"context"
	"html/template"

	"github.com/a-h/templ"
	"github.com/frostyard/site/internal/content"
	"github.com/frostyard/site/templates/components"
	"github.com/frostyard/site/templates/layouts"
)

// RenderDocsPage renders a documentation page to an HTML string.
func RenderDocsPage(page *content.Page, site *content.Site) (string, error) {
	meta := layouts.PageMeta{
		Title:       page.Title,
		Description: page.Description,
		Path:        page.Path,
		SiteName:    "Frostyard",
	}
	sidebar := buildSidebar(site.Sections)
	toc := buildTOC(page.Headings)

	contentComponent := templ.Raw(string(page.Content))
	component := layouts.Docs(meta, sidebar, toc)

	return renderComponent(component, contentComponent)
}

// RenderBlogPost renders a blog post to an HTML string.
func RenderBlogPost(page *content.Page) (string, error) {
	meta := layouts.PageMeta{
		Title:       page.Title,
		Description: page.Description,
		Path:        page.Path,
		SiteName:    "Frostyard",
	}

	contentComponent := templ.Raw(string(page.Content))
	component := layouts.Blog(meta)

	return renderComponent(component, contentComponent)
}

// RenderLandingPage renders the home page to an HTML string.
func RenderLandingPage(pageContent template.HTML) (string, error) {
	meta := layouts.PageMeta{
		SiteName: "Frostyard",
		Path:     "/",
	}

	contentComponent := templ.Raw(string(pageContent))
	component := layouts.Landing(meta)

	return renderComponent(component, contentComponent)
}

func renderComponent(wrapper templ.Component, children ...templ.Component) (string, error) {
	var buf bytes.Buffer
	// If we have children, we need to use a wrapper approach
	if len(children) > 0 {
		combined := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			return wrapper.Render(templ.WithChildren(ctx, children[0]), w)
		})
		if err := combined.Render(context.Background(), &buf); err != nil {
			return "", err
		}
	} else {
		if err := wrapper.Render(context.Background(), &buf); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func buildSidebar(sections []*content.Section) []components.SidebarSection {
	var result []components.SidebarSection
	for _, sec := range sections {
		result = append(result, convertSection(sec))
	}
	return result
}

func convertSection(sec *content.Section) components.SidebarSection {
	s := components.SidebarSection{
		Title: sec.Title,
		Path:  sec.Path,
	}
	for _, p := range sec.Pages {
		s.Pages = append(s.Pages, components.SidebarLink{
			Title: p.Title,
			Path:  p.Path,
		})
	}
	for _, sub := range sec.Subsections {
		s.Subsections = append(s.Subsections, convertSection(sub))
	}
	return s
}

func buildTOC(headings []content.Heading) []components.TOCHeading {
	var result []components.TOCHeading
	for _, h := range headings {
		result = append(result, components.TOCHeading{
			Level: h.Level,
			ID:    h.ID,
			Text:  h.Text,
		})
	}
	return result
}
```

**Step 2: Create the build orchestrator**

Create `internal/build/build.go`:

```go
package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frostyard/site/internal/content"
	"github.com/frostyard/site/internal/render"
)

// Config holds build configuration.
type Config struct {
	ContentDir string // Path to content/ directory
	StaticDir  string // Path to static/ directory
	OutputDir  string // Path to dist/ output directory
}

// Build runs the full static site generation pipeline.
func Build(cfg Config) error {
	// Clean output directory
	if err := os.RemoveAll(cfg.OutputDir); err != nil {
		return fmt.Errorf("cleaning output dir: %w", err)
	}

	// Load all content
	site, err := content.LoadContent(cfg.ContentDir)
	if err != nil {
		return fmt.Errorf("loading content: %w", err)
	}

	fmt.Printf("Loaded %d pages, %d blog posts\n", len(site.Pages), len(site.Posts))

	// Render each page
	for _, page := range site.Pages {
		if err := renderPage(page, site, cfg.OutputDir); err != nil {
			return fmt.Errorf("rendering %s: %w", page.Path, err)
		}
	}

	// Copy static assets
	if err := copyDir(cfg.StaticDir, filepath.Join(cfg.OutputDir)); err != nil {
		return fmt.Errorf("copying static assets: %w", err)
	}

	// Generate sitemap
	if err := generateSitemap(site, cfg.OutputDir); err != nil {
		return fmt.Errorf("generating sitemap: %w", err)
	}

	fmt.Printf("Build complete: %s\n", cfg.OutputDir)
	return nil
}

func renderPage(page *content.Page, site *content.Site, outputDir string) error {
	var html string
	var err error

	switch {
	case strings.HasPrefix(page.Path, "/blog/posts/"):
		html, err = render.RenderBlogPost(page)
	case page.Path == "/":
		html, err = render.RenderLandingPage(page.Content)
	default:
		html, err = render.RenderDocsPage(page, site)
	}
	if err != nil {
		return err
	}

	// Write HTML to output directory
	outPath := filepath.Join(outputDir, page.Path, "index.html")
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(html), 0o644)
}

func copyDir(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil // No static dir is fine
	}
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, 0o644)
	})
}
```

**Step 3: Create build_test.go**

Create `internal/build/build_test.go`:

```go
package build

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	dir := t.TempDir()

	// Create minimal content
	files := map[string]string{
		"content/docs/_index.md": `---
title: "Docs"
---
Welcome to the docs.
`,
		"content/docs/intro.md": `---
title: "Introduction"
weight: 1
---

# Introduction

Hello world.
`,
	}

	for path, body := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	outputDir := filepath.Join(dir, "dist")
	cfg := Config{
		ContentDir: filepath.Join(dir, "content"),
		StaticDir:  filepath.Join(dir, "static"),
		OutputDir:  outputDir,
	}

	if err := Build(cfg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify output files exist
	indexHTML := filepath.Join(outputDir, "docs", "intro", "index.html")
	if _, err := os.Stat(indexHTML); os.IsNotExist(err) {
		t.Errorf("expected %s to exist", indexHTML)
	}
}
```

**Step 4: Run tests**

```bash
templ generate && go test ./internal/build/ -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/build/ internal/render/
git commit -m "feat: add build pipeline and HTML renderer"
```

---

## Task 6: Sitemap and RSS Generation

**Files:**
- Create: `internal/build/sitemap.go`
- Create: `internal/build/rss.go`

**Step 1: Implement sitemap generation**

Create `internal/build/sitemap.go`:

```go
package build

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/frostyard/site/internal/content"
)

const siteURL = "https://frostyard.github.io"

type urlSet struct {
	XMLName xml.Name  `xml:"urlset"`
	XMLNS   string    `xml:"xmlns,attr"`
	URLs    []siteURL_ `xml:"url"`
}

type siteURL_ struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

func generateSitemap(site *content.Site, outputDir string) error {
	us := urlSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	now := time.Now().Format("2006-01-02")
	for _, page := range site.Pages {
		us.URLs = append(us.URLs, siteURL_{
			Loc:     siteURL + page.Path,
			LastMod: now,
		})
	}

	data, err := xml.MarshalIndent(us, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling sitemap: %w", err)
	}

	out := append([]byte(xml.Header), data...)
	return os.WriteFile(filepath.Join(outputDir, "sitemap.xml"), out, 0o644)
}
```

**Step 2: Implement RSS feed generation**

Create `internal/build/rss.go`:

```go
package build

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/frostyard/site/internal/content"
)

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate,omitempty"`
}

func generateRSS(site *content.Site, outputDir string) error {
	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:       "Frostyard Blog",
			Link:        siteURL + "/blog/",
			Description: "Updates from the Frostyard ecosystem",
		},
	}

	for _, post := range site.Posts {
		item := rssItem{
			Title:       post.Title,
			Link:        siteURL + post.Path,
			Description: post.Description,
		}
		if !post.ParsedDate.IsZero() {
			item.PubDate = post.ParsedDate.Format("Mon, 02 Jan 2006 15:04:05 -0700")
		}
		feed.Channel.Items = append(feed.Channel.Items, item)
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling RSS: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(outputDir, "blog"), 0o755); err != nil {
		return err
	}

	out := append([]byte(xml.Header), data...)
	return os.WriteFile(filepath.Join(outputDir, "blog", "feed.xml"), out, 0o644)
}
```

**Step 3: Wire RSS into the build pipeline**

In `internal/build/build.go`, add after the sitemap generation call:

```go
// In the Build function, after generateSitemap:
if err := generateRSS(site, cfg.OutputDir); err != nil {
    return fmt.Errorf("generating RSS: %w", err)
}
```

**Step 4: Run tests and verify build**

```bash
go test ./internal/build/ -v
```

**Step 5: Commit**

```bash
git add internal/build/sitemap.go internal/build/rss.go internal/build/build.go
git commit -m "feat: add sitemap.xml and RSS feed generation"
```

---

## Task 7: CLI Entry Point

Create the main CLI with `build`, `serve`, and `new` subcommands.

**Files:**
- Create: `cmd/frostyard/main.go`

**Step 1: Implement the CLI**

Create `cmd/frostyard/main.go`:

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/frostyard/site/internal/build"
	"github.com/frostyard/site/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Find project root (directory containing go.mod)
	root, err := findProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		cfg := build.Config{
			ContentDir: filepath.Join(root, "content"),
			StaticDir:  filepath.Join(root, "static"),
			OutputDir:  filepath.Join(root, "dist"),
		}
		if err := build.Build(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
			os.Exit(1)
		}

	case "serve":
		addr := ":3000"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		cfg := server.Config{
			ContentDir: filepath.Join(root, "content"),
			StaticDir:  filepath.Join(root, "static"),
			OutputDir:  filepath.Join(root, "dist"),
			Addr:       addr,
			Root:       root,
		}
		if err := server.Serve(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
			os.Exit(1)
		}

	case "new":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: frostyard new <page|post> <path-or-title>")
			os.Exit(1)
		}
		if err := scaffoldContent(root, os.Args[2], os.Args[3]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: frostyard <command>

Commands:
  build              Build the static site
  serve [addr]       Start dev server (default :3000)
  new page <path>    Create a new page
  new post <title>   Create a new blog post`)
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (no go.mod found)")
		}
		dir = parent
	}
}

func scaffoldContent(root, kind, arg string) error {
	switch kind {
	case "page":
		return scaffoldPage(root, arg)
	case "post":
		return scaffoldPost(root, arg)
	default:
		return fmt.Errorf("unknown content type %q (use 'page' or 'post')", kind)
	}
}

func scaffoldPage(root, relPath string) error {
	if !strings.HasPrefix(relPath, "content/") {
		relPath = filepath.Join("content", relPath)
	}
	fullPath := filepath.Join(root, relPath)

	title := strings.TrimSuffix(filepath.Base(relPath), ".md")
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.Title(title)

	content := fmt.Sprintf(`---
title: "%s"
description: ""
weight: 0
draft: false
---

`, title)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return err
	}
	fmt.Printf("Created: %s\n", relPath)
	return nil
}

func scaffoldPost(root, title string) error {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s.md", date, slug)
	relPath := filepath.Join("content", "blog", "posts", filename)
	fullPath := filepath.Join(root, relPath)

	content := fmt.Sprintf(`---
title: "%s"
date: "%s"
author: ""
tags: []
description: ""
draft: true
---

`, title, date)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return err
	}
	fmt.Printf("Created: %s\n", relPath)
	return nil
}
```

**Step 2: Verify it compiles**

```bash
templ generate && go build ./cmd/frostyard/
```

**Step 3: Commit**

```bash
git add cmd/frostyard/
git commit -m "feat: add frostyard CLI with build, serve, and new commands"
```

---

## Task 8: Dev Server with Live Reload

**Files:**
- Create: `internal/server/server.go`

**Step 1: Implement the dev server**

Create `internal/server/server.go`:

```go
package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/frostyard/site/internal/build"
)

// Config holds dev server configuration.
type Config struct {
	ContentDir string
	StaticDir  string
	OutputDir  string
	Addr       string
	Root       string
}

// Serve starts the development server with file watching and live reload.
func Serve(cfg Config) error {
	// Initial build
	buildCfg := build.Config{
		ContentDir: cfg.ContentDir,
		StaticDir:  cfg.StaticDir,
		OutputDir:  cfg.OutputDir,
	}
	if err := build.Build(buildCfg); err != nil {
		return fmt.Errorf("initial build: %w", err)
	}

	// SSE clients for live reload
	var (
		mu      sync.Mutex
		clients []chan struct{}
	)

	// File watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating watcher: %w", err)
	}
	defer watcher.Close()

	// Watch content, templates, and static directories
	for _, dir := range []string{cfg.ContentDir, filepath.Join(cfg.Root, "templates"), cfg.StaticDir} {
		if err := watchRecursive(watcher, dir); err != nil {
			log.Printf("Warning: could not watch %s: %v", dir, err)
		}
	}

	// Rebuild on file changes (debounced)
	go func() {
		var timer *time.Timer
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(200*time.Millisecond, func() {
						log.Println("Rebuilding...")
						if err := build.Build(buildCfg); err != nil {
							log.Printf("Rebuild error: %v", err)
							return
						}
						log.Println("Rebuild complete")
						// Notify all SSE clients
						mu.Lock()
						for _, ch := range clients {
							select {
							case ch <- struct{}{}:
							default:
							}
						}
						mu.Unlock()
					})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	// SSE endpoint for live reload
	http.HandleFunc("/_reload", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch := make(chan struct{}, 1)
		mu.Lock()
		clients = append(clients, ch)
		mu.Unlock()

		defer func() {
			mu.Lock()
			for i, c := range clients {
				if c == ch {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			mu.Unlock()
		}()

		for {
			select {
			case <-ch:
				fmt.Fprintf(w, "data: reload\n\n")
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	})

	// Inject live reload script into HTML responses
	fs := http.FileServer(http.Dir(cfg.OutputDir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve the file
		path := r.URL.Path
		if path == "/" || strings.HasSuffix(path, "/") {
			path = filepath.Join(path, "index.html")
		}
		fullPath := filepath.Join(cfg.OutputDir, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// Try with index.html
			fullPath = filepath.Join(cfg.OutputDir, r.URL.Path, "index.html")
		}

		if strings.HasSuffix(fullPath, ".html") {
			data, err := os.ReadFile(fullPath)
			if err != nil {
				fs.ServeHTTP(w, r)
				return
			}
			// Inject live reload script before </body>
			script := `<script>
const es = new EventSource('/_reload');
es.onmessage = () => location.reload();
es.onerror = () => setTimeout(() => location.reload(), 1000);
</script>`
			html := strings.Replace(string(data), "</body>", script+"</body>", 1)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(html))
			return
		}

		fs.ServeHTTP(w, r)
	})

	fmt.Printf("Dev server running at http://localhost%s\n", cfg.Addr)
	fmt.Println("Watching for changes...")
	return http.ListenAndServe(cfg.Addr, nil)
}

func watchRecursive(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}
```

**Step 2: Install fsnotify and verify compile**

```bash
go mod tidy && templ generate && go build ./cmd/frostyard/
```

**Step 3: Commit**

```bash
git add internal/server/ go.mod go.sum
git commit -m "feat: add dev server with file watching and SSE live reload"
```

---

## Task 9: Tailwind CSS Setup

**Files:**
- Create: `input.css`
- Create: `tailwind.config.js`
- Modify: `internal/build/build.go` (add Tailwind step)

**Step 1: Download Tailwind standalone CLI**

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64
mv tailwindcss-linux-x64 tailwindcss
```

**Step 2: Create Tailwind config**

Create `tailwind.config.js`:

```js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.templ",
    "./templates/**/*_templ.go",
    "./dist/**/*.html",
  ],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        frost: {
          50: "#f0f9ff",
          100: "#e0f2fe",
          200: "#bae6fd",
          300: "#7dd3fc",
          400: "#38bdf8",
          500: "#0ea5e9",
        },
      },
    },
  },
  plugins: [
    require("@tailwindcss/typography"),
  ],
};
```

Note: Tailwind v4 uses CSS-based config. If the standalone CLI is v4, we'll use `input.css` with `@import` directives instead. Adjust accordingly during implementation.

**Step 3: Create input.css**

Create `input.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

**Step 4: Add Tailwind build step to the build pipeline**

In `internal/build/build.go`, add a function to run Tailwind and call it from `Build()`:

```go
func runTailwind(root, outputDir string) error {
	cssDir := filepath.Join(outputDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		return err
	}
	cmd := exec.Command(
		filepath.Join(root, "tailwindcss"),
		"-i", filepath.Join(root, "input.css"),
		"-o", filepath.Join(cssDir, "style.css"),
		"--minify",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
```

Add `Root string` to `build.Config` and call `runTailwind(cfg.Root, cfg.OutputDir)` in the Build function after rendering pages.

**Step 5: Verify Tailwind generates CSS**

```bash
./tailwindcss -i input.css -o dist/css/style.css
```

**Step 6: Commit**

```bash
git add input.css tailwind.config.js internal/build/build.go
git commit -m "feat: add Tailwind CSS setup and build integration"
```

---

## Task 10: Starter Content and Static Pages

Create the home page, downloads page, community page, and seed docs content.

**Files:**
- Create: `templates/pages/home.templ`
- Create: `templates/pages/downloads.templ`
- Create: `templates/pages/community.templ`
- Create: `content/docs/_index.md`
- Create: `content/docs/getting-started/_index.md`
- Create: `content/docs/getting-started/quickstart.md`
- Create: `content/docs/images/_index.md`
- Create: `content/docs/tools/_index.md`
- Create: `content/docs/faq.md`
- Create: `content/blog/_index.md`

**Step 1: Create the home page template**

Create `templates/pages/home.templ`:

```templ
package pages

import "github.com/frostyard/site/templates/layouts"

templ Home() {
	@layouts.Landing(layouts.PageMeta{SiteName: "Frostyard", Path: "/"}) {
		<!-- Hero -->
		<section class="py-20 px-4">
			<div class="max-w-4xl mx-auto text-center">
				<h1 class="text-5xl font-bold text-slate-100 mb-4">
					<span class="text-sky-400">Frost</span>yard
				</h1>
				<p class="text-xl text-slate-400 mb-8 max-w-2xl mx-auto">
					Secure, reproducible OCI-based Linux images built with bootc.
					Immutable by design. Atomic updates. Your infrastructure, containerized.
				</p>
				<div class="flex gap-4 justify-center">
					<a href="/docs/getting-started/" class="px-6 py-3 bg-sky-500 hover:bg-sky-400 text-white rounded-lg font-medium transition-colors">
						Get Started
					</a>
					<a href="https://github.com/frostyard" class="px-6 py-3 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded-lg font-medium transition-colors border border-slate-700">
						GitHub
					</a>
				</div>
			</div>
		</section>
		<!-- Feature cards -->
		<section class="py-16 px-4 border-t border-slate-800">
			<div class="max-w-5xl mx-auto grid md:grid-cols-3 gap-8">
				<div class="p-6 rounded-lg bg-slate-800/50 border border-slate-700">
					<h3 class="text-lg font-semibold text-slate-100 mb-2">Images</h3>
					<p class="text-slate-400 text-sm">Pre-built bootc container images: Snow, Snowfield, Cayo, and more. Ready to deploy.</p>
					<a href="/docs/images/" class="text-sky-400 text-sm mt-3 inline-block hover:text-sky-300">Browse images &rarr;</a>
				</div>
				<div class="p-6 rounded-lg bg-slate-800/50 border border-slate-700">
					<h3 class="text-lg font-semibold text-slate-100 mb-2">Tools</h3>
					<p class="text-slate-400 text-sm">nbc for disk installs, Chairlift for system management, First Setup for onboarding.</p>
					<a href="/docs/tools/" class="text-sky-400 text-sm mt-3 inline-block hover:text-sky-300">Explore tools &rarr;</a>
				</div>
				<div class="p-6 rounded-lg bg-slate-800/50 border border-slate-700">
					<h3 class="text-lg font-semibold text-slate-100 mb-2">Atomic Updates</h3>
					<p class="text-slate-400 text-sm">A/B partition scheme with automatic rollback. Updates are transactional and safe.</p>
					<a href="/docs/getting-started/" class="text-sky-400 text-sm mt-3 inline-block hover:text-sky-300">Learn more &rarr;</a>
				</div>
			</div>
		</section>
	}
}
```

**Step 2: Create seed markdown content**

Create all the `_index.md` and starter content files listed above with appropriate frontmatter and brief placeholder content. The key docs structure:

- `content/docs/_index.md` — title: "Documentation"
- `content/docs/getting-started/_index.md` — title: "Getting Started", weight: 1
- `content/docs/getting-started/quickstart.md` — title: "Quickstart", weight: 1
- `content/docs/images/_index.md` — title: "Images", weight: 2
- `content/docs/tools/_index.md` — title: "Tools", weight: 3
- `content/docs/faq.md` — title: "FAQ", weight: 99
- `content/blog/_index.md` — title: "Blog"

**Step 3: Create downloads and community Templ pages**

Create `templates/pages/downloads.templ` and `templates/pages/community.templ` as simple static pages rendered through the base layout with hardcoded content (links to GitHub releases, container registries, community channels).

**Step 4: Verify full build**

```bash
templ generate && go run ./cmd/frostyard build
```

**Step 5: Commit**

```bash
git add content/ templates/pages/ static/
git commit -m "feat: add starter content, home page, downloads, and community pages"
```

---

## Task 11: Justfile and Tooling Setup

**Files:**
- Modify: `Justfile`

**Step 1: Write the new Justfile**

Replace the existing Justfile:

```just
# Frostyard Site
# Go static site generator with Templ + Tailwind

default:
    just --list --unsorted

# Generate templ files and build the static site
build:
    templ generate
    go run ./cmd/frostyard build

# Start dev server with live reload
serve:
    templ generate
    go run ./cmd/frostyard serve

# Run all tests
test:
    go test ./... -v

# Generate templ Go code
generate:
    templ generate

# Deploy to GitHub Pages
deploy: build
    @echo "Deploying to GitHub Pages..."
    ghp-import -p -f -b pages dist

# Create a new docs page
new-page path:
    go run ./cmd/frostyard new page {{ path }}

# Create a new blog post
new-post title:
    go run ./cmd/frostyard new post "{{ title }}"

# Run Pagefind to build search index (post-build)
search-index:
    pagefind --site dist

# Clean build artifacts
clean:
    rm -rf dist
```

**Step 2: Verify Justfile works**

```bash
just build
```

**Step 3: Commit**

```bash
git add Justfile
git commit -m "chore: replace MkDocs Justfile with Go SSG build commands"
```

---

## Task 12: Port nbc Documentation

Migrate the nbc documentation from the saved submodule content into the new content structure.

**Files:**
- Create: `content/docs/tools/nbc/_index.md`
- Create: `content/docs/tools/nbc/` (multiple files ported from /tmp/frostyard-nbc-docs/)

**Step 1: Review saved nbc docs**

```bash
ls /tmp/frostyard-nbc-docs/
```

**Step 2: Create nbc section index and port each doc file**

For each nbc doc file, create a corresponding file under `content/docs/tools/nbc/` with proper frontmatter added. The main feature docs (AB-UPDATES.md, ENCRYPTION.md, etc.) become pages. The CLI reference files go under `content/docs/tools/nbc/cli/`.

Port the content as-is with frontmatter headers added. Don't rewrite the docs — just add the YAML frontmatter block and adjust any internal links.

**Step 3: Verify build with ported content**

```bash
just build
```

**Step 4: Commit**

```bash
git add content/docs/tools/nbc/
git commit -m "feat: port nbc documentation from submodule"
```

---

## Task 13: Pagefind Search Integration

**Files:**
- Modify: `internal/build/build.go` (add Pagefind step)

**Step 1: Install Pagefind**

```bash
# Install via npm (or download binary)
npx pagefind --version
# Or download standalone binary
```

**Step 2: Add Pagefind as post-build step**

In `internal/build/build.go`, add a function to run Pagefind after the HTML is generated:

```go
func runPagefind(root, outputDir string) error {
	cmd := exec.Command("pagefind", "--site", outputDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
```

Call it at the end of `Build()`. If pagefind is not found, log a warning and skip (it's optional for dev builds).

**Step 3: Verify search works**

```bash
just build && just search-index && just serve
```

Open browser, test search.

**Step 4: Commit**

```bash
git add internal/build/build.go
git commit -m "feat: integrate Pagefind for client-side search"
```

---

## Task 14: GitHub Actions Workflow

**Files:**
- Create: `.github/workflows/deploy.yml`

**Step 1: Create the deployment workflow**

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.26"

      - name: Install Templ
        run: go install github.com/a-h/templ/cmd/templ@latest

      - name: Install Tailwind CSS
        run: |
          curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
          chmod +x tailwindcss-linux-x64
          mv tailwindcss-linux-x64 tailwindcss

      - name: Install Pagefind
        run: npx pagefind --version || npm install -g pagefind

      - name: Build site
        run: |
          templ generate
          go run ./cmd/frostyard build

      - name: Build search index
        run: pagefind --site dist

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: dist

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

**Step 2: Commit**

```bash
git add .github/workflows/deploy.yml
git commit -m "ci: add GitHub Actions workflow for building and deploying to GitHub Pages"
```

---

## Task 15: End-to-End Verification

**Step 1: Clean build from scratch**

```bash
just clean && just build
```

**Step 2: Verify output structure**

```bash
find dist -type f | head -30
```

Expected: HTML files in proper directory structure, CSS, sitemap.xml, RSS feed.

**Step 3: Start dev server and test manually**

```bash
just serve
```

Verify in browser:
- Home page renders with frost theme
- Navigation works
- Docs pages have sidebar and TOC
- Blog section exists
- Dark mode toggle works
- Search box appears (empty until Pagefind runs)

**Step 4: Run all tests**

```bash
just test
```

Expected: ALL PASS

**Step 5: Final commit if any fixes needed**

```bash
git add -A
git commit -m "fix: end-to-end verification fixes"
```
