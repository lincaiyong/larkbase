package larkfield

type Field interface {
	SetSelf(Field)
	Id() string
	SetId(string)
	Name() string
	SetName(string)
	TypeStr() string
	Type() Type
	SetType(Type)
	UnderlayValue() any
	SetUnderlayValue(any)
	SetUnderlayValueNoDirty(any)
	Dirty() bool
	SetDirty(v bool)
	StringValue() string

	Fork() Field

	Parse(v any) error
	Build() any
}
