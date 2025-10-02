package main

import (
	"fmt"
	"github.com/lincaiyong/larkbase"
	"os"
)

type DemoRecord struct {
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
	Formula    larkbase.FormulaField      `lark:"公式"`
	Link       larkbase.UrlField          `lark:"超链接"`
	Phone      larkbase.PhoneField        `lark:"电话"`
	Progress   larkbase.NumberField       `lark:"进度"`
	Email      larkbase.TextField         `lark:"Email"`
	Code       larkbase.TextField         `lark:"条码"`
	Test       larkbase.NumberField       `lark:"货币"`
}

var (
	larkAppId     = os.Getenv("LARK_APP_ID")
	larkAppSecret = os.Getenv("LARK_APP_SECRET")
)

func main() {
	conn, err := larkbase.Connect[DemoRecord](larkAppId, larkAppSecret)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var record DemoRecord
	err = conn.FindOne(&record, conn.Filter().Name.Is("andy"))
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

	record.Age.SetIntValue(123456)
	err = conn.UpdateOne(&record)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var records []DemoRecord
	err = conn.FindAll(&records, conn.Filter().Name.IsNot("andy"))
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
