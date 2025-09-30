package larkbase

func (c *Client) Count() (int, error) {
	if err := c.checkCurrent(); err != nil {
		return 0, err
	}
	total := 0
	pageToken := ""
	for {
		var err error
		pageToken, err = c.getTableRecordByPage(pageToken, &total)
		if err != nil {
			return 0, err
		}
		if pageToken == "" {
			break
		}
	}
	return total, nil
}

func (c *Client) Where() *Client {
	if err := c.checkCurrent(); err != nil {
		return c
	}
	return c
}

func (c *Client) Records() ([]*Record, error) {
	if err := c.checkCurrent(); err != nil {
		return nil, err
	}
	return nil, nil
}

//}
//
//func (c *Client) searchRecordsByPage(fields map[string]IField, pageToken string, pageSize int, filters ...*larkbitable.Condition) ([]*Record, string) {
//	if fields == nil {
//		fields = c.GetFields(nil)
//	}
//	fieldNames := sortedMapKeys(fields)
//	bodyBuilder := larkbitable.NewSearchAppTableRecordReqBodyBuilder()
//	bodyBuilder.FieldNames(fieldNames)
//	if len(filters) > 0 {
//		bodyBuilder.Filter(larkbitable.NewFilterInfoBuilder().
//			Conjunction(`and`).
//			Conditions(filters).
//			Build())
//	}
//	bodyBuilder.AutomaticFields(true)
//	req := larkbitable.NewSearchAppTableRecordReqBuilder().
//		AppToken(c.appToken).
//		TableId(c.table.Id).
//		PageToken(pageToken).
//		PageSize(pageSize).
//		Body(bodyBuilder.Build()).Build()
//	resp, err := c.Bitable.V1.AppTableRecord.Search(context.Background(), req)
//	if err != nil {
//		log.FatalLog("fail to call bitable search table: %v", err)
//	}
//	if !resp.Success() {
//		log.FatalLog("unexpected response error: %s", larkcore.Prettify(resp.CodeError))
//	}
//	result := make([]*Record, 0)
//	for _, item := range resp.Data.Items {
//		record := &Record{
//			Id:     *item.RecordId,
//			Fields: make(map[string]IField),
//		}
//		for name, fi := range item.Fields {
//			f := fields[name]
//			field := f.Parse(fi)
//			record.Fields[name] = field
//		}
//		result = append(result, record)
//	}
//	if *resp.Data.HasMore {
//		return result, *resp.Data.PageToken
//	}
//	return result, ""
//}
