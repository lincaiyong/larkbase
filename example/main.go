package main

import (
	"fmt"
	"github.com/lincaiyong/larkbase"
	"os"
)

type Demo struct {
	larkbase.Meta `lark:"https://bytedance.larkoffice.com/base/RB31bsA7Pa3f5JsKDlhcoTYdnue?table=tblRyfYXwEhFVX9y"`

	No         larkbase.AutoNumberField   `lark:"no"`
	Name       larkbase.TextField         `lark:"name"`
	Age        larkbase.NumberField       `lark:"age"`
	Attachment larkbase.MediaField        `lark:"附件"`
	Date       larkbase.DateField         `lark:"日期"`
	Multi      larkbase.MultiSelectField  `lark:"multi"`
	Single     larkbase.SingleSelectField `lark:"单选"`
	Person     larkbase.PersonField       `lark:"人员"`
	Check      larkbase.CheckboxField     `lark:"check"`
}

func main() {
	demo := &Demo{}
	conn, err := larkbase.Connect(os.Getenv("LARK_APP_ID"), os.Getenv("LARK_APP_SECRET"), demo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	count, err := conn.CountRecords()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(count)

	var records []Demo
	err = conn.QueryRecords(&records, demo.Name.FilterIs("andy"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s, err := larkbase.Marshal(records)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(s)

	err = conn.QueryRecords(&records, demo.Name.FilterIsNot("andy"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s, err = larkbase.Marshal(records)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(s)
}
