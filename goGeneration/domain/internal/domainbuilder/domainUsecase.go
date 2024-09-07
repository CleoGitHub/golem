package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/stringtool"
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
