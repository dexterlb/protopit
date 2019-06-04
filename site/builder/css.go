package builder

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os/exec"

	"github.com/wellington/go-libsass"
)

var sassSite *Site

type cssRenderer struct {
	styleDir string
}

func (c *cssRenderer) Render() ([]byte, error) {
	dummy := []byte(`@import 'main'`)
	w := &bytes.Buffer{}
	comp, err := libsass.New(w, bytes.NewReader(dummy))
	noerr("cannot set sass include dir", comp.Option(libsass.IncludePaths(
		[]string{c.styleDir},
	)))

	if err != nil {
		return nil, fmt.Errorf("cannot initialise style: %s", err)
	}
	err = comp.Run()
	if err != nil {
		return nil, fmt.Errorf("cannot compile stype: %s", err)
	}

	return w.Bytes(), nil
}

func (c *cssRenderer) Extension() string {
	return ".css"
}

func (c *cssRenderer) MimeType() string {
	return "text/css"
}

func (c *cssRenderer) HashData(w io.Writer) error {
	cmd := exec.Command("find", c.styleDir, "-printf", "%p %s %t \\n")
	cmd.Stdout = w
    stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		edata := string(stderr.Bytes())
		if edata != "" {
			err = fmt.Errorf("command failed: %s", edata)
		}
	}
	if err != nil {
		return fmt.Errorf("cannot list files in %s: %s", c.styleDir, err)
	}

	return nil
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

func (s *Site) RenderCss() {
	libsass.RegisterSassFunc("image($name)", sassImage)
	sassSite = s

	url, err := s.Media.GetFile(
		&cssRenderer{styleDir: s.StyleDir},
		s.MediaOutDir, s.MediaUrlBase,
	)
	noerr("cannot render css", err)

	s.CssTag = template.HTML(fmt.Sprintf(
		`<link rel="stylesheet" type="text/css" href="%s">`, url,
	))
}
