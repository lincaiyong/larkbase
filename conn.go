package larkbase

import (
	"context"
	"errors"
	"fmt"
	"github.com/lincaiyong/larkbase/larkfield"
	lark "github.com/lincaiyong/larkbase/larksuite"
	larkbitable "github.com/lincaiyong/larkbase/larksuite/service/bitable/v1"
	"github.com/lincaiyong/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/bitable-overview
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/search

type Filter = larkbitable.FilterInfo
type ViewFilter = larkbitable.AppTableViewPropertyFilterInfo
type Condition = larkfield.Condition

const modifiedTimeFieldName = "modified_time"

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
			sb.WriteString(fmt.Sprintf("    %s larkbase.%sField `lark:\"%s\"`\n", field.Name(), field.Type(), field.Name()))
		}
	}
	sb.WriteString("}")
	return sb.String(), nil
}

func Connect[T any](ctx context.Context, appId, appSecret string) (*Connection[T], error) {
	structPtr := new(T)
	conn := &Connection[T]{ctx: ctx, condition: structPtr}
	if err := conn.checkStructPtr(structPtr); err != nil {
		return nil, err
	}
	var err error
	conn.tableUrl, conn.appToken, conn.structName, conn.tableId, conn.fields, err = conn.extractAndFillConditionInstance(structPtr)
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
}

func (c *Connection[T]) Sort() *T {
	return c.condition
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

func (c *Connection[T]) Find(structPtr *T, filter *larkbitable.FilterInfo, sorts ...*larkbitable.Sort) error {
	if structPtr == nil {
		return errors.New("structPtr is nil")
	}
	if err := c.fillStructPtr(structPtr); err != nil {
		return err
	}
	var err error
	records := make([]*Record, 0)
	records, _, err = c.queryRecordsByPage(filter, sorts, "", 1, records)
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

func (c *Connection[T]) FindAll(structPtrSlicePtr *[]*T, filter *larkbitable.FilterInfo, sorts ...*larkbitable.Sort) error {
	if structPtrSlicePtr == nil {
		return errors.New("structSlicePtr is nil")
	}
	if err := c.fillStructPtrSlice(*structPtrSlicePtr); err != nil {
		return err
	}
	records := make([]*Record, 0)
	if err := queryAllPages(func(pageToken string) (newPageToken string, err error) {
		records, newPageToken, err = c.queryRecordsByPage(filter, sorts, pageToken, 0, records)
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

func (c *Connection[T]) SyncToDatabase(db *gorm.DB, batchSize int) error {
	items := make([]string, 0)
	for _, name := range c.fieldNames {
		items = append(items, fmt.Sprintf("`%s` TEXT", name))
	}
	tableName := strings.ToLower(c.structName)
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (`record_id` VARCHAR(255) PRIMARY KEY, %s)",
		tableName, strings.Join(items, ", "))
	log.InfoLog("sql: %s", sql)
	if err := db.WithContext(c.ctx).Exec(sql).Error; err != nil {
		return err
	}
	var count int64
	sql = fmt.Sprintf("SELECT COUNT(1) FROM `%s`", tableName)
	if err := db.WithContext(c.ctx).Raw(sql).Scan(&count).Error; err != nil {
		return err
	}
	log.InfoLog("count: %d", count)
	var latestModifiedTime *time.Time
	if count > 0 {
		sql = fmt.Sprintf("SELECT MAX(modified_time) FROM `%s`", tableName)
		var latestModifiedTimeStr string
		result := db.WithContext(c.ctx).Raw(sql).Scan(&latestModifiedTimeStr)
		if result.Error != nil {
			return result.Error
		}
		t, err := larkfield.BeijingDateTimeStrToTime(latestModifiedTimeStr)
		if err != nil {
			return err
		}
		latestModifiedTime = &t
		log.InfoLog("latestModifiedTime: %s", latestModifiedTimeStr)
	}
	var filter *larkbitable.FilterInfo
	if latestModifiedTime != nil {
		field := c.fieldMap[modifiedTimeFieldName].(*larkfield.ModifiedTimeField)
		filter = c.FilterAnd(field.IsGreater(*latestModifiedTime))
	}
	var rawRecords []*T
	err := c.FindAll(&rawRecords, filter)
	if err != nil {
		return err
	}
	log.InfoLog("rawRecords count: %d", len(rawRecords))
	var records []*Record
	for _, rawRecord := range rawRecords {
		record, convErr := c.convertStructPtrToRecord(rawRecord)
		if convErr != nil {
			return convErr
		}
		if latestModifiedTime != nil {
			if !record.ModifiedTime.After(*latestModifiedTime) {
				continue
			}
		}
		records = append(records, record)
	}
	log.InfoLog("records count: %d", len(records))
	if len(records) > 0 {
		if batchSize == 0 || batchSize > len(records) {
			batchSize = len(records)
		}
		for i := 0; i < len(records)/batchSize; i++ {
			var batchRecords []*Record
			if (i+1)*batchSize > len(records) {
				batchRecords = records[i*batchSize:]
			} else {
				batchRecords = records[i*batchSize : (i+1)*batchSize]
			}
			columns := []string{"`record_id`"}
			for _, name := range c.fieldNames {
				columns = append(columns, fmt.Sprintf("`%s`", name))
			}
			values := make([]any, 0)
			valuesPlaceHolders := make([]string, 0, len(batchRecords))
			for _, record := range batchRecords {
				placeholders := []string{"?"}
				values = append(values, record.Id)
				for _, fieldName := range c.fieldNames {
					placeholders = append(placeholders, "?")
					values = append(values, record.Fields[fieldName].StringValue())
				}
				valuesPlaceHolders = append(valuesPlaceHolders, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
			}
			var updateItems []string
			for _, name := range c.fieldNames {
				updateItems = append(updateItems, fmt.Sprintf("`%s` = VALUES(`%s`)", name, name))
			}
			sql = fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s ON DUPLICATE KEY UPDATE %s",
				tableName,
				strings.Join(columns, ", "),
				strings.Join(valuesPlaceHolders, ", "),
				strings.Join(updateItems, ", "))
			log.InfoLog("sql: %s", sql)
			if err = db.WithContext(c.ctx).Exec(sql, values...).Error; err != nil {
				return err
			}
			log.InfoLog("insert or update %d records", len(batchRecords))
		}
	}
	log.InfoLog("successfully sync larkbase to database")
	return nil
}
