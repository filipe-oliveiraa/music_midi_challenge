package conf

import (
	"testing"
	"time"

	"crossjoin.com/gorxestra/util/conf/typ"
	"github.com/stretchr/testify/assert"
)

type DefaultValues struct {
	Uint8       uint8        `conf:"default:1"`
	Uint16      uint16       `conf:"default:2"`
	Uint32      uint32       `conf:"default:3"`
	Uint64      uint64       `conf:"default:4"`
	Int8        uint8        `conf:"default:5"`
	Int16       uint16       `conf:"default:6"`
	Int32       uint32       `conf:"default:7"`
	Int64       uint64       `conf:"default:8"`
	Float32     float32      `conf:"default:9"`
	Float64     float32      `conf:"default:10"`
	String      string       `conf:"default:thisisastring"`
	Bool        bool         `conf:"default:true"`
	StringSlice []string     `conf:"default:first;second"`
	Duration    typ.Duration `conf:"default:10s"`
}

var testDuration = typ.Duration(time.Second * 10)

var ExpectedDefaultValue = DefaultValues{
	Uint8:       1,
	Uint16:      2,
	Uint32:      3,
	Uint64:      4,
	Int8:        5,
	Int16:       6,
	Int32:       7,
	Int64:       8,
	Float32:     9,
	Float64:     10,
	String:      "thisisastring",
	Bool:        true,
	StringSlice: []string{"first", "second"},
	Duration:    testDuration,
}

type TestGetFieldsStruct struct {
	TestUint32 uint32
	TestString string
	TestBool   bool
	TestSlice  []string
	TestStruct struct {
		StructField1 string
		StructField2 string
		AStruct      struct {
			SubStructField1 string
		}
	}
}

func TestSetDefaultValues(t *testing.T) {
	var v DefaultValues
	fields, err := getFields(&v)
	assert.Nil(t, err)
	err = setFieldDefaultValues(fields)
	assert.Nil(t, err)
	assert.Equal(t, ExpectedDefaultValue, v)
}

func TestGetFields(t *testing.T) {
	v := TestGetFieldsStruct{}

	fs, err := getFields(&v)
	assert.Nil(t, err)
	assert.Equal(t, 7, len(fs))
}

func TestParseTags(t *testing.T) {
	options := parseTag("default:10,required,env:ENV_UINT,hide,hello:world")

	assert.Equal(t, options, FieldOptions{
		Required:     true,
		DefaultValue: "10",
		EnvName:      "ENV_UINT",
		Hide:         true,
		Raw: map[string]string{
			RequiredTag:     "",
			DefaultValueTag: "10",
			EnvNameTag:      "ENV_UINT",
			HideTag:         "",
			"hello":         "world",
		},
	})
}

func TestGetValueByName(t *testing.T) {
	v := ExpectedDefaultValue

	value := getValueByName(&v, "Int8")
	assert.EqualValues(t, v.Int8, value.Uint())
}
