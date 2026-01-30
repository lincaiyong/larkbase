package larkbase

import (
	"github.com/lincaiyong/larkbase/larkfield"
	"time"
)

func NewRecord() *Record {
	return new(Record)
}

type Record struct {
	Id           string
	ModifiedTime time.Time
	Fields       map[string]larkfield.Field
}

func (r *Record) buildForLarkSuite() (map[string]any, error) {
	fields := make(map[string]any)
	for name, field := range r.Fields {
		if field.Dirty() {
			fields[name] = field.Build()
		}
	}
	return fields, nil
}

type AnyRecord struct {
	Meta
	Id     NumberField `lark:"id"`
	Data   map[string]string
	update map[string]string
}

func (r *AnyRecord) Update(k, v string) {
	if r.update == nil {
		r.update = make(map[string]string)
	}
	r.update[k] = v
	r.Data[k] = v
}
