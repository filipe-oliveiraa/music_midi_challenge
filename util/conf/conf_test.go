package conf

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct1 struct {
	Field1 uint `conf:"default:10,env:ENV_FIELD1,hide" json:"field1"`
	Field2 struct {
		InnerField1 uint     `conf:"default:11,env:ENV_INNER_FIELD1" json:"innerfield1"`
		InnerField2 []string `conf:"default:a;b;c,env:ENV_INNER_FIELD2" json:"innerfield2"`
	} `json:"field2"`
}

const JsonTestStruct1 = `{
    "field1": 10,
    "field2": {
        "innerfield1": 11,
        "innerfield2": [
            "a",
            "b",
            "c"
        ]
    }
}
`

func TestOutputTo(t *testing.T) {
	tmp := os.TempDir()
	path := path.Join(tmp, "f.json")

	var s TestStruct1
	err := OutputTo(&s, NewJsonOutputter(path))
	assert.ErrorIs(t, err, nil)

	equalFileContentWith(t, path, JsonTestStruct1)
}

func equalFileContentWith(t *testing.T, path, content string) {
	f, err := os.Open(path)
	assert.ErrorIs(t, err, nil)

	bs, err := io.ReadAll(f)
	assert.ErrorIs(t, err, nil)

	assert.Equal(t, content, string(bs))
}
