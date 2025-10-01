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

var (
	larkAppId     = os.Getenv("LARK_APP_ID")
	larkAppSecret = os.Getenv("LARK_APP_SECRET")
)

func main() {
	demo := &Demo{}
	conn, err := larkbase.Connect(larkAppId, larkAppSecret, demo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var record Demo
	err = conn.FindOne(&record, demo.Name.FilterIs("andy"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s, err := larkbase.Marshal(record)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(s)

	record.Age.Value = "123456"
	err = conn.UpdateOne(&record)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var records []Demo
	err = conn.FindAll(&records, demo.Name.FilterIsNot("andy"))
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
