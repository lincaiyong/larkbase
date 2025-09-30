package main

import (
	"fmt"
	"github.com/lincaiyong/larkbase"
	"os"
)

type Demo struct {
	larkbase.Meta `lark:"https://bytedance.larkoffice.com/base/RB31bsA7Pa3f5JsKDlhcoTYdnue?table=tblRyfYXwEhFVX9y"`

	No   larkbase.AutoNumberField `lark:"no"`
	Name larkbase.TextField       `lark:"name"`
	Age  larkbase.NumberField     `lark:"age"`
}

func main() {
	client := larkbase.NewClient(os.Getenv("LARK_APP_ID"), os.Getenv("LARK_APP_SECRET"))

	name, err := client.Table(Demo{}).Name()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(name)

	fields, err := client.Table(Demo{}).Fields()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, field := range fields {
		fmt.Println(field.Name(), field.Type().String())
	}

	count, err := client.Table(Demo{}).Count()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(count)

	records, err := client.Table(Demo{}).Where().Records()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, record := range records {
		fmt.Println(record.Id)
		for _, field := range record.Fields {
			fmt.Println(field.Name(), field.Type().String(), field.Value())
		}
	}
}
