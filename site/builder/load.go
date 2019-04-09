package builder

import (
	"fmt"
	"html/template"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mattn/go-zglob"
)

type Site struct {
	ContentDir string
	Variant    string
	Pages      map[string]*Page
	OutputDir  string
	Template   *template.Template
	CssTag     template.HTML
	StyleDir   string
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
		Variant:    variant,
		Template:   template,
		Pages:      make(map[string]*Page),
	}
}

func (s *Site) LoadPages() {
	names, err := zglob.Glob(fmt.Sprintf("%s/**/*.%s.md", s.ContentDir, s.Variant))
	noerr("cannot find pages", err)
	for i := range names {
		filename, err := filepath.Abs(names[i])
		noerr("cannot get file path", err)
		name := regexp.MustCompile(
			`^(.*)\.[^.]*\.[^.]*$`,
		).FindStringSubmatch(strings.TrimPrefix(filename, s.ContentDir))[1]

		s.ParsePage(filename, name)
	}
}
