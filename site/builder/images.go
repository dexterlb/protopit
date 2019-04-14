package builder

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type imageRenderer struct {
	filename string
	sizeSpec string
}

func (i *imageRenderer) Extension() string {
	return filepath.Ext(i.filename)
}

func (i *imageRenderer) MimeType() string {
	data, err := ioutil.ReadFile(i.filename)
	noerr("cannot read image", err)

	ctype := http.DetectContentType(data)

	if strings.HasPrefix(ctype, "text/xml") && i.Extension() == ".svg" {
		ctype = strings.ReplaceAll(ctype, "text/xml", "image/svg+xml")
	}

	if !strings.HasPrefix(ctype, "image/") {
		noerr("invalid image", fmt.Errorf("%s has type %s which is not image", i.filename, ctype))
	}

	return ctype
}

func (i *imageRenderer) Render() ([]byte, error) {
	return ioutil.ReadFile(i.filename)
}

func (i *imageRenderer) HashData(w io.Writer) error {
	fmt.Fprintf(w, i.filename)
	fmt.Fprintf(w, i.sizeSpec)

	stat, err := os.Stat(i.filename)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%d %s", stat.Size(), stat.ModTime())
	return nil
}

func (s *Site) GetImage(name string, sizeSpec string) []byte {
	data, err := s.Media.Get(&imageRenderer{
		filename: filepath.Join(s.MediaDir, name),
		sizeSpec: "",
	})
	noerr("cannot render image", err)
	return data
}

func (s *Site) GetImageData(name string, sizeSpec string) string {
	data, err := s.Media.GetBase64(&imageRenderer{
		filename: filepath.Join(s.MediaDir, name),
		sizeSpec: "",
	})
	noerr("cannot render image", err)
	return data
}

func (s *Site) GetImageFile(name string, sizeSpec string) string {
	data, err := s.Media.GetFile(&imageRenderer{
		filename: filepath.Join(s.MediaDir, name),
		sizeSpec: "",
	}, s.MediaOutDir, s.MediaUrlBase)
	noerr("cannot render image", err)
	return data
}
