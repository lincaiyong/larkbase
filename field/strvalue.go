package field

import (
	"fmt"
	"strings"
)

func (f *TextField) StringValue() string {
	return f.value.(string)
}

func (f *NumberField) StringValue() string {
	return fmt.Sprintf("%g", f.value)
}

func (f *SingleSelectField) StringValue() string {
	return ""
}

func (f *MultiSelectField) StringValue() string {
	return ""
}

func (f *DateField) StringValue() string {
	return ""
}

func (f *CheckboxField) StringValue() string {
	return ""
}

func (f *PersonField) StringValue() string {
	return strings.Join(f.value.([]string), ",")
}

func (f *PhoneField) StringValue() string {
	return ""
}

func (f *UrlField) StringValue() string {
	return ""
}

func (f *MediaField) StringValue() string {
	return ""
}

func (f *SingleLinkField) StringValue() string {
	return ""
}

func (f *LookupField) StringValue() string {
	return ""
}

func (f *FormulaField) StringValue() string {
	return ""
}

func (f *DuplexLinkField) StringValue() string {
	return ""
}

func (f *LocationField) StringValue() string {
	return ""
}

func (f *GroupField) StringValue() string {
	return ""
}

func (f *WorkflowField) StringValue() string {
	return ""
}

func (f *CreatedTimeField) StringValue() string {
	return ""
}

func (f *ModifiedTimeField) StringValue() string {
	return ""
}

func (f *CreatePersonField) StringValue() string {
	return ""
}

func (f *ModifyPersonField) StringValue() string {
	return ""
}

func (f *AutoNumberField) StringValue() string {
	return fmt.Sprintf("%d", f.value)
}

func (f *ButtonField) StringValue() string {
	return ""
}
