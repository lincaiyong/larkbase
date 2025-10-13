package larkfield

import (
	"fmt"
	"time"
)

func (f *TextField) Build() any {
	return f.value
}

func (f *NumberField) Build() any {
	return f.value
}

func (f *SingleSelectField) Build() any {
	return f.value
}

func (f *MultiSelectField) Build() any {
	return f.value
}

func (f *DateField) Build() any {
	return f.value.(time.Time).UnixMilli()
}

func (f *CheckboxField) Build() any {
	return f.value
}

func (f *UrlField) Build() any {
	return map[string]string{
		"link": f.value.(string),
		"text": f.value.(string),
	}
}

func (f *AutoNumberField) Build() any {
	return fmt.Sprintf("%d", f.Value())
}

func (f *ModifiedTimeField) Build() any {
	return f.value.(time.Time).UnixMilli()
}
