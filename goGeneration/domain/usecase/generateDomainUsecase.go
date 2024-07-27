package usecase

import (
	"context"

	"github.com/cleoGitHub/golem/coredomaindefinition"
	"github.com/cleoGitHub/golem/goGeneration/domain/consts"
	"github.com/cleoGitHub/golem/goGeneration/domain/internal/domainbuilder"
	"github.com/cleoGitHub/golem/pkg/merror"
)

// GenerateDomainUsecase implements GenerationUsecase.
func (g *GenerationUsecaseImpl) GenerateDomainUsecase(ctx context.Context, domainDefinition coredomaindefinition.Domain, path string) error {
	domainBuilder := domainbuilder.NewDomainBuilder(
		domainDefinition.Name,
		consts.DefaultModelFields,
	)

	for _, m := range domainDefinition.Models {
		domainBuilder.WithModel(ctx, m)
	}

	for _, r := range domainDefinition.Relations {
		domainBuilder.WithRelation(ctx, r)
	}

	for _, r := range domainDefinition.Repositories {
		domainBuilder.WithRepository(ctx, r)
	}

	for _, r := range domainDefinition.CRUDs {
		domainBuilder.WithCRUD(ctx, r)
	}

	for _, r := range domainDefinition.Usecases {
		domainBuilder.WithUsecase(ctx, r)
	}

	domain, err := domainBuilder.Build(ctx)
	if err != nil {
		return merror.Stack(err)
	}

	if err := g.InitDomainUsecase(ctx, domain, path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteUsecaseValidatorInterfaceUsecase(ctx, domain, domainBuilder.GetValidator(ctx), path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteDomainUsecasesCRUDImplUsecase(ctx, domain, domain.UsecasesCRUDImpl, path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteDomainUsecasesValidatorImplUsecase(ctx, domain, domainBuilder.GetUsecaseValidatorImpl(ctx), path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteDomainRepositoryUsecase(ctx, domain, domain.DomainRepository, path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteRepositoryTransactionInterfaceUsecase(ctx, domain, domainBuilder.GetTransation(ctx), path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteRepositoryPaginationUsecase(ctx, domain, domainBuilder.GetPagination(ctx), path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteRepositoryOrderingUsecase(ctx, domain, domainBuilder.GetOrdering(ctx), path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteDomainUsecasesUsecase(ctx, domain, domainBuilder.GetDomainUsecase(ctx), path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteUsecasesStructsUsecase(ctx, domain, path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteRepositoryErrorsUsecase(ctx, domain, path); err != nil {
		return merror.Stack(err)
	}

	for _, m := range domain.Models {
		if err := g.WriteModelUsecase(ctx, domain, m, path); err != nil {
			return merror.Stack(err)
		}
	}

	for _, r := range domain.Repositories {
		if err := g.WriteRepositoryUsecase(ctx, domain, r, path); err != nil {
			return merror.Stack(err)
		}
	}

	for _, c := range domain.Controllers {
		if err := g.WriteControllerUsecase(ctx, domain, c, path); err != nil {
			return merror.Stack(err)
		}
	}

	if err := g.WriteGormTransactionUsecase(ctx, domain, domain.GormTransaction, path); err != nil {
		return merror.Stack(err)
	}

	for _, m := range domain.GormModels {
		if err := g.WriteGormModelUsecase(ctx, domain, m, path); err != nil {
			return merror.Stack(err)
		}
	}

	if err := g.WriteGormDomainRepositoryUsecase(ctx, domain, domain.GormDomainRepository, path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteServiceUsecase(ctx, domain, domain.Service, path); err != nil {
		return merror.Stack(err)
	}

	if err := g.WriteHttpServiceUsecase(ctx, domain, domain.HttpService, path); err != nil {
		return merror.Stack(err)
	}

	if err := g.FormatDomainUsecase(ctx, domain, path); err != nil {
		return merror.Stack(err)
	}

	if err := g.GoTidyDomainUsecase(ctx, path, domain); err != nil {
		return merror.Stack(err)
	}

	return nil
}
