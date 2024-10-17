package model

type Param struct {
	Name string
	Type Type
}

func (p *Param) Copy() *Param {
	return &Param{
		Name: p.Name,
		Type: Copy(p.Type),
	}
}

type Params []*Param

func (p Params) Copy() Params {
	cp := make(Params, len(p))
	for i, v := range p {
		cp[i] = v.Copy()
	}
	return cp
}
