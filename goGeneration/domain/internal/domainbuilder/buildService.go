package domainbuilder

import (
	"context"
	"strings"

	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/stringtool"
)

func (b *domainBuilder) buildService(ctx context.Context) *domainBuilder {
	if b.Err != nil {
		return b
	}

	service := &model.Interface{
		Name:    stringtool.UpperFirstLetter(b.Domain.Name) + "Service",
		Methods: []*model.Function{},
	}

	for _, usecase := range b.Domain.Usecases {
		m := usecase.Function.Copy().(*model.Function)
		m.Name = strings.TrimSuffix(m.Name, "Usecase")
		m.Args[1].Type = &model.PointerType{
			Type: m.Args[1].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.Copy(),
		}
		m.Results[0].Type = &model.PointerType{
			Type: m.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.Copy(),
		}
		service.Methods = append(service.Methods, m)

	}

	b.Domain.Service = service

	return b
}
