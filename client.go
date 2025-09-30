package larkbase

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"reflect"
	"regexp"
	"strings"
	"unicode"
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

func (c *Client) Connect(arg any) (*Connection, error) {
	meta, fields, err := c.checkTableStruct(arg)
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

	structType := reflect.TypeOf(arg).Elem()
	structValue := reflect.ValueOf(arg).Elem()
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)
		if fieldType.Name == "Meta" {
			continue
		}
		f := fields[fieldType.Tag.Get("lark")]
		instance := reflect.New(fieldType.Type)
		instanceAsField := instance.Convert(reflect.TypeOf(&f))
		ff := instanceAsField.Interface().(*Field)
		ff.Name = f.Name
		ff.Type = f.Type
		fieldValue.Set(instance.Elem())
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

func (c *Client) connectLarkTable(meta string, fields map[string]Field) (appToken string, table *Table, err error) {
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
		expectField, ok := table.fields[name]
		if !ok {
			continue
		}
		expectType := expectField.Type
		if expectType == "" {
			continue
		}
		actualType := FieldType(*item.Type)
		if expectType != actualType.String() {
			return "", fmt.Errorf("type of field \"%s\" mismatch, expect: %s, actual: %s(%d)", name,
				expectType, actualType.String(), actualType)
		}
	}
	hasMore := *resp.Data.HasMore
	if hasMore {
		return *resp.Data.PageToken, nil
	}
	return "", nil
}

func (c *Client) checkIsPointerOfStruct(arg any) error {
	structPtrType := reflect.TypeOf(arg)
	if structPtrType.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid argument: expect a pointer, not %v", arg)
	}
	structType := structPtrType.Elem()
	if structType.Kind() != reflect.Struct {
		return fmt.Errorf("invalid argument: expect a pointer of struct, not %v", arg)
	}
	return nil
}

func (c *Client) checkStructField(name string, fieldType reflect.StructField, fieldValue reflect.Value) error {
	type_ := fieldType.Type.String()
	if !strings.HasPrefix(type_, "larkbase.") || !strings.HasSuffix(type_, "Field") {
		return fmt.Errorf("invalid struct field type: %s %s, expect larkbase.XxxField", name, type_)
	}
	if !unicode.IsUpper(rune(name[0])) {
		return fmt.Errorf("invalid struct field name: %s, should start with an uppercase letter", name)
	}
	if !fieldValue.CanSet() {
		return fmt.Errorf("invalid struct field: %s %s", name, type_)
	}
	if fieldType.Type.Kind() != reflect.Struct {
		return fmt.Errorf("invalid struct field type: %s %s, expect larkbase.XxxField", name, type_)
	}
	return nil
}

func (c *Client) checkTableStruct(arg any) (meta string, fields map[string]Field, err error) {
	if err = c.checkIsPointerOfStruct(arg); err != nil {
		return
	}
	structType := reflect.TypeOf(arg).Elem()
	structValue := reflect.ValueOf(arg).Elem()
	fields = make(map[string]Field)
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		name := fieldType.Name
		tag := fieldType.Tag.Get("lark")
		type_ := fieldType.Type.String()
		if name == "Meta" && type_ == "larkbase.Meta" {
			meta = tag
			continue
		}
		fieldValue := structValue.Field(i)
		if err = c.checkStructField(name, fieldType, fieldValue); err != nil {
			return
		}
		newFieldInstance := reflect.New(fieldType.Type).Elem()
		newFieldInstanceAsField := newFieldInstance.Convert(reflect.TypeOf(Field{}))
		f := newFieldInstanceAsField.Interface().(Field)
		f.Name = tag
		f.Type = type_[len("larkbase.") : len(type_)-len("Field")]
		fieldValue.Set(newFieldInstance)
		fields[tag] = f
	}
	if meta == "" {
		err = fmt.Errorf("invalid struct: missing Meta field, %s", structType.Name())
		return
	}
	return
}
