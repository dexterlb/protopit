package builder

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

func markdownInspect(s *Site, p *Page) func(io.Writer, ast.Node, bool) (ast.WalkStatus, bool) {
	return func(w io.Writer, anyNode ast.Node, entering bool) (ast.WalkStatus, bool) {
		switch node := anyNode.(type) {
		case *ast.CodeBlock:
			if string(node.Info) == "meta" {
				p.Meta = ParseMetaData(node.Literal)
			}
		}
		return ast.GoToNext, false // do nothing
	}
}

func (s *Site) ParsePage(filename string, name string) {
	page := &Page{
		Name:     name,
		Filename: filename,
		Type:     "normal.html",
	}
	s.ParseMarkdownPage(page)
	if page.Meta == nil {
		noerr("page has no metadata", fmt.Errorf("no metadata in %s", filename))
	}
	s.Pages[page.Name] = page
}

func (s *Site) ParseMarkdownPage(page *Page) {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: markdownInspect(s, page),
	}
	renderer := html.NewRenderer(opts)
	data, err := ioutil.ReadFile(page.Filename)
	noerr("cannot read file", err)
	html := markdown.ToHTML(data, nil, renderer)
	page.Html = html
}
