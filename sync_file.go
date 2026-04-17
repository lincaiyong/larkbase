package larkbase

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lincaiyong/larkbase/larkfield"
	"github.com/lincaiyong/larkbase/larksuite/bitable"
	larkcore "github.com/lincaiyong/larkbase/larksuite/core"
	"github.com/lincaiyong/log"
)

type fileCache struct {
	SyncedAt string                      `json:"synced_at"`
	Records  map[string]*fileCacheRecord `json:"records"`
}

type fileCacheRecord struct {
	RecordId         string         `json:"record_id"`
	LastModifiedTime int64          `json:"last_modified_time"`
	Fields           map[string]any `json:"fields"`
}

// SyncToFile 从飞书拉取记录并缓存到本地 JSON 文件。
// fullSync=false: 增量同步（只拉 modified_time > 上次同步时间的记录）
// fullSync=true:  全量同步（拉取所有记录，替换文件内容，可检测飞书侧删除）
// 文件不存在时自动按全量处理。
func (c *Connection[T]) SyncToFile(filePath string) error {
	cache := &fileCache{Records: make(map[string]*fileCacheRecord)}
	larkRecordIds := make(map[string]bool)
	err := queryAllPages(func(pageToken string) (string, error) {
		items, nextPageToken, err := c.queryRawRecordsByPage(nil, pageToken)
		if err != nil {
			return "", err
		}
		for _, item := range items {
			if item.RecordId == nil {
				continue
			}
			recordId := *item.RecordId
			larkRecordIds[recordId] = true

			var modifiedTime int64
			if item.LastModifiedTime != nil {
				modifiedTime = *item.LastModifiedTime
			}

			cache.Records[recordId] = &fileCacheRecord{
				RecordId:         recordId,
				LastModifiedTime: modifiedTime,
				Fields:           item.Fields,
			}
		}
		return nextPageToken, nil
	})
	if err != nil {
		return fmt.Errorf("fail to sync records from lark: %v", err)
	}

	for id := range cache.Records {
		if !larkRecordIds[id] {
			delete(cache.Records, id)
		}
	}
	cache.SyncedAt = larkfield.TimeToBeijingDateTimeStr(time.Now())
	log.InfoLog("sync to file: %s, records: %d", filePath, len(cache.Records))
	return writeFileCache(filePath, cache)
}

// ReadFromFile 从本地 JSON 缓存文件读取所有记录，返回类型安全的 []*T。
func (c *Connection[T]) ReadFromFile(filePath string) ([]*T, error) {
	cache, err := readFileCache(filePath)
	if err != nil {
		return nil, err
	}
	if cache == nil || len(cache.Records) == 0 {
		return nil, nil
	}

	records := make([]*Record, 0, len(cache.Records))
	for _, cached := range cache.Records {
		recordId := cached.RecordId
		lastModifiedTime := cached.LastModifiedTime
		item := &bitable.AppTableRecord{
			RecordId:         &recordId,
			LastModifiedTime: &lastModifiedTime,
			Fields:           cached.Fields,
		}
		record, err := c.parseAppTableRecord(item)
		if err != nil {
			return nil, fmt.Errorf("fail to parse cached record %s: %v", cached.RecordId, err)
		}
		records = append(records, record)
	}

	var result []*T
	if err := c.convertRecordsToStructPtrSlicePtr(records, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Connection[T]) queryRawRecordsByPage(filter *bitable.FilterInfo, pageToken string) ([]*bitable.AppTableRecord, string, error) {
	bodyBuilder := bitable.NewSearchAppTableRecordReqBodyBuilder()
	bodyBuilder.FieldNames(c.fieldNames)
	if filter != nil {
		bodyBuilder.Filter(filter)
	}
	bodyBuilder.AutomaticFields(true)
	req := bitable.NewSearchAppTableRecordReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		PageToken(pageToken).
		PageSize(100).
		Body(bodyBuilder.Build()).
		Build()

	var resp *bitable.SearchAppTableRecordResp
	var err error
	err = c.retry(func() error {
		resp, err = c.client.Bitable.V1.AppTableRecord.Search(c.ctx, req)
		return err
	})
	if err != nil {
		return nil, "", fmt.Errorf("fail to call bitable search table: %v", err)
	}
	if !resp.Success() {
		return nil, "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	if *resp.Data.HasMore {
		return resp.Data.Items, *resp.Data.PageToken, nil
	}
	return resp.Data.Items, "", nil
}

func readFileCache(filePath string) (*fileCache, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("fail to read cache file: %v", err)
	}
	var cache fileCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("fail to parse cache file: %v", err)
	}
	if cache.Records == nil {
		cache.Records = make(map[string]*fileCacheRecord)
	}
	return &cache, nil
}

func writeFileCache(filePath string, cache *fileCache) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("fail to create cache directory: %v", err)
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("fail to marshal cache: %v", err)
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("fail to write cache file: %v", err)
	}
	return nil
}
