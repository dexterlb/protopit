package builder

import "fmt"

func (s *Site) SanityCheck() {
	noerr("event inconsistency", s.checkEvents())
}

func (s *Site) checkEvents() error {
	for _, page := range s.Pages {
		if page.Meta.EventData != nil {
			var fine bool
			for _, tag := range page.Meta.Tags {
				if tag == "events" {
					fine = true
				}
			}
			if !fine {
				return fmt.Errorf("page %s has an event date but no events tag", page.Url)
			}

            for _, s2 := range s.AllVariants {
                if page2, ok := s2.Pages[page.Name]; ok {
                    if page2.Meta.EventData == nil {
                        return fmt.Errorf("page %s is not an event across all variants", page.Name)
                    }
                    if page2.Meta.EventData.RuleRaw != page.Meta.EventData.RuleRaw {
                        return fmt.Errorf("page %s has different event rules across variants", page.Name)
                    }
                }
            }
		}
	}
	return nil
}
