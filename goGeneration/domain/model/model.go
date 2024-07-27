package model

type Model struct {
	// Struct related to model
	JsonName  string
	Struct    *Struct
	Relations []*Relation
	Activable bool
}

func (m *Model) GetType(typeOpts ...GetTypeOpt) string {
	return m.Struct.Name
}

func (m *Model) SubTypes() []Type {
	return []Type{}
}

func (m *Model) Copy() Type {
	return &Model{Struct: m.Struct.Copy().(*Struct), Relations: m.Relations}
}
