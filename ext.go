package larkbase

import (
	"context"
	"fmt"
	"github.com/lincaiyong/larkbase/larkfield"
	lark "github.com/lincaiyong/larkbase/larksuite"
	"github.com/lincaiyong/larkbase/larksuite/bitable"
	larkcore "github.com/lincaiyong/larkbase/larksuite/core"
	"strings"
)

//	func CreateTable(ctx context.Context, name string) error {
//		client := lark.NewClient(appId, appSecret)
//		req := bitable.NewCreateAppReqBuilder().
//			ReqApp(bitable.NewReqAppBuilder().Name(name).FolderToken(`fldcnqquW1svRIYVT2Np6Iabcef`).Build()).
//			Build()
//		resp, err := client.Bitable.V1.App.Create(context.Background(), req)
//		if err != nil {
//			return err
//		}
//		if !resp.Success() {
//			return fmt.Errorf("error response: %s", larkcore.Prettify(resp.CodeError))
//		}
//		resp.Data.App.Url
//		return nil
//	}

func CreateTable(ctx context.Context, name string, fields []string) (string, error) {
	client := lark.NewClient(appId, appSecret)
	var buildFields []*bitable.AppTableCreateHeader
	for _, field := range fields {
		buildFields = append(buildFields, bitable.NewAppTableCreateHeaderBuilder().FieldName(field).Type(int(larkfield.TypeText)).Build())
	}
	req := bitable.NewCreateAppTableReqBuilder().
		Body(bitable.NewCreateAppTableReqBodyBuilder().
			Table(bitable.NewReqTableBuilder().
				Name(name).
				DefaultViewName(`表格`).
				Fields(buildFields).
				Build()).
			Build()).
		Build()
	resp, err := client.Bitable.V1.AppTable.Create(ctx, req)
	if err != nil {
		return "", err
	}
	if !resp.Success() {
		return "", fmt.Errorf("error response: %s", larkcore.Prettify(resp.CodeError))
	}
	return *resp.Data.TableId, nil
}

func CreateAll(ctx context.Context, url string, data []map[string]string, tosFields string) error {
	tosFieldMap := make(map[string]struct{})
	if tosFields != "" {
		for _, field := range strings.Split(tosFields, ",") {
			tosFieldMap[field] = struct{}{}
		}
	}
	conn, err := ConnectAny(ctx, url)
	if err != nil {
		return err
	}
	var fields []string
	var records []*AnyRecord
	for i, row := range data {
		var r AnyRecord
		r.Data = make(map[string]string)
		for k, v := range row {
			if _, ok := tosFieldMap[k]; ok {
				v, err = tosPutFn(ctx, []byte(v))
				if err != nil {
					return err
				}
				k = fmt.Sprintf("tos_%s", k)
			}
			r.Data[k] = v
			r.Update(k, v)
			if i == 0 {
				fields = append(fields, k)
			}
		}
		records = append(records, &r)
	}
	_, err = conn.CreateAllAny(fields, records)
	if err != nil {
		return err
	}
	return nil
}
