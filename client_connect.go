package larkbase

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"reflect"
	"regexp"
)

func (c *Client) Connect(obj any) (*Connection, error) {
	meta, fields, err := checkTableStruct(obj)
	if err != nil {
		return nil, err
	}

	appToken, table, err := connectLarkApp(c, meta, fields)
	if err != nil {
		return nil, err
	}

	pageToken := ""
	for {
		pageToken, err = checkFieldsByPage(c, appToken, table, pageToken)
		if err != nil {
			return nil, err
		}
		if pageToken == "" {
			break
		}
	}
	ret := NewConnection(c, appToken, table)
	return ret, nil
}

func checkTableStruct(table any) (meta string, fields map[string]IField, err error) {
	obj := reflect.TypeOf(table)
	for obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}
	if obj.Kind() != reflect.Struct {
		err = fmt.Errorf("invalid argument: %v", table)
		return
	}
	fields = make(map[string]IField)
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		name := field.Name
		tag := field.Tag.Get("lark")
		type_ := field.Type.String()
		if name == "Meta" && type_ == "larkbase.Meta" {
			meta = tag
			continue
		}
		if field.Type.Kind() != reflect.Struct {
			err = fmt.Errorf("invalid struct field: %s %s", name, type_)
			return
		}
		newInstance := reflect.New(field.Type)
		tableField := newInstance.Interface().(IField)
		tableField.SetName(tag)
		fields[tag] = tableField
	}
	if meta == "" {
		err = fmt.Errorf("invalid table: missing Meta, %s", obj.Name())
		return
	}
	return
}

func connectLarkApp(client *Client, meta string, fields map[string]IField) (appToken string, table *Table, err error) {
	re := regexp.MustCompile(`^https://bytedance\.larkoffice\.com/base/(\w+)\?table=(\w+)`)
	match := re.FindStringSubmatch(meta)
	if len(match) != 3 {
		err = fmt.Errorf("invalid table url: %s", meta)
		return
	}
	appToken = match[1]
	tableId := match[2]
	err = client.connectLarkApp(appToken)
	if err != nil {
		return
	}
	tableName := client.appTables[appToken][tableId]
	if tableName == "" {
		err = fmt.Errorf("fail to find table: %s", tableId)
		return
	}
	table = NewTable(tableId, tableName, fields)
	return
}

func checkFieldsByPage(client *Client, appToken string, table *Table, pageToken string) (string, error) {
	const pageSize = 10
	req := larkbitable.NewListAppTableFieldReqBuilder().
		AppToken(appToken).
		TableId(table.id).
		PageSize(pageSize).
		PageToken(pageToken).
		Build()
	resp, err := client.Bitable.V1.AppTableField.List(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("fail to call bitable list field: %v", err)
	}
	if !resp.Success() {
		return "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	for _, item := range resp.Data.Items {
		name := *item.FieldName
		expectField := table.fields[name]
		if expectField == nil {
			continue
		}
		expectType := expectField.Type()
		if expectType == 0 {
			continue
		}
		actualType := FieldType(*item.Type)
		if expectType != actualType {
			return "", fmt.Errorf("type of field \"%s\" mismatch, expect: %s(%d), actual: %s(%d)", name,
				expectType.String(), expectType, actualType.String(), actualType)
		}
	}
	hasMore := *resp.Data.HasMore
	if hasMore {
		return *resp.Data.PageToken, nil
	}
	return "", nil
}
