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

func (f *PersonField) Build() any {
	return nil
}

func (f *PhoneField) Build() any {
	return nil
}

func (f *UrlField) Build() any {
	return nil
}

func (f *MediaField) Build() any {
	return nil
}

func (f *SingleLinkField) Build() any {
	return nil
}

func (f *LookupField) Build() any {
	return nil
}

func (f *FormulaField) Build() any {
	return nil
}

func (f *DuplexLinkField) Build() any {
	return nil
}

func (f *LocationField) Build() any {
	return nil
}

func (f *GroupField) Build() any {
	return nil
}

func (f *WorkflowField) Build() any {
	return nil
}

func (f *AutoNumberField) Build() any {
	return fmt.Sprintf("%d", f.Value())
}

func (f *ButtonField) Build() any {
	return nil
}

func (f *ModifiedTimeField) Build() any {
	return f.value.(time.Time).UnixMilli()
}
