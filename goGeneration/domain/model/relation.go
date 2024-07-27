package model

import "github.com/cleoGitHub/golem/coredomaindefinition"

type RelationType string

const (
	RelationMultiple        RelationType = "multiple"
	RelationSingleOptionnal RelationType = "singleOptionnal"
	RelationSingleMandatory RelationType = "singleMandatory"
)

type Relation struct {
	On         *Model
	Type       RelationType
	Definition *coredomaindefinition.Relation
}
