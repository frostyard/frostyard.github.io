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
