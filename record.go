package larkbase

func NewRecord() *Record {
	return new(Record)
}

type Record struct {
	Id     string
	Fields map[string]Field
}

func (r *Record) buildForLarkSuite() (map[string]any, error) {
	fields := make(map[string]any)
	for name, field := range r.Fields {
		if field.Dirty() {
			var err error
			fields[name], err = field.BuildForLarkSuite()
			if err != nil {
				return nil, err
			}
		}
	}
	return fields, nil
}
