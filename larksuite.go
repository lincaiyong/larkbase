package larkbase

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/lincaiyong/larkbase/larkfield"
)

func queryAllPages(f func(pageToken string) (newPageToken string, err error)) error {
	pageToken := ""
	for {
		var err error
		pageToken, err = f(pageToken)
		if err != nil {
			return err
		}
		if pageToken == "" {
			break
		}
	}
	return nil
}

// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/search
func (c *Connection[T]) queryRecordsByPage(filters []*larkbitable.Condition, pageToken string, pageSize int, records []*Record) ([]*Record, string, error) {
	if pageSize == 0 {
		pageSize = 100
	}
	bodyBuilder := larkbitable.NewSearchAppTableRecordReqBodyBuilder()
	bodyBuilder.FieldNames(c.fieldNames)
	if len(filters) > 0 {
		bodyBuilder.Filter(larkbitable.NewFilterInfoBuilder().
			Conjunction(`and`).
			Conditions(filters).
			Build())
	}
	bodyBuilder.AutomaticFields(true)
	reqBuilder := larkbitable.NewSearchAppTableRecordReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		PageToken(pageToken).
		PageSize(pageSize).
		Body(bodyBuilder.Build())
	req := reqBuilder.Build()
	var resp, err = c.client.Bitable.V1.AppTableRecord.Search(context.Background(), req)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call bitable search table: %v", err)
	}
	if !resp.Success() {
		return nil, "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	for _, item := range resp.Data.Items {
		record := &Record{
			Id:     *item.RecordId,
			Fields: make(map[string]larkfield.Field),
		}
		for name, value := range item.Fields {
			structField := c.fieldMap[name].Fork()
			err = structField.Parse(value)
			if err != nil {
				return nil, "", err
			}
			record.Fields[name] = structField
		}
		records = append(records, record)
	}
	if *resp.Data.HasMore {
		return records, *resp.Data.PageToken, nil
	}
	return records, "", nil
}

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-record/update
func (c *Connection[T]) updateRecord(record *Record) error {
	fields, err := record.buildForLarkSuite()
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return nil
	}
	req := larkbitable.NewUpdateAppTableRecordReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		RecordId(record.Id).
		AppTableRecord(larkbitable.NewAppTableRecordBuilder().
			Fields(fields).
			Build()).
		Build()
	resp, err := c.client.Bitable.V1.AppTableRecord.Update(context.Background(), req)
	if err != nil {
		return fmt.Errorf("fail to call bitable update table: %v", err)
	}
	if !resp.Success() {
		return fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	return nil
}

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-record/batch_update
func (c *Connection[T]) updateRecords(records []*Record) error {
	reqRecords := make([]*larkbitable.AppTableRecord, 0, len(records))
	for _, record := range records {
		fields, err := record.buildForLarkSuite()
		if err != nil {
			return err
		}
		if len(fields) == 0 {
			continue
		}
		reqRecords = append(reqRecords, larkbitable.NewAppTableRecordBuilder().
			Fields(fields).
			RecordId(record.Id).
			Build())
	}
	if len(reqRecords) == 0 {
		return nil
	}
	req := larkbitable.NewBatchUpdateAppTableRecordReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		Body(larkbitable.NewBatchUpdateAppTableRecordReqBodyBuilder().
			Records(reqRecords).
			Build()).Build()
	resp, err := c.client.Bitable.V1.AppTableRecord.BatchUpdate(context.Background(), req)
	if err != nil {
		return fmt.Errorf("fail to call bitable update table: %v", err)
	}
	if !resp.Success() {
		return fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	return nil
}

func (c *Connection[T]) queryFieldsByPage(pageToken string, fields map[string]larkfield.Type) (string, error) {
	pageSize := 100
	req := larkbitable.NewListAppTableFieldReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		PageToken(pageToken).
		PageSize(pageSize).
		Build()
	resp, err := c.client.Bitable.V1.AppTableField.List(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("fail to call bitable list field: %v", err)
	}
	if !resp.Success() {
		return "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	for _, item := range resp.Data.Items {
		fields[*item.FieldName] = larkfield.Type(*item.Type)
	}
	hasMore := *resp.Data.HasMore
	if hasMore {
		return *resp.Data.PageToken, nil
	}
	return "", nil
}
