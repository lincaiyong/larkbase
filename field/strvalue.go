package field

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (f *TextField) StringValue() string {
	if f.value == nil {
		return ""
	}
	return f.value.(string)
}

func (f *NumberField) StringValue() string {
	if f.value == nil {
		return ""
	}
	return fmt.Sprintf("%g", f.value)
}

func (f *SingleSelectField) StringValue() string {
	if f.value == nil {
		return ""
	}
	return f.value.(string)
}

func (f *MultiSelectField) StringValue() string {
	if f.value == nil {
		return ""
	}
	b, _ := json.Marshal(f.value)
	return string(b)
}

func (f *DateField) StringValue() string {
	if f.value == nil {
		return ""
	}
	return timeToBeijingDateTimeStr(f.value.(time.Time))
}

func (f *CheckboxField) StringValue() string {
	if f.value == nil || !f.value.(bool) {
		return ""
	}
	return "1"
}

func (f *PersonField) StringValue() string {
	if f.value == nil {
		return ""
	}
	return strings.Join(f.value.([]string), ",")
}

func (f *PhoneField) StringValue() string {
	if f.value == nil {
		return ""
	}
	return f.value.(string)
}

func (f *UrlField) StringValue() string {
	if f.value == nil {
		return ""
	}
	return f.value.(string)
}

func (f *MediaField) StringValue() string {
	if f.value == nil {
		return ""
	}
	return strings.Join(f.value.([]string), ",")
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

func (f *AutoNumberField) StringValue() string {
	if f.value == nil {
		return ""
	}
	return fmt.Sprintf("%d", f.value)
}

func (f *ButtonField) StringValue() string {
	return ""
}
