package larkbase

func NewRecord() *Record {
	return new(Record)
}

type Record struct {
	Id     string
	Fields map[string]Field
}

func (r *Record) buildForLarkSuite() map[string]any {
	fields := make(map[string]any)
	for name, field := range r.Fields {
		if field.Value != "" {
			fields[name] = field.buildForLarkSuite()
		}
	}
	return fields
}
