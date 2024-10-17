package usecase

import (
	"context"

	"github.com/cleogithub/golem/coredomaindefinition"
)

type GenerationUsecase interface {
	GenerateDomainUsecase(ctx context.Context, domainDefinition coredomaindefinition.Domain, path string) error
}
