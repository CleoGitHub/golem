package model

import "reflect"

type Var struct {
	Name    string
	Type    Type
	Value   interface{}
	IsConst bool
}

func (c *Var) Copy() *Var {
	// use Copy on c.Type with refelction
	v := reflect.ValueOf(c.Type)
	results := v.MethodByName("Copy").Call([]reflect.Value{})
	typeCopied := results[0].Interface().(Type)

	return &Var{
		Name:  c.Name,
		Type:  typeCopied,
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
