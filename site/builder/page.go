package builder

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

type Page struct {
	Filename string
	Meta     *MetaData
	Name     string
	Html     []byte
	Url		 string
}

func (p *Page) Content() template.HTML {
	return template.HTML(p.Html)
}

func (s *Site) Page(name string) *Page {
	site := s
	var ok bool
	var page *Page

	sname := strings.SplitN(name, ":", 2)
	if len(sname) >= 2 {
		site, ok = s.AllVariants[sname[1]]
		name = sname[0]
		if !ok {
			noerr("cannot decode page url", fmt.Errorf("no such variant: %s", sname[1]))
		}
	}

	if page, ok = site.Pages[name]; !ok {
		noerr("cannot get page url", fmt.Errorf("no such page: %s:%s", name, site.Variant))
	}

	return page
}

func (s *Site) PageUrl(name string) string {
	return s.Page(name).Url
}

func (s *Site) AbsPageUrl(name string) string {
	return filepath.Join("/", s.PageUrl(name))
}
