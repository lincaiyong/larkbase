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

func extractTableAndFillBasicInfo(structPtr any) (tableUrl, appToken, tableId string, fields []Field, err error) {
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

func checkUserStructType(structType reflect.Type) error {
	// TODO: cache
	if structType.Kind() != reflect.Struct {
		return fmt.Errorf("structPtr is not a pointer to struct, got %s", structType.Kind().String())
	}
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

func checkUserStructPtr(structPtr any) error {
	structPtrValue := reflect.ValueOf(structPtr)
	if structPtrValue.Kind() != reflect.Ptr {
		return fmt.Errorf("structPtr is not a pointer to struct, got %T", structPtr)
	}
	return checkUserStructType(structPtrValue.Elem().Type())
}

func checkUserStructSlicePtr(structSlicePtr any) error {
	structSlicePtrType := reflect.TypeOf(structSlicePtr)
	if structSlicePtrType.Kind() != reflect.Ptr {
		return fmt.Errorf("structPtr is not a pointer to slice of struct, got %T", structSlicePtr)
	}
	structSliceType := structSlicePtrType.Elem()
	if structSliceType.Kind() != reflect.Slice {
		return fmt.Errorf("structPtr is not a pointer to slice of struct, got %T", structSlicePtr)
	}
	return checkUserStructType(structSliceType.Elem())
}

func convertUserStructToRecord(structPtr any) (record *Record, err error) {
	structPtrValue := reflect.ValueOf(structPtr)
	if structPtrValue.Kind() != reflect.Ptr {
		err = fmt.Errorf("structPtr is not a pointer to struct, got %T", structPtr)
		return
	}
	if structPtrValue.IsNil() {
		err = fmt.Errorf("structPtr is a nil pointer, expect address of a struct")
		return
	}
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

func convertRecordToUserStruct(record *Record, structPtr any) error {
	structPtrValue := reflect.ValueOf(structPtr)
	if structPtrValue.Kind() != reflect.Ptr {
		return fmt.Errorf("structPtr is not a pointer to struct, got %T", structPtr)
	}
	if structPtrValue.IsNil() {
		return fmt.Errorf("structPtr is a nil pointer, expect address of a struct")
	}
	structValue := structPtrValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return fmt.Errorf("structPtr is not a pointer to struct, got %T", structPtr)
	}
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

func convertRecordsToUserStructSlicePtr(records []*Record, structSlicePtr any) error {
	structSlicePtrValue := reflect.ValueOf(structSlicePtr)
	StructSliceValue := structSlicePtrValue.Elem()
	newSlice := reflect.MakeSlice(StructSliceValue.Type(), len(records), len(records))
	for i, record := range records {
		if err := convertRecordToUserStruct(record, newSlice.Index(i).Addr().Interface()); err != nil {
			return err
		}
	}
	StructSliceValue.Set(newSlice)
	return nil
}
