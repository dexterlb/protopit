package builder

import (
	"fmt"
	"os"
	"path/filepath"
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
	outDir := filepath.Join(s.OutputDir, p.Name, s.Variant)
	noerr("cannot create dir", os.MkdirAll(outDir, os.ModePerm))
	f, err := os.Create(filepath.Join(outDir, "index.html"))
	defer f.Close()
	noerr("cannot create output file", err)
	err = s.Template.ExecuteTemplate(f, fmt.Sprintf("%s.html", p.Meta.Type), TemplateData{
		Page: p,
		Site: s,
	})
	noerr("cannot render page", err)
}
