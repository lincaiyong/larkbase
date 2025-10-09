package larkbase

import (
	"context"
	"errors"
	"fmt"
	"github.com/lincaiyong/larkbase/larkfield"
	lark "github.com/lincaiyong/larkbase/larksuite"
	larkcore "github.com/lincaiyong/larkbase/larksuite/core"
	larkbitable "github.com/lincaiyong/larkbase/larksuite/service/bitable/v1"
	"github.com/lincaiyong/log"
	"time"
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
	if c.fieldMap[larkfield.ModifiedTimeFieldName] == nil || c.fieldMap[larkfield.ModifiedTimeFieldName].Type() != "Number" {
		return fmt.Errorf("field \"%s\" with Number type is required in struct: %s", larkfield.ModifiedTimeFieldName, c.structName)
	}
	fields := make(map[string]larkfield.Field)
	err := queryAllPages(func(pageToken string) (newPageToken string, err error) {
		return c.queryFieldsByPage(pageToken, fields)
	})
	if err != nil {
		return err
	}
	if c.fieldMap[larkfield.ModifiedTimeFieldName].Type() != "Number" {
		return fmt.Errorf("field \"%s\" with Number type is required in larkbase table: %s", larkfield.ModifiedTimeFieldName, c.tableUrl)
	}
	for name, structField := range c.fieldMap {
		f, ok := fields[name]
		if !ok {
			return fmt.Errorf("field %s is not found in larkbase table: %s", name, c.tableUrl)
		}
		if structField.Type() != f.Type() {
			return fmt.Errorf("field %s in larkbase table %s has type %s, not %s", name, c.tableUrl, f.Type(), structField.Type())
		}
		structField.SetId(f.Id())
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
func (c *Connection[T]) queryRecordsByPage(filter *larkbitable.FilterInfo, sorts []*larkbitable.Sort, pageToken string, pageSize int, records []*Record) ([]*Record, string, error) {
	if pageSize == 0 {
		pageSize = 100
	}
	bodyBuilder := larkbitable.NewSearchAppTableRecordReqBodyBuilder()
	bodyBuilder.FieldNames(c.fieldNames)
	if filter != nil {
		bodyBuilder.Filter(filter)
	}
	if sorts != nil {
		bodyBuilder.Sort(sorts)
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
	var newRecords []*Record
	newRecords, err = c.parseAppTableRecords(resp.Data.Items)
	if err != nil {
		return nil, "", err
	}
	for _, record := range newRecords {
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

func queryFieldsByPage(client *lark.Client, appToken, tableId, pageToken string, fields map[string]larkfield.Field) (string, error) {
	pageSize := 100
	req := larkbitable.NewListAppTableFieldReqBuilder().
		AppToken(appToken).
		TableId(tableId).
		PageToken(pageToken).
		PageSize(pageSize).
		Build()
	resp, err := client.Bitable.V1.AppTableField.List(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("fail to call bitable list field: %v", err)
	}
	if !resp.Success() {
		return "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	for _, item := range resp.Data.Items {
		id := *item.FieldId
		name := *item.FieldName
		type_ := larkfield.Type(*item.Type).String()
		fields[*item.FieldName] = larkfield.NewBaseField(id, name, type_)
	}
	hasMore := *resp.Data.HasMore
	if hasMore {
		return *resp.Data.PageToken, nil
	}
	return "", nil
}

func (c *Connection[T]) queryFieldsByPage(pageToken string, fields map[string]larkfield.Field) (string, error) {
	return queryFieldsByPage(c.client, c.appToken, c.tableId, pageToken, fields)
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
	fields[larkfield.ModifiedTimeFieldName] = larkfield.TimeToModifiedTime(time.Now())
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

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-view/create
func (c *Connection[T]) createView(name string) (string, error) {
	req := larkbitable.NewCreateAppTableViewReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		ReqView(larkbitable.NewReqViewBuilder().
			ViewName(name).
			ViewType(`grid`).
			Build()).
		Build()
	resp, err := c.client.Bitable.V1.AppTableView.Create(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("fail to call bitable create view: %v", err)
	}
	if !resp.Success() {
		return "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	viewId := *resp.Data.View.ViewId
	return viewId, nil
}

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-view/patch
func (c *Connection[T]) updateView(viewId, viewName string, filter *ViewFilter) error {
	req := larkbitable.NewPatchAppTableViewReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		ViewId(viewId).
		Body(larkbitable.NewPatchAppTableViewReqBodyBuilder().
			ViewName(viewName).
			Property(larkbitable.NewAppTableViewPropertyBuilder().
				FilterInfo(filter).
				Build()).
			Build()).
		Build()
	resp, err := c.client.Bitable.V1.AppTableView.Patch(context.Background(), req)
	if err != nil {
		return fmt.Errorf("fail to call bitable update view: %v", err)
	}
	if !resp.Success() {
		return fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	return nil
}

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-field/list?appId=cli_a8d592a8e236d00b
func (c *Connection[T]) listFields() (map[string]larkfield.Type, error) {
	req := larkbitable.NewListAppTableFieldReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		PageSize(100).
		Build()
	resp, err := c.client.Bitable.V1.AppTableField.List(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("fail to call bitable list field: %v", err)
	}
	if !resp.Success() {
		return nil, fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	if *resp.Data.HasMore {
		log.WarnLog("unexpected more than 100 fields, ignored")
	}
	ret := make(map[string]larkfield.Type)
	for _, item := range resp.Data.Items {
		ret[*item.FieldName] = larkfield.Type(*item.Type)
	}
	return ret, nil
}

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-field/create
func (c *Connection[T]) createField(name string, type_ larkfield.Type) error {
	req := larkbitable.NewCreateAppTableFieldReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		AppTableField(larkbitable.NewAppTableFieldBuilder().
			FieldName(name).
			Type(int(type_)).
			Build()).
		Build()
	resp, err := c.client.Bitable.V1.AppTableField.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("fail to call bitable create field: %v", err)
	}
	if !resp.Success() {
		return fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	return nil
}
