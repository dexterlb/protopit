package builder

import (
	"fmt"
	"html/template"
	"os"
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
	MediaDir   string
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

func (s *Site) Clean() {
	noerr("cannot clean output dir", os.RemoveAll(s.OutputDir))
	noerr("cannot create output dir", os.MkdirAll(s.OutputDir, os.ModePerm))
}

func (s *Site) LoadPages() {
	names, err := zglob.Glob(fmt.Sprintf("%s/**/*.%s.md", s.ContentDir, s.Variant))
	noerr("cannot find pages", err)
	for i := range names {
		filename, err := filepath.Abs(names[i])
		noerr("cannot get file path", err)
		name := regexp.MustCompile(
			`^(.*)\.[^.]*\.[^.]*$`,
		).FindStringSubmatch(strings.TrimPrefix(filename, fmt.Sprintf("%s/", s.ContentDir)))[1]

		s.ParsePage(filename, name)
	}
}
