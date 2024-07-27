package model

type Var struct {
	Name    string
	Type    Type
	Value   interface{}
	IsConst bool
}

func (c *Var) Copy() *Var {
	return &Var{
		Name:  c.Name,
		Type:  c.Type.Copy(),
		Value: c.Value,
	}
}

type ArrayConsts []*Var

func (c ArrayConsts) Copy() []*Var {
	newConsts := make([]*Var, len(c))
	for i, v := range c {
		newConsts[i] = v.Copy()
	}
	return newConsts
}
