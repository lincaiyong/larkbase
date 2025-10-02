package field

type Field interface {
	Name() string
	SetName(string)
	Type() string
	SetType(string)
	UnderlayValue() any
	SetUnderlayValue(any)
	SetUnderlayValueNoDirty(any)
	Dirty() bool
	StringValue() string

	Fork() Field

	Parse(v any)
	Build() any
}
