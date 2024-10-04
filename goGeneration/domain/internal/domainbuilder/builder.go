package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem/coredomaindefinition"
)

type Builder interface {
	WithModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) Builder
	WithRepository(ctx context.Context, repositoryDefinition *coredomaindefinition.Repository) Builder
	WithRelation(ctx context.Context, relationDefinition *coredomaindefinition.Relation) Builder
	Build(ctx context.Context) error
}
