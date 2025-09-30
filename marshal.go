package larkbase

import (
	"encoding/json"
	"errors"
	"reflect"
)

func Marshal(obj any) (string, error) {
	objType := reflect.TypeOf(obj)
	if objType.Kind() != reflect.Slice {
		return "", errors.New("not slice")
	}
	sliceValue := reflect.ValueOf(obj)
	s := make([]map[string]string, 0)
	for i := 0; i < sliceValue.Len(); i++ {
		elemValue := sliceValue.Index(i)
		if elemValue.Kind() != reflect.Struct {
			return "", errors.New("not struct")
		}
		m := make(map[string]string)
		for j := 0; j < elemValue.NumField(); j++ {
			fieldValue := elemValue.Field(j)
			fieldType := fieldValue.Type()
			if fieldType.Name() == "Meta" {
				continue
			}
			fieldPtr := fieldValue.Addr()
			ifield := fieldPtr.Interface().(IField)
			m[ifield.Name()] = ifield.Value()
		}
		s = append(s, m)
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b), nil
}
