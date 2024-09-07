package model

import "github.com/cleogithub/golem-common/pkg/stringtool"

type Struct struct {
	// Consts referenced in this struct
	Consts  []*Var
	Name    string
	Fields  []*Field
	Methods []*Function
}

func (s *Struct) GetType(typeOpts ...GetTypeOpt) string {
	return s.Name
}

func (s *Struct) SubTypes() []Type {
	return []Type{}
}

func (s *Struct) GetMethodName() string {
	return stringtool.LowerFirstLetter(s.Name)
}

func (s *Struct) Copy() Type {
	return &Struct{
		Consts:  ArrayConsts(s.Consts).Copy(),
		Name:    s.Name,
		Fields:  Fields(s.Fields).Copy(),
		Methods: s.Methods,
	}
}
