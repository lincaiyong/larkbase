package larkbase

import (
	"reflect"
	"regexp"
)

func (c *Client) Name() (string, error) {
	if err := c.checkCurrent(); err != nil {
		return "", err
	}
	return c.current.table.name, nil
}

func (c *Client) Id() (string, error) {
	if err := c.checkCurrent(); err != nil {
		return "", err
	}
	return c.current.table.id, nil
}

func (c *Client) Table(table any) *Client {
	ret := &Client{
		Client:    c.Client,
		appTables: c.appTables,
		current:   &ClientContext{},
	}
	ret = ret.realTable(table)

	return ret
}

func (c *Client) realTable(table any) *Client {
	obj := reflect.TypeOf(table)
	for obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}
	if obj.Kind() != reflect.Struct {
		c.failCurrent("invalid argument: %v", table)
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
		if field.Type.Kind() != reflect.Struct {
			c.failCurrent("invalid struct field: %s %s", name, type_)
			return c
		}

		newInstance := reflect.New(field.Type)
		typeMethod := newInstance.MethodByName("Type")
		if !typeMethod.IsValid() || typeMethod.Type().NumIn() != 0 {
			c.failCurrent("invalid struct field: %s %s", name, type_)
			return c
		}
		results := typeMethod.Call(nil)
		if len(results) == 0 {
			c.failCurrent("invalid struct field: %s %s", name, type_)
			return c
		}
		t := results[0].Int()
		fields[tag] = FieldType(t)
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

	c.checkFields()

	return c
}
