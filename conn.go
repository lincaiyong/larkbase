package larkbase

import (
	"errors"
	"github.com/lincaiyong/larkbase/larkfield"
	lark "github.com/lincaiyong/larkbase/larksuite"
	larkbitable "github.com/lincaiyong/larkbase/larksuite/service/bitable/v1"
)

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/bitable-overview
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/search

type Filter = larkbitable.FilterInfo
type ViewFilter = larkbitable.AppTableViewPropertyFilterInfo
type Condition = larkfield.Condition

func Connect[T any](appId, appSecret string) (*Connection[T], error) {
	structPtr := new(T)
	conn := &Connection[T]{condition: structPtr}
	if err := conn.checkStructPtr(structPtr); err != nil {
		return nil, err
	}
	var err error
	conn.tableUrl, conn.appToken, conn.tableId, conn.fields, err = conn.extractAndFillConditionInstance(structPtr)
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

type Connection[T any] struct {
	client *lark.Client

	condition *T

	tableUrl   string
	appToken   string
	tableId    string
	fields     []larkfield.Field
	fieldNames []string
	fieldMap   map[string]larkfield.Field
}

func (c *Connection[T]) Condition() *T {
	return c.condition
}

func (c *Connection[T]) FilterAnd(conditions ...*Condition) *Filter {
	s := make([]*larkbitable.Condition, len(conditions))
	for i, condition := range conditions {
		s[i] = condition.ToLarkCondition()
	}
	return larkbitable.NewFilterInfoBuilder().
		Conjunction(`and`).
		Conditions(s).Build()
}

func (c *Connection[T]) FilterOr(conditions ...*Condition) *Filter {
	s := make([]*larkbitable.Condition, len(conditions))
	for i, condition := range conditions {
		s[i] = condition.ToLarkCondition()
	}
	return larkbitable.NewFilterInfoBuilder().
		Conjunction(`or`).
		Conditions(s).Build()
}

func (c *Connection[T]) ViewFilterAnd(conditions ...*Condition) *ViewFilter {
	s := make([]*larkbitable.AppTableViewPropertyFilterInfoCondition, len(conditions))
	for i, condition := range conditions {
		s[i] = condition.ToLarkViewCondition()
	}
	return larkbitable.NewAppTableViewPropertyFilterInfoBuilder().
		Conjunction(`and`).
		Conditions(s).Build()
}

func (c *Connection[T]) ViewFilterOr(conditions ...*Condition) *ViewFilter {
	s := make([]*larkbitable.AppTableViewPropertyFilterInfoCondition, len(conditions))
	for i, condition := range conditions {
		s[i] = condition.ToLarkViewCondition()
	}
	return larkbitable.NewAppTableViewPropertyFilterInfoBuilder().
		Conjunction(`or`).
		Conditions(s).Build()
}

var errorNotFound = errors.New("record not found")

func (c *Connection[T]) IsNotFoundError(err error) bool {
	return errors.Is(err, errorNotFound)
}

func (c *Connection[T]) Find(structPtr *T, filter *larkbitable.FilterInfo) error {
	if structPtr == nil {
		return errors.New("structPtr is nil")
	}
	if err := c.fillStructPtr(structPtr); err != nil {
		return err
	}
	var err error
	records := make([]*Record, 0)
	records, _, err = c.queryRecordsByPage(filter, "", 1, records)
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

func (c *Connection[T]) FindAll(structPtrSlicePtr *[]*T, filter *larkbitable.FilterInfo) error {
	if structPtrSlicePtr == nil {
		return errors.New("structSlicePtr is nil")
	}
	if err := c.fillStructPtrSlice(*structPtrSlicePtr); err != nil {
		return err
	}
	records := make([]*Record, 0)
	if err := queryAllPages(func(pageToken string) (newPageToken string, err error) {
		records, newPageToken, err = c.queryRecordsByPage(filter, pageToken, 0, records)
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
