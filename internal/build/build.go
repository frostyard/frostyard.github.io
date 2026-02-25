package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/frostyard/site/internal/content"
	"github.com/frostyard/site/internal/render"
)

// Config holds the build configuration.
type Config struct {
	ContentDir string // Path to content directory (e.g., "content")
	StaticDir  string // Path to static assets directory (e.g., "static")
	OutputDir  string // Path to output directory (e.g., "dist")
	Root       string // Project root directory
}

// Build orchestrates the full site build: load content, render HTML, copy static assets.
func Build(cfg Config) error {
	// Clean output directory
	if err := os.RemoveAll(cfg.OutputDir); err != nil {
		return fmt.Errorf("cleaning output directory: %w", err)
	}

	// Load content
	site, err := content.LoadContent(cfg.ContentDir)
	if err != nil {
		return fmt.Errorf("loading content: %w", err)
	}

	fmt.Printf("Loaded %d pages, %d blog posts\n", len(site.Pages), len(site.Posts))

	// Render each page to HTML
	for _, page := range site.Pages {
		if err := renderPage(page, site, cfg.OutputDir); err != nil {
			return fmt.Errorf("rendering %s: %w", page.Path, err)
		}
	}

	// Copy static assets
	if err := copyDir(cfg.StaticDir, cfg.OutputDir); err != nil {
		return fmt.Errorf("copying static assets: %w", err)
	}

	// Generate sitemap
	if err := generateSitemap(site, cfg.OutputDir); err != nil {
		return fmt.Errorf("generating sitemap: %w", err)
	}

	// Generate RSS feed
	if err := generateRSS(site, cfg.OutputDir); err != nil {
		return fmt.Errorf("generating RSS feed: %w", err)
	}

	fmt.Printf("Build complete: %s\n", cfg.OutputDir)
	return nil
}

// renderPage renders a single page to HTML and writes it to the output directory.
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
		return fmt.Errorf("rendering page: %w", err)
	}

	// Determine output file path: outputDir/page.Path/index.html
	outPath := filepath.Join(outputDir, page.Path, "index.html")

	// Create directory
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", outPath, err)
	}

	// Write HTML file
	if err := os.WriteFile(outPath, []byte(html), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", outPath, err)
	}

	return nil
}

// copyDir copies all files from src to dst recursively.
// Skips silently if src does not exist.
func copyDir(src, dst string) error {
	// Skip if source doesn't exist
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path from src
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		return copyFile(path, destPath)
	})
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
