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
