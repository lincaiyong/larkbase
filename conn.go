package larkbase

import (
	"context"
	"errors"
	"fmt"
	"github.com/lincaiyong/larkbase/larkfield"
	lark "github.com/lincaiyong/larkbase/larksuite"
	"strings"
	"unicode"
)

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/bitable-overview
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/search

func DescribeTable(ctx context.Context, appId, appSecret, url string) (string, error) {
	appToken, tableId := extractAppTokenTableIdFromUrl(url)
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
		switch field.Type() {
		case "Text", "Number", "SingleSelect", "MultiSelect", "Date", "Checkbox", "Url", "AutoNumber", "ModifiedTime":
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

func ConnectWithOpts[T any](ctx context.Context, appId, appSecret, tableUrl string, fieldNameMapping map[string]string) (*Connection[T], error) {
	structPtr := new(T)
	conn := &Connection[T]{ctx: ctx, condition: structPtr, fieldNameMapping: fieldNameMapping}
	if err := conn.checkStructPtr(structPtr, tableUrl); err != nil {
		return nil, err
	}
	var err error
	conn.tableUrl, conn.appToken, conn.structName, conn.tableId, conn.fields, err = conn.extractAndFillConditionInstance(structPtr, tableUrl)
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

func ConnectWithUrl[T any](ctx context.Context, appId, appSecret, tableUrl string) (*Connection[T], error) {
	return ConnectWithOpts[T](ctx, appId, appSecret, tableUrl, nil)
}

func Connect[T any](ctx context.Context, appId, appSecret string) (*Connection[T], error) {
	return ConnectWithUrl[T](ctx, appId, appSecret, "")
}

type Connection[T any] struct {
	ctx    context.Context
	client *lark.Client

	condition *T

	tableUrl   string
	appToken   string
	structName string
	tableId    string
	fields     []larkfield.Field
	fieldNames []string
	fieldMap   map[string]larkfield.Field

	fieldNameMapping map[string]string
}

func (c *Connection[T]) fieldRealName(fieldName string) string {
	if c.fieldNameMapping != nil {
		if ret, ok := c.fieldNameMapping[fieldName]; ok {
			return ret
		}
	}
	return fieldName
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

func (c *Connection[T]) UpdateAll(structPtrSlice []*T) error {
	if err := c.fillStructPtrSlice(structPtrSlice); err != nil {
		return err
	}
	records, err := c.convertStructPtrSliceToRecords(structPtrSlice)
	if err != nil {
		return err
	}
	return c.updateRecords(records)
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
	records, err = c.createRecords(records)
	if err != nil {
		return nil, err
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
