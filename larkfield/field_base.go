package larkfield

import (
	"encoding/json"
	"github.com/lincaiyong/larkbase/larksuite/bitable"
	"time"
)

func NewBaseField(self Field, id, name string, type_ Type) BaseField {
	return BaseField{
		self:    self,
		id:      id,
		name:    name,
		typeStr: type_.String(),
		type_:   type_,
	}
}

type BaseField struct {
	self    Field
	id      string
	name    string
	typeStr string
	type_   Type
	value   any
	dirty   bool
}

func (f *BaseField) SetSelf(field Field) {
	f.self = field
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

func (f *BaseField) TypeStr() string {
	return f.typeStr
}

func (f *BaseField) Type() Type {
	return f.type_
}

func (f *BaseField) SetType(t Type) {
	f.type_ = t
	f.typeStr = t.String()
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
	return f.type_.CreateField(f.id, f.name, f.type_)
}

func (f *BaseField) Parse(_ any) error {
	return f.self.Parse(f.value)
}

func (f *BaseField) Build() any {
	return f.self.Build()
}

func (f *BaseField) Asc() *bitable.Sort {
	builder := &bitable.SortBuilder{}
	return builder.FieldName(f.name).Build()
}

func (f *BaseField) Desc() *bitable.Sort {
	builder := &bitable.SortBuilder{}
	return builder.FieldName(f.name).Desc(true).Build()
}
