package builder

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/net/html"
)

type TemplateData struct {
	Site *Site
	Page *Page
}

func (s *Site) Render() {
	for _, page := range s.Pages {
		s.RenderPage(page)
	}
}

func (s *Site) RenderPage(p *Page) {
	// render html from template
	node := s.renderTemplate("main.html", TemplateData{
		Page: p,
		Site: s,
	})

	s.transformHtml(p, node)

	// create output file
	outDir := filepath.Join(s.OutputDir, s.PageUrl(p.Name))
	noerr("cannot create dir", os.MkdirAll(outDir, os.ModePerm))
	f, err := os.Create(filepath.Join(outDir, "index.html"))
	defer f.Close()
	noerr("cannot create output file", err)

	// render html
	err = html.Render(f, node)
	noerr("cannot render html", err)
}

func (s *Site) AbsPageUrl(name string) string {
	return filepath.Join("/", s.PageUrl(name))
}

func (s *Site) PageUrl(name string) string {
	if _, ok := s.Pages[name]; !ok {
		noerr("cannot get page url", fmt.Errorf("no such page: %s", name))
	}
	return filepath.Join(name, s.Variant)
}

func (s *Site) renderTemplate(name string, data interface{}) *html.Node {
	w := &bytes.Buffer{}
	err := s.Template.ExecuteTemplate(w, name, data)
	noerr("cannot render page", err)

	node, err := html.Parse(w)
	noerr("invalid HTML", err)

	return node
}

func (s *Site) renderTemplateAt(name string, data interface{}, context *html.Node) []*html.Node {
	w := &bytes.Buffer{}
	err := s.Template.ExecuteTemplate(w, name, data)
	noerr("cannot render page", err)

	nodes, err := html.ParseFragment(w, context)
	noerr("invalid HTML", err)

	return nodes
}
