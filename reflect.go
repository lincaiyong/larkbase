package larkbase

import (
	"fmt"
	"github.com/lincaiyong/larkbase/field"
	"reflect"
	"regexp"
	"strings"
)

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
	return s[len("field.") : len(s)-len("Field")] // a little bit hacking
}

func (c *Connection[T]) extractAndFillFilterInstance(structPtr *T) (tableUrl, appToken, tableId string, fields []Field, err error) {
	structValue := reflect.ValueOf(structPtr).Elem()
	metaField := structValue.Type().Field(0)
	tableUrl = metaField.Tag.Get("lark")
	appToken, tableId = extractAppTokenTableIdFromUrl(tableUrl)
	for i := 1; i < structValue.NumField(); i++ {
		structField := structValue.Type().Field(i)
		fieldValue := structValue.Field(i)
		f := reflect.New(structField.Type).Interface().(Field)
		f.SetName(structField.Tag.Get("lark"))
		f.SetType(convertToFieldType(structField.Type.String()))
		fields = append(fields, f)
		fieldValue.Set(reflect.ValueOf(f).Elem())
	}
	return
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
		if !strings.HasPrefix(typeName, "field.") || !strings.HasSuffix(typeName, "Field") || structField.Type.Kind() != reflect.Struct {
			return fmt.Errorf("field type of user struct should be field.XxxField, got %s", typeName)
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
	record.Fields = make(map[string]Field)
	for i := 0; i < structValue.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)
		if structField.Name == "Meta" {
			record.Id = fieldValue.Interface().(Meta).RecordId
			continue
		}
		baseField := fieldValue.Field(0)
		hack := baseField.Convert(reflect.TypeOf(field.HackBaseField{})).Interface().(field.HackBaseField)
		f := reflect.New(structField.Type).Interface().(Field)
		tag := structField.Tag.Get("lark")
		f.SetName(hack.Name())
		f.SetType(hack.Type())
		f.SetUnderlayValueNoDirty(hack.Value())
		f.SetDirty(hack.Dirty())
		record.Fields[tag] = f
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
			meta := Meta{RecordId: record.Id}
			fieldValue.Set(reflect.ValueOf(meta))
			continue
		}
		tag := structField.Tag.Get("lark")
		f, ok := record.Fields[tag]
		if !ok {
			continue
		}
		value := f.UnderlayValue()
		ff := reflect.New(structField.Type).Interface().(Field)
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
