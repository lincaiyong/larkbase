package main

import (
	"context"
	"encoding/json"
	"github.com/lincaiyong/larkbase"
	"github.com/lincaiyong/log"
	"os"
)

func main() {
	appId, appSecret := os.Getenv("LARK_APP_ID"), os.Getenv("LARK_APP_SECRET")
	tableUrl := "https://bytedance.larkoffice.com/base/RB31bsA7Pa3f5JsKDlhcoTYdnue?table=tblRyfYXwEhFVX9y"
	conn, err := larkbase.ConnectAny(context.Background(), appId, appSecret, tableUrl)
	if err != nil {
		log.ErrorLog("fail to connect: %v", err)
		return
	}
	var records []*larkbase.AnyRecord
	err = conn.FindAll(&records, nil)
	if err != nil {
		log.ErrorLog("fail to find: %v", err)
		return
	}
	for _, record := range records {
		b, _ := json.MarshalIndent(record.Data, "", "  ")
		log.InfoLog(string(b))
	}
	log.InfoLog("done")
}
