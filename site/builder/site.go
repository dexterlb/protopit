package builder

import (
	"html/template"
	"path/filepath"

	"github.com/mattn/go-zglob"
)

type Site struct {
	ContentDir  string
	Variant     string
	Pages       map[string]*Page
	PagesByDate []*Page
	OutputDir   string
	Template    *template.Template
	CssTag      template.HTML
	StyleDir    string
	MediaDir    string
	AllVariants map[string]*Site
}

func Init(variant string, contentDir string) *Site {
	properContentDir, err := filepath.Abs(contentDir)
	noerr("cannot get content dir path", err)

	templateNames, err := zglob.Glob("templates/**/*.html")
	noerr("cannot get template names", err)

	template, err := template.ParseFiles(templateNames...)
	noerr("cannot load templates", err)

	return &Site{
		ContentDir: properContentDir,
		OutputDir:  "output",
		StyleDir:   "styles",
		MediaDir:   "media",
		Variant:    variant,
		Template:   template,
		Pages:      make(map[string]*Page),
	}
}
