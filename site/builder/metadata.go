package builder

import (
	"fmt"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/teambition/rrule-go"
)

type MetaData struct {
	Title       string
	Description string
	Type        string
	Thumb       string
	Date        time.Time
	Url         *string
	Tags        []string
    Event       string
    NoFeed      bool
	EventData   *EventData
	location    *time.Location
}

func ParseMetaData(data []byte, loc *time.Location, variant string) *MetaData {
	var meta MetaData
	_, err := toml.Decode(string(data), &meta)
	noerr("cannot parse metadata", err)

	if meta.Type == "" {
		meta.Type = "normal"
	}
	if meta.Date.IsZero() && meta.Type == "normal" {
		noerr("cannot parse metadata", fmt.Errorf("no date!"))
	}
	meta.location = loc
	meta.EventData = parseEventData(meta.Event, loc, variant)
	return &meta
}

func (m *MetaData) HasTags() bool {
	return len(m.Tags) != 0
}

type EventData struct {
	Rule          *rrule.Set
	RuleRaw       string
	Next          time.Time
	Prev          time.Time
	HumanReadable string
	location      *time.Location
}

func (e *EventData) After(t time.Time) time.Time {
	return e.Rule.After(t, false).In(e.location)
}

func (e *EventData) Before(t time.Time) time.Time {
	return e.Rule.Before(t, false).In(e.location)
}

func parseEventData(s string, loc *time.Location, variant string) *EventData {
	if s == "" {
		return nil
	}
	ss := strings.SplitN(s, "#", 2)
	ed := &EventData{}
	ed.location = loc
	if len(ss) == 2 {
		ed.HumanReadable = ss[0]
		s = ss[1]
	}
	ed.RuleRaw = s

	ruleSpecs := strings.Split(s, "|")
	rule, err := rrule.StrSliceToRRuleSetInLoc(ruleSpecs, loc)
	if err != nil {
		var err2 error
		var t time.Time
		for _, format := range []string{
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
			"2006-01-02T15:04",
			"2006-01-02",
		} {
			t, err2 = time.ParseInLocation(format, s, time.Local)
			if err2 == nil {
				err = nil
				singleRule, err2 := rrule.NewRRule(rrule.ROption{
					Dtstart: t,
					Count:   1,
				})
				noerr("cannot construct rule from date", err2)
				rule = &rrule.Set{}
				rule.RRule(singleRule)
				break
			}
		}
	}
	noerr("cannot parse event date rule", err)
	ed.Rule = rule
	ed.Next = ed.After(time.Now())
	ed.Prev = ed.Before(time.Now())
	return ed
}
