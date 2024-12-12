package conf

import (
	"reflect"
	"strconv"
	"strings"

	"crossjoin.com/gorxestra/util/conf/typ"
)

const (
	TagSeparator      = ","
	SliceSeparator    = ";"
	KeyValueSeparator = ":"
)

const (
	ConfTag         = "conf"
	FlagTag         = "flag"
	RequiredTag     = "required"
	EnvNameTag      = "env"
	DefaultValueTag = "default"
	HideTag         = "hide"
)

type Field struct {
	Name string

	Value reflect.Value

	ParentsName []string

	Options FieldOptions
}

type FieldOptions struct {
	Required bool

	EnvName      string
	DefaultValue string

	Hide bool

	Raw map[string]string
}

func getFields(v any) ([]Field, error) {
	fields, err := getFieldsRec(v, []string{})
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func getFieldsRec(v any, parents []string) ([]Field, error) {
	fields := make([]Field, 0)
	value := getValueElm(reflect.ValueOf(v))

	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			sfield := value.Type().Field(i)

			if !field.CanSet() {
				return fields, nil
			}

			if field.Kind() == reflect.Struct {
				structFields, err := getFieldsRec(field.Addr().Interface(), append(parents, sfield.Name))
				if err != nil {
					return nil, err
				}

				fields = append(fields, structFields...)
				continue
			}

			tag := sfield.Tag.Get(ConfTag)
			options := parseTag(tag)
			fields = append(fields, Field{
				Name:        sfield.Name,
				Value:       field,
				ParentsName: parents,
				Options:     options,
			})
		}
	}

	return fields, nil
}

func getValueByName(v any, name string) reflect.Value {
	value := getValueElm(reflect.ValueOf(v))
	return value.FieldByName(name)
}

func getValueElm(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Interface && !value.IsNil() {
		elm := value.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			return elm
		}
	}

	if value.Kind() == reflect.Pointer {
		return value.Elem()
	}

	return value
}

// parseTag returns a FieldOption when parsing the tag
func parseTag(tag string) FieldOptions {
	raw := make(map[string]string)

	parts := strings.Split(tag, TagSeparator)

	for i := range parts {
		kv := strings.SplitN(parts[i], KeyValueSeparator, 2)

		if len(kv) > 1 {
			raw[kv[0]] = kv[1]
		} else {
			raw[kv[0]] = ""
		}
	}

	_, required := raw[RequiredTag]
	_, hide := raw[HideTag]
	return FieldOptions{
		Required:     required,
		EnvName:      raw[EnvNameTag],
		DefaultValue: raw[DefaultValueTag],
		Hide:         hide,
		Raw:          raw,
	}
}

func setFieldDefaultValues(fields []Field) error {
	for i := range fields {
		err := setFieldValue(fields[i], fields[i].Options.DefaultValue)
		if err != nil {
			return err
		}
	}

	return nil
}

func setFieldValue(field Field, strValue string) error {
	value := field.Value
	kind := value.Kind()

	// Verify if is a custom type
	v := value.Interface()
	custom, ok := v.(typ.Type)
	if ok {
		customV, err := custom.GetValue(strValue)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(customV))
		return nil
	}

	switch kind {
	case reflect.Bool:
		b, err := strconv.ParseBool(strValue)
		if err != nil {
			return err
		}
		value.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return err
		}
		value.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(strValue, 10, 64)
		if err != nil {
			return err
		}
		value.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return err
		}
		value.SetFloat(v)
	case reflect.String:
		value.SetString(strValue)
	case reflect.Slice:
		if err := setSliceField(field, strValue); err != nil {
			return err
		}
	}

	return nil
}

func setSliceField(field Field, strValue string) error {
	value := field.Value
	kind := value.Type().String()

	if kind == "[]string" {
		v := parseStringSlice(strValue)
		value.Set(reflect.ValueOf(v))
	}

	return nil
}

func parseStringSlice(v string) []string {
	return strings.Split(v, SliceSeparator)
}
