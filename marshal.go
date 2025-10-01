package larkbase

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func marshalSliceOfRecord(arg any) (string, error) {
	sliceValue := reflect.ValueOf(arg)
	s := make([]map[string]string, 0)
	for i := 0; i < sliceValue.Len(); i++ {
		elemValue := sliceValue.Index(i)
		var m map[string]string
		var err error
		if elemValue.Kind() == reflect.Struct {
			m, err = marshalRecord(elemValue)
		} else if elemValue.Kind() == reflect.Ptr {
			m, err = marshalRecord(elemValue.Elem())
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

func marshalRecord(recordValue reflect.Value) (map[string]string, error) {
	m := make(map[string]string)
	//for j := 0; j < recordValue.NumField(); j++ {
	//	fieldValue := recordValue.Field(j)
	//	fieldType := fieldValue.Type()
	//	if fieldType.Name() == "Meta" {
	//		meta := fieldValue.Convert(reflect.TypeOf(Meta{})).Interface().(Meta)
	//		m["_record_id"] = meta.RecordId
	//		continue
	//	}
	//	f := fieldValue.Convert(reflect.TypeOf(Field{})).Interface().(Field)
	//	n := f.Name()
	//	v := f.Value()
	//	if v != "" {
	//		m[n] = v
	//	}
	//}
	return m, nil
}

func Marshal(obj any) (string, error) {
	objType := reflect.TypeOf(obj)
	if objType.Kind() == reflect.Slice {
		return marshalSliceOfRecord(obj)
	}
	var m map[string]string
	var err error
	if objType.Kind() == reflect.Ptr {
		m, err = marshalRecord(reflect.ValueOf(obj).Elem())
	} else if objType.Kind() == reflect.Struct {
		m, err = marshalRecord(reflect.ValueOf(obj))
	} else {
		return "", fmt.Errorf("invalid argument: %v", obj)
	}
	if err != nil {
		return "", err
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b), nil
}
