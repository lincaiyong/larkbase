package larkfield

import (
	"fmt"
	"strconv"
	"strings"
)

func parseError(type_, expect string, v any) error {
	return fmt.Errorf("fail to parse value for %s field: expect %s, got %T", type_, expect, v)
}

func (f *TextField) Parse(v any) error {
	if v1, ok1 := v.([]any); ok1 {
		items := make([]string, 0)
		for _, v2 := range v1 {
			if v3, ok3 := v2.(map[string]any); ok3 {
				type_, okT := v3["type"].(string)
				if !okT {
					return parseError(f.type_, "string", v3["type"])
				}
				switch type_ {
				case "text", "url":
					value, okV := v3["text"].(string)
					if !okV {
						return parseError(f.type_, "string", v3["value"])
					}
					items = append(items, value)
				default:
					return fmt.Errorf("fail to handle text field type: %s", type_)
				}

			}
		}
		if len(items) == 1 {
			f.value = items[0]
		} else if len(items) > 1 {
			f.value = strings.Join(items, "")
		}
		return nil
	}
	if v1, ok1 := v.(string); ok1 {
		f.value = v1
		return nil
	}
	return fmt.Errorf("fail to handle value of type %T", v)
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
	if v1, ok1 := v.([]any); ok1 {
		items := make([]string, 0)
		for _, v2 := range v1 {
			if v3, ok3 := v2.(string); ok3 {
				items = append(items, v3)
			} else {
				return parseError(f.type_, "string", v2)
			}
		}
		f.value = items
		return nil
	} else {
		return parseError(f.type_, "[]any", v)
	}
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
	if v1, ok1 := v.(map[string]any); ok1 {
		if v1["type"] != nil && v1["type"].(string) != "url" {
			return fmt.Errorf("fail to handle url field type: %s", v1["type"])
		}
		link, ok3 := v1["link"].(string)
		if !ok3 {
			return parseError(f.type_, "string", v1["link"])
		}
		f.value = link
		return nil
	} else {
		return parseError(f.type_, "map[string]any", v)
	}
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
