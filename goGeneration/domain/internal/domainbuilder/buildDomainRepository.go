package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	REPOSITORY_MIGRATE           = "Migrate"
	REPOSITORY_BEGIN_TRANSACTION = "BeginTransaction"
)

func (b *domainBuilder) buildDomainRepository(ctx context.Context) (*model.File, error) {
	if b.Err != nil {
		return nil, b.Err
	}

	domainRepository := &model.Interface{
		Name: GetDomainRepositoryName(ctx, b.Definition),
		Methods: []*model.Function{
			{
				Name: REPOSITORY_MIGRATE,
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PkgReference{
							Pkg: consts.CommonPkgs["context"],
							Reference: &model.ExternalType{
								Type: "Context",
							},
						},
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
			{
				Name: REPOSITORY_BEGIN_TRANSACTION,
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PkgReference{
							Pkg: consts.CommonPkgs["context"],
							Reference: &model.ExternalType{
								Type: "Context",
							},
						},
					},
				},
				Results: []*model.Param{
					{
						Type: &model.PkgReference{
							Pkg: b.Domain.Architecture.RepositoryPkg,
							Reference: &model.ExternalType{
								Type: TRANSACTION_NAME,
							},
						},
					},
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
		},
	}

	for _, repositoryBuilder := range b.RepositoryBuilders {
		domainRepository.Methods = append(domainRepository.Methods, repositoryBuilder.Methods...)
	}

	return &model.File{
		Name: GetDomainRepositoryName(ctx, b.Definition),
		Pkg:  b.GetRepositoryPackage(),
		Elements: []interface{}{
			domainRepository,
		},
	}, nil
}
