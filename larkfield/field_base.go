package larkfield

import (
	"encoding/json"
	larkbitable "github.com/lincaiyong/larkbase/larksuite/service/bitable/v1"
	"time"
)

func NewBaseField(id, name, type_ string) *BaseField {
	return &BaseField{
		id:    id,
		name:  name,
		type_: type_,
	}
}

type BaseField struct {
	id    string
	name  string
	type_ string
	value any
	dirty bool
}

func (f *BaseField) Id() string {
	return f.id
}

func (f *BaseField) SetId(v string) {
	f.id = v
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

func (f *BaseField) UnderlayValue() any {
	return f.value
}

func (f *BaseField) SetUnderlayValue(value any) {
	old := f.StringValue()
	new_ := (&BaseField{value: value}).StringValue()
	if old != new_ {
		f.dirty = true
		f.value = value
	}
}

func (f *BaseField) SetUnderlayValueNoDirty(value any) {
	f.value = value
}

func (f *BaseField) Dirty() bool {
	return f.dirty
}

func (f *BaseField) SetDirty(v bool) {
	f.dirty = v
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	if v, ok := value.(string); ok {
		return v
	}
	if v, ok := value.(time.Time); ok {
		return TimeToBeijingDateTimeStr(v)
	}
	b, _ := json.Marshal(value)
	return string(b)
}

func (f *BaseField) StringValue() string {
	return stringValue(f.value)
}

func (f *BaseField) Fork() Field {
	panic("should not happen")
}

func (f *BaseField) Parse(_ any) error {
	panic("should not happen")
	return nil
}

func (f *BaseField) Build() any {
	panic("should not happen")
}

func (f *BaseField) Asc() *larkbitable.Sort {
	builder := &larkbitable.SortBuilder{}
	return builder.FieldName(f.name).Build()
}

func (f *BaseField) Desc() *larkbitable.Sort {
	builder := &larkbitable.SortBuilder{}
	return builder.FieldName(f.name).Desc(true).Build()
}

type HackBaseField BaseField

func (f HackBaseField) Id() string {
	return f.id
}

func (f HackBaseField) Name() string {
	return f.name
}

func (f HackBaseField) Type() string {
	return f.type_
}

func (f HackBaseField) Value() any {
	return f.value
}

func (f HackBaseField) StringValue() string {
	return stringValue(f.value)
}

func (f HackBaseField) Dirty() bool {
	return f.dirty
}
