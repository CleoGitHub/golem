package coredomaindefinition

type Field struct {
	Name        string
	Type        Type
	Validations []*Validation
}
