package larkbase

import (
	"context"
	"fmt"
	"github.com/lincaiyong/larkbase/larkfield"
	"github.com/lincaiyong/larkbase/larksuite/bitable"
	"github.com/lincaiyong/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

const modifiedTimeFieldName = "modified_time"
const updatedTimeFieldName = "updated_time"
const idFieldName = "id"

func (c *Connection[T]) SyncToDatabase(db *gorm.DB, batchSize int) error {
	if c.fieldMap[modifiedTimeFieldName] == nil || c.fieldMap[modifiedTimeFieldName].Type() != "ModifiedTime" {
		return fmt.Errorf("field \"%s\" with ModifiedTime type is required in struct: %s", modifiedTimeFieldName, c.structName)
	}
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
	var filter *bitable.FilterInfo
	if latestModifiedTime != nil {
		field := c.fieldMap[modifiedTimeFieldName].(*larkfield.ModifiedTimeField)
		filter = c.FilterAnd(field.IsGreater(*latestModifiedTime))
	}
	var rawRecords []*T
	err := c.FindAll(&rawRecords, NewFindOption(filter))
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

func (c *Connection[T]) SyncFromApi(getUpdatedFunc func(context.Context, string, time.Time) []map[string]any) error {
	if c.fieldMap[updatedTimeFieldName] == nil || c.fieldMap[updatedTimeFieldName].Type() != "Date" {
		return fmt.Errorf("field \"%s\" with Date type is required in struct: %s", updatedTimeFieldName, c.structName)
	}
	if c.fieldMap[idFieldName] == nil || c.fieldMap[idFieldName].Type() != "Text" {
		return fmt.Errorf("field \"%s\" with Text type is required in struct: %s", idFieldName, c.structName)
	}
	//var latestUpdatedAt time.Time
	//var latest T
	//err := c.Find(&latest, NewFindOption(nil, c.fieldMap[updatedTimeFieldName].(*larkfield.BaseField).Desc()))
	//if err != nil {
	//	if c.IsNotFoundError(err) {
	//		latestUpdatedAt = time.Unix(0, 0)
	//	} else {
	//		log.ErrorLog("fail to find record: %v", err)
	//		return err
	//	}
	//} else {
	//	log.InfoLog("record: %s", c.MarshalIgnoreError(&latest))
	//}
	//
	//memos, err = getUpdatedFunc(ctx, slug, updatedAt)
	//if err != nil {
	//	log.ErrorLog("fail to get updated memo: %v", err)
	//}
	//log.InfoLog("memos count: %d", len(memos))
	//if len(memos) > 0 {
	//	memoMap := make(map[string]*flomo.Memo)
	//	for _, memo := range memos {
	//		memoMap[memo.Slug] = memo
	//	}
	//	conditions := make([]*larkbase.Condition, len(memos))
	//	for i, m := range memos {
	//		conditions[i] = conn.Condition().Id.Is(m.Slug)
	//	}
	//	var existsRecords []*FlomoMemo
	//	err = conn.FindAll(&existsRecords, conn.FilterOr(conditions...))
	//	if len(existsRecords) > 0 {
	//		for _, record := range existsRecords {
	//			m := memoMap[record.Id.StringValue()]
	//			delete(memoMap, m.Slug)
	//			record.Content.SetValue(m.Content)
	//			record.Tags.SetValue(m.Tags)
	//			record.CreatedAt.SetValue(m.CreatedAt)
	//			record.UpdatedAt.SetValue(m.UpdatedAt)
	//		}
	//		err = conn.UpdateAll(existsRecords)
	//		if err != nil {
	//			log.ErrorLog("fail to update records: %v", err)
	//			return
	//		}
	//		log.InfoLog("updated records count: %d", len(existsRecords))
	//	}
	//	if len(memoMap) > 0 {
	//		var records []*FlomoMemo
	//		for _, m := range memoMap {
	//			var record FlomoMemo
	//			record.Id.SetValue(m.Slug)
	//			record.Content.SetValue(m.Content)
	//			record.Tags.SetValue(m.Tags)
	//			record.CreatedAt.SetValue(m.CreatedAt)
	//			record.UpdatedAt.SetValue(m.UpdatedAt)
	//			records = append(records, &record)
	//		}
	//		_, err = conn.CreateAll(records)
	//		if err != nil {
	//			log.ErrorLog("fail to create records: %v", err)
	//			return
	//		}
	//		log.InfoLog("create records count: %d", len(records))
	//	}
	//}
	return nil
}
