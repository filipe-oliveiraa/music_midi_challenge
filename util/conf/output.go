package conf

import (
	"encoding/json"
	"os"
)

type Outputter interface {
	Output(v any) error
}

type JsonOutputter struct {
	path string
}

func NewJsonOutputter(path string) Outputter {
	return JsonOutputter{
		path: path,
	}
}

func (j JsonOutputter) Output(v any) error {
	f, err := os.Create(j.path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")

	return enc.Encode(v)
}
