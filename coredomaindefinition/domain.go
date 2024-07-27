package coredomaindefinition

type Domain struct {
	Name          string
	Configuration *DomainConfiguration
	Models        []*Model
	Relations     []*Relation
	Repositories  []*Repository
	Usecases      []*Usecase
	CRUDs         []*CRUD
}
