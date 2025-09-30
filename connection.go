package larkbase

import (
	"fmt"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"reflect"
)

func NewConnection(client *Client, appToken string, table *Table) *Connection {
	return &Connection{
		client:   client,
		appToken: appToken,
		table:    *table,
	}
}

type Connection struct {
	client *Client

	appToken string
	table    Table

	filters []*larkbitable.Condition
}

func (c *Connection) TableName() string {
	return c.table.name
}

func (c *Connection) TableId() string {
	return c.table.id
}

func (c *Connection) TableFields() map[string]IField {
	return c.table.fields
}

func (c *Connection) TableFieldKeys() []string {
	return c.table.fieldKeys
}

func (c *Connection) CountRecords() (int, error) {
	total := 0
	pageToken := ""
	for {
		var err error
		pageToken, err = c.queryTableRecordByPage(pageToken, &total)
		if err != nil {
			return 0, err
		}
		if pageToken == "" {
			break
		}
	}
	return total, nil
}

func (c *Connection) Where(filters ...*larkbitable.Condition) *Connection {
	c.filters = filters
	return c
}

func (c *Connection) QueryRecords(ptr any) error {
	ptrValue := reflect.ValueOf(ptr)
	if ptrValue.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid argument: expect a pointer, %v", ptr)
	}
	sliceValue := ptrValue.Elem()
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("invalid argument: expect a pointer of slice, %v", ptr)
	}
	sliceElemType := sliceValue.Type().Elem()
	if sliceElemType.Kind() != reflect.Struct {
		return fmt.Errorf("invalid argument: expect a pointer of slice of struct, %v", ptr)
	}

	records := make([]*Record, 0)
	pageToken := ""
	for {
		var err error
		records, pageToken, err = c.queryRecordsByPage(pageToken, records)
		if err != nil {
			return err
		}
		if pageToken == "" {
			break
		}
	}
	elemType := sliceValue.Type().Elem()
	newSlice := reflect.MakeSlice(sliceValue.Type(), len(records), len(records))
	for i, record := range records {
		elemValue := reflect.New(elemType).Elem()
		err := c.convertRecordToStruct(record, elemValue)
		if err != nil {
			return err
		}
		newSlice.Index(i).Set(elemValue)
	}
	sliceValue.Set(newSlice)

	return nil
}

func (c *Connection) convertRecordToStruct(record *Record, structValue reflect.Value) error {
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		if !fieldValue.CanSet() {
			continue
		}
		if field.Name == "Meta" && field.Type.String() == "larkbase.Meta" {
			continue
		}
		tag := field.Tag.Get("lark")
		if tag == "" {
			continue
		}
		var value string
		recordField := record.Fields[tag]
		if recordField != nil {
			value = recordField.Value()
		}
		if field.Type.Kind() != reflect.Struct {
			return fmt.Errorf("invalid struct field: %s %s", field.Name, field.Type.String())
		}
		newInstance := reflect.New(field.Type)
		tableField := newInstance.Interface().(IField)
		tableField.SetName(tag)
		tableField.SetValue(value)
		fieldValue.Set(newInstance.Elem())
	}
	return nil
}
