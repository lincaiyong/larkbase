package larkbase

import (
	"github.com/lincaiyong/larkbase/larkfield"
	"time"
)

func NewRecord() *Record {
	return new(Record)
}

type Record struct {
	Id     string
	Fields map[string]larkfield.Field

	CreatedTime  time.Time
	ModifiedTime time.Time
	CreatePerson string
	ModifyPerson string
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
