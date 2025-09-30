package larkbase

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"reflect"
)

func NewConnection(client *lark.Client, table *Table) *Connection {
	return &Connection{client, table}
}

type Connection struct {
	client *lark.Client
	table  *Table
}

func (c *Connection) CountRecords() (int, error) {
	total := 0
	pageToken := ""
	for {
		var err error
		pageToken, err = c.queryTableRecordByPage(pageToken, &total)
		if err != nil {
			return 0, err
		}
		if pageToken == "" {
			break
		}
	}
	return total, nil
}

func (c *Connection) QueryRecords(arg any, filters ...*larkbitable.Condition) error {
	slicePtrValue := reflect.ValueOf(arg)
	if slicePtrValue.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid argument: expect a pointer, not %v", arg)
	}
	sliceValue := slicePtrValue.Elem()
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("invalid argument: expect a pointer of slice, not %v", arg)
	}
	sliceElemType := sliceValue.Type().Elem()
	if sliceElemType.Kind() != reflect.Struct {
		return fmt.Errorf("invalid argument: expect a pointer of slice with struct element (like []Demo), %v", arg)
	}

	records := make([]*Record, 0)
	pageToken := ""
	for {
		var err error
		records, pageToken, err = c.queryRecordsByPage(filters, pageToken, records)
		if err != nil {
			return err
		}
		if pageToken == "" {
			break
		}
	}
	elemType := sliceValue.Type().Elem()
	newSlice := reflect.MakeSlice(sliceValue.Type(), len(records), len(records))
	for i, record := range records {
		elemValue := reflect.New(elemType).Elem()
		err := c.convertRecordToStruct(record, elemValue)
		if err != nil {
			return err
		}
		newSlice.Index(i).Set(elemValue)
	}
	sliceValue.Set(newSlice)

	return nil
}

func (c *Connection) convertRecordToStruct(record *Record, structValue reflect.Value) error {
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		if !fieldValue.CanSet() {
			continue
		}
		if field.Name == "Meta" && field.Type.String() == "larkbase.Meta" {
			continue
		}
		tag := field.Tag.Get("lark")
		if tag == "" {
			continue
		}
		var value string
		recordField, ok := record.Fields[tag]
		if ok {
			value = recordField.Value
		}
		if field.Type.Kind() != reflect.Struct {
			return fmt.Errorf("invalid struct field: %s %s", field.Name, field.Type.String())
		}
		newInstance := reflect.New(field.Type)
		newInstanceAsField := newInstance.Convert(reflect.TypeOf(&Field{})).Interface().(*Field)
		newInstanceAsField.Name = tag
		newInstanceAsField.Value = value
		fieldValue.Set(newInstance.Elem())
	}
	return nil
}

func (c *Connection) queryRecordsByPage(filters []*larkbitable.Condition, pageToken string, records []*Record) ([]*Record, string, error) {
	const pageSize = 100
	bodyBuilder := larkbitable.NewSearchAppTableRecordReqBodyBuilder()
	bodyBuilder.FieldNames(c.table.FieldNames())
	if len(filters) > 0 {
		bodyBuilder.Filter(larkbitable.NewFilterInfoBuilder().
			Conjunction(`and`).
			Conditions(filters).
			Build())
	}
	bodyBuilder.AutomaticFields(true)
	req := larkbitable.NewSearchAppTableRecordReqBuilder().
		AppToken(c.table.AppToken()).
		TableId(c.table.TableId()).
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
			field := c.table.GetField(name)
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
		AppToken(c.table.AppToken()).
		TableId(c.table.TableId()).
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

//func (c *Client) UpdateRecord(record *Record) {
//	req := larkbitable.NewUpdateAppTableRecordReqBuilder().
//		AppToken(c.appToken).
//		TableId(c.table.Id).
//		RecordId(record.Id).
//		AppTableRecord(larkbitable.NewAppTableRecordBuilder().
//			Fields(record.Build()).
//			Build()).
//		Build()
//	resp, err := c.Bitable.V1.AppTableRecord.Update(context.Background(), req)
//	if err != nil {
//		log.FatalLog("fail to call bitable update table: %v", err)
//	}
//	if !resp.Success() {
//		log.FatalLog("unexpected response error: %s", larkcore.Prettify(resp.CodeError))
//	}
//}
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
