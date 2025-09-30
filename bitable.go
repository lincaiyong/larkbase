package larkbase

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/bitable-overview
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/search

//func ConnectByUrl(url string) *Client {
//	re := regexp.MustCompile(`^https://bytedance\.larkoffice\.com/base/(\w+)\?table=(\w+)`)
//	ret := re.FindStringSubmatch(url)
//	if len(ret) != 3 {
//		log.FatalLog("invalid bitable url: %s", url)
//	}
//	appToken := ret[1]
//	tableId := ret[2]
//	client := &Client{Client: lark.NewClient(gAppId, gAppSecret), appToken: appToken}
//	tables := client.listTables()
//	table, ok := tables[tableId]
//	if !ok {
//		log.FatalLog("fail to connect to table: '%s' not found", tableId)
//	}
//	client.table = table
//	return client
//}

//func (c *Client) SearchRecords(fields map[string]IField, limit int, filters ...*larkbitable.Condition) []*Record {
//	if limit == -1 {
//		limit = 5000
//	}
//	records := make([]*Record, 0)
//	var pageToken string
//	for {
//		var tmp []*Record
//		pageSize := limit
//		if pageSize > 100 {
//			pageSize = 100
//			limit -= 100
//		}
//		tmp, pageToken = c.searchRecords(fields, pageToken, pageSize, filters...)
//		records = append(records, tmp...)
//		if pageToken == "" || len(records) >= limit {
//			break
//		}
//	}
//	return records
//}
//
//func sortedMapKeys(m interface{}) []string {
//	v := reflect.ValueOf(m)
//	ret := make([]string, 0)
//	if v.Kind() == reflect.Map {
//		for _, val := range v.MapKeys() {
//			ret = append(ret, val.String())
//		}
//	}
//	sort.Strings(ret)
//	return ret
//}

//
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

//
//func (c *Client) SyncLarkToDb() {
//	table := c.name
//	fields := c.GetFields(map[string]int{
//		"updated_time": FieldTypeUpdatedTime,
//	})
//	query := buildSql(table, fields)
//	db := mysqldb.ConnectBitable()
//	r := db.Exec(query)
//	if r.Error != nil {
//		log.FatalLog("fail to execute query: %s, %v", query, r.Error)
//	}
//	var count64 int64
//	r = db.Table(table).Count(&count64)
//	if r.Error != nil {
//		log.FatalLog("fail to count table: %s, %v", table, r.Error)
//	}
//	var filters []*larkbitable.Condition
//	if count64 > 0 {
//		query = fmt.Sprintf(`SELECT _updated_time FROM %s ORDER BY _updated_time DESC LIMIT 1;`, table)
//		var updatedTime time.Time
//		r = db.Raw(query).Scan(&updatedTime)
//		if r.Error != nil {
//			log.FatalLog("fail to execute query: %s, %v", query, r.Error)
//		}
//		updatedTimeStr := util.TimestampToDateTimeStr(updatedTime.Unix())
//		filter := FilterDateIsGreater(fields["updated_time"], updatedTimeStr)
//		filters = append(filters, filter)
//	}
//
//	records := c.SearchRecords(nil, -1, filters...)
//	cli.Info("search record results: %d", len(records))
//	if len(records) == 0 {
//		return
//	}
//	for _, record := range records {
//		query = fmt.Sprintf("INSERT IGNORE INTO %s (_record_id) VALUES (?)", table)
//		r = db.Exec(query, record.Id)
//		if r.Error != nil {
//			log.FatalLog("fail to execute query: %s, %v", query, r.Error)
//		}
//		toUpdate := make(map[string]any)
//		for _, field := range record.Fields {
//			toUpdate[field.Name()] = field.Value()
//		}
//		r = db.Table(table).Where("_record_id = ?", record.Id).Updates(toUpdate)
//		if r.Error != nil {
//			log.FatalLog("fail to execute query: %v", r.Error)
//		}
//	}
//}

//func buildSql(table string, fields map[string]IField) string {
//	var fieldsSql []string
//	for _, name := range util.SortedMapKeys(fields) {
//		field := fields[name]
//		fieldSql := fmt.Sprintf("%s VARCHAR(1024)", field.Name())
//		fieldsSql = append(fieldsSql, fieldSql)
//	}
//	ret := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
//_record_id VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT '',
//_created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
//_updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
//index idx_record_id (_record_id),
//%s
//)`,
//		table, strings.Join(fieldsSql, ",\n"))
//	return ret
//}
//
//func (c *Client) AddFields(fields map[string]int) {
//	for k, v := range fields {
//		builder := larkbitable.NewCreateAppTableFieldReqBuilder().
//			AppToken(c.appToken).
//			TableId(c.table.Id).
//			AppTableField(larkbitable.NewAppTableFieldBuilder().
//				FieldName(k).Type(v).Build())
//		req := builder.Build()
//		resp, err := c.Bitable.V1.AppTableField.Create(context.Background(), req)
//		if err != nil {
//			log.FatalLog("fail to add table field: %v", err)
//		}
//		if !resp.Success() {
//			log.FatalLog("fail to add table field: %v", resp.Error())
//		}
//	}
//}
//
//func (c *Client) DeleteFields(fields map[string]IField) {
//	for _, v := range fields {
//		if v.IsPrimary() {
//			continue
//		}
//		builder := larkbitable.NewDeleteAppTableFieldReqBuilder().
//			AppToken(c.appToken).
//			TableId(c.table.Id).
//			FieldId(v.Id())
//		req := builder.Build()
//		resp, err := c.Bitable.V1.AppTableField.Delete(context.Background(), req)
//		if err != nil {
//			log.FatalLog("fail to delete table field: %v", err)
//		}
//		if !resp.Success() {
//			log.FatalLog("fail to delete table field: %v", resp.Error())
//		}
//	}
//}
//
//func (c *Client) InitTable(fields map[string]int) {
//	oldFields := c.GetFields(nil)
//	c.DeleteFields(oldFields)
//	c.AddFields(fields)
//	records := c.SearchRecords(nil, -1)
//	c.DeleteRecords(records)
//}
//
//func DownloadMedia(fileToken string) []byte {
//	session := dbutil.QuerySecret("lark_session_id")
//	url := fmt.Sprintf("https://internal-api-drive-stream.larkoffice.com/space/api/box/stream/download/all/%s/", fileToken)
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		log.FatalLog("fail to create request: %v", err)
//	}
//	req.Header.Add("Cookie", "session="+session)
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		log.FatalLog("fail to execute request: %v", err)
//	}
//	defer func() { _ = resp.Body.Close() }()
//	if resp.StatusCode != 200 {
//		log.FatalLog("fail to execute request: %v", resp.Status)
//	}
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.FatalLog("fail to read body: %v", err)
//	}
//	return body
//}

/*
	GET  HTTP/1.1
	Host:
	Cookie: session=XN0YXJ0-d81g2e71-7cf3-45c9-8424-d1040d2405a5-WVuZA
*/
