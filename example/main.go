package main

import (
	"context"
	"fmt"
	"github.com/lincaiyong/larkbase"
	"os"
	"time"
)

type DemoRecord struct {
	larkbase.Meta `lark:"https://bytedance.larkoffice.com/base/RB31bsA7Pa3f5JsKDlhcoTYdnue?table=tblRyfYXwEhFVX9y"`

	No           larkbase.AutoNumberField   `lark:"no"`
	Name         larkbase.TextField         `lark:"name"`
	Age          larkbase.NumberField       `lark:"age"`
	Date         larkbase.DateField         `lark:"日期"`
	Multi        larkbase.MultiSelectField  `lark:"multi"`
	Single       larkbase.SingleSelectField `lark:"单选"`
	Check        larkbase.CheckboxField     `lark:"check"`
	Link         larkbase.UrlField          `lark:"超链接"`
	Progress     larkbase.ProgressField     `lark:"进度"`
	Email        larkbase.EmailField        `lark:"Email"`
	Code         larkbase.BarcodeField      `lark:"条码"`
	Currency     larkbase.CurrencyField     `lark:"货币"`
	Rating       larkbase.RatingField       `lark:"评分"`
	Lookup       larkbase.LookupField       `lark:"lookup"`
	Formula      larkbase.FormulaField      `lark:"formula"`
	ModifiedTime larkbase.ModifiedTimeField `lark:"modified_time"`
}

var (
	larkAppId     = os.Getenv("LARK_APP_ID")
	larkAppSecret = os.Getenv("LARK_APP_SECRET")
)

func testBatch(conn *larkbase.Connection[DemoRecord]) {
	records := make([]*DemoRecord, 0)
	for i := 0; i < 10; i++ {
		record := &DemoRecord{}
		record.Name.SetValue(fmt.Sprintf("test-%d", i))
		records = append(records, record)
	}
	var err error
	records, err = conn.CreateAll(records)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, record := range records {
		fmt.Println(record.RecordId)
	}
	var results []*DemoRecord
	conditions := make([]*larkbase.Condition, 0)
	for i := range records {
		conditions = append(conditions, conn.Condition().Name.Is(fmt.Sprintf("test-%d", i)))
	}
	err = conn.FindAll(&results, larkbase.NewFindOption(conn.FilterOr(conditions...)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(results)
	time.Sleep(3 * time.Second)
	for i, record := range records {
		record.Age.SetIntValue(20 + i)
	}
	err = conn.UpdateAll(records)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	time.Sleep(3 * time.Second)
	err = conn.DeleteAll(records)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	//fmt.Println(larkfield.TimeToModifiedTime(time.Now()))
	//fmt.Println(larkfield.ModifiedTimeToTime(2509091209).Format(time.DateTime))

	conn, err := larkbase.Connect[DemoRecord](context.Background(), larkAppId, larkAppSecret)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//dsn := fmt.Sprintf(os.Getenv("DATABASE_DSN"))
	//db, err := gorm.Open(mysql.Open(dsn))
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//err = conn.SyncToDatabase(db, 0)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}

	var records []*DemoRecord
	err = conn.FindAll(&records, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(conn.MarshalIgnoreError(records))

	s, err := larkbase.DescribeTable(context.Background(), larkAppId, larkAppSecret, "https://bytedance.larkoffice.com/base/RB31bsA7Pa3f5JsKDlhcoTYdnue?table=tblRyfYXwEhFVX9y")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(s)
	return

	//conn, err := larkbase.Connect[DemoRecord](larkAppId, larkAppSecret)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//var record DemoRecord
	//err = conn.Find(&record, conn.FilterAnd(conn.Condition().Name.IsEmpty()))
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(record)
	//
	//err = conn.CreateView("empty", conn.ViewFilterAnd(conn.Condition().Name.IsNotEmpty()))
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//return
	//
	//testBatch(conn)
	//
	//var r DemoRecord
	//r.Name.SetValue("test")
	//err = conn.Create(&r)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//time.Sleep(3 * time.Second)
	//s, err := conn.MarshalRecord(&r)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(s)
	//
	//r.Name.SetValue("test2")
	//err = conn.Update(&r)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//err = conn.Delete(&r)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//var record DemoRecord
	//err = conn.Find(&record, conn.FilterAnd(conn.Condition().Name.Is("andy")))
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//s, err = conn.MarshalRecord(&record)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(s)
	//
	//record.Age.SetIntValue(123456)
	//record.Date.SetValue(time.Now())
	//err = conn.Update(&record)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//var records []*DemoRecord
	//err = conn.FindAll(&records, conn.FilterAnd(conn.Condition().Name.IsNot("andy")))
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//s, err = conn.MarshalRecords(records)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(s)
	//
	//for _, r := range records {
	//	r.Date.SetValue(time.Now())
	//}
	//err = conn.UpdateAll(records)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
}
