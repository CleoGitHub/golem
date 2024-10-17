package coredomaindefinition

type Repository struct {
	On *Model
	// Name of the table for the model
	TableName string
	// Optionnal: updatedAt will be used
	DefaultOrderBy string
	// Method to define in repository
	Methods []*RepositoryMethod
}

type RepositoryMethod struct {
	Name      string
	Params    []*Param
	Results   []Type
	Paginable bool
}
