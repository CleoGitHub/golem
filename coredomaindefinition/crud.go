package coredomaindefinition

type CRUD struct {
	On            *Model
	Create        CRUDAction
	Get           CRUDAction
	GetActive     CRUDAction
	List          CRUDAction
	ListActive    CRUDAction
	Update        CRUDAction
	Delete        CRUDAction
	RelationCRUDs []*RelationCRUD
}

type RelationCRUD struct {
	Relation *Relation
	Roles    []string
	Add      CRUDAction
	Remove   CRUDAction
	List     CRUDAction
}

type CRUDAction struct {
	Active bool
	Roles  []string
}
