package builder

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type MetaData struct {
	Title string
	Type  string
	Thumb string
	Date  time.Time
	Url   *string
	Tags []string
}

func ParseMetaData(data []byte) *MetaData {
	var meta MetaData
	_, err := toml.Decode(string(data), &meta)
	noerr("cannot parse metadata", err)

	if meta.Type == "" {
		meta.Type = "normal"
	}
	if meta.Date.IsZero() && meta.Type == "normal" {
		noerr("cannot parse metadata", fmt.Errorf("no date!"))
	}
	return &meta
}

func (m *MetaData) HasTags() bool {
    return len(m.Tags) != 0
}
