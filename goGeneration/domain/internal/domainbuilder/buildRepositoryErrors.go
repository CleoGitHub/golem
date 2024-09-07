package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (b *domainBuilder) buildRepositoryErrors(ctx context.Context) *domainBuilder {
	if b.Err != nil {
		return b
	}

	b.RepositoryErrors["notFound"] = &model.Var{
		Name: "ErrNotFound",
		Value: &model.PkgReference{
			Pkg: consts.CommonPkgs["fmt"],
			Reference: &model.ExternalType{
				Type: `Errorf("not found")`,
			},
		},
	}

	b.Domain.RepositoryErrors = []*model.Var{}

	for _, err := range b.RepositoryErrors {
		b.Domain.RepositoryErrors = append(b.Domain.RepositoryErrors, err)
	}

	return b
}
