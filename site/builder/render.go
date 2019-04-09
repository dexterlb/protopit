package builder

import (
	"os"
	"path/filepath"
)

type TemplateData struct {
	Site *Site
	Page *Page
}

func (s *Site) Render() {
    noerr("cannot clean output dir", os.RemoveAll(s.OutputDir))
    for _, page := range s.Pages {
        s.RenderPage(page)
    }
}

func (s *Site) RenderPage(p *Page) {
	outDir := filepath.Join(s.OutputDir, p.Name, s.Variant)
	noerr("cannot create dir", os.MkdirAll(outDir, os.ModePerm))
	f, err := os.Create(filepath.Join(outDir, "index.html"))
	noerr("cannot create output file", err)
	err = s.Template.ExecuteTemplate(f, p.Type, TemplateData{
		Page: p,
		Site: s,
	})
	noerr("cannot render page", err)
}
