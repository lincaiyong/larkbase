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

	current *ClientContext
}

type ClientContext struct {
	appToken string
	table    *Table
	error    error
}

func (c *Client) checkCurrent() error {
	if c.current == nil {
		return fmt.Errorf("current table is not set")
	}
	return c.current.error
}

func (c *Client) failCurrent(msg string, args ...any) {
	c.current.error = fmt.Errorf(msg, args...)
}

func (c *Client) Error() error {
	if c.current != nil {
		return c.current.error
	}
	return nil
}

func (c *Client) Table(table any) *Client {
	ret := &Client{
		Client:    c.Client,
		appTables: c.appTables,
		current:   &ClientContext{},
	}
	return ret.realTable(table)
}

func (c *Client) realTable(table any) *Client {
	obj := reflect.TypeOf(table)
	for obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}
	if obj.Kind() != reflect.Struct {
		c.failCurrent("invalid argument to Table(): %v", table)
	}
	meta := ""
	fields := make(map[string]FieldType)
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		name := field.Name
		tag := field.Tag.Get("lark")
		type_ := field.Type.String()
		if name == "Meta" && type_ == "larkbase.Meta" {
			meta = tag
			continue
		}
		var t FieldType
		switch type_ {
		case "larkbase.TextField":
			t = FieldTypeText
		case "larkbase.NumberField":
			t = FieldTypeNumber
		case "larkbase.SingleSelectField":
			t = FieldTypeSingleSelect
		case "larkbase.MultiSelectField":
			t = FieldTypeMultiSelect
		case "larkbase.DateField":
			t = FieldTypeDate
		case "larkbase.CheckboxField":
			t = FieldTypeCheckbox
		case "larkbase.PersonField":
			t = FieldTypePerson
		case "larkbase.UrlField":
			t = FieldTypeUrl
		case "larkbase.MediaField":
			t = FieldTypeMedia
		case "larkbase.UpdatedTimeField":
			t = FieldTypeUpdatedTime
		default:
			c.failCurrent("invalid table field: %s %s `lark:\"%s\"`", name, type_, tag)
			return c
		}
		fields[tag] = t
	}
	if meta == "" {
		c.failCurrent("invalid table: missing Meta, %s", obj.Name())
		return c
	}
	re := regexp.MustCompile(`^https://bytedance\.larkoffice\.com/base/(\w+)\?table=(\w+)`)
	match := re.FindStringSubmatch(meta)
	if len(match) != 3 {
		c.failCurrent("invalid table url: %s", meta)
		return c
	}
	c.current.appToken = match[1]
	tableId := match[2]
	err := c.connectApp(c.current.appToken)
	if err != nil {
		c.current.error = err
		return c
	}
	tableName := c.appTables[c.current.appToken][tableId]
	if table == "" {
		c.failCurrent("fail to find table: %s", tableId)
		return c
	}
	c.current.table = NewTable(tableId, tableName, fields)
	return c
}

func (c *Client) Count(total *int) *Client {
	if c.checkCurrent() != nil {
		return c
	}
	pageToken := ""
	for {
		var err error
		pageToken, err = c.getTableRecordByPage(pageToken, total)
		if err != nil {
			c.current.error = err
			return c
		}
		if pageToken == "" {
			break
		}
	}
	return c
}

func (c *Client) connectApp(appToken string) error {
	if c.appTables[appToken] != nil {
		return nil
	}
	c.appTables[appToken] = map[string]string{}
	pageToken := ""
	for {
		var err error
		pageToken, err = c.getAppTablesByPage(appToken, pageToken)
		if err != nil {
			return err
		}
		if pageToken == "" {
			break
		}
	}
	return nil
}

func (c *Client) getAppTablesByPage(appToken, pageToken string) (string, error) {
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

func (c *Client) getTableRecordByPage(pageToken string, total *int) (string, error) {
	const pageSize = 100
	req := larkbitable.NewSearchAppTableRecordReqBuilder().
		AppToken(c.current.appToken).
		TableId(c.current.table.id).
		PageToken(pageToken).
		PageSize(pageSize).
		Body(larkbitable.NewSearchAppTableRecordReqBodyBuilder().Build()).Build()
	resp, err := c.Bitable.V1.AppTableRecord.Search(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("fail to call bitable search table: %v", err)
	}
	if !resp.Success() {
		return "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	*total += *resp.Data.Total
	if *resp.Data.HasMore {
		return *resp.Data.PageToken, nil
	}
	return "", nil
}
