package builder

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s *Site) GetImage(name string) []byte {
	f, err := os.Open(filepath.Join(s.MediaDir, name))
	noerr("cannot open image", err)
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	noerr("cannot load image", err)
	return data
}

func (s *Site) GetImageData(name string) string {
	data := s.GetImage(name)
	ctype := http.DetectContentType(data)

	if strings.HasPrefix(ctype, "text/xml") && strings.HasSuffix(name, ".svg") {
		ctype = strings.ReplaceAll(ctype, "text/xml", "image/svg+xml")
	}

	if !strings.HasPrefix(ctype, "image/") {
		noerr("invalid image", fmt.Errorf("%s has type %s which is not image", name, ctype))
	}

	b64data := base64.StdEncoding.EncodeToString(data)

	return fmt.Sprintf("data:%s;base64,%s", ctype, b64data)
}
