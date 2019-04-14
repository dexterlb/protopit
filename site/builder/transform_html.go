package builder

import (
	"bytes"
	"fmt"
	"html/template"

	"golang.org/x/net/html"
)

func (s *Site) transformHtml(p *Page, node *html.Node) *html.Node {
	newNodes := s.replaceNode(p, node)
	if newNodes != nil {
		parent := node.Parent
		next := node.NextSibling
		parent.RemoveChild(node)
		for i := range newNodes {
			parent.InsertBefore(newNodes[i], next)
		}
		if len(newNodes) > 0 {
			return newNodes[0]
		} else {
			return nil
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		child = s.transformHtml(p, child)
	}
	return node
}

type TransformData struct {
	Site *Site
	Page *Page
	Node *html.Node
}

func (s *Site) replaceNode(p *Page, node *html.Node) []*html.Node {
	if node.Type == html.ElementNode {
		transformer := fmt.Sprintf("transform_%s.html", string(node.Data))
		if s.Template.Lookup(transformer) != nil {
			return s.renderTemplateAt(transformer, &TransformData{
				Site: s,
				Page: p,
				Node: node,
			}, node.Parent)
		}
	}
	return nil
}

func (t *TransformData) Attr(name string) string {
		for _, a := range t.Node.Attr {
			if a.Key == name {
				return a.Val
			}
		}
		return ""
}

func (t *TransformData) Content() template.HTML {
	w := &bytes.Buffer{}
	for child := t.Node.FirstChild; child != nil; child = child.NextSibling {
		err := html.Render(w, child)
		noerr("cannot render html node", err)
	}
	return template.HTML(w.Bytes())
}
