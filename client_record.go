package larkbase

import (
	"context"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func (c *Client) Count() (int, error) {
	if err := c.checkCurrent(); err != nil {
		return 0, err
	}
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

func (c *Client) Where(filters ...*larkbitable.Condition) *Client {
	if err := c.checkCurrent(); err != nil {
		return c
	}
	c.current.filters = filters
	return c
}

func (c *Client) Records() ([]*Record, error) {
	if err := c.checkCurrent(); err != nil {
		return nil, err
	}
	ret := make([]*Record, 0)
	pageToken := ""
	for {
		ret, pageToken = c.queryRecordsByPage(pageToken, ret)
		if c.current.error != nil {
			return nil, c.current.error
		}
		if pageToken == "" {
			break
		}
	}
	return ret, nil
}

func (c *Client) queryRecordsByPage(pageToken string, records []*Record) ([]*Record, string) {
	const pageSize = 100
	bodyBuilder := larkbitable.NewSearchAppTableRecordReqBodyBuilder()
	bodyBuilder.FieldNames(c.current.table.fieldKeys)
	if len(c.current.filters) > 0 {
		bodyBuilder.Filter(larkbitable.NewFilterInfoBuilder().
			Conjunction(`and`).
			Conditions(c.current.filters).
			Build())
	}
	bodyBuilder.AutomaticFields(true)
	req := larkbitable.NewSearchAppTableRecordReqBuilder().
		AppToken(c.current.appToken).
		TableId(c.current.table.id).
		PageToken(pageToken).
		PageSize(pageSize).
		Body(bodyBuilder.Build()).Build()
	resp, err := c.Bitable.V1.AppTableRecord.Search(context.Background(), req)
	if err != nil {
		c.failCurrent("fail to call bitable search table: %v", err)
		return nil, ""
	}
	if !resp.Success() {
		c.failCurrent("get response with error: %s", larkcore.Prettify(resp.CodeError))
		return nil, ""
	}
	result := make([]*Record, 0)
	for _, item := range resp.Data.Items {
		record := &Record{
			Id:     *item.RecordId,
			Fields: make(map[string]IField),
		}
		for name, fi := range item.Fields {
			f := c.current.table.fields[name]
			field := f.Parse(fi)
			record.Fields[name] = field
		}
		result = append(result, record)
	}
	if *resp.Data.HasMore {
		return result, *resp.Data.PageToken
	}
	return result, ""
}
