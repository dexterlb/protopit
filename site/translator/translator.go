package translator

import (
	"encoding/json"
	"fmt"
	"os"
)

type Translator struct {
	texts map[string](map[string]string)
}

func New() *Translator {
	return &Translator{
		texts: make(map[string](map[string]string)),
	}
}

func Load(filename string) (*Translator, error) {
	tran := New()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return tran, nil
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %s", err)
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(&tran.texts)
	if err != nil {
		return nil, err
	}
	return tran, nil
}

func (t *Translator) Store(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("cannot open file: %s", err)
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	return enc.Encode(t.texts)
}

func (t *Translator) Get(source string, variant string) string {
	if _, ok := t.texts[source]; !ok {
		t.texts[source] = make(map[string]string)
	}
	if _, ok := t.texts[source][variant]; !ok {
		t.texts[source][variant] = source
	}
	return t.texts[source][variant]
}
