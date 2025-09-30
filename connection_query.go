package larkbase

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

type Record struct {
	Id     string
	Fields map[string]Field
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
	for _, item := range resp.Data.Items {
		record := &Record{
			Id:     *item.RecordId,
			Fields: make(map[string]Field),
		}
		for name, fi := range item.Fields {
			field := c.table.fields[name]
			field.parseValue(fi)
			record.Fields[name] = field
		}
		records = append(records, record)
	}
	if *resp.Data.HasMore {
		return records, *resp.Data.PageToken, nil
	}
	return records, "", nil
}

func (c *Connection) queryTableRecordByPage(pageToken string, total *int) (string, error) {
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
