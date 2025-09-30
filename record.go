package larkbase

type Record struct {
	Id     string            `json:"id"`
	Fields map[string]IField `json:"fields"`
}

func (r *Record) Build() map[string]any {
	fields := make(map[string]any)
	for name, field := range r.Fields {
		if field.Value() != "" {
			fields[name] = field.Build()
		}
	}
	return fields
}
