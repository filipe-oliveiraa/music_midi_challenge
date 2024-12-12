package errors

import "reflect"

func IsHashable(err error) bool {
	v := reflect.ValueOf(err)
	return isValueHashable(v)
}

func isValueHashable(v reflect.Value) bool {
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return true
		}

		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map, reflect.Slice:
		return false
	case reflect.Struct:
		for i := range v.NumField() {
			if !isValueHashable(v.Field(i)) {
				return false
			}
		}
	}

	return true
}
