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
}

func (c *Client) connectLarkApp(appToken string) error {
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
