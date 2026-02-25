package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/frostyard/site/internal/build"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	root, err := findProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch cmd {
	case "build":
		cfg := build.Config{
			ContentDir: filepath.Join(root, "content"),
			StaticDir:  filepath.Join(root, "static"),
			OutputDir:  filepath.Join(root, "dist"),
			Root:       root,
		}
		if err := build.Build(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
			os.Exit(1)
		}

	case "serve":
		fmt.Println("Server not yet implemented")

	case "new":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: frostyard new <page|post> <args>\n")
			os.Exit(1)
		}
		subCmd := os.Args[2]
		switch subCmd {
		case "page":
			if len(os.Args) < 4 {
				fmt.Fprintf(os.Stderr, "Usage: frostyard new page <path>\n")
				os.Exit(1)
			}
			relPath := os.Args[3]
			if err := scaffoldPage(root, relPath); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating page: %v\n", err)
				os.Exit(1)
			}
		case "post":
			if len(os.Args) < 4 {
				fmt.Fprintf(os.Stderr, "Usage: frostyard new post <title>\n")
				os.Exit(1)
			}
			title := strings.Join(os.Args[3:], " ")
			if err := scaffoldPost(root, title); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating post: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown new subcommand: %s\n", subCmd)
			printUsage()
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: frostyard <command> [args]

Commands:
  build              Build the site to dist/
  serve              Start a local development server
  new page <path>    Create a new page (e.g., docs/guides/setup)
  new post <title>   Create a new blog post`)
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working directory: %w", err)
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

func scaffoldPage(root, relPath string) error {
	// Ensure .md extension
	if !strings.HasSuffix(relPath, ".md") {
		relPath = relPath + ".md"
	}

	fullPath := filepath.Join(root, "content", relPath)

	// Don't overwrite existing files
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("file already exists: %s", fullPath)
	}

	// Derive title from the filename
	base := filepath.Base(relPath)
	title := strings.TrimSuffix(base, ".md")
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.Title(title) //nolint:staticcheck

	content := fmt.Sprintf(`---
title: "%s"
---
`, title)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("Created page: %s\n", fullPath)
	return nil
}

func scaffoldPost(root, title string) error {
	slug := slugify(title)
	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s.md", date, slug)

	postsDir := filepath.Join(root, "content", "blog", "posts")
	fullPath := filepath.Join(postsDir, filename)

	// Don't overwrite existing files
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("file already exists: %s", fullPath)
	}

	content := fmt.Sprintf(`---
title: "%s"
date: "%s"
author: ""
description: ""
tags: []
---
`, title, date)

	if err := os.MkdirAll(postsDir, 0o755); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("Created post: %s\n", fullPath)
	return nil
}

// slugify converts a title to a URL-friendly slug.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' {
			return r
		}
		return -1
	}, s)
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`[\s-]+`).ReplaceAllString(s, "-")
	return s
}
