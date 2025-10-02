package larkbase

import (
	"errors"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/lincaiyong/larkbase/field"
	"reflect"
)

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/bitable-overview
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/search

func Connect[T any](appId, appSecret string) (*Connection[T], error) {
	structPtr := reflect.New(reflect.TypeOf((*T)(nil)).Elem()).Interface().(*T)
	conn := &Connection[T]{filter: structPtr}
	if err := conn.checkStructPtr(structPtr); err != nil {
		return nil, err
	}
	var err error
	conn.tableUrl, conn.appToken, conn.tableId, conn.fields, err = conn.extractAndFillFilterInstance(structPtr)
	if err != nil {
		return nil, err
	}
	conn.client = lark.NewClient(appId, appSecret)

	conn.fieldMap = make(map[string]Field)
	for _, structField := range conn.fields {
		conn.fieldNames = append(conn.fieldNames, structField.Name())
		conn.fieldMap[structField.Name()] = structField
	}
	err = conn.checkFields()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type Connection[T any] struct {
	client *lark.Client

	filter *T

	tableUrl   string
	appToken   string
	tableId    string
	fields     []Field
	fieldNames []string
	fieldMap   map[string]Field
}

func (c *Connection[T]) Filter() *T {
	return c.filter
}

var errorNotFound = errors.New("record not found")

func (c *Connection[T]) IsNotFoundError(err error) bool {
	return errors.Is(err, errorNotFound)
}

func (c *Connection[T]) FindOne(structPtr *T, filters ...*larkbitable.Condition) error {
	if structPtr == nil {
		return errors.New("structPtr is nil")
	}
	if err := c.checkStructPtr(structPtr); err != nil {
		return err
	}
	var err error
	records := make([]*Record, 0)
	records, _, err = c.queryRecordsByPage(filters, "", 1, records)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return errorNotFound
	}
	record := records[0]
	err = c.convertRecordToStructPtr(record, structPtr)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection[T]) FindAll(structPtrSlicePtr *[]*T, filters ...*larkbitable.Condition) error {
	if structPtrSlicePtr == nil {
		return errors.New("structSlicePtr is nil")
	}
	if err := c.checkStructPtrSlicePtr(structPtrSlicePtr); err != nil {
		return err
	}
	records := make([]*Record, 0)
	if err := queryAllPages(func(pageToken string) (newPageToken string, err error) {
		records, newPageToken, err = c.queryRecordsByPage(filters, pageToken, 0, records)
		return
	}); err != nil {
		return err
	}
	return c.convertRecordsToStructPtrSlicePtr(records, structPtrSlicePtr)
}

func (c *Connection[T]) UpdateOne(structPtr *T) error {
	if structPtr == nil {
		return errors.New("structPtr is nil")
	}
	if err := c.checkStructPtr(structPtr); err != nil {
		return err
	}
	record, err := c.convertStructPtrToRecord(structPtr)
	if err != nil {
		return err
	}
	return c.updateRecord(record)
}

func (c *Connection[T]) UpdateAll(structPtrSlice []*T) error {
	if err := c.checkStructPtrSlice(structPtrSlice); err != nil {
		return err
	}
	records, err := c.convertStructPtrSliceToRecords(structPtrSlice)
	if err != nil {
		return err
	}
	return c.updateRecords(records)
}

func (c *Connection[T]) checkFields() error {
	fields := make(map[string]field.Type)
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

//
//func (c *Client) DeleteRecords(records []*Record) {
//	for _, record := range records {
//		builder := larkbitable.NewDeleteAppTableRecordReqBuilder().
//			AppToken(c.appToken).TableId(c.table.Id).
//			RecordId(record.Id)
//		req := builder.Build()
//		resp, err := c.Bitable.V1.AppTableRecord.Delete(context.Background(), req)
//		if err != nil {
//			log.FatalLog("fail to call bitable delete table: %v", err)
//		}
//		if !resp.Success() {
//			log.FatalLog("unexpected response error: %s", larkcore.Prettify(resp.CodeError))
//		}
//	}
//}
//
//func (c *Client) AddRecord(fields map[string]IField) {
//	record := Record{Fields: fields}
//	req := larkbitable.NewCreateAppTableRecordReqBuilder().
//		AppToken(c.appToken).
//		TableId(c.table.Id).
//		AppTableRecord(larkbitable.NewAppTableRecordBuilder().
//			Fields(record.Build()).
//			Build()).
//		Build()
//	resp, err := c.Bitable.V1.AppTableRecord.Create(context.Background(), req)
//	if err != nil {
//		log.FatalLog("fail to call bitable create record: %v", err)
//	}
//	if !resp.Success() {
//		log.FatalLog("unexpected response error: %s", larkcore.Prettify(resp.CodeError))
//	}
//}
//
//func (c *Client) UploadFile(filePath string) string {
//	stat, err := os.Stat(filePath)
//	if err != nil {
//		log.FatalLog("fail to stat file: %v", err)
//	}
//	file, err := os.Open(filePath)
//	if err != nil {
//		log.FatalLog("fail to open file: %v", err)
//	}
//	defer func() { _ = file.Close() }()
//	req := larkdrive.NewUploadAllMediaReqBuilder().
//		Body(larkdrive.NewUploadAllMediaReqBodyBuilder().
//			FileName(path.Base(filePath)).
//			ParentType(`bitable_file`).
//			ParentNode(c.appToken).
//			Size(int(stat.Size())).
//			File(file).
//			Build()).
//		Build()
//
//	resp, err := c.Drive.V1.Media.UploadAll(context.Background(), req)
//	if err != nil {
//		log.FatalLog("fail to call bitable upload table: %v", err)
//	}
//	if !resp.Success() {
//		log.FatalLog("unexpected response error: %s", larkcore.Prettify(resp.CodeError))
//	}
//	return *resp.Data.FileToken
//}
