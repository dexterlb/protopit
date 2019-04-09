package builder

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/wellington/go-libsass"
)

func (s *Site) RenderCss() {
	dummy := []byte(`@import 'main'`)
	w := &bytes.Buffer{}
	comp, err := libsass.New(w, bytes.NewReader(dummy))
	noerr("cannot set sass include dir", comp.Option(libsass.IncludePaths(
		[]string{s.StyleDir},
	)))

	noerr("cannot initialise style", err)
	noerr("cannot compile style", comp.Run())

	s.CssTag = template.HTML(fmt.Sprintf("<style type=\"text/css\">\n%s\n</style>", w))
}
