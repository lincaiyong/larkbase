package larkbase

import (
	"encoding/json"
	"fmt"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/lincaiyong/log"
	"strconv"
	"strings"
	"time"
)

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-field/guide

const FieldTypeText = 1
const FieldTypeNumber = 2
const FieldTypeSingleSelect = 3
const FieldTypeMultiSelect = 4
const FieldTypeDate = 5
const FieldTypeCheckbox = 7
const FieldTypePerson = 11
const FieldTypePhone = 13
const FieldTypeUrl = 15
const FieldTypeMedia = 17
const FieldTypeSingleLink = 18
const FieldTypeLookup = 19
const FieldTypeFormula = 20
const FieldTypeDuplexLink = 21
const FieldTypeLocation = 22
const FieldTypeGroup = 23
const FieldTypeWorkflow = 24
const FieldTypeCreatedTime = 1001
const FieldTypeModifiedTime = 1002
const FieldTypeCreatePerson = 1003
const FieldTypeModifyPerson = 1004
const FieldTypeAutoNumber = 1005
const FieldTypeButton = 3001

type FieldType int

func (t FieldType) String() string {
	switch t {
	case FieldTypeText:
		return "Text"
	case FieldTypeNumber:
		return "Number"
	case FieldTypeSingleSelect:
		return "SingleSelect"
	case FieldTypeMultiSelect:
		return "MultiSelect"
	case FieldTypeDate:
		return "Date"
	case FieldTypeCheckbox:
		return "Checkbox"
	case FieldTypePerson:
		return "Person"
	case FieldTypePhone:
		return "Phone"
	case FieldTypeUrl:
		return "Url"
	case FieldTypeMedia:
		return "Media"
	case FieldTypeSingleLink:
		return "SingleLink"
	case FieldTypeLookup:
		return "Lookup"
	case FieldTypeFormula:
		return "Formula"
	case FieldTypeDuplexLink:
		return "DuplexLink"
	case FieldTypeLocation:
		return "Location"
	case FieldTypeGroup:
		return "Group"
	case FieldTypeWorkflow:
		return "Workflow"
	case FieldTypeCreatedTime:
		return "CreatedTime"
	case FieldTypeModifiedTime:
		return "ModifiedTime"
	case FieldTypeCreatePerson:
		return "CreatePerson"
	case FieldTypeModifyPerson:
		return "ModifyPerson"
	case FieldTypeAutoNumber:
		return "AutoNumber"
	case FieldTypeButton:
		return "Button"
	default:
		return "?"
	}
}

type Field struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (f *Field) parseValue(v any) {
	switch f.Type {
	case "Text":
		if v2, ok := v.([]any); ok && len(v2) == 1 {
			if v3, ok2 := v2[0].(map[string]any); ok2 {
				f.Value = v3["text"].(string)
			}
		}
	case "Number":
		f.Value = fmt.Sprintf("%g", v.(float64))
	case "SingleSelect":
		f.Value = v.(string)
	case "MultiSelect":
		if vv, ok := v.([]any); ok {
			items := make([]string, 0)
			for _, v2 := range vv {
				items = append(items, v2.(string))
			}
			b, _ := json.Marshal(items)
			f.Value = string(b)
		}
	case "Date":
		f.Value = unixSecondsToBeijingDateTimeStr(int64(v.(float64) / 1000))
	case "Checkbox":
		f.Value = map[bool]string{true: "1", false: ""}[v.(bool)]
	case "Person":
		if vv, ok := v.([]any); ok {
			items := make([]string, 0)
			for _, v2 := range vv {
				if v3, ok2 := v2.(map[string]any); ok2 {
					items = append(items, v3["name"].(string))
				}
			}
			f.Value = strings.Join(items, ",")
		}
	case "Phone":
		f.Value = v.(string)
	case "Url":
		if v2, ok := v.(map[string]any); ok {
			f.Value = v2["link"].(string)
		}
	case "Media":
		if vv, ok := v.([]any); ok {
			items := make([]string, 0)
			for _, v2 := range vv {
				if v3, ok2 := v2.(map[string]any); ok2 {
					items = append(items, v3["file_token"].(string))
				}
			}
			f.Value = strings.Join(items, ",")
		}
	case "SingleLink":
		f.Value = v.(string)
	case "Lookup":
		f.Value = v.(string)
	case "Formula":
		f.Value = v.(string)
	case "DuplexLink":
		f.Value = v.(string)
	case "Location":
		f.Value = v.(string)
	case "Group":
		f.Value = v.(string)
	case "Workflow":
		f.Value = v.(string)
	case "CreatedTime":
		f.Value = unixSecondsToBeijingDateTimeStr(int64(v.(float64) / 1000))
	case "ModifiedTime":
		f.Value = unixSecondsToBeijingDateTimeStr(int64(v.(float64) / 1000))
	case "CreatePerson":
		f.Value = v.(string)
	case "ModifyPerson":
		f.Value = v.(string)
	case "AutoNumber":
		f.Value = v.(string)
	case "Button":
		f.Value = v.(string)
	}
}

// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/bitable-record-data-structure-overview
func (f *Field) buildForLarkSuite() any {
	switch f.Type {
	case "AutoNumber":
		return f.Value
	case "Checkbox":
		return map[string]bool{"1": true}[f.Value]
	case "Date":
		us, _ := beijingDateTimeStrToUnixSeconds(f.Value)
		return us * 1000
	case "Media":
		items := strings.Split(f.Value, ",")
		tmp := make([]any, len(items))
		for i, sel := range items {
			tmp[i] = map[string]string{
				"file_token": sel,
			}
		}
		return tmp
	case "MultiSelect":
		var items []string
		_ = json.Unmarshal([]byte(f.Value), &items)
		tmp := make([]any, len(items))
		for i, sel := range items {
			tmp[i] = sel
		}
		return tmp
	case "Number":
		ff, _ := strconv.ParseFloat(f.Value, 64)
		return ff
	case "SingleSelect", "Text":
		return f.Value
	case "ModifiedTime", "CreatedTime":
		t, _ := beijingDateTimeStrToTime(f.Value)
		return &t
	case "Url":
		return map[string]any{
			"text": f.Value,
			"link": f.Value,
		}
	default:
		log.WarnLog("unsupported field type to build: %s, ignored", f.Type)
		return nil
	}
}

type TextField Field

func (f *TextField) FilterIs(value string) *larkbitable.Condition {
	return filterIs(f.Name, value)
}
func (f *TextField) FilterIsNot(value string) *larkbitable.Condition {
	return filterIsNot(f.Name, value)
}
func (f *TextField) FilterContains(value string) *larkbitable.Condition {
	return filterContains(f.Name, value)
}
func (f *TextField) FilterDoesNotContains(value string) *larkbitable.Condition {
	return filterDoesNotContains(f.Name, value)
}
func (f *TextField) FilterIsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.Name)
}
func (f *TextField) FilterIsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.Name)
}

type NumberField Field

func (f *NumberField) FilterIs(value string) *larkbitable.Condition {
	return filterIs(f.Name, value)
}
func (f *NumberField) FilterIsNot(value string) *larkbitable.Condition {
	return filterIsNot(f.Name, value)
}
func (f *NumberField) FilterIsGreater(value string) *larkbitable.Condition {
	return filterIsGreater(f.Name, value)
}
func (f *NumberField) FilterIsGreaterEqual(value string) *larkbitable.Condition {
	return filterIsGreaterEqual(f.Name, value)
}
func (f *NumberField) FilterIsLess(value string) *larkbitable.Condition {
	return filterIsLess(f.Name, value)
}
func (f *NumberField) FilterIsLessEqual(value string) *larkbitable.Condition {
	return filterIsLessEqual(f.Name, value)
}
func (f *NumberField) FilterIsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.Name)
}
func (f *NumberField) FilterIsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.Name)
}

type SingleSelectField TextField

type MultiSelectField TextField

type DateField Field

func (f *DateField) FilterIsToday() *larkbitable.Condition {
	return filterDateIsToday(f.Name)
}
func (f *DateField) FilterIsTomorrow() *larkbitable.Condition {
	return filterDateIsTomorrow(f.Name)
}
func (f *DateField) FilterIsYesterday() *larkbitable.Condition {
	return filterDateIsYesterday(f.Name)
}
func (f *DateField) FilterIs(time time.Time) *larkbitable.Condition {
	return filterDateIs(f.Name, time)
}
func (f *DateField) FilterIsGreaterThanToday() *larkbitable.Condition {
	return filterDateIsGreaterThanToday(f.Name)
}
func (f *DateField) FilterIsGreaterThanTomorrow() *larkbitable.Condition {
	return filterDateIsGreaterThanTomorrow(f.Name)
}
func (f *DateField) FilterIsGreaterThanYesterday() *larkbitable.Condition {
	return filterDateIsGreaterThanYesterday(f.Name)
}
func (f *DateField) FilterIsGreater(time time.Time) *larkbitable.Condition {
	return filterDateIsGreater(f.Name, time)
}
func (f *DateField) FilterIsLessThanToday() *larkbitable.Condition {
	return filterDateIsLessThanToday(f.Name)
}
func (f *DateField) FilterIsLessThanTomorrow() *larkbitable.Condition {
	return filterDateIsLessThanTomorrow(f.Name)
}
func (f *DateField) FilterIsLessThanYesterday() *larkbitable.Condition {
	return filterDateIsLessThanYesterday(f.Name)
}

func (f *DateField) FilterIsLess(time time.Time) *larkbitable.Condition {
	return filterDateIsLess(f.Name, time)
}
func (f *DateField) FilterIsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.Name)
}
func (f *DateField) FilterIsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.Name)
}

func (f *DateField) FilterIsCurrentWeek() *larkbitable.Condition {
	return filterDateIsCurrentWeek(f.Name)
}

func (f *DateField) FilterIsLastWeek() *larkbitable.Condition {
	return filterDateIsLastWeek(f.Name)
}
func (f *DateField) FilterIsCurrentMonth() *larkbitable.Condition {
	return filterDateIsCurrentMonth(f.Name)
}

func (f *DateField) FilterIsLastMonth() *larkbitable.Condition {
	return filterDateIsLastMonth(f.Name)
}

func (f *DateField) FilterIsTheLastWeek() *larkbitable.Condition {
	return filterDateIsTheLastWeek(f.Name)
}

func (f *DateField) FilterTheNextWeek() *larkbitable.Condition {
	return filterDateTheNextWeek(f.Name)
}

func (f *DateField) FilterIsTheLastMonth() *larkbitable.Condition {
	return filterDateIsTheLastMonth(f.Name)
}

func (f *DateField) FilterTheNextMonth() *larkbitable.Condition {
	return filterDateTheNextMonth(f.Name)
}

type CheckboxField Field

func (f *CheckboxField) FilterIs(value bool) *larkbitable.Condition {
	return filterIs(f.Name, map[bool]string{true: "1", false: ""}[value])
}

type PersonField TextField

type PhoneField TextField

type UrlField TextField

type MediaField Field

func (f *MediaField) FilterIsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.Name)
}
func (f *MediaField) FilterIsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.Name)
}

type SingleLinkField Field
type LookupField Field
type FormulaField Field
type DuplexLinkField Field
type LocationField Field
type GroupField Field
type WorkflowField Field
type CreatedTimeField DateField
type ModifiedTimeField DateField
type CreatePersonField PersonField
type ModifyPersonField PersonField
type AutoNumberField NumberField
type ButtonField Field
