package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (b *domainBuilder) GetDomainUsecase(ctx context.Context) *model.Interface {
	if b.DomainUsecase != nil {
		return b.DomainUsecase
	}

	b.DomainUsecase = &model.Interface{
		Name: stringtool.UpperFirstLetter(b.Domain.Name) + "Usecases",
	}

	return b.DomainUsecase
}
