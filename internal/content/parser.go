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
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"gopkg.in/yaml.v3"
)

// ParsePage parses a markdown file with YAML frontmatter and returns a Page.
// sourcePath is the filesystem path relative to the project root (e.g., "content/docs/tools/nbc/install.md").
func ParsePage(data []byte, sourcePath string) (*Page, error) {
	fm, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("splitting frontmatter: %w", err)
	}

	var page Page
	if len(fm) > 0 {
		if err := yaml.Unmarshal(fm, &page); err != nil {
			return nil, fmt.Errorf("parsing frontmatter: %w", err)
		}
	}

	html, headings, err := renderMarkdown(body)
	if err != nil {
		return nil, fmt.Errorf("rendering markdown: %w", err)
	}

	page.Content = template.HTML(html)
	page.Headings = headings
	page.SourcePath = sourcePath
	page.Path = computePath(sourcePath)
	page.Slug = computeSlug(sourcePath)
	page.IsIndex = strings.HasSuffix(sourcePath, "_index.md")

	if page.Date != "" {
		parsed, err := parseDate(page.Date)
		if err == nil {
			page.ParsedDate = parsed
		}
	}

	return &page, nil
}

// splitFrontmatter splits YAML frontmatter (delimited by ---) from the markdown body.
// Returns empty frontmatter if no frontmatter delimiters are found.
func splitFrontmatter(data []byte) (frontmatter, body []byte, err error) {
	content := string(data)

	// Frontmatter must start with ---
	if !strings.HasPrefix(content, "---") {
		return nil, data, nil
	}

	// Find the closing ---
	rest := content[3:]
	// Skip the newline after opening ---
	if idx := strings.Index(rest, "\n"); idx >= 0 {
		rest = rest[idx+1:]
	} else {
		return nil, data, nil
	}

	endIdx := strings.Index(rest, "\n---")
	if endIdx < 0 {
		return nil, data, nil
	}

	frontmatter = []byte(rest[:endIdx])
	body = []byte(rest[endIdx+4:]) // skip \n---

	// Skip the newline after closing ---
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	}

	return frontmatter, body, nil
}

// computePath derives the URL path from a source file path.
// It strips the "content/" prefix, removes .md extension, handles _index.md files,
// and ensures leading and trailing slashes.
func computePath(sourcePath string) string {
	// Normalize to forward slashes
	p := filepath.ToSlash(sourcePath)

	// Strip "content/" prefix
	p = strings.TrimPrefix(p, "content/")

	// Remove .md extension
	p = strings.TrimSuffix(p, ".md")

	// Handle _index files: remove the _index part
	if strings.HasSuffix(p, "/_index") {
		p = strings.TrimSuffix(p, "/_index")
	} else if p == "_index" {
		p = ""
	}

	// Ensure leading slash
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	// Ensure trailing slash
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}

	return p
}

// computeSlug returns the URL-friendly name derived from the filename.
func computeSlug(sourcePath string) string {
	base := filepath.Base(sourcePath)
	slug := strings.TrimSuffix(base, ".md")
	return slug
}

// renderMarkdown converts markdown source to HTML and extracts headings.
// Uses goldmark with GFM extension, syntax highlighting (Chroma dracula style with CSS classes),
// and auto heading IDs.
func renderMarkdown(source []byte) (string, []Heading, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithCSSWriter(nil),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	// Parse to AST to extract headings
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	headings := extractHeadings(doc, source)

	// Render to HTML
	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, source, doc); err != nil {
		return "", nil, fmt.Errorf("rendering markdown: %w", err)
	}

	return buf.String(), headings, nil
}

// extractHeadings walks the AST and extracts all headings with their level, ID, and text.
func extractHeadings(node ast.Node, source []byte) []Heading {
	var headings []Heading

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if heading, ok := n.(*ast.Heading); ok {
			h := Heading{
				Level: heading.Level,
			}

			// Get the auto-generated ID
			if id, ok := heading.AttributeString("id"); ok {
				h.ID = string(id.([]byte))
			}

			// Get the text content of the heading
			var textBuf bytes.Buffer
			for child := heading.FirstChild(); child != nil; child = child.NextSibling() {
				if t, ok := child.(*ast.Text); ok {
					textBuf.Write(t.Segment.Value(source))
				}
			}
			h.Text = textBuf.String()

			headings = append(headings, h)
		}

		return ast.WalkContinue, nil
	})

	return headings
}

// parseDate tries several common date formats.
func parseDate(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"January 2, 2006",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date %q", s)
}
