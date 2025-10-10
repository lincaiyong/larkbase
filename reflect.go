package larkbase

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lincaiyong/larkbase/larkfield"
	"reflect"
	"regexp"
	"strings"
)

func (c *Connection[T]) convertToHack(fieldValue reflect.Value) larkfield.HackBaseField {
	baseFieldPtr := fieldValue.Field(0)
	return baseFieldPtr.Convert(reflect.TypeOf(larkfield.HackBaseField{})).Interface().(larkfield.HackBaseField)
}

func extractAppTokenTableIdFromUrl(url string) (string, string) {
	re := regexp.MustCompile(`^https://bytedance\.larkoffice\.com/base/(\w+)\?table=(\w+)`)
	match := re.FindStringSubmatch(url)
	if len(match) != 3 {
		return "", ""
	}
	appToken := match[1]
	tableId := match[2]
	return appToken, tableId
}

func convertToFieldType(s string) string {
	return s[len("larkfield.") : len(s)-len("Field")] // a little bit hacking
}

func (c *Connection[T]) extractAndFillConditionInstance(structPtr *T) (tableUrl, appToken, structName, tableId string, fields []larkfield.Field, err error) {
	structValue := reflect.ValueOf(structPtr).Elem()
	structType := structValue.Type()
	metaField := structType.Field(0)
	structName = structType.Name()
	tableUrl = metaField.Tag.Get("lark")
	appToken, tableId = extractAppTokenTableIdFromUrl(tableUrl)
	err = c.fillStructPtr(structPtr)
	for i := 1; i < structType.NumField(); i++ {
		fieldValue := structValue.Field(i)
		field := fieldValue.Addr().Interface().(larkfield.Field)
		fields = append(fields, field)
	}
	return
}

func (c *Connection[T]) fillStructPtr(structPtr *T) error {
	structValue := reflect.ValueOf(structPtr).Elem()
	structType := structValue.Type()
	for i := 1; i < structValue.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)
		field := fieldValue.Addr().Interface().(larkfield.Field)
		field.SetName(structField.Tag.Get("lark"))
		field.SetType(convertToFieldType(structField.Type.String()))
	}
	return nil
}

func (c *Connection[T]) fillStructPtrSlice(structPtrSlicePtr []*T) error {
	for _, structPtr := range structPtrSlicePtr {
		err := c.fillStructPtr(structPtr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Connection[T]) checkStructType(structType reflect.Type) error {
	if structType.NumField() == 0 {
		return fmt.Errorf("user struct has no field")
	}
	metaField := structType.Field(0)
	tableUrl := metaField.Tag.Get("lark")
	if !metaField.Anonymous || metaField.Type.String() != "larkbase.Meta" || tableUrl == "" {
		return fmt.Errorf("first field of user struct should be larkbase.Meta with tableUrl in lark tag")
	}
	appToken, tableId := extractAppTokenTableIdFromUrl(tableUrl)
	if appToken == "" || tableId == "" {
		return fmt.Errorf("tableUrl is invalid: %s", tableUrl)
	}
	for i := 1; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		typeName := structField.Type.String()
		if !strings.HasPrefix(typeName, "larkfield.") || !strings.HasSuffix(typeName, "Field") || structField.Type.Kind() != reflect.Struct {
			return fmt.Errorf("field type of user struct should be larkfield.XxxField, got %s", typeName)
		}
		tag := structField.Tag.Get("lark")
		if tag == "" {
			return fmt.Errorf("field type of user struct should have lark tag like`lark:\"xxx\"`: %s %s", structField.Name, typeName)
		}
	}
	return nil
}

func (c *Connection[T]) checkStructPtr(structPtr *T) error {
	return c.checkStructType(reflect.TypeOf(structPtr).Elem())
}

func (c *Connection[T]) checkStructPtrSlicePtr(structPtrSlicePtr *[]*T) error {
	return c.checkStructType(reflect.TypeOf(structPtrSlicePtr).Elem().Elem().Elem())
}

func (c *Connection[T]) checkStructPtrSlice(structPtrSlice []*T) error {
	structSliceType := reflect.TypeOf(structPtrSlice)
	return c.checkStructType(structSliceType.Elem().Elem())
}

func (c *Connection[T]) convertStructPtrToRecord(structPtr *T) (record *Record, err error) {
	structPtrValue := reflect.ValueOf(structPtr)
	structValue := structPtrValue.Elem()
	if structValue.Kind() != reflect.Struct {
		err = fmt.Errorf("structPtr is not a pointer to struct, got %T", structPtr)
		return
	}
	structType := structValue.Type()
	record = NewRecord()
	record.Fields = make(map[string]larkfield.Field)
	for i := 0; i < structValue.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)
		if structField.Name == "Meta" {
			meta := fieldValue.Interface().(Meta)
			record.Id = meta.RecordId
			record.ModifiedTime = meta.ModifiedTime
			continue
		}
		hack := c.convertToHack(fieldValue)
		field := reflect.New(structField.Type).Interface().(larkfield.Field)
		tag := structField.Tag.Get("lark")
		field.SetName(hack.Name())
		field.SetType(hack.Type())
		field.SetUnderlayValueNoDirty(hack.Value())
		field.SetDirty(hack.Dirty())
		record.Fields[tag] = field
	}
	return
}

func (c *Connection[T]) convertStructPtrSliceToRecords(structPtrSlice []*T) (records []*Record, err error) {
	for _, structPtr := range structPtrSlice {
		var record *Record
		record, err = c.convertStructPtrToRecord(structPtr)
		if err != nil {
			return
		}
		records = append(records, record)
	}
	return
}

func (c *Connection[T]) convertRecordToStructPtr(record *Record, structPtr *T) error {
	structValue := reflect.ValueOf(structPtr).Elem()
	structType := structValue.Type()
	for i := 0; i < structValue.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)
		if structField.Name == "Meta" {
			meta := Meta{RecordId: record.Id, ModifiedTime: record.ModifiedTime}
			fieldValue.Set(reflect.ValueOf(meta))
			continue
		}
		tag := structField.Tag.Get("lark")
		field, ok := record.Fields[tag]
		if !ok {
			continue
		}
		value := field.UnderlayValue()
		ff := reflect.New(structField.Type).Interface().(larkfield.Field)
		ff.SetName(tag)
		ff.SetType(convertToFieldType(structField.Type.String()))
		ff.SetUnderlayValueNoDirty(value)
		fieldValue.Set(reflect.ValueOf(ff).Elem())
	}
	return nil
}

func (c *Connection[T]) convertRecordsToStructPtrSlicePtr(records []*Record, structPtrSlicePtr *[]*T) error {
	ret := make([]*T, len(records))
	for i, record := range records {
		structPtr := new(T)
		if err := c.convertRecordToStructPtr(record, structPtr); err != nil {
			return err
		}
		ret[i] = structPtr
	}
	*structPtrSlicePtr = ret
	return nil
}

func (c *Connection[T]) marshalStructPtr(structPtr *T) (map[string]string, error) {
	if structPtr == nil {
		return nil, errors.New("structPtr is nil")
	}
	m := make(map[string]string)
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
		hack := baseFieldValue.Convert(reflect.TypeOf(larkfield.HackBaseField{})).Interface().(larkfield.HackBaseField)
		name := hack.Name()
		value := hack.StringValue()
		if value != "" {
			m[name] = value
		}
	}
	return m, nil
}

func (c *Connection[T]) MarshalRecords(structPtrSlice []*T) (string, error) {
	s := make([]map[string]string, 0)
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

func (c *Connection[T]) MarshalIgnoreError(x any) string {
	if record, ok := x.(*T); ok {
		s, _ := c.MarshalRecord(record)
		return s
	} else if records, ok2 := x.([]*T); ok2 {
		s, _ := c.MarshalRecords(records)
		return s
	}
	return ""
}
