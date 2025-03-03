package reqvalidator

import (
	"errors"
	"reflect"
)

func Validate(target any) error {
	val := reflect.ValueOf(target)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.IsZero() {
			return errors.New("field '" + typ.Field(i).Name + "' is empty")
		}
	}

	return nil
}
