package larkbase

import (
	"encoding/json"
	"errors"
	"github.com/lincaiyong/larkbase/field"
	"reflect"
)

func (c *Connection[T]) marshalStructPtr(structPtr *T) (map[string]any, error) {
	if structPtr == nil {
		return nil, errors.New("structPtr is nil")
	}
	m := make(map[string]any)
	structValue := reflect.ValueOf(structPtr).Elem()
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

func (c *Connection[T]) MarshalRecords(structPtrSlice []*T) (string, error) {
	s := make([]map[string]any, 0)
	for _, structPtr := range structPtrSlice {
		m, err := c.marshalStructPtr(structPtr)
		if err != nil {
			return "", err
		}
		s = append(s, m)
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b), nil
}

func (c *Connection[T]) MarshalRecord(structPtr *T) (string, error) {
	if structPtr == nil {
		return "", errors.New("structPtr is nil")
	}
	m, err := c.marshalStructPtr(structPtr)
	if err != nil {
		return "", err
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b), nil
}
