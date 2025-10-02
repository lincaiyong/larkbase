package field

import (
	"github.com/lincaiyong/log"
	"strconv"
	"strings"
)

func parseError(type_, expect string, v any) {
	log.FatalLog("fail to parse value for number field: expect float64, got %T", v)
}

func (f *TextField) Parse(v any) {
	if v1, ok1 := v.([]any); ok1 {
		items := make([]string, 0)
		for _, v2 := range v1 {
			if v3, ok3 := v2.(map[string]any); ok3 {
				type_, okT := v3["type"].(string)
				if !okT {
					parseError(f.type_, "string", v3["type"])
				}
				if type_ != "text" {
					log.FatalLog("fail to handle text field type: %s", type_)
				}
				value, okV := v3["text"].(string)
				if !okV {
					parseError(f.type_, "string", v3["value"])
				}
				items = append(items, value)
			}
		}
		if len(items) > 1 {
			log.FatalLog("fail to handle text field with more than 1 item: %s", strings.Join(items, ","))
		}
		if len(items) == 1 {
			f.value = items[0]
		}
	}
}

func (f *NumberField) Parse(v any) {
	if vv, ok := v.(float64); ok {
		f.value = vv
	} else {
		parseError(f.type_, "float64", v)
	}
}

func (f *SingleSelectField) Parse(v any) {

}

func (f *MultiSelectField) Parse(v any) {

}

func (f *DateField) Parse(v any) {

}

func (f *CheckboxField) Parse(v any) {
	if val, ok := v.(bool); ok {
		f.value = val
	} else {
		parseError(f.type_, "bool", v)
	}
}

func (f *PersonField) Parse(v any) {
	if v1, ok1 := v.([]any); ok1 {
		items := make([]string, 0)
		for _, v2 := range v1 {
			if v3, ok3 := v2.(map[string]any); ok3 {
				if v4, ok4 := v3["name"].(string); ok4 {
					items = append(items, v4)
				} else {
					parseError(f.type_, "string", v3["name"])
				}
			} else {
				parseError(f.type_, "map[string]any", v2)
			}
		}
		f.value = items
	} else {
		parseError(f.type_, "[]any", v)
	}
}

func (f *PhoneField) Parse(v any) {

}

func (f *UrlField) Parse(v any) {

}

func (f *MediaField) Parse(v any) {

}

func (f *SingleLinkField) Parse(v any) {

}

func (f *LookupField) Parse(v any) {

}

func (f *FormulaField) Parse(v any) {

}

func (f *DuplexLinkField) Parse(v any) {

}

func (f *LocationField) Parse(v any) {

}

func (f *GroupField) Parse(v any) {

}

func (f *WorkflowField) Parse(v any) {

}

func (f *CreatedTimeField) Parse(v any) {

}

func (f *ModifiedTimeField) Parse(v any) {

}

func (f *CreatePersonField) Parse(v any) {

}

func (f *ModifyPersonField) Parse(v any) {

}

func (f *AutoNumberField) Parse(v any) {
	if val, err := strconv.ParseInt(v.(string), 0, 64); err == nil {
		f.value = int(val)
	} else {
		parseError(f.type_, "integer string", v)
	}
}

func (f *ButtonField) Parse(v any) {

}
