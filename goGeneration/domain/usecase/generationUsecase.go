package usecase

import (
	"context"

	"github.com/cleoGitHub/golem/coredomaindefinition"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
)

type GenerationUsecase interface {
	// Generation
	GenerateDomainUsecase(ctx context.Context, domainDefinition coredomaindefinition.Domain, path string) error

	// Format
	FormatDomainUsecase(ctx context.Context, domain *model.Domain, path string) error

	GoTidyDomainUsecase(ctx context.Context, path string, domain *model.Domain) error
	InitDomainUsecase(ctx context.Context, domain *model.Domain, path string) error
}
