package coredomaindefinition

type Repository struct {
	On *Model
	// Name of the table for the model
	TableName string
	// Optionnal: updatedAt will be used
	DefaultOrderBy string
	// Method to define in repository
	Methods []*RepositoryMethod
	// Element to preload with repository
	Withables []*Model
	// Default get method generated, no need to defined if true, require on to be defined
	GetMethod bool
	// Default list method generated, no need to defined if true
	ListMethod bool
	// Default create method generated, no need to defined if true
	CreateMethod bool
	// Default update method generated, no need to defined if true
	UpdateMethod bool
	// Default delete method generated, no need to defined if true
	DeleteMethod bool
}

type RepositoryMethod struct {
	Name    string
	Params  []*Param
	Results []Type
}
