package coredomaindefinition

type Model struct {
	Name       string
	Fields     []*Field
	Activable  bool
	Archivable bool
}

func (m Model) GetType() string {
	return m.Name
}

func NewModel(name string) *Model {
	return &Model{
		Name:       name,
		Archivable: true,
	}
}
