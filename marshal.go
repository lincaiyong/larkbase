package larkbase

import (
	"encoding/json"
	"fmt"
	"github.com/lincaiyong/larkbase/field"
	"reflect"
)

func marshalSliceOfUserStruct(arg any) (string, error) {
	sliceValue := reflect.ValueOf(arg)
	s := make([]map[string]any, 0)
	for i := 0; i < sliceValue.Len(); i++ {
		elemValue := sliceValue.Index(i)
		var m map[string]any
		var err error
		if elemValue.Kind() == reflect.Struct {
			m, err = marshalUserStruct(elemValue)
		} else if elemValue.Kind() == reflect.Ptr {
			m, err = marshalUserStruct(elemValue.Elem())
		} else {
			return "", fmt.Errorf("invalid argument: %v", arg)
		}
		if err != nil {
			return "", err
		}
		s = append(s, m)
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b), nil
}

func marshalUserStruct(structValue reflect.Value) (map[string]any, error) {
	m := make(map[string]any)
	for j := 0; j < structValue.NumField(); j++ {
		fieldValue := structValue.Field(j)
		fieldType := fieldValue.Type()
		if fieldType.Name() == "Meta" {
			meta := fieldValue.Convert(reflect.TypeOf(Meta{})).Interface().(Meta)
			m["_record_id"] = meta.RecordId
			continue
		}
		baseFieldValue := fieldValue.Field(0)
		hack := baseFieldValue.Convert(reflect.TypeOf(field.HackBaseField{})).Interface().(field.HackBaseField)
		name := hack.Name()
		value := hack.Value()
		if value != nil {
			m[name] = value
		}
	}
	return m, nil
}

func Marshal(obj any) (string, error) {
	objType := reflect.TypeOf(obj)
	if objType.Kind() == reflect.Slice {
		return marshalSliceOfUserStruct(obj)
	}
	var m map[string]any
	var err error
	if objType.Kind() == reflect.Ptr {
		m, err = marshalUserStruct(reflect.ValueOf(obj).Elem())
	} else if objType.Kind() == reflect.Struct {
		m, err = marshalUserStruct(reflect.ValueOf(obj))
	} else {
		return "", fmt.Errorf("invalid argument: %v", obj)
	}
	if err != nil {
		return "", err
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b), nil
}
