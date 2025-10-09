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

func (f *PersonField) Fork() Field {
	return &PersonField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *PhoneField) Fork() Field {
	return &PhoneField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *UrlField) Fork() Field {
	return &UrlField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *MediaField) Fork() Field {
	return &MediaField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *SingleLinkField) Fork() Field {
	return &SingleLinkField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *LookupField) Fork() Field {
	return &LookupField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *FormulaField) Fork() Field {
	return &FormulaField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *DuplexLinkField) Fork() Field {
	return &DuplexLinkField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *LocationField) Fork() Field {
	return &LocationField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *GroupField) Fork() Field {
	return &GroupField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *WorkflowField) Fork() Field {
	return &WorkflowField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *AutoNumberField) Fork() Field {
	return &AutoNumberField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *ButtonField) Fork() Field {
	return &ButtonField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}

func (f *ModifiedTimeField) Fork() Field {
	return &ModifiedTimeField{BaseField{id: f.id, name: f.name, type_: f.type_, value: f.value}}
}
