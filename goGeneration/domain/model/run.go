package model

import "github.com/cleogithub/golem/coredomaindefinition"

type Run struct {
	Domain           *Domain
	DomainDefinition coredomaindefinition.Domain

	ModelToDefinition map[*Model]*coredomaindefinition.Model
	DefinitionToModel map[*coredomaindefinition.Model]*Model
	ModelToRepository map[*Model]*Repository
	RepositoryToModel map[*Repository]*Model
}
