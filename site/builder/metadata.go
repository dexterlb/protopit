package builder

import (
	"fmt"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/DexterLB/protopit/site/builder/translator"
	"github.com/teambition/rrule-go"
)

type MetaData struct {
	Title     string
	Type      string
	Thumb     string
	Date      time.Time
	Url       *string
	Tags      []string
	Event     string
	EventData *EventData
	location  *time.Location
}

func ParseMetaData(data []byte, loc *time.Location, variant string, trans *translator.Translator) *MetaData {
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
	meta.EventData = parseEventData(meta.Event, loc, variant, trans)
	return &meta
}

func (m *MetaData) HasTags() bool {
	return len(m.Tags) != 0
}

type EventData struct {
	Rule          *rrule.Set
	HumanReadable string
}

func parseEventData(s string, loc *time.Location, variant string, trans *translator.Translator) *EventData {
	if s == "" {
		return nil
	}
	ss := strings.SplitN(s, "#", 2)
	ed := &EventData{}
	if len(ss) == 2 {
		ed.HumanReadable = trans.Get(ss[0], variant)
		s = ss[1]
	}

	ruleSpecs := strings.Split(s, "|")
	rule, err := rrule.StrSliceToRRuleSetInLoc(ruleSpecs, loc)
	noerr("cannot parse event date rule", err)
	ed.Rule = rule
	return ed
}
