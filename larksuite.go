package larkbase

import (
	"context"
	"errors"
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

func (c *Connection[T]) checkFields() error {
	fields := make(map[string]larkfield.Type)
	err := queryAllPages(func(pageToken string) (newPageToken string, err error) {
		return c.queryFieldsByPage(pageToken, fields)
	})
	if err != nil {
		return err
	}
	for name, structField := range c.fieldMap {
		f, ok := fields[name]
		if !ok {
			return fmt.Errorf("field %s is not found in larkbase table: %s", name, c.tableUrl)
		}
		if structField.Type() != f.String() {
			return fmt.Errorf("field %s in larkbase table %s has type %s, not %s", name, c.tableUrl, f.String(), structField.Type())
		}
	}
	return nil
}

func (c *Connection[T]) parseAppTableRecord(item *larkbitable.AppTableRecord) (*Record, error) {
	record := &Record{
		Id:     *item.RecordId,
		Fields: make(map[string]larkfield.Field),
	}
	if item.LastModifiedTime != nil {
		record.ModifiedTime = larkfield.UnixSecondsToTime((*item.LastModifiedTime) / 1000)
	}
	for name, value := range item.Fields {
		structField := c.fieldMap[name].Fork()
		err := structField.Parse(value)
		if err != nil {
			return nil, err
		}
		record.Fields[name] = structField
	}
	return record, nil
}

func (c *Connection[T]) parseAppTableRecords(items []*larkbitable.AppTableRecord) ([]*Record, error) {
	records := make([]*Record, 0, len(items))
	for _, item := range items {
		record, err := c.parseAppTableRecord(item)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
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
	records, err = c.parseAppTableRecords(resp.Data.Items)
	if err != nil {
		return nil, "", err
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
	if record.Id == "" {
		return fmt.Errorf("record id is empty")
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

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-record/create
func (c *Connection[T]) createRecord(record *Record) (*Record, error) {
	fields, err := record.buildForLarkSuite()
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, errors.New("fail to create record: empty fields")
	}
	req := larkbitable.NewCreateAppTableRecordReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		AppTableRecord(larkbitable.NewAppTableRecordBuilder().
			Fields(fields).
			Build()).
		Build()
	resp, err := c.client.Bitable.V1.AppTableRecord.Create(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("fail to call bitable create record: %v", err)
	}
	if !resp.Success() {
		return nil, fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	var ret *Record
	ret, err = c.parseAppTableRecord(resp.Data.Record)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-record/batch_create
func (c *Connection[T]) createRecords(records []*Record) ([]*Record, error) {
	reqRecords := make([]*larkbitable.AppTableRecord, 0, len(records))
	for _, record := range records {
		fields, err := record.buildForLarkSuite()
		if err != nil {
			return nil, err
		}
		if len(fields) == 0 {
			continue
		}
		reqRecords = append(reqRecords, larkbitable.NewAppTableRecordBuilder().
			Fields(fields).
			Build())
	}
	if len(reqRecords) == 0 {
		return nil, fmt.Errorf("fail to create records: empty records")
	}
	req := larkbitable.NewBatchCreateAppTableRecordReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		Body(larkbitable.NewBatchCreateAppTableRecordReqBodyBuilder().Records(reqRecords).Build()).
		Build()
	resp, err := c.client.Bitable.V1.AppTableRecord.BatchCreate(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("fail to call bitable create record: %v", err)
	}
	if !resp.Success() {
		return nil, fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	var ret []*Record
	ret, err = c.parseAppTableRecords(resp.Data.Records)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-record/delete
func (c *Connection[T]) deleteRecord(record *Record) error {
	if record.Id == "" {
		return fmt.Errorf("record id is empty")
	}
	builder := larkbitable.NewDeleteAppTableRecordReqBuilder().
		AppToken(c.appToken).TableId(c.tableId).
		RecordId(record.Id)
	req := builder.Build()
	resp, err := c.client.Bitable.V1.AppTableRecord.Delete(context.Background(), req)
	if err != nil {
		return fmt.Errorf("fail to call bitable delete record: %v", err)
	}
	if !resp.Success() {
		return fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	return nil
}

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-record/batch_delete
func (c *Connection[T]) deleteRecords(records []*Record) error {
	recordIds := make([]string, 0, len(records))
	for _, record := range records {
		if record.Id == "" {
			continue
		}
		recordIds = append(recordIds, record.Id)
	}
	if len(recordIds) == 0 {
		return fmt.Errorf("fail to delete records: empty records")
	}
	req := larkbitable.NewBatchDeleteAppTableRecordReqBuilder().
		AppToken(c.appToken).TableId(c.tableId).
		Body(larkbitable.NewBatchDeleteAppTableRecordReqBodyBuilder().
			Records(recordIds).
			Build()).Build()
	resp, err := c.client.Bitable.V1.AppTableRecord.BatchDelete(context.Background(), req)
	if err != nil {
		return fmt.Errorf("fail to call bitable delete record: %v", err)
	}
	if !resp.Success() {
		return fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	return nil
}
