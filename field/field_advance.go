package field

type SingleLinkField struct{ BaseField }
type LookupField struct{ BaseField }
type FormulaField struct{ BaseField }
type DuplexLinkField struct{ BaseField }
type LocationField struct{ BaseField }
type GroupField struct{ BaseField }
type WorkflowField struct{ BaseField }
type CreatedTimeField DateField
type ModifiedTimeField DateField
type CreatePersonField PersonField
type ModifyPersonField PersonField
type AutoNumberField NumberField
type ButtonField struct{ BaseField }
