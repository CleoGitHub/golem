package coredomaindefinition

type CRUD struct {
	On     *Model
	Create *CRUDAction
	Get    *CRUDAction
	List   *CRUDAction
	Update *CRUDAction
	Delete *CRUDAction
}

type CRUDAction struct {
	Active bool
	Roles  []string
	// Roles required for retriving active element
	RolesForActive []string
}
