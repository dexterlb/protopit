package builder

import (
	"html/template"
	"path/filepath"

	"github.com/DexterLB/protopit/site/builder/media"
	"github.com/DexterLB/protopit/site/builder/translator"
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
	Media       *media.Media
	MediaDir    string
	AllVariants map[string]*Site
	Translator  *translator.Translator
}

func Init(variant string, contentDir string, translator *translator.Translator) *Site {
	properContentDir, err := filepath.Abs(contentDir)
	noerr("cannot get content dir path", err)

	templateNames, err := zglob.Glob("templates/**/*.html")
	noerr("cannot get template names", err)

	funcs := template.FuncMap{
		"translate": func(text string) string {
			return translator.Get(text, variant)
		},
	}

	templ := template.New("").Funcs(funcs)

	templ, err = templ.ParseFiles(templateNames...)
	noerr("cannot load templates", err)

	return &Site{
		ContentDir: properContentDir,
		OutputDir:  "output",
		StyleDir:   "styles",
		MediaDir:   "media",
		Media:      media.New("cache"),
		Variant:    variant,
		Template:   templ,
		Pages:      make(map[string]*Page),
		Translator: translator,
	}
}
