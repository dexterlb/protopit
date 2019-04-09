package builder

import (
	"io"
	"io/ioutil"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

// return (ast.GoToNext, true) to tell html renderer to skip rendering this node
// (because you've rendered it)
func renderHookDropCodeBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	// skip all nodes that are not CodeBlock nodes
	if _, ok := node.(*ast.CodeBlock); !ok {
		return ast.GoToNext, false
	}
	// custom rendering logic for ast.CodeBlock. By doing nothing it won't be
	// present in the output
	return ast.GoToNext, true
}

func (s *Site) ParsePage(filename string, name string) {
	page := &Page{
		Name:     name,
		Filename: filename,
		Type:     "normal.html",
	}
	s.ParseMarkdownPage(page)
	s.Pages[page.Name] = page
}

func (s *Site) ParseMarkdownPage(page *Page) {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: renderHookDropCodeBlock,
	}
	renderer := html.NewRenderer(opts)
	data, err := ioutil.ReadFile(page.Filename)
	noerr("cannot read file", err)
	html := markdown.ToHTML(data, nil, renderer)
	page.Html = html
}
