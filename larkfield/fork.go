package larkfield

func (f *TextField) Fork() Field {
	return &TextField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *NumberField) Fork() Field {
	return &NumberField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *SingleSelectField) Fork() Field {
	return &SingleSelectField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *MultiSelectField) Fork() Field {
	return &MultiSelectField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *DateField) Fork() Field {
	return &DateField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *CheckboxField) Fork() Field {
	return &CheckboxField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *UrlField) Fork() Field {
	return &UrlField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *AutoNumberField) Fork() Field {
	return &AutoNumberField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *ModifiedTimeField) Fork() Field {
	return &ModifiedTimeField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}
