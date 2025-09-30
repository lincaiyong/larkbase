package larkbase

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func NewClient(appId, appSecret string) *Client {
	return &Client{
		Client:    lark.NewClient(appId, appSecret),
		appTables: map[string]map[string]string{},
	}
}

type Client struct {
	*lark.Client
	appTables map[string]map[string]string // key is app token, table id

	current *ClientContext
}

type ClientContext struct {
	appToken string
	table    *Table
	error    error
	filters  []*larkbitable.Condition
}

func (c *Client) checkCurrent() error {
	if c.current == nil {
		return fmt.Errorf("current table is not set")
	}
	return c.current.error
}

func (c *Client) failCurrent(msg string, args ...any) {
	c.current.error = fmt.Errorf(msg, args...)
}

func (c *Client) connectApp(appToken string) error {
	if c.appTables[appToken] != nil {
		return nil
	}
	c.appTables[appToken] = map[string]string{}
	pageToken := ""
	for {
		var err error
		pageToken, err = c.getAppTablesByPage(appToken, pageToken)
		if err != nil {
			return err
		}
		if pageToken == "" {
			break
		}
	}
	return nil
}

func (c *Client) getAppTablesByPage(appToken, pageToken string) (string, error) {
	const pageSize = 10
	req := larkbitable.NewListAppTableReqBuilder().
		AppToken(appToken).
		PageSize(pageSize).
		PageToken(pageToken).
		Build()
	resp, err := c.Bitable.V1.AppTable.List(context.Background(), req)
	if err != nil {
		return "", err
	}
	if !resp.Success() {
		return "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	data := resp.Data
	tables := c.appTables[appToken]
	for _, item := range data.Items {
		tables[*item.TableId] = *item.Name
	}
	if *data.HasMore {
		return *data.PageToken, nil
	}
	return "", nil
}

func (c *Client) getTableRecordByPage(pageToken string, total *int) (string, error) {
	const pageSize = 100
	req := larkbitable.NewSearchAppTableRecordReqBuilder().
		AppToken(c.current.appToken).
		TableId(c.current.table.id).
		PageToken(pageToken).
		PageSize(pageSize).
		Body(larkbitable.NewSearchAppTableRecordReqBodyBuilder().Build()).Build()
	resp, err := c.Bitable.V1.AppTableRecord.Search(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("fail to call bitable search table: %v", err)
	}
	if !resp.Success() {
		return "", fmt.Errorf("get response with error: %s", larkcore.Prettify(resp.CodeError))
	}
	*total += *resp.Data.Total
	if *resp.Data.HasMore {
		return *resp.Data.PageToken, nil
	}
	return "", nil
}
