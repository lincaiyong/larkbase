package larkbase

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func (c *Connection) CountRecords() (int, error) {
	total := 0
	pageToken := ""
	for {
		var err error
		pageToken, err = c.getTableRecordByPage(pageToken, &total)
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

func (c *Connection) queryRecordsByPage(pageToken string, records []*Record) ([]*Record, string, error) {
	const pageSize = 100
	bodyBuilder := larkbitable.NewSearchAppTableRecordReqBodyBuilder()
	bodyBuilder.FieldNames(c.table.fieldKeys)
	if len(c.filters) > 0 {
		bodyBuilder.Filter(larkbitable.NewFilterInfoBuilder().
			Conjunction(`and`).
			Conditions(c.filters).
			Build())
	}
	bodyBuilder.AutomaticFields(true)
	req := larkbitable.NewSearchAppTableRecordReqBuilder().
		AppToken(c.appToken).
		TableId(c.table.id).
		PageToken(pageToken).
		PageSize(pageSize).
		Body(bodyBuilder.Build()).Build()
	resp, err := c.client.Bitable.V1.AppTableRecord.Search(context.Background(), req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call bitable search table: %v", err)
	}
	if !resp.Success() {
		return nil, "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	result := make([]*Record, 0)
	for _, item := range resp.Data.Items {
		record := &Record{
			Id:     *item.RecordId,
			Fields: make(map[string]IField),
		}
		for name, fi := range item.Fields {
			f := c.table.fields[name]
			field := f.Parse(fi)
			record.Fields[name] = field
		}
		result = append(result, record)
	}
	if *resp.Data.HasMore {
		return result, *resp.Data.PageToken, nil
	}
	return result, "", nil
}

func (c *Connection) getTableRecordByPage(pageToken string, total *int) (string, error) {
	const pageSize = 100
	req := larkbitable.NewSearchAppTableRecordReqBuilder().
		AppToken(c.appToken).
		TableId(c.table.id).
		PageToken(pageToken).
		PageSize(pageSize).
		Body(larkbitable.NewSearchAppTableRecordReqBodyBuilder().Build()).Build()
	resp, err := c.client.Bitable.V1.AppTableRecord.Search(context.Background(), req)
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
