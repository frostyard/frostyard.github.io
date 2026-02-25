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

// RenderDocsPage renders a docs page with sidebar navigation and table of contents.
func RenderDocsPage(page *content.Page, site *content.Site) (string, error) {
	meta := layouts.PageMeta{
		Title:    page.Title,
		Path:     page.Path,
		SiteName: "Frostyard",
	}
	if page.Description != "" {
		meta.Description = page.Description
	}

	sidebar := buildSidebar(site.Sections)
	toc := buildTOC(page.Headings)
	rawContent := templ.Raw(string(page.Content))

	wrapper := layouts.Docs(meta, sidebar, toc)
	return renderWithChildren(wrapper, rawContent)
}

// RenderBlogPost renders a blog post page.
func RenderBlogPost(page *content.Page) (string, error) {
	meta := layouts.PageMeta{
		Title:    page.Title,
		Path:     page.Path,
		SiteName: "Frostyard",
	}
	if page.Description != "" {
		meta.Description = page.Description
	}

	rawContent := templ.Raw(string(page.Content))
	wrapper := layouts.Blog(meta)
	return renderWithChildren(wrapper, rawContent)
}

// RenderLandingPage renders the home/landing page.
func RenderLandingPage(pageContent template.HTML) (string, error) {
	meta := layouts.PageMeta{
		Path:     "/",
		SiteName: "Frostyard",
	}

	rawContent := templ.Raw(string(pageContent))
	wrapper := layouts.Landing(meta)
	return renderWithChildren(wrapper, rawContent)
}

// renderWithChildren renders a templ component with children using templ.WithChildren context.
func renderWithChildren(wrapper templ.Component, child templ.Component) (string, error) {
	var buf bytes.Buffer
	ctx := templ.WithChildren(context.Background(), child)
	if err := wrapper.Render(ctx, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// buildSidebar converts the content section tree into sidebar component data.
func buildSidebar(sections []*content.Section) []components.SidebarSection {
	result := make([]components.SidebarSection, 0, len(sections))
	for _, sec := range sections {
		result = append(result, convertSection(sec))
	}
	return result
}

// convertSection recursively converts a content.Section to a components.SidebarSection.
func convertSection(sec *content.Section) components.SidebarSection {
	ss := components.SidebarSection{
		Title: sec.Title,
		Path:  sec.Path,
	}

	for _, p := range sec.Pages {
		ss.Pages = append(ss.Pages, components.SidebarLink{
			Title: p.Title,
			Path:  p.Path,
		})
	}

	for _, sub := range sec.Subsections {
		ss.Subsections = append(ss.Subsections, convertSection(sub))
	}

	return ss
}

// buildTOC converts content headings into TOC component data.
func buildTOC(headings []content.Heading) []components.TOCHeading {
	result := make([]components.TOCHeading, 0, len(headings))
	for _, h := range headings {
		result = append(result, components.TOCHeading{
			Level: h.Level,
			ID:    h.ID,
			Text:  h.Text,
		})
	}
	return result
}
