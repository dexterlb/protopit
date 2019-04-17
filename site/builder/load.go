package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-zglob"
)

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

	for _, page := range s.Pages {
		s.PagesByDate = append(s.PagesByDate, page)
	}
	sort.Slice(s.PagesByDate, func(i, j int) bool {
		return s.PagesByDate[i].Meta.Date.After(s.PagesByDate[j].Meta.Date)
	})

}
