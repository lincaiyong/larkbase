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
	//X          larkbase.CheckboxField     `lark:"公式"`
}

func main() {
	client := larkbase.NewClient(os.Getenv("LARK_APP_ID"), os.Getenv("LARK_APP_SECRET"))
	demo := Demo{}
	demoConn, err := client.Connect(&demo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(demoConn.TableName(), demoConn.TableId())
	for _, field := range demoConn.TableFields() {
		fmt.Println(field.Name(), field.Type().String())
	}

	count, err := demoConn.CountRecords()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(count)

	var records []Demo
	err = demoConn.Where(demo.Name.FilterIs("andy")).QueryRecords(&records)
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
}
