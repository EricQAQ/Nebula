package core

import (
	"reflect"
)

func GetFieldValue(stru interface{}, field string) reflect.Value {
	v := reflect.ValueOf(stru)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.FieldByName(field)
}
