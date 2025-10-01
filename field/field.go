package field

type Field interface {
	Name() string
	SetName(string)
	Type() string
	SetType(string)
	Value() string
	SetValue(string)
	SetValueNoDirty(string)
	Dirty() bool

	ParseFromLarkSuite(v any)
	BuildForLarkSuite() (any, error)
}
