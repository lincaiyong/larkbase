package main

import (
	"fmt"
	"github.com/lincaiyong/larkbase"
	"os"
)

type Bar struct {
	larkbase.Meta `lark:"https://bytedance.larkoffice.com/base/RB31bsA7Pa3f5JsKDlhcoTYdnue?table=tblRyfYXwEhFVX9y"`

	No   larkbase.NumberField `lark:"no"`
	Name larkbase.TextField   `lark:"name"`
	Age  larkbase.NumberField `lark:"age"`
}

func main() {
	client := larkbase.NewClient(os.Getenv("LARK_APP_ID"), os.Getenv("LARK_APP_SECRET"))
	var count int
	client.Table(Bar{}).Count(&count)
	if client.Error() != nil {
		fmt.Println(client.Error())
		os.Exit(1)
	}
	fmt.Println(count)
}
