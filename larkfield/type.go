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

const TypeText = 1
const TypeNumber = 2
const TypeSingleSelect = 3
const TypeMultiSelect = 4
const TypeDate = 5
const TypeCheckbox = 7
const TypePerson = 11
const TypePhone = 13
const TypeUrl = 15
const TypeMedia = 17
const TypeSingleLink = 18
const TypeLookup = 19
const TypeFormula = 20
const TypeDuplexLink = 21
const TypeLocation = 22
const TypeGroup = 23
const TypeWorkflow = 24
const TypeCreatedTime = 1001
const TypeModifiedTime = 1002
const TypeCreatePerson = 1003
const TypeModifyPerson = 1004
const TypeAutoNumber = 1005
const TypeButton = 3001

type Type int

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
	case TypePerson:
		return "Person"
	case TypePhone:
		return "Phone"
	case TypeUrl:
		return "Url"
	case TypeMedia:
		return "Media"
	case TypeSingleLink:
		return "SingleLink"
	case TypeLookup:
		return "Lookup"
	case TypeFormula:
		return "Formula"
	case TypeDuplexLink:
		return "DuplexLink"
	case TypeLocation:
		return "Location"
	case TypeGroup:
		return "Group"
	case TypeWorkflow:
		return "Workflow"
	case TypeCreatedTime:
		return "CreatedTime"
	case TypeModifiedTime:
		return "ModifiedTime"
	case TypeCreatePerson:
		return "CreatePerson"
	case TypeModifyPerson:
		return "ModifyPerson"
	case TypeAutoNumber:
		return "AutoNumber"
	case TypeButton:
		return "Button"
	default:
		return "?"
	}
}
