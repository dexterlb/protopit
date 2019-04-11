package builder

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type MetaData struct {
	Title string
	Type  string
	Date  time.Time
	Url   *string
}

func ParseMetaData(data []byte) *MetaData {
	var meta MetaData
	_, err := toml.Decode(string(data), &meta)
	noerr("cannot parse metadata", err)

	if meta.Type == "" {
		meta.Type = "normal"
	}
	if meta.Date.IsZero() {
		noerr("cannot parse metadata", fmt.Errorf("no date!"))
	}
	return &meta
}
