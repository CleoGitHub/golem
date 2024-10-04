package model

type Enum struct {
	Name   string
	Type   Type
	Values map[string]interface{}
}

func (e *Enum) GetType(opts ...GetTypeOpt) string {
	return e.Name
}

func (e *Enum) Copy() Type {
	return &Enum{
		Name:   e.Name,
		Type:   e.Type.Copy(),
		Values: e.Values,
	}
}

func (e *Enum) SubTypes() []Type {
	return []Type{e.Type}
}
