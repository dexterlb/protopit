package builder

import "github.com/BurntSushi/toml"

type MetaData struct {
	Title string
}

func ParseMetaData(data []byte) *MetaData {
	var meta MetaData
	_, err := toml.Decode(string(data), &meta)
	noerr("cannot parse metadata", err)

	return &meta
}
