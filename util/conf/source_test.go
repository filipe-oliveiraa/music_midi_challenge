package conf

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestParseFlags struct {
	DataDir string
	Field1  string `conf:"flag:f1,env:F1" json:"field1"`
	Field2  uint   `conf:"flag:f2,env:F2" json:"field2"`
	Field3  bool   `conf:"flag:f3,env:F3" json:"field3"`
	Field4  string `conf:"" json:"field4"`
}

func TestFlagSource(t *testing.T) {
	{
		args := []string{"--f1=World", "--f2=5", "-f3", "--field4=test"}
		source := NewFlagSource("", args)
		expected := TestParseFlags{
			Field1: "World",
			Field2: 5,
			Field3: true,
			Field4: "test",
		}
		var actual TestParseFlags

		fields, err := getFields(&actual)
		assert.Nil(t, err)

		source.Apply(actual, fields)

		assert.Equal(t, expected, actual)
	}

	{
		args := []string{"--help"}
		source := NewFlagSource("", args)
		err := source.Apply(nil, nil)
		assert.ErrorIs(t, err, ErrHelp)
	}
}

func TestEnvSource(t *testing.T) {
	{
		os.Setenv("F1", "World")
		os.Setenv("F2", "5")
		os.Setenv("F3", "true")
		os.Setenv("ENV_FIELD4", "test")

		source := NewEnvSource("ENV")
		expected := TestParseFlags{
			Field1: "World",
			Field2: 5,
			Field3: true,
			Field4: "test",
		}
		var actual TestParseFlags
		fields, err := getFields(&actual)
		assert.Nil(t, err)
		err = source.Apply(actual, fields)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	}

}

const TestJsonFile = `
{
	"field1": "hello",
	"field2": 10,
	"field3": true,
	"field4":"world"
}
`

var TestJsonStruct = TestParseFlags{
	DataDir: "mockPath",
	Field1:  "hello",
	Field2:  10,
	Field3:  true,
	Field4:  "world",
}

func TestJsonSource(t *testing.T) {
	var s TestParseFlags
	s.DataDir = "mockPath"

	fields, err := getFields(&s)
	assert.Nil(t, err)
	json := NewJsonSourceWithOpen("", true, mockOpen)
	err = json.Apply(&s, fields)
	assert.Nil(t, err)
	assert.Equal(t, TestJsonStruct, s)
}

func mockOpen(name string) (io.Reader, error) {
	if name != "mockPath" {
		return nil, errors.New("file not found")
	}

	return strings.NewReader(TestJsonFile), nil
}
