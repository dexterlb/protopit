package builder

import "fmt"

func (s *Site) SanityCheck() {
	s.checkEvents()
}

func (s *Site) checkEvents() {
	for _, page := range s.Pages {
		if page.Meta.EventData != nil {
			var fine bool
			for _, tag := range page.Meta.Tags {
				if tag == "events" {
					fine = true
				}
			}
			if !fine {
				noerr("event inconsistency", fmt.Errorf("page %s has an event date but no events tag", page.Url))
			}
		}
	}
}
