package larkbase

import (
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
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

func (c *Connection) QueryRecords() ([]*Record, error) {
	ret := make([]*Record, 0)
	pageToken := ""
	for {
		var err error
		ret, pageToken, err = c.queryRecordsByPage(pageToken, ret)
		if err != nil {
			return nil, err
		}
		if pageToken == "" {
			break
		}
	}
	return ret, nil
}
