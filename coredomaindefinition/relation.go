package coredomaindefinition

type Relation struct {
	Source        *Model
	Target        *Model
	Type          RelationType
	IgnoreReverse bool
}

type RelationType string

const (
	// Relation one to one not required
	RelationTypeOneToOne RelationType = "oneToOne"

	// Relation one to many not required
	RelationTypeOneToMany RelationType = "oneToMany"

	// Relation many to one not required
	RelationTypeManyToOne RelationType = "manyToOne"

	// Relation many to many not required
	RelationTypeManyToMany RelationType = "manyToMany"

	// Relation belongs to, can not be null. Required.
	RelationTypeBelongsTo RelationType = "belongsTo"

	// Relation belongs to, can not be null. Required.
	RelationTypeSubresourcesOf RelationType = "subresourcesOf"
)
