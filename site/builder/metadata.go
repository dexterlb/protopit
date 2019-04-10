package builder

import "github.com/BurntSushi/toml"

type MetaData struct {
	Title string
	Type  string
}

func ParseMetaData(data []byte) *MetaData {
	var meta MetaData
	_, err := toml.Decode(string(data), &meta)
	noerr("cannot parse metadata", err)

	if meta.Type == "" {
		meta.Type = "normal"
	}
	return &meta
}
