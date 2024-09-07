package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (b *domainBuilder) GetTransation(ctx context.Context) *model.Interface {
	if b.Err != nil {
		return nil
	}

	if b.Transaction != nil {
		return b.Transaction
	}

	transtactionIntf := &model.Interface{
		Name: "Transaction",
		Methods: []*model.Function{
			{
				Name: "Get",
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
						Type: model.PrimitiveTypeInterface,
					},
				},
			},
			{
				Name: "Commit",
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
				Name: "Rollback",
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
		},
	}

	b.Transaction = transtactionIntf
	b.Domain.RepositoryTransaction = transtactionIntf

	return transtactionIntf
}
