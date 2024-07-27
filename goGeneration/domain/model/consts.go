package model

type Consts struct {
	Name   string
	Values []interface{}
}

func (c Consts) Copy() Consts {
	newConsts := Consts{
		Name:   c.Name,
		Values: c.Values,
	}
	// for i, v := range c.Values {
	// 	newConsts.Values[i] = v.Copy().(*Var)
	// }
	return newConsts
}
