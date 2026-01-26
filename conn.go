package larkbase

import (
	"context"
	"errors"
	"fmt"
	"github.com/lincaiyong/larkbase/larkfield"
	lark "github.com/lincaiyong/larkbase/larksuite"
	"os"
	"strings"
	"unicode"
)

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/bitable-overview
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/search

func DescribeTable(ctx context.Context, url string) (string, error) {
	appToken, tableId, _ := extractAppTokenTableIdViewIdFromUrl(url)
	if appToken == "" || tableId == "" {
		return "", fmt.Errorf("invalid table url: %s", url)
	}
	client := lark.NewClient(appId, appSecret)
	fields := make(map[string]larkfield.Field)
	err := queryAllPages(func(pageToken string) (newPageToken string, err error) {
		return queryFieldsByPage(ctx, client, appToken, tableId, pageToken, fields)
	})
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString("type Record struct {\n")
	sb.WriteString(fmt.Sprintf("    larkbase.Meta `lark:\"%s\"`\n\n", url))
	for _, field := range fields {
		switch field.TypeStr() {
		case "Text", "Number", "SingleSelect", "MultiSelect", "Date", "Checkbox", "Url", "AutoNumber", "ModifiedTime", "Lookup", "Formula":
			sb.WriteString(fmt.Sprintf("    %s larkbase.%sField `lark:\"%s\"`\n", toCamelCase(field.Name()), field.Type(), field.Name()))
		}
	}
	sb.WriteString("}")
	return sb.String(), nil
}

func toCamelCase(s string) string {
	var sb strings.Builder
	shouldUpper := true
	for _, r := range s {
		if r == '_' { // 下划线分隔符，将下一个字符转换为大写字母
			shouldUpper = true
		} else {
			if shouldUpper && unicode.IsLetter(r) { // 需要转换为大写字母
				sb.WriteRune(unicode.ToUpper(r))
				shouldUpper = false
			} else {
				sb.WriteRune(unicode.ToLower(r))
			}
		}
	}
	return sb.String()
}

var appId, appSecret string

func init() {
	appId = os.Getenv("LARK_APP_ID")
	appSecret = os.Getenv("LARK_APP_SECRET")
}

func SetAppIdSecret(appId_, appSecret_ string) {
	appId = appId_
	appSecret = appSecret_
}

func ConnectAny(ctx context.Context, tableUrl string) (*Connection[AnyRecord], error) {
	conn := &Connection[AnyRecord]{
		ctx:         ctx,
		tableUrl:    tableUrl,
		isAnyRecord: true,
		condition:   &AnyRecord{},
	}
	err := conn.fillStructPtr(conn.condition)
	if err != nil {
		return nil, err
	}
	conn.appToken, conn.tableId, conn.viewId = extractAppTokenTableIdViewIdFromUrl(tableUrl)
	if conn.appToken == "" || conn.tableId == "" {
		return nil, fmt.Errorf("invalid table url: %s", tableUrl)
	}
	conn.client = lark.NewClient(appId, appSecret)
	fields := make(map[string]larkfield.Field)
	err = queryAllPages(func(pageToken string) (newPageToken string, err error) {
		return conn.queryFieldsByPage(pageToken, fields)
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query fields: %v", err)
	}
	conn.fieldMap = make(map[string]larkfield.Field)
	for _, field := range fields {
		conn.fieldMap[field.Name()] = field
		conn.fieldNames = append(conn.fieldNames, field.Name())
	}
	return conn, nil
}

func ConnectUrl[T any](ctx context.Context, tableUrl string) (*Connection[T], error) {
	structPtr := new(T)
	conn := &Connection[T]{ctx: ctx, condition: structPtr}
	if err := conn.checkStructPtr(structPtr, tableUrl); err != nil {
		return nil, err
	}
	var err error
	conn.tableUrl, conn.appToken, conn.structName, conn.tableId, conn.viewId, conn.fields, err = conn.extractAndFillConditionInstance(structPtr, tableUrl)
	if err != nil {
		return nil, err
	}
	conn.client = lark.NewClient(appId, appSecret)
	conn.fieldMap = make(map[string]larkfield.Field)
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

func Connect[T any](ctx context.Context) (*Connection[T], error) {
	return ConnectUrl[T](ctx, "")
}

type Connection[T any] struct {
	ctx    context.Context
	client *lark.Client

	condition *T

	tableUrl   string
	appToken   string
	tableId    string
	viewId     string
	structName string
	fields     []larkfield.Field
	fieldNames []string
	fieldMap   map[string]larkfield.Field

	isAnyRecord bool

	batchSize int
}

func (c *Connection[T]) Context() context.Context {
	return c.ctx
}

func (c *Connection[T]) SetTosValue(f *larkfield.TextField, b []byte) error {
	key, err := tosPutFn(c.ctx, b)
	if err != nil {
		return err
	}
	f.SetValue(key)
	return nil
}

func (c *Connection[T]) GetTosValue(f *larkfield.TextField) ([]byte, error) {
	hash := f.StringValue()
	if hash == "" {
		return nil, nil
	}
	if len(hash) != 32 {
		return nil, fmt.Errorf("invalid md5: %s", hash)
	}
	return tosGetFn(c.ctx, hash)
}

func (c *Connection[T]) TableUrl() string {
	return c.tableUrl
}

func (c *Connection[T]) AppToken() string {
	return c.appToken
}

func (c *Connection[T]) TableId() string {
	return c.tableId
}

func (c *Connection[T]) ViewId() string {
	return c.viewId
}

func (c *Connection[T]) StructName() string {
	return c.structName
}

var errorNotFound = errors.New("record not found")

func (c *Connection[T]) IsNotFoundError(err error) bool {
	return errors.Is(err, errorNotFound)
}

func (c *Connection[T]) Find(structPtr *T, opt *FindOption) error {
	if structPtr == nil {
		return errors.New("structPtr is nil")
	}
	if err := c.fillStructPtr(structPtr); err != nil {
		return err
	}
	var err error
	records := make([]*Record, 0)
	if opt == nil {
		opt = &FindOption{}
	}
	records, _, err = c.queryRecordsByPage(opt.viewId, opt.filter, opt.sorts, "", 1, records)
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

func (c *Connection[T]) FindAll(structPtrSlicePtr *[]*T, opt *FindOption) error {
	if structPtrSlicePtr == nil {
		return errors.New("structSlicePtr is nil")
	}
	if err := c.fillStructPtrSlice(*structPtrSlicePtr); err != nil {
		return err
	}
	records := make([]*Record, 0)
	pageSize := 0
	if opt == nil {
		opt = &FindOption{}
	}
	if opt.limit > 0 {
		pageSize = opt.limit
	}
	if err := queryAllPages(func(pageToken string) (newPageToken string, err error) {
		records, newPageToken, err = c.queryRecordsByPage(opt.viewId, opt.filter, opt.sorts, pageToken, pageSize, records)
		if opt.limit > 0 && len(records) >= opt.limit {
			newPageToken = ""
		}
		return
	}); err != nil {
		return err
	}
	return c.convertRecordsToStructPtrSlicePtr(records, structPtrSlicePtr)
}

func (c *Connection[T]) Update(structPtr *T) error {
	if structPtr == nil {
		return errors.New("structPtr is nil")
	}
	if err := c.fillStructPtr(structPtr); err != nil {
		return err
	}
	record, err := c.convertStructPtrToRecord(structPtr)
	if err != nil {
		return err
	}
	return c.updateRecord(record)
}

func (c *Connection[T]) SetBatchSize(batchSize int) {
	c.batchSize = batchSize
}

func (c *Connection[T]) UpdateAll(structPtrSlice []*T) error {
	if err := c.fillStructPtrSlice(structPtrSlice); err != nil {
		return err
	}
	records, err := c.convertStructPtrSliceToRecords(structPtrSlice)
	if err != nil {
		return err
	}
	if c.batchSize == 0 {
		return c.updateRecords(records)
	} else {
		for i := 0; i < len(records); i += c.batchSize {
			end := i + c.batchSize
			if end > len(records) {
				end = len(records)
			}
			batchRecords := records[i:end]
			err = c.updateRecords(batchRecords)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (c *Connection[T]) Create(structPtr *T) error {
	if structPtr == nil {
		return errors.New("structPtr is nil")
	}
	if err := c.fillStructPtr(structPtr); err != nil {
		return err
	}
	record, err := c.convertStructPtrToRecord(structPtr)
	if err != nil {
		return err
	}
	record, err = c.createRecord(record)
	if err != nil {
		return err
	}
	return c.convertRecordToStructPtr(record, structPtr)
}

func (c *Connection[T]) CreateAll(structPtrSlice []*T) ([]*T, error) {
	if err := c.fillStructPtrSlice(structPtrSlice); err != nil {
		return nil, err
	}
	records, err := c.convertStructPtrSliceToRecords(structPtrSlice)
	if err != nil {
		return nil, err
	}
	if c.batchSize == 0 {
		records, err = c.createRecords(records)
		if err != nil {
			return nil, err
		}
	} else {
		resultRecords := make([]*Record, 0, len(records))
		for i := 0; i < len(records); i += c.batchSize {
			end := i + c.batchSize
			if end > len(records) {
				end = len(records)
			}
			batchRecords := records[i:end]
			batchRecords, err = c.createRecords(batchRecords)
			if err != nil {
				return nil, err
			}
			resultRecords = append(resultRecords, batchRecords...)
		}
		records = resultRecords
	}
	err = c.convertRecordsToStructPtrSlicePtr(records, &structPtrSlice)
	if err != nil {
		return nil, err
	}
	return structPtrSlice, nil
}

func (c *Connection[T]) Delete(structPtr *T) error {
	if structPtr == nil {
		return errors.New("structPtr is nil")
	}
	if err := c.fillStructPtr(structPtr); err != nil {
		return err
	}
	record, err := c.convertStructPtrToRecord(structPtr)
	if err != nil {
		return err
	}
	err = c.deleteRecord(record)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection[T]) DeleteAll(structPtrSlice []*T) error {
	if err := c.fillStructPtrSlice(structPtrSlice); err != nil {
		return err
	}
	records, err := c.convertStructPtrSliceToRecords(structPtrSlice)
	if err != nil {
		return err
	}
	err = c.deleteRecords(records)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection[T]) CreateView(name string, filter *ViewFilter) error {
	viewId, err := c.createView(name)
	if err != nil {
		return err
	}
	err = c.updateView(viewId, name, filter)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection[T]) ListFields() (map[string]larkfield.Type, error) {
	return c.listFields()
}

func (c *Connection[T]) CreateField(name string, type_ larkfield.Type) error {
	return c.createField(name, type_)
}
