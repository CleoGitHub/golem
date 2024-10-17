package model

type Interface struct {
	Name    string
	Methods []*Function
}

func (i *Interface) GetType(...GetTypeOpt) string {
	return i.Name
}

func (i *Interface) SubTypes() []Type {
	return []Type{}
}

func (i *Interface) Copy() *Interface {
	return &Interface{
		Name:    i.Name,
		Methods: Functions(i.Methods).Copy(),
	}
}
