package larkfield

type Field interface {
	Id() string
	SetId(string)
	Name() string
	SetName(string)
	Type() string
	SetType(string)
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
