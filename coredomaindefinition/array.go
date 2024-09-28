package coredomaindefinition

type Array struct {
	Type Type
}

func (t Array) GetType() string {
	return "array"
}
