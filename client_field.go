package larkbase

import (
	"context"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func (c *Client) Fields() (map[string]IField, error) {
	if err := c.checkCurrent(); err != nil {
		return nil, err
	}
	return c.current.table.fields, nil
}

func (c *Client) checkFields() {
	pageToken := ""
	for {
		pageToken = c.checkFieldsByPage(pageToken)
		if pageToken == "" {
			break
		}
	}
}

func (c *Client) checkFieldsByPage(pageToken string) string {
	const pageSize = 10
	req := larkbitable.NewListAppTableFieldReqBuilder().
		AppToken(c.current.appToken).
		TableId(c.current.table.id).
		PageSize(pageSize).
		PageToken(pageToken).
		Build()
	resp, err := c.Bitable.V1.AppTableField.List(context.Background(), req)
	if err != nil {
		c.failCurrent("fail to call bitable list field: %v", err)
		return ""
	}
	if !resp.Success() {
		c.failCurrent("get response with error: %s", larkcore.Prettify(resp.CodeError))
		return ""
	}
	for _, item := range resp.Data.Items {
		name := *item.FieldName
		expectField := c.current.table.fields[name]
		if expectField == nil {
			continue
		}
		expectType := expectField.Type()
		if expectType == 0 {
			continue
		}
		actualType := FieldType(*item.Type)
		if expectType != actualType {
			c.failCurrent("type of field \"%s\" mismatch, expect: %s(%d), actual: %s(%d)", name,
				expectType.String(), expectType, actualType.String(), actualType)
			return ""
		}
	}
	hasMore := *resp.Data.HasMore
	if hasMore {
		return *resp.Data.PageToken
	}
	return ""
}
