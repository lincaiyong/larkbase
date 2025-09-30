package larkbase

/*
https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-field/guide

文本：原值展示，不支持 markdown 语法，"xxx"
单选：填写选项值，对于新的选项值，将会创建一个新的选项，"xxx"
数字：填写数字格式的值，100
日期：填写毫秒级时间戳，1674206443000
多选：填写多个选项值，对于新的选项值，将会创建一个新的选项。如果填写多个相同的新选项值，将会创建多个相同的选项，[]any{"xx"}
复选框：填写 true 或 false
超链接：参考以下示例，text 为文本值，link 为 URL 链接
	map[string]interface{}{`text`: `飞书多维表格官网`, `link`: `https://www.feishu.cn/product/base`}
附件：参考 [{"file_token": "Vl3FbVkvnowlgpxpqsAbBrtFcrd"}]
*/

const FieldTypeText = 1
const FieldTypeNumber = 2
const FieldTypeSingleSelect = 3
const FieldTypeMultiSelect = 4
const FieldTypeDate = 5
const FieldTypeCheckbox = 7
const FieldTypePerson = 11
const FieldTypeUrl = 15
const FieldTypeMedia = 17
const FieldTypeUpdatedTime = 1002

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
	case FieldTypeUrl:
		return "Url"
	case FieldTypeMedia:
		return "Media"
	default:
		return "UpdatedTime"
	}
}

type IField interface {
	Id() string
	Name() string
	IsPrimary() bool
	Type() FieldType
	Value() string
	SetValue(v any) error
	Build() any
	Parse(v any) IField
}

type BaseField struct {
	id        string
	name      string
	isPrimary bool
}

func (f *BaseField) Id() string {
	return f.id
}

func (f *BaseField) Name() string {
	return f.name
}

func (f *BaseField) IsPrimary() bool {
	return f.isPrimary
}

func (f *BaseField) Type() FieldType {
	return 0
}

func (f *BaseField) Value() string {
	return ""
}

func (f *BaseField) SetValue(_ any) error {
	return nil
}

func (f *BaseField) Build() any {
	return ""
}

func (f *BaseField) Parse(v any) IField {
	return nil
}
