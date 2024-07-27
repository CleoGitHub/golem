package coredomaindefinition

type Model struct {
	Name      string
	Fields    []*Field
	Activable bool
}

func (m Model) GetType() string {
	return m.Name
}
