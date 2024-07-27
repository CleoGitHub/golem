package coredomaindefinition

type Usecase struct {
	Name    string
	Args    []*Param
	Results []*Param
	Roles   []string
}
