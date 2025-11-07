package larkfield

/*
https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-field/guide#46001acf

1：文本（默认值）、条码（需声明 "ui_type": "Barcode"）、邮箱（需声明"ui_type": "Email")
2：数字（默认值）、进度（需声明 "ui_type": "Progress"）、货币（需声明 "ui_type": "Currency"）、评分（需声明 "ui_type": "Rating")
3：单选
4：多选
5：日期
7：复选框
11：人员
13：电话号码
15：超链接
17：附件
18：单向关联
19：查找引用
20：公式
21：双向关联
22：地理位置
23：群组
24：流程（不支持通过写接口新增或编辑，仅支持读接口）
1001：创建时间
1002：最后更新时间
1003：创建人
1004：修改人
1005：自动编号
3001：按钮（不支持通过写接口新增或编辑，仅支持读接口）
*/

type Type int

const TypeUnknown Type = 0
const TypeText Type = 1
const TypeNumber Type = 2
const TypeSingleSelect Type = 3
const TypeMultiSelect Type = 4
const TypeDate Type = 5
const TypeCheckbox Type = 7
const TypeUrl Type = 15
const TypeLookup = 19
const TypeFormula = 20
const TypeModifiedTime Type = 1002
const TypeAutoNumber Type = 1005

func (t Type) String() string {
	switch t {
	case TypeText:
		return "Text"
	case TypeNumber:
		return "Number"
	case TypeSingleSelect:
		return "SingleSelect"
	case TypeMultiSelect:
		return "MultiSelect"
	case TypeDate:
		return "Date"
	case TypeCheckbox:
		return "Checkbox"
	case TypeUrl:
		return "Url"
	case TypeLookup:
		return "Lookup"
	case TypeFormula:
		return "Formula"
	case TypeModifiedTime:
		return "ModifiedTime"
	case TypeAutoNumber:
		return "AutoNumber"
	default:
		return "?"
	}
}

func (t Type) CreateField(id, name string, type_ Type) Field {
	switch t {
	case TypeText:
		ret := &TextField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	case TypeNumber:
		ret := &NumberField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	case TypeSingleSelect:
		ret := &SingleSelectField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	case TypeMultiSelect:
		ret := &MultiSelectField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	case TypeDate:
		ret := &DateField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	case TypeCheckbox:
		ret := &CheckboxField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	case TypeUrl:
		ret := &UrlField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	case TypeLookup:
		ret := &LookupField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	case TypeFormula:
		ret := &FormulaField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	case TypeModifiedTime:
		ret := &ModifiedTimeField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	case TypeAutoNumber:
		ret := &AutoNumberField{}
		ret.BaseField = NewBaseField(ret, id, name, type_)
		return ret
	default:
		return nil
	}
}

func TypeFromString(s string) Type {
	switch s {
	case "Text":
		return TypeText
	case "Number":
		return TypeNumber
	case "SingleSelect":
		return TypeSingleSelect
	case "MultiSelect":
		return TypeMultiSelect
	case "Date":
		return TypeDate
	case "Checkbox":
		return TypeCheckbox
	case "Url":
		return TypeUrl
	case "Lookup":
		return TypeLookup
	case "Formula":
		return TypeFormula
	case "ModifiedTime":
		return TypeModifiedTime
	case "AutoNumber":
		return TypeAutoNumber
	default:
		return TypeUnknown
	}
}
