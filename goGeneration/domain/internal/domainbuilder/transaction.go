package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	TRANSACTION_NAME     = "Transaction"
	TRANSACTION_GET      = "Get"
	TRANSACTION_COMMIT   = "Commit"
	TRANSACTION_ROLLBACK = "Rollback"
)

var TRANSACTION = &model.Interface{
	Name: TRANSACTION_NAME,
	Methods: []*model.Function{
		{
			Name: TRANSACTION_GET,
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
			Name: TRANSACTION_COMMIT,
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
			Name: TRANSACTION_ROLLBACK,
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

func (b *domainBuilder) addTransaction(ctx context.Context) *domainBuilder {
	if b.Err != nil {
		return b
	}

	b.Domain.Ports = append(b.Domain.Ports, &model.File{
		Name: TRANSACTION_NAME,
		Pkg:  b.GetRepositoryPackage(),
		Elements: []interface{}{
			TRANSACTION,
		},
	})

	return b
}
