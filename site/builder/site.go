package builder

import (
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/DexterLB/protopit/site/builder/media"
	"github.com/DexterLB/protopit/site/builder/translator"
	"github.com/mattn/go-zglob"
)

type Site struct {
	ContentDir     string
	Variant        string
	Pages          map[string]*Page
	PagesByUrl     map[string]*Page
	PagesByTag     map[string][]*Page
	PagesByDate    []*Page
	UpcomingEvents []*Page
	OutputDir      string
	Template       *template.Template
	CssTag         template.HTML
	HasFeeds       bool
	StyleDir       string
	Media          *media.Media
	MediaDir       string
	MediaOutDir    string
	MediaUrlBase   string
	BaseUrl        string
	AllVariants    map[string]*Site
	Translator     *translator.Translator
	Location       *time.Location
}

func Init(variant string, contentDir string, locationSpec string, translator *translator.Translator, url string) *Site {
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
		"mcall_url": func(obj reflect.Value, method string, args ...reflect.Value) template.URL {
			raw := obj.MethodByName(method).Call(args)[0].String()
			return template.URL(raw)
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
		"err": func(msg string, args ...interface{}) interface{} {
			noerr("cannot render template", fmt.Errorf(msg, args...))
			return 42
		},
		"never": func(t time.Time) bool {
			return t.IsZero()
		},
		"take": func(n int, items reflect.Value) reflect.Value {
			l := items.Len()
			if n < l {
				l = n
			}
			return items.Slice(0, l)
		},
		"tformat": func(option string, t time.Time) string {
			switch option {
			case "date":
				return t.Format("2006-01-02")
			case "time":
				return t.Format("15:04")
			case "datetime":
				return t.Format("2006-01-02 15:04")
			default:
				noerr("cannot format time", fmt.Errorf("unknown format: %s", option))
				return ""
			}
		},
		"strip": func(s string) string {
			return strings.TrimSpace(s)
		},
		"join": func(sep string, s ...string) string {
			return strings.Join(s, sep)
		},
	}

	templ := template.New("").Funcs(funcs)

	templ, err = templ.ParseFiles(templateNames...)
	noerr("cannot load templates", err)

	loc, err := time.LoadLocation(locationSpec)
	noerr("cannot determine time zone", err)

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
		PagesByTag:   make(map[string][]*Page),
		PagesByUrl:   make(map[string]*Page),
		Translator:   translator,
		Location:     loc,
		BaseUrl:      url,
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

func (s *Site) GetPagesByTag(tag string) []*Page {
	return s.PagesByTag[tag]
}

func (s *Site) CanonicalUrl(suffix string) string {
	if filepath.Join(suffix) == "/" {
		return s.BaseUrl
	}

	u, err := url.Parse(s.BaseUrl)
	noerr("invalid url", err)

	u.Path = filepath.Join(u.Path, suffix)
	return u.String()
}
