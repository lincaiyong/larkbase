package larkbase

import (
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"reflect"
	"regexp"
	"strings"
)

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/bitable-overview
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/search

func Connect(appId, appSecret string, dataPtr any) (*Connection, error) {
	table, err := checkTable(dataPtr)
	if err != nil {
		return nil, err
	}
	client := lark.NewClient(appId, appSecret)
	return NewConnection(client, table), nil
}

func checkTable(dataPtr any) (*Table, error) {
	dataPtrType := reflect.TypeOf(dataPtr)
	if dataPtrType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("table is not a pointer: %v", dataPtr)
	}
	dataStructType := dataPtrType.Elem()
	if dataStructType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("table is not a struct: %v", dataPtr)
	}
	if dataStructType.NumField() == 0 {
		return nil, fmt.Errorf("table has no fields: %v", dataPtr)
	}

	dataStructValue := reflect.ValueOf(dataPtr).Elem()
	metaField := dataStructType.Field(0)
	tableUrl := metaField.Tag.Get("lark")
	if !metaField.Anonymous || metaField.Type.String() != "larkbase.Meta" || tableUrl == "" {
		return nil, fmt.Errorf("table has no larkbase.Meta field with tableUrl tag")
	}
	re := regexp.MustCompile(`^https://bytedance\.larkoffice\.com/base/(\w+)\?table=(\w+)`)
	match := re.FindStringSubmatch(tableUrl)
	if len(match) != 3 {
		return nil, fmt.Errorf("invalid table url: %s", tableUrl)
	}
	appToken := match[1]
	tableId := match[2]

	fields := make([]Field, 0)
	for i := 1; i < dataStructType.NumField(); i++ {
		fieldValue := dataStructValue.Field(i)
		structField := dataStructType.Field(i)
		tag := structField.Tag.Get("lark")
		name := structField.Type.Name()
		typeName := structField.Type.String()
		if strings.HasPrefix(typeName, "*larkbase.") {
			return nil, fmt.Errorf("table has invalid field type: %s, should be a struct, not porinter", typeName)
		}
		if !strings.HasPrefix(typeName, "larkbase.") || !strings.HasSuffix(typeName, "Field") || structField.Type.Kind() != reflect.Struct {
			return nil, fmt.Errorf("table has invalid field type: %s, expect type like larkbase.XxxField", typeName)
		}
		if !fieldValue.CanSet() {
			return nil, fmt.Errorf("table has unexported field: %s %s", name, typeName)
		}
		if tag == "" {
			return nil, fmt.Errorf("table has empty field tag: %s %s, expect `lark:\"xxx\"`", name, typeName)
		}
		f := fieldValue.Addr().Convert(reflect.TypeOf(&Field{})).Interface().(*Field)
		f.Name = tag
		f.Type = typeName[9 : len(typeName)-5]
		fields = append(fields, *f)
	}
	table := NewTable(tableUrl, appToken, tableId, fields)
	return table, nil
}
