package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
)

type Source interface {
	Apply(v any, f []Field) error
	Help(f Field) string
	HasHelp() bool
}

type EnvSource struct {
	prefix string
}

func NewEnvSource(prefix string) Source {
	return &EnvSource{prefix: prefix}
}

// Apply implements Source.
func (e *EnvSource) Apply(v any, fields []Field) error {
	for i := range fields {
		env := fields[i].Options.Raw[EnvNameTag]
		if env == "" {
			env = joinParents(fields[i].ParentsName, e.prefix, fields[i].Name, "_", UpperCase)
		}

		v, ok := os.LookupEnv(env)
		if ok {
			_ = setFieldValue(fields[i], v)
		}
	}

	return nil
}

func (e *EnvSource) HasHelp() bool {
	return true
}

func (e *EnvSource) Help(f Field) string {
	env := joinParents(f.ParentsName, e.prefix, f.Name, "_", UpperCase)
	return fmt.Sprintf("$%s", env)
}

type FlagSource struct {
	prefix string
	args   []string
	flags  map[string]string
}

func NewFlagSource(prefix string, args []string) Source {
	return &FlagSource{
		prefix: prefix,
		args:   args,
		flags:  make(map[string]string),
	}
}

// Apply implements Source.
func (e *FlagSource) Apply(v any, fields []Field) error {
	// Parse flags
	for {
		ok, err := e.parseOne()
		if err != nil {
			return err
		}

		if !ok {
			break
		}
	}

	_, ok := e.flags["help"]
	if ok {
		return ErrHelp
	}

	for i := range fields {
		_ = e.applyOne(fields[i])
	}

	return nil
}

func (e *FlagSource) HasHelp() bool {
	return true
}

// based on go flag package
func (e *FlagSource) parseOne() (bool, error) {
	// no more arguments
	if len(e.args) == 0 {
		return false, nil
	}

	// arg does not start with -
	s := e.args[0]
	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}

	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			return false, nil
		}
	}

	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, fmt.Errorf("bad flag syntax: %s", s)
	}

	// it's a flag. does it have an argument?
	value := ""
	e.args = e.args[1:]
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			name = name[0:i]
			break
		}
	}

	e.flags[name] = value
	return true, nil
}

func (e *FlagSource) applyOne(field Field) error {
	flagName := field.Options.Raw[FlagTag]

	if flagName == "" {
		flagName = joinParents(field.ParentsName, e.prefix, field.Name, "-", LowerCase)
	}

	v, ok := e.flags[flagName]
	if ok {
		if v == "" {
			return setFieldValue(field, "true")
		}
		return setFieldValue(field, v)
	}

	return nil
}

func (e *FlagSource) Help(f Field) string {
	env := f.Options.Raw[FlagTag]
	if env == "" {
		env = joinParents(f.ParentsName, e.prefix, f.Name, "-", LowerCase)
	}

	return fmt.Sprintf("--%s", env)
}

const DefaultDataDirFieldName = "DataDir"

type Open = func(name string) (io.Reader, error)

type JsonSource struct {
	filename  string
	seekPath  bool
	fieldName string
	open      Open
}

// NewJsonSource attempts to seek path under the value (DataDir)
// if it finds DataDir value is used, otherwise uses path
func NewJsonSource(filename string, seekPath bool) Source {
	return &JsonSource{
		filename:  filename,
		seekPath:  seekPath,
		fieldName: DefaultDataDirFieldName,
		open:      openWrapper(os.Open),
	}
}

func NewJsonSourceWithOpen(filename string, seekPath bool, open Open) Source {
	return &JsonSource{
		filename:  filename,
		seekPath:  seekPath,
		fieldName: DefaultDataDirFieldName,
		open:      open,
	}
}

func openWrapper(f func(name string) (*os.File, error)) Open {
	return func(name string) (io.Reader, error) {
		f, err := f(name)
		return f, err
	}
}

// Apply implements Source.
func (j *JsonSource) Apply(v any, fields []Field) error {
	path := j.getConfPath(v)

	f, err := j.open(path)
	if err != nil {
		// Ignore open errors
		return nil
	}

	bs, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bs, v); err != nil {
		return err
	}

	return nil
}

func (j *JsonSource) getConfPath(v any) string {
	if !j.seekPath {
		return path.Join(".", j.fieldName)
	}

	vPath := getValueByName(v, j.fieldName)
	if vPath.Kind() == reflect.String {
		return path.Join(vPath.String(), j.filename)
	} else {
		return path.Join(".", j.filename)
	}
}

func (e *JsonSource) HasHelp() bool {
	return false
}

func (j *JsonSource) Help(f Field) string {
	return ""
}
