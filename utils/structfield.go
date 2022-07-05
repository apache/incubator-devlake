package utils

import (
	"reflect"
)

// WalkFiled get the field data by tag
func WalkFields(t reflect.Type, filter func(field *reflect.StructField) bool) (f []reflect.StructField) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			f = append(f, WalkFields(field.Type, filter)...)
		} else {
			if filter == nil {
				f = append(f, field)
			} else if filter(&field) {
				f = append(f, field)
			}
		}
	}
	return f
}
