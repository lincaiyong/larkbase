package larkbase

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"reflect"
	"regexp"
)

func NewClient(appId, appSecret string) *Client {
	return &Client{
		Client:    lark.NewClient(appId, appSecret),
		appTables: map[string]map[string]string{},
	}
}

type Client struct {
	*lark.Client
	appTables map[string]map[string]string // key is app token, table id
}

func (c *Client) Connect(ptr any) (*Connection, error) {
	meta, fields, err := c.checkTableStruct(ptr)
	if err != nil {
		return nil, err
	}

	appToken, table, err := c.connectLarkTable(meta, fields)
	if err != nil {
		return nil, err
	}

	pageToken := ""
	for {
		pageToken, err = c.checkFieldsByPage(appToken, table, pageToken)
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

func (c *Client) connectLarkApp(appToken string) error {
	if c.appTables[appToken] != nil {
		return nil
	}
	c.appTables[appToken] = map[string]string{}
	pageToken := ""
	for {
		var err error
		pageToken, err = c.queryAppTablesByPage(appToken, pageToken)
		if err != nil {
			return err
		}
		if pageToken == "" {
			break
		}
	}
	return nil
}

func (c *Client) queryAppTablesByPage(appToken, pageToken string) (string, error) {
	const pageSize = 10
	req := larkbitable.NewListAppTableReqBuilder().
		AppToken(appToken).
		PageSize(pageSize).
		PageToken(pageToken).
		Build()
	resp, err := c.Bitable.V1.AppTable.List(context.Background(), req)
	if err != nil {
		return "", err
	}
	if !resp.Success() {
		return "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	data := resp.Data
	tables := c.appTables[appToken]
	for _, item := range data.Items {
		tables[*item.TableId] = *item.Name
	}
	if *data.HasMore {
		return *data.PageToken, nil
	}
	return "", nil
}

func (c *Client) connectLarkTable(meta string, fields map[string]IField) (appToken string, table *Table, err error) {
	re := regexp.MustCompile(`^https://bytedance\.larkoffice\.com/base/(\w+)\?table=(\w+)`)
	match := re.FindStringSubmatch(meta)
	if len(match) != 3 {
		err = fmt.Errorf("invalid table url: %s", meta)
		return
	}
	appToken = match[1]
	tableId := match[2]
	err = c.connectLarkApp(appToken)
	if err != nil {
		return
	}
	tableName := c.appTables[appToken][tableId]
	if tableName == "" {
		err = fmt.Errorf("fail to find table: %s", tableId)
		return
	}
	table = NewTable(tableId, tableName, fields)
	return
}

func (c *Client) checkFieldsByPage(appToken string, table *Table, pageToken string) (string, error) {
	const pageSize = 10
	req := larkbitable.NewListAppTableFieldReqBuilder().
		AppToken(appToken).
		TableId(table.id).
		PageSize(pageSize).
		PageToken(pageToken).
		Build()
	resp, err := c.Bitable.V1.AppTableField.List(context.Background(), req)
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

func (c *Client) checkTableStruct(ptr any) (meta string, fields map[string]IField, err error) {
	objValue := reflect.ValueOf(ptr).Elem()
	obj := reflect.TypeOf(ptr)
	if obj.Kind() != reflect.Ptr {
		err = fmt.Errorf("invalid argument: expect pointer, %v", ptr)
		return
	}
	obj = obj.Elem()
	if obj.Kind() != reflect.Struct {
		err = fmt.Errorf("invalid argument: %v", ptr)
		return
	}
	fields = make(map[string]IField)
	for i := 0; i < obj.NumField(); i++ {
		fieldValue := objValue.Field(i)
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

		fieldValue.Set(newInstance.Elem())
	}
	if meta == "" {
		err = fmt.Errorf("invalid table: missing Meta, %s", obj.Name())
		return
	}
	return
}
