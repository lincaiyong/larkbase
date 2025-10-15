package larkfield

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"strings"
)

func parseError(type_, expect string, v any) error {
	return fmt.Errorf("fail to parse value for %s field: expect %s, got %T", type_, expect, v)
}

func (f *TextField) Parse(v any) error {
	if v1, ok1 := v.(string); ok1 {
		f.value = v1
		return nil
	}
	var data []struct {
		Type string `mapstructure:"type"`
		Text string `mapstructure:"text"`
	}
	err := mapstructure.WeakDecode(v, &data)
	if err != nil {
		return fmt.Errorf("fail to parse text field: %w, %v", err, v)
	}
	items := make([]string, 0)
	for _, item := range data {
		items = append(items, item.Text)
	}
	f.value = strings.Join(items, "")
	return nil
}

func (f *NumberField) Parse(v any) error {
	if vv, ok := v.(float64); ok {
		f.value = vv
		return nil
	} else {
		return parseError(f.type_, "float64", v)
	}
}

func (f *SingleSelectField) Parse(v any) error {
	if v1, ok1 := v.(string); ok1 {
		f.value = v1
		return nil
	} else {
		return parseError(f.type_, "string", v)
	}
}

func (f *MultiSelectField) Parse(v any) error {
	var data []string
	err := mapstructure.Decode(v, &data)
	if err != nil {
		return fmt.Errorf("fail to parse multi select field: %w, %v", err, v)
	}
	f.value = data
	return nil
}

func (f *DateField) Parse(v any) error {
	if v1, ok1 := v.(float64); ok1 {
		f.value = UnixSecondsToTime(int64(v1 / 1000))
		return nil
	} else {
		return parseError(f.type_, "float64", v)
	}
}

func (f *CheckboxField) Parse(v any) error {
	if val, ok := v.(bool); ok {
		f.value = val
		return nil
	} else {
		return parseError(f.type_, "bool", v)
	}
}

//
//func (f *PersonField) Parse(v any) error {
//	if v1, ok1 := v.([]any); ok1 {
//		items := make([]string, 0)
//		for _, v2 := range v1 {
//			if v3, ok3 := v2.(map[string]any); ok3 {
//				if v4, ok4 := v3["name"].(string); ok4 {
//					items = append(items, v4)
//				} else {
//					return parseError(f.type_, "string", v3["name"])
//				}
//			} else {
//				return parseError(f.type_, "map[string]any", v2)
//			}
//		}
//		f.value = items
//		return nil
//	} else {
//		return parseError(f.type_, "[]any", v)
//	}
//}
//
//func (f *PhoneField) Parse(v any) error {
//	if v1, ok1 := v.(string); ok1 {
//		f.value = v1
//		return nil
//	} else {
//		return parseError(f.type_, "string", v)
//	}
//}

func (f *UrlField) Parse(v any) error {
	var data struct {
		Type string `mapstructure:"type"`
		Link string `mapstructure:"link"`
	}
	err := mapstructure.Decode(v, &data)
	if err != nil {
		return fmt.Errorf("fail to parse url field: %w, %v", err, v)
	}
	f.value = data.Link
	return nil
}

//func (f *MediaField) Parse(v any) error {
//	if v1, ok1 := v.([]any); ok1 {
//		items := make([]string, 0)
//		for _, v2 := range v1 {
//			if v3, ok3 := v2.(map[string]any); ok3 {
//				fileToken, ok4 := v3["file_token"].(string)
//				if !ok4 {
//					return parseError(f.type_, "string", v3["file_token"])
//				}
//				items = append(items, fileToken)
//			} else {
//				return parseError(f.type_, "map[string]any", v2)
//			}
//		}
//		return nil
//	} else {
//		return parseError(f.type_, "[]any", v)
//	}
//}

func (f *AutoNumberField) Parse(v any) error {
	if val, err := strconv.ParseInt(v.(string), 0, 64); err == nil {
		f.value = int(val)
		return nil
	} else {
		return parseError(f.type_, "integer string", v)
	}
}

func (f *ModifiedTimeField) Parse(v any) error {
	if v1, ok1 := v.(float64); ok1 {
		f.value = UnixSecondsToTime(int64(v1 / 1000))
		return nil
	} else {
		return parseError(f.type_, "float64", v)
	}
}

func parseAnyTypeField(v any) (string, error) {
	var data struct {
		Type  Type `mapstructure:"type"`
		Value any  `mapstructure:"value"`
	}
	err := mapstructure.Decode(v, &data)
	if err != nil {
		return "", err
	}
	switch data.Type {
	case TypeText:
		f := &TextField{}
		err = f.Parse(data.Value)
		if err != nil {
			return "", err
		}
		return f.StringValue(), nil
	case TypeNumber:
		f := &NumberField{}
		err = f.Parse(data.Value)
		if err != nil {
			return "", err
		}
		return f.StringValue(), nil
	case TypeSingleSelect:
		f := &SingleSelectField{}
		err = f.Parse(data.Value)
		if err != nil {
			return "", err
		}
		return f.StringValue(), nil
	default:
		return "", fmt.Errorf("unsupported type field for lookup or formula: %s", data.Type.String())
	}
}

func (f *LookupField) Parse(v any) error {
	val, err := parseAnyTypeField(v)
	if err != nil {
		return fmt.Errorf("fail to parse lookup field: %w, %v", err, v)
	}
	f.value = val
	return nil
}

func (f *FormulaField) Parse(v any) error {
	val, err := parseAnyTypeField(v)
	if err != nil {
		return fmt.Errorf("fail to parse formula field: %w, %v", err, v)
	}
	f.value = val
	return nil
}
