package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleoGitHub/golem/goGeneration/domain/consts"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
	"github.com/cleoGitHub/golem/pkg/stringtool"
)

func (b *domainBuilder) GetUsecaseValidatorImpl(ctx context.Context) *model.Struct {
	if b.Err != nil {
		return nil
	}

	if b.UsecaseValidatorImpl != nil {
		return b.UsecaseValidatorImpl
	}

	usecaseValidator := &model.Struct{
		Name: stringtool.UpperFirstLetter(b.Domain.Name) + "UsecasesValidatorImpl",
		Fields: []*model.Field{
			{
				Name: "Validator",
				Type: b.GetValidator(ctx),
			},
			{
				Name: "Usecases",
				Type: b.DomainUsecase,
			},
		},
	}

	for _, usecase := range b.Domain.Usecases {
		f := usecase.Function.Copy().(*model.Function)
		f.Content = func() (string, []*model.GoPkg) {
			str := fmt.Sprintf(
				"if err := %s.%s.%s(ctx, %s); err != nil {",
				usecaseValidator.GetMethodName(), usecaseValidator.Fields[0].Name, usecaseValidator.Fields[0].Type.(*model.Interface).Methods[0].Name, f.Args[1].Name,
			) + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN
			str += fmt.Sprintf("return %s.%s.%s(ctx, %s)", usecaseValidator.GetMethodName(), usecaseValidator.Fields[1].Name, f.Name, f.Args[1].Name) + consts.LN
			return str, nil
		}
		usecaseValidator.Methods = append(usecaseValidator.Methods, f)
	}

	b.UsecaseValidatorImpl = usecaseValidator

	return usecaseValidator
}
