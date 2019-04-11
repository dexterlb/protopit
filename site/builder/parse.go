package builder

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func markdownInspect(s *Site, p *Page) func(io.Writer, ast.Node, bool) (ast.WalkStatus, bool) {
	return func(w io.Writer, anyNode ast.Node, entering bool) (ast.WalkStatus, bool) {
		switch node := anyNode.(type) {
		case *ast.CodeBlock:
			if string(node.Info) == "meta" {
				p.Meta = ParseMetaData(node.Literal)
				return ast.GoToNext, true // skip
			}
		}
		return ast.GoToNext, false // do nothing
	}
}

func (s *Site) ParsePage(filename string, name string) {
	page := &Page{
		Name:     name,
		Filename: filename,
	}
	s.ParseMarkdownPage(page)
	if page.Meta == nil {
		noerr("page has no metadata", fmt.Errorf("no metadata in %s", filename))
	}
	if page.Meta.NameOverride != nil {
		page.Name = *page.Meta.NameOverride
	}
	s.Pages[page.Name] = page
}

func (s *Site) ParseMarkdownPage(page *Page) {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: markdownInspect(s, page),
	}
	renderer := html.NewRenderer(opts)

	extensions := parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.SuperSubscript |
		parser.Footnotes

	p := parser.NewWithExtensions(extensions)

	data, err := ioutil.ReadFile(page.Filename)
	noerr("cannot read file", err)
	html := markdown.ToHTML(data, p, renderer)
	page.Html = html
}
