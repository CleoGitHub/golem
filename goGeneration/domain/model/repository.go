package model

type Repository struct {
	On              *Model
	Name            string
	TableName       string
	AllowedWheres   Consts
	AllowedOrderBys Consts
	DefaultOrderBy  string
	FieldToColumn   Map
	Methods         []*RepositoryMethod
	Functions       []*Function
}

func (r *Repository) GetType(typeOpts ...GetTypeOpt) string {
	return r.Name
}

func (r *Repository) SubTypes() []Type {
	return []Type{}
}

func (r *Repository) Copy() Type {
	return &Repository{
		On:              r.On,
		Name:            r.Name,
		TableName:       r.TableName,
		AllowedWheres:   r.AllowedWheres.Copy(),
		AllowedOrderBys: r.AllowedOrderBys.Copy(),
		DefaultOrderBy:  r.DefaultOrderBy,
		FieldToColumn:   r.FieldToColumn,
		Methods:         RepositoryMethods(r.Methods).Copy(),
		Functions:       Functions(r.Functions).Copy(),
	}
}

type RepositoryMethod struct {
	Context  *Struct
	Opt      *TypeDefinition
	Opts     []*Function
	Function *Function
}

func (r *RepositoryMethod) Copy() *RepositoryMethod {
	return &RepositoryMethod{
		Context:  r.Context,
		Opt:      r.Opt,
		Opts:     Functions(r.Opts).Copy(),
		Function: r.Function,
	}
}

type RepositoryMethods []*RepositoryMethod

func (r RepositoryMethods) Copy() []*RepositoryMethod {
	newMethods := make([]*RepositoryMethod, len(r))
	for i, v := range r {
		newMethods[i] = v.Copy()
	}
	return newMethods
}
