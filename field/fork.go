package field

func (f *TextField) Fork() Field {
	return &TextField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *NumberField) Fork() Field {
	return &NumberField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *SingleSelectField) Fork() Field {
	return &SingleSelectField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *MultiSelectField) Fork() Field {
	return &MultiSelectField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *DateField) Fork() Field {
	return &DateField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *CheckboxField) Fork() Field {
	return &CheckboxField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *PersonField) Fork() Field {
	return &PersonField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *PhoneField) Fork() Field {
	return &PhoneField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *UrlField) Fork() Field {
	return &UrlField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *MediaField) Fork() Field {
	return &MediaField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *SingleLinkField) Fork() Field {
	return &SingleLinkField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *LookupField) Fork() Field {
	return &LookupField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *FormulaField) Fork() Field {
	return &FormulaField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *DuplexLinkField) Fork() Field {
	return &DuplexLinkField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *LocationField) Fork() Field {
	return &LocationField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *GroupField) Fork() Field {
	return &GroupField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *WorkflowField) Fork() Field {
	return &WorkflowField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *AutoNumberField) Fork() Field {
	return &AutoNumberField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}

func (f *ButtonField) Fork() Field {
	return &ButtonField{BaseField{name: f.name, type_: f.type_, value: f.value}}
}
