package larkfield

// https://open.larkoffice.com/document/server-docs/docs/bitable-v1/app-table-field/guide

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
