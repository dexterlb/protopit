package media

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Media struct {
	cacheDir string
}

func New(cacheDir string) *Media {
	return &Media{
		cacheDir: cacheDir,
	}
}

type Renderable interface {
	HashData(w io.Writer) error
	Extension() string
	MimeType() string
	Render() ([]byte, error)
}

func calcHash(r Renderable) (string, error) {
	b := &bytes.Buffer{}
	err := r.HashData(b)
	if err != nil {
		return "", err
	}
	return dataHash(b.Bytes()), nil
}

func dataHash(data []byte) string {
	hashRaw := sha256.Sum256(data)
	hash := make([]byte, 128)
	base64.RawURLEncoding.Encode(hash, hashRaw[:])
    return string(hash[0:24])
}

func (m *Media) Get(r Renderable) ([]byte, error) {
	hash, err := calcHash(r)
	if err != nil {
		return nil, err
	}

    err = os.MkdirAll(m.cacheDir, os.ModePerm)
    if err != nil {
        return nil, fmt.Errorf("cannot create cache dir: %s", err)
    }

	filename := filepath.Join(m.cacheDir, fmt.Sprintf("%s%s", hash, r.Extension()))

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		data, err := r.Render()
		if err != nil {
			return nil, fmt.Errorf("error while rendering media: %s", err)
		}

		err = ioutil.WriteFile(filename, data, 0666)
		if err != nil {
			return nil, fmt.Errorf("unable to write to cached file: %s", err)
		}

		return data, nil
	}
	return ioutil.ReadFile(filename)
}

func (m *Media) GetBase64(r Renderable) (string, error) {
	data, err := m.Get(r)
	if err != nil {
		return "", err
	}

	b64data := base64.StdEncoding.EncodeToString(data)

	return fmt.Sprintf("data:%s;base64,%s", r.MimeType(), b64data), nil
}

func (m *Media) GetFile(r Renderable, rootDir string, urlBase string) (string, error) {
	data, err := m.Get(r)
	if err != nil {
		return "", err
	}

    filename := fmt.Sprintf("%s%s", dataHash(data), r.Extension())
    outFile := filepath.Join(rootDir, filename)

    err = ioutil.WriteFile(outFile, data, 0666)
    if err != nil {
        return "", fmt.Errorf("unable to write file: %s", err)
    }

    return filepath.Join(urlBase, filename), nil
}
