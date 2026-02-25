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
	PubDate     string `xml:"pubDate"`
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
		return fmt.Errorf("marshaling RSS feed: %w", err)
	}

	out := xml.Header + string(data)

	blogDir := filepath.Join(outputDir, "blog")
	if err := os.MkdirAll(blogDir, 0o755); err != nil {
		return fmt.Errorf("creating blog directory: %w", err)
	}

	outPath := filepath.Join(blogDir, "feed.xml")
	if err := os.WriteFile(outPath, []byte(out), 0o644); err != nil {
		return fmt.Errorf("writing RSS feed: %w", err)
	}

	return nil
}
