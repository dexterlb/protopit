package builder

import (
	"html/template"
	"path/filepath"
	"reflect"

	"github.com/DexterLB/protopit/site/builder/media"
	"github.com/DexterLB/protopit/site/builder/translator"
	"github.com/mattn/go-zglob"
)

type Site struct {
	ContentDir   string
	Variant      string
	Pages        map[string]*Page
	PagesByDate  []*Page
	OutputDir    string
	Template     *template.Template
	CssTag       template.HTML
	StyleDir     string
	Media        *media.Media
	MediaDir     string
	MediaOutDir  string
	MediaUrlBase string
	AllVariants  map[string]*Site
	Translator   *translator.Translator
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
		"mcall": func(obj reflect.Value, method string, args ...reflect.Value) reflect.Value {
			return obj.MethodByName(method).Call(args)[0]
		},
		"set": func(v map[string]interface{}, key string, obj interface{}) interface{} {
			v[key] = obj
			return obj
		},
		"get": func(v map[string]interface{}, key string) interface{} {
			return v[key]
		},
		"make": func() interface{} {
			return make(map[string]interface{})
		},
	}

	templ := template.New("").Funcs(funcs)

	templ, err = templ.ParseFiles(templateNames...)
	noerr("cannot load templates", err)

	return &Site{
		ContentDir:   properContentDir,
		OutputDir:    "output",
		StyleDir:     "styles",
		MediaDir:     "media",
		MediaOutDir:  "output/media",
		MediaUrlBase: "/media",
		Media:        media.New("cache"),
		Variant:      variant,
		Template:     templ,
		Pages:        make(map[string]*Page),
		Translator:   translator,
	}
}

func (s *Site) PagesByType(t string) []*Page {
	var result []*Page
	for _, p := range s.PagesByDate {
		if p.Meta.Type == t {
			result = append(result, p)
		}
	}
	return result
}
