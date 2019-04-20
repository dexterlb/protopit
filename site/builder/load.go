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
		if page.Meta.EventData != nil && !page.Meta.EventData.Next.IsZero() {
			s.UpcomingEvents = append(s.UpcomingEvents, page)
		}
	}
	sort.Slice(s.PagesByDate, func(i, j int) bool {
		return s.PagesByDate[i].Meta.Date.After(s.PagesByDate[j].Meta.Date)
	})
	sort.Slice(s.UpcomingEvents, func(i, j int) bool {
		return s.UpcomingEvents[i].Meta.EventData.Next.After(
			s.UpcomingEvents[j].Meta.EventData.Next,
		)
	})

	for _, page := range s.PagesByDate {
		if page.Meta.Type == "tag" {
			s.PagesByTag[page.Name] = nil
		}
	}
	for _, page := range s.PagesByDate {
		for _, tag := range page.Meta.Tags {
			if _, ok := s.PagesByTag[tag]; !ok {
				noerr("cannot build tag table", fmt.Errorf("invalid tag: %s", tag))
			}
			s.PagesByTag[tag] = append(s.PagesByTag[tag], page)
		}
	}
}
