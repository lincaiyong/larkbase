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

func extractAppTokenTableIdViewIdFromUrl(url string) (string, string, string) {
	re := regexp.MustCompile(`^https://bytedance\.larkoffice\.com/(?:base|wiki)/(\w+)\?table=(\w+)(?:&view=(\w+))?`)
	match := re.FindStringSubmatch(url)
	if len(match) != 4 {
		return "", "", ""
	}
	appToken := match[1]
	tableId := match[2]
	viewId := match[3]
	return appToken, tableId, viewId
}

func (c *Connection[T]) fieldTypeFromStructField(structField reflect.StructField) larkfield.Type {
	s := structField.Type.String()
	s = s[len("larkfield.") : len(s)-len("Field")]
	return larkfield.TypeFromString(s)
}

func (c *Connection[T]) fieldNameFromStructField(structField reflect.StructField) string {
	return structField.Tag.Get("lark")
}

func (c *Connection[T]) fieldFromStructField(structField reflect.StructField) (larkfield.Field, error) {
	ft := c.fieldTypeFromStructField(structField)
	if ft == larkfield.TypeUnknown {
		return nil, fmt.Errorf("field type of %s is not supported", structField.Type.String())
	}
	name := c.fieldNameFromStructField(structField)
	return ft.CreateField("", name, ft), nil
}

func (c *Connection[T]) extractAndFillConditionInstance(structPtr *T, tableUrl_ string) (tableUrl, appToken, structName, tableId, viewId string, fields []larkfield.Field, err error) {
	structValue := reflect.ValueOf(structPtr).Elem()
	structType := structValue.Type()
	metaField := structType.Field(0)
	structName = structType.Name()
	tableUrl = tableUrl_
	if tableUrl == "" {
		tableUrl = metaField.Tag.Get("lark")
	}
	appToken, tableId, viewId = extractAppTokenTableIdViewIdFromUrl(tableUrl)
	err = c.fillStructPtr(structPtr)
	for i := 1; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)
		if !fieldValue.CanAddr() || !fieldValue.Addr().CanInterface() {
			err = fmt.Errorf("%s is not exported", structType.Field(i).Name)
			return
		}
		var field larkfield.Field
		field, err = c.fieldFromStructField(structField)
		if err != nil {
			return
		}
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
		if !fieldValue.CanAddr() || !fieldValue.Addr().CanInterface() {
			return fmt.Errorf("%s is not exported", structType.Field(i).Name)
		}
		field := fieldValue.Addr().Interface().(larkfield.Field)
		field.SetSelf(field)
		field.SetName(c.fieldNameFromStructField(structField))
		field.SetType(c.fieldTypeFromStructField(structField))
		if c.isAnyRecord {
			break
		}
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

func (c *Connection[T]) checkStructType(structType reflect.Type, tableUrl string) error {
	if structType.NumField() == 0 {
		return fmt.Errorf("user struct has no field")
	}
	metaField := structType.Field(0)
	if tableUrl == "" {
		tableUrl = metaField.Tag.Get("lark")
	}
	if !metaField.Anonymous || metaField.Type.String() != "larkbase.Meta" || tableUrl == "" {
		return fmt.Errorf("first field of user struct should be larkbase.Meta with tableUrl in lark tag")
	}
	appToken, tableId, _ := extractAppTokenTableIdViewIdFromUrl(tableUrl)
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

func (c *Connection[T]) checkStructPtr(structPtr *T, tableUrl string) error {
	return c.checkStructType(reflect.TypeOf(structPtr).Elem(), tableUrl)
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
		if c.isAnyRecord {
			break
		}
		field := fieldValue.Addr().Interface().(larkfield.Field)
		record.Fields[field.Name()] = field
	}
	if c.isAnyRecord {
		if anyRecord, ok := any(structPtr).(*AnyRecord); ok {
			if anyRecord.update != nil {
				for k, v := range anyRecord.update {
					f := new(TextField)
					f.SetValue(v)
					record.Fields[k] = f
				}
			}
		}
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
	meta := Meta{RecordId: record.Id, ModifiedTime: record.ModifiedTime}

	if c.isAnyRecord {
		metaField := structValue.Field(0)
		metaField.Set(reflect.ValueOf(meta))
		dataField := structValue.Field(2)
		data := make(map[string]string)
		for k, f := range record.Fields {
			data[k] = f.StringValue()
		}
		dataField.Set(reflect.ValueOf(data))
		return nil
	}

	for i := 0; i < structValue.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)
		if structField.Name == "Meta" {
			fieldValue.Set(reflect.ValueOf(meta))
			continue
		}
		name := c.fieldNameFromStructField(structField)
		field, ok := record.Fields[name]
		if !ok {
			continue
		}
		newField, err := c.fieldFromStructField(structField)
		if err != nil {
			return err
		}
		value := field.UnderlayValue()
		newField.SetUnderlayValueNoDirty(value)
		fieldValue.Set(reflect.ValueOf(newField).Elem())
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
		field := fieldValue.Addr().Interface().(larkfield.Field)
		name := field.Name()
		value := field.StringValue()
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
