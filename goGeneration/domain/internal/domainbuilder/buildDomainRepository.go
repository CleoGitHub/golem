package domainbuilder

import (
	"context"

	"github.com/cleoGitHub/golem/goGeneration/domain/consts"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
	"github.com/cleoGitHub/golem/pkg/stringtool"
)

func (b *domainBuilder) buildDomainRepository(ctx context.Context) *domainBuilder {
	if b.Err != nil {
		return b
	}

	domainRepository := &model.Interface{
		Name: stringtool.UpperFirstLetter(b.Domain.Name) + "Repository",
		Methods: []*model.Function{
			{
				Name: "Migrate",
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
				Name: "BeginTransaction",
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
							Pkg:       b.Domain.Architecture.RepositoryPkg,
							Reference: b.GetTransation(ctx),
						},
					},
				},
			},
		},
	}

	b.Domain.DomainRepository = domainRepository

	return b
}
