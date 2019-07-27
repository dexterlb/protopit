package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
)

func (s *Site) globalFeed() *feeds.Feed {
	homepage, ok := s.PagesByUrl[""]
	if !ok {
		noerr("cannot find homepage", fmt.Errorf("no page with url ''"))
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       homepage.Meta.Title,
		Link:        &feeds.Link{Href: s.CanonicalUrl("/")},
		Description: homepage.Meta.Description,
		Created:     now,
	}

	for _, page := range s.Pages {
		if page.Meta.NoFeed {
			continue
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Title:       page.Meta.Title,
			Link:        &feeds.Link{Href: s.CanonicalUrl(page.AbsUrl())},
			Description: page.Meta.Description,
			Created:     page.Meta.Date,
		})
	}

	if len(feed.Items) == 0 {
		return nil
	}

	return feed
}

func (s *Site) RenderFeeds() {
	noerr("cannot create dir", os.MkdirAll(filepath.Join(s.OutputDir, s.feedDir()), os.ModePerm))

	globalFeed := s.globalFeed()
	if globalFeed == nil {
		s.HasFeeds = false
		return
	}
	s.HasFeeds = true

	atomFeed, err := globalFeed.ToAtom()
    noerr("cannot generate atom feed", err)

	rssFeed, err := globalFeed.ToRss()
    noerr("cannot generate rss feed", err)

	jsonFeed, err := globalFeed.ToJSON()
	noerr("cannot generate json feed", err)

	noerr("cannot create global atom feed", ioutil.WriteFile(
		filepath.Join(s.OutputDir, s.getFeed("atom")),
		[]byte(atomFeed),
		0777,
	))
	noerr("cannot create global rss feed", ioutil.WriteFile(
		filepath.Join(s.OutputDir, s.getFeed("rss")),
		[]byte(rssFeed),
		0777,
	))
	noerr("cannot create global json feed", ioutil.WriteFile(
		filepath.Join(s.OutputDir, s.getFeed("json")),
		[]byte(jsonFeed),
		0777,
	))
}

func (s *Site) feedDir() string {
	if s.Variant == "any" {
		return "feeds"
	} else {
		return filepath.Join("feeds", s.Variant)
	}
}

func (s *Site) FeedUrl(format string) string {
    return filepath.Join("/", s.getFeed(format))
}

func (s *Site) getFeed(format string) string {
	dir := s.feedDir()

	switch format {
	case "atom":
		return filepath.Join(dir, "atom.xml")
	case "rss":
		return filepath.Join(dir, "rss.xml")
	case "json":
		return filepath.Join(dir, "feed.json")
	}

	return ""
}
