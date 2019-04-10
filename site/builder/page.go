package builder

import "html/template"

type Page struct {
	Filename string
	Meta     *MetaData
	Name     string
	Html     []byte
}

func (p *Page) Content() template.HTML {
	return template.HTML(p.Html)
}
