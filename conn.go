package larkbase

import (
	"errors"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/bitable-overview
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/search

func Connect(appId, appSecret string, structPtr any) (*Connection, error) {
	if err := checkUserStructPtr(structPtr); err != nil {
		return nil, err
	}
	table, err := extractTableAndFillBasicInfo(structPtr)
	if err != nil {
		return nil, err
	}
	client := lark.NewClient(appId, appSecret)
	return &Connection{client, table}, nil
}

type Connection struct {
	client *lark.Client
	table  *Table
}

var errorNotFound = errors.New("record not found")

func (c *Connection) IsNotFoundError(err error) bool {
	return errors.Is(err, errorNotFound)
}

func (c *Connection) FindOne(structPtr any, filters ...*larkbitable.Condition) error {
	if err := checkUserStructPtr(structPtr); err != nil {
		return err
	}

	var err error
	records := make([]*Record, 0)
	records, _, err = c.queryRecordsByPage(filters, "", 1, records)
	if err != nil {
		return err
	}
	if len(records) != 1 {
		return errorNotFound
	}
	record := records[0]
	err = convertRecordToUserStruct(record, structPtr)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) FindAll(structSlicePtr any, filters ...*larkbitable.Condition) error {
	if err := checkUserStructSlicePtr(structSlicePtr); err != nil {
		return err
	}

	records := make([]*Record, 0)
	if err := queryAllPages(func(pageToken string) (newPageToken string, err error) {
		records, newPageToken, err = c.queryRecordsByPage(filters, pageToken, 0, records)
		return
	}); err != nil {
		return err
	}
	return convertRecordsToUserStructSlicePtr(records, structSlicePtr)
}

func (c *Connection) UpdateOne(structPtr any) error {
	if err := checkUserStructPtr(structPtr); err != nil {
		return err
	}
	record, err := convertUserStructToRecord(structPtr)
	if err != nil {
		return err
	}
	return c.updateRecord(record)
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
