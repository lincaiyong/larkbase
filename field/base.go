package field

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type BaseField struct {
	name  string
	type_ string
	value string
	dirty bool
}

func (f *BaseField) Name() string {
	return f.name
}

func (f *BaseField) SetName(name string) {
	f.name = name
}

func (f *BaseField) Type() string {
	return f.type_
}

func (f *BaseField) SetType(type_ string) {
	f.type_ = type_
}

func (f *BaseField) Value() string {
	return f.value
}

func (f *BaseField) SetValue(value string) {
	if f.value != value {
		f.dirty = true
		f.value = value
	}
}

func (f *BaseField) SetValueNoDirty(value string) {
	f.value = value
}

func (f *BaseField) Dirty() bool {
	return f.dirty
}

func (f *BaseField) Fork() Field {
	return &BaseField{
		name:  f.name,
		type_: f.type_,
		value: f.value,
	}
}

func (f *BaseField) ParseFromLarkSuite(v any) {
	switch f.Type() {
	case "Text":
		if v2, ok := v.([]any); ok && len(v2) == 1 {
			if v3, ok2 := v2[0].(map[string]any); ok2 {
				f.value = v3["text"].(string)
			}
		}
	case "Number":
		f.value = fmt.Sprintf("%g", v.(float64))
	case "SingleSelect":
		f.value = v.(string)
	case "MultiSelect":
		if vv, ok := v.([]any); ok {
			items := make([]string, 0)
			for _, v2 := range vv {
				items = append(items, v2.(string))
			}
			b, _ := json.Marshal(items)
			f.value = string(b)
		}
	case "Date":
		f.value = unixSecondsToBeijingDateTimeStr(int64(v.(float64) / 1000))
	case "Checkbox":
		f.value = map[bool]string{true: "1", false: ""}[v.(bool)]
	case "Person":
		if vv, ok := v.([]any); ok {
			items := make([]string, 0)
			for _, v2 := range vv {
				if v3, ok2 := v2.(map[string]any); ok2 {
					items = append(items, v3["name"].(string))
				}
			}
			f.value = strings.Join(items, ",")
		}
	case "Phone":
		f.value = v.(string)
	case "Url":
		if v2, ok := v.(map[string]any); ok {
			f.value = v2["link"].(string)
		}
	case "Media":
		if vv, ok := v.([]any); ok {
			items := make([]string, 0)
			for _, v2 := range vv {
				if v3, ok2 := v2.(map[string]any); ok2 {
					items = append(items, v3["file_token"].(string))
				}
			}
			f.value = strings.Join(items, ",")
		}
	case "SingleLink":
		f.value = v.(string)
	case "Lookup":
		f.value = v.(string)
	case "Formula":
		f.value = v.(string)
	case "DuplexLink":
		f.value = v.(string)
	case "Location":
		f.value = v.(string)
	case "Group":
		f.value = v.(string)
	case "Workflow":
		f.value = v.(string)
	case "CreatedTime":
		f.value = unixSecondsToBeijingDateTimeStr(int64(v.(float64) / 1000))
	case "ModifiedTime":
		f.value = unixSecondsToBeijingDateTimeStr(int64(v.(float64) / 1000))
	case "CreatePerson":
		f.value = v.(string)
	case "ModifyPerson":
		f.value = v.(string)
	case "AutoNumber":
		f.value = v.(string)
	case "Button":
		f.value = v.(string)
	}
}

// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/bitable-record-data-structure-overview

func (f *BaseField) BuildForLarkSuite() (ret any, err error) {
	switch f.Type() {
	case "AutoNumber":
		ret = f.value
	case "Checkbox":
		ret = map[string]bool{"1": true}[f.value]
	case "Date":
		var us int64
		us, err = beijingDateTimeStrToUnixSeconds(f.value)
		ret = us * 1000
	case "Media":
		items := strings.Split(f.value, ",")
		tmp := make([]any, len(items))
		for i, sel := range items {
			tmp[i] = map[string]string{
				"file_token": sel,
			}
		}
		ret = tmp
	case "MultiSelect":
		var items []string
		err = json.Unmarshal([]byte(f.value), &items)
		tmp := make([]any, len(items))
		for i, sel := range items {
			tmp[i] = sel
		}
		ret = tmp
	case "Number":
		var ff float64
		ff, err = strconv.ParseFloat(f.value, 64)
		ret = ff
	case "SingleSelect", "Text":
		ret = f.value
	case "ModifiedTime", "CreatedTime":
		var t time.Time
		t, err = beijingDateTimeStrToTime(f.value)
		ret = &t
	case "Url":
		ret = map[string]any{
			"text": f.value,
			"link": f.value,
		}
	default:
		err = fmt.Errorf("unsupported field type to build: %s", f.Type())
	}
	return
}
