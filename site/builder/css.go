package builder

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/wellington/go-libsass"
)

var sassSite *Site

func (s *Site) RenderCss() {
	libsass.RegisterSassFunc("image($name)", sassImage)
	sassSite = s

	dummy := []byte(`@import 'main'`)
	w := &bytes.Buffer{}
	comp, err := libsass.New(w, bytes.NewReader(dummy))
	noerr("cannot set sass include dir", comp.Option(libsass.IncludePaths(
		[]string{s.StyleDir},
	)))

	noerr("cannot initialise style", err)
	noerr("cannot compile style", comp.Run())

	s.writeCss(w.Bytes())
}

func sassImage(ctx context.Context, usv libsass.SassValue) (*libsass.SassValue, error) {
    // TODO: if 2 arguments given, use the second as sizeSpec
	args := []interface{}{""}
	noerr("cannot process sass image() function", libsass.Unmarshal(usv, &args))

	res, err := libsass.Marshal(
		fmt.Sprintf("url('%s')", sassSite.GetImageData(args[0].(string), "")),
	)
	return &res, err
}

func (s *Site) writeCss(data []byte) {
	hashRaw := sha256.Sum256(data)
	hash := make([]byte, 128)
	base64.URLEncoding.Encode(hash, hashRaw[:])
	name := fmt.Sprintf("style.%s.css", string(hash[0:8]))

	f, err := os.Create(filepath.Join(s.OutputDir, name))
	noerr("cannot create output css file", err)
	defer f.Close()
	_, err = f.Write(data)
	noerr("cannot write css", err)

	s.CssTag = template.HTML(fmt.Sprintf(
		`<link rel="stylesheet" type="text/css" href="/%s">`, name,
	))
}
