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
	meta, fields := c.checkTableFields(table)
	if c.current.error != nil {
		return c
	}
	c.checkConnect(meta, fields)
	if c.current.error != nil {
		return c
	}
	c.checkFields()
	return c
}

func (c *Client) checkTableFields(table any) (string, map[string]IField) {
	obj := reflect.TypeOf(table)
	for obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}
	if obj.Kind() != reflect.Struct {
		c.failCurrent("invalid argument: %v", table)
	}
	meta := ""
	fields := make(map[string]IField)
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
			return "", nil
		}
		newInstance := reflect.New(field.Type)
		tableField := newInstance.Interface().(IField)
		tableField.SetName(tag)
		fields[tag] = tableField
	}
	if meta == "" {
		c.failCurrent("invalid table: missing Meta, %s", obj.Name())
		return "", nil
	}
	return meta, fields
}

func (c *Client) checkConnect(meta string, fields map[string]IField) *Client {
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
	if tableName == "" {
		c.failCurrent("fail to find table: %s", tableId)
		return c
	}
	c.current.table = NewTable(tableId, tableName, fields)
	return c
}
