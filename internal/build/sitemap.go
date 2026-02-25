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
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []urlEntry
}

type urlEntry struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
	LastMod string   `xml:"lastmod"`
}

func generateSitemap(site *content.Site, outputDir string) error {
	today := time.Now().Format("2006-01-02")

	set := urlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	for _, page := range site.Pages {
		set.URLs = append(set.URLs, urlEntry{
			Loc:     siteURL + page.Path,
			LastMod: today,
		})
	}

	data, err := xml.MarshalIndent(set, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling sitemap: %w", err)
	}

	out := xml.Header + string(data)
	outPath := filepath.Join(outputDir, "sitemap.xml")

	if err := os.WriteFile(outPath, []byte(out), 0o644); err != nil {
		return fmt.Errorf("writing sitemap: %w", err)
	}

	return nil
}
