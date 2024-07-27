package model

type TypeDefinition struct {
	Name string
	Type Type
}

func (t *TypeDefinition) GetType(opt ...GetTypeOpt) string {
	return t.Name
}

func (t *TypeDefinition) SubTypes() []Type {
	return []Type{}
}

func (t *TypeDefinition) Copy() Type {
	return &TypeDefinition{
		Name: t.Name,
		Type: t.Type.Copy(),
	}
}
