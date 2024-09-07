package model

type Domain struct {
	Name                  string
	Architecture          *Architecture
	RepositoryTransaction *Interface
	Models                []*Model
	UsecaseStructs        []*Struct
	Usecases              []*Usecase
	UsecasesCRUDImpl      *Struct
	DomainRepository      *Interface
	Repositories          []*Repository
	RepositoryErrors      []*Var
	Controllers           []*Struct
	PortImplementations   map[string][]*Struct
	GormTransaction       *Struct
	GormModels            []*GormModel
	GormDomainRepository  *Struct
	HttpService           *Struct
	Service               *Interface
	Pagination            *Struct
	Ordering              *Struct
}
