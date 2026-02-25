package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LoadContent walks contentDir, parses all .md files, skips drafts,
// separates blog posts, and builds a section tree.
func LoadContent(contentDir string) (*Site, error) {
	var allPages []*Page
	var posts []*Page

	// The parent of contentDir — sourcePaths should be relative to this
	// so that they start with "content/".
	baseDir := filepath.Dir(contentDir)

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .md files
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		// Compute sourcePath relative to baseDir so it starts with "content/"
		sourcePath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return fmt.Errorf("computing relative path for %s: %w", path, err)
		}
		// Normalize to forward slashes for consistent path handling
		sourcePath = filepath.ToSlash(sourcePath)

		page, err := ParsePage(data, sourcePath)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", sourcePath, err)
		}

		// Skip drafts
		if page.Draft {
			return nil
		}

		allPages = append(allPages, page)

		// Identify blog posts: path starts with /blog/posts/ and is not an index page
		if strings.HasPrefix(page.Path, "/blog/posts/") && !page.IsIndex {
			posts = append(posts, page)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking content directory: %w", err)
	}

	// Sort posts by date descending
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ParsedDate.After(posts[j].ParsedDate)
	})

	// Build section tree
	sections := buildSectionTree(allPages)

	return &Site{
		Pages:    allPages,
		Posts:    posts,
		Sections: sections,
	}, nil
}

// buildSectionTree organizes pages into a hierarchical section tree.
// Two-pass algorithm:
//  1. Create sections from _index.md pages (keyed by page.Path).
//  2. Assign non-index pages to their parent section, sort pages by weight,
//     build hierarchy by assigning subsections to parents, sort subsections by weight,
//     and return root sections.
func buildSectionTree(pages []*Page) []*Section {
	sectionMap := make(map[string]*Section)

	// Pass 1: Create sections from _index.md pages
	for _, p := range pages {
		if !p.IsIndex {
			continue
		}
		sec := &Section{
			Title:       p.Title,
			Description: p.Description,
			Icon:        p.Icon,
			Path:        p.Path,
			IndexPage:   p,
			Weight:      p.Weight,
		}
		sectionMap[p.Path] = sec
	}

	// Pass 2: Assign non-index pages to their parent section
	for _, p := range pages {
		if p.IsIndex {
			continue
		}
		// Parent path is the directory of the page's URL path, with trailing slash
		parentPath := filepath.Dir(strings.TrimSuffix(p.Path, "/"))
		if !strings.HasSuffix(parentPath, "/") {
			parentPath += "/"
		}

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

	// Build hierarchy: assign subsections to their parent sections
	for path, sec := range sectionMap {
		parentPath := filepath.Dir(strings.TrimSuffix(path, "/"))
		if !strings.HasSuffix(parentPath, "/") {
			parentPath += "/"
		}

		if parent, ok := sectionMap[parentPath]; ok {
			parent.Subsections = append(parent.Subsections, sec)
		}
	}

	// Sort subsections by weight
	for _, sec := range sectionMap {
		sort.Slice(sec.Subsections, func(i, j int) bool {
			return sec.Subsections[i].Weight < sec.Subsections[j].Weight
		})
	}

	// Collect root sections: those whose parent is not in the map
	var roots []*Section
	for path, sec := range sectionMap {
		parentPath := filepath.Dir(strings.TrimSuffix(path, "/"))
		if !strings.HasSuffix(parentPath, "/") {
			parentPath += "/"
		}

		if _, ok := sectionMap[parentPath]; !ok {
			roots = append(roots, sec)
		}
	}

	// Sort root sections by weight
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].Weight < roots[j].Weight
	})

	return roots
}
