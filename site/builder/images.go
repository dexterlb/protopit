package builder

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type imageRenderer struct {
	filename string
	sizeSpec string
}

func (i *imageRenderer) Extension() string {
	opts := parseImageOptions(i.sizeSpec)
	if opts != nil && opts.format != "" {
		return fmt.Sprintf(".%s", opts.format)
	}

	return filepath.Ext(i.filename)
}

func (i *imageRenderer) MimeType() string {
	if i.Extension() == ".svg" {
		return "image/svg+xml"
	}

	return fmt.Sprintf("image/%s", strings.TrimPrefix(i.Extension(), "."))
}

func (i *imageRenderer) Render() ([]byte, error) {
	data, err := ioutil.ReadFile(i.filename)
	if err != nil {
		return nil, err
	}

	opts := parseImageOptions(i.sizeSpec)
	if opts == nil {
		return data, nil
	}

	return convertImage(data, strings.TrimPrefix(filepath.Ext(i.filename), "."), opts)
}

func (i *imageRenderer) HashData(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", i.filename)
	fmt.Fprintf(w, "%s\n", i.sizeSpec)

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
		sizeSpec: sizeSpec,
	})
	noerr("cannot render image", err)
	return data
}

func (s *Site) GetImageFile(name string, sizeSpec string) string {
	url, err := s.Media.GetFile(&imageRenderer{
		filename: filepath.Join(s.MediaDir, name),
		sizeSpec: sizeSpec,
	}, s.MediaOutDir, s.MediaUrlBase)
	noerr("cannot render image", err)
	return url
}

type imageOptions struct {
	width  int
	height int
	method string
	format string
}

func parseImageOptions(s string) *imageOptions {
	if s == "" {
		return nil
	}
	m := regexp.MustCompile(
		`^([0-9]*)x([0-9]*)(\*([0-9.]*))?(:([a-z]+))?(\.([a-z]+))?`,
	).FindStringSubmatch(s)
	if len(m) != 9 {
		noerr("cannot parse image settings", fmt.Errorf("`%s` is invalid (gr %d", s, len(m)))
	}
	var opts imageOptions

	if m[1] != "" {
		width, err := strconv.ParseInt(m[1], 10, 64)
		noerr(fmt.Sprintf("`%s` is invalid image options", s), err)
		opts.width = int(width)
	}

	if m[2] != "" {
		height, err := strconv.ParseInt(m[2], 10, 64)
		noerr(fmt.Sprintf("`%s` is invalid image options", s), err)
		opts.height = int(height)
	}

	if m[4] != "" {
		scale, err := strconv.ParseFloat(m[4], 64)
		noerr(fmt.Sprintf("`%s` is invalid image options", s), err)
		opts.width = int(float64(opts.width) * scale)
		opts.height = int(float64(opts.height) * scale)
	}

	opts.method = m[6]
	opts.format = m[8]
	return &opts
}

func convertImage(data []byte, inFormat string, opts *imageOptions) ([]byte, error) {
	shell := &strings.Builder{}
	outFormat := inFormat
	if opts.format != "" {
		outFormat = opts.format
	}

	width := opts.width
	height := opts.height

	fmt.Fprintf(shell, `magick convert -background none -antialias %s:- `, inFormat)
	if opts.method == "" || opts.method == "fit" || opts.method == "scale" {
		if width == 0 {
			width = 99999
		}
		if height == 0 {
			height = 99999
		}
		if opts.method == "scale" {
			fmt.Fprintf(shell, ` -resize %dx%d `, width, height)
		} else {
			fmt.Fprintf(shell, ` -resize %dx%d\> `, width, height)
		}
	} else if opts.method == "crop" {
		if width == 0 || height == 0 {
			noerr("cannot convert image", fmt.Errorf("crop* size must have both width and height"))
		}
		fmt.Fprintf(shell, ` -resize %dx%d^ `, width, height)
		fmt.Fprintf(shell, ` -gravity center `)
		fmt.Fprintf(shell, ` -extent %dx%d`, width, height)
	} else {
		noerr("cannot convert image", fmt.Errorf("unknown method: %s", opts.method))
	}
	fmt.Fprintf(shell, ` %s:-`, outFormat)

	cmd := exec.Command("bash", "-c", shell.String())
	cmd.Stdin = bytes.NewReader(data)

	out := &bytes.Buffer{}
	cmd.Stdout = out

	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("cannot convert image: %s\nstderr:\n%s", err, string(stderr.Bytes()))
	}

	return out.Bytes(), nil
}
