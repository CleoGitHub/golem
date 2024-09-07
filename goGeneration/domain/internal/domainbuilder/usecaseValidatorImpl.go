package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/internal/stringbuilder"
	"github.com/cleogithub/golem/goGeneration/domain/model"
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
			strct := usecaseValidator.GetMethodName()
			validator := usecaseValidator.Fields[0].Name
			validate := usecaseValidator.Fields[0].Type.(*model.Interface).Methods[0].Name
			validateMimeTypes := usecaseValidator.Fields[0].Type.(*model.Interface).Methods[4].Name

			strBuilder := stringbuilder.NewStringBuilder()
			strBuilder.Append(&stringbuilder.If{
				Condition: fmt.Sprintf("err := %s.%s.%s(ctx, %s); err != nil",
					strct, validator, validate, f.Args[1].Name,
				),
				Then: &stringbuilder.Chainable{
					Elem: stringbuilder.String("return nil, err"),
				},
			})

			// str := fmt.Sprintf(
			// 	"if err := %s.%s.%s(ctx, %s); err != nil {",
			// 	usecaseValidator.GetMethodName(), usecaseValidator.Fields[0].Name, usecaseValidator.Fields[0].Type.(*model.Interface).Methods[0].Name, f.Args[1].Name,
			// ) + consts.LN
			// str += "return nil, err" + consts.LN
			// str += "}" + consts.LN

			for _, field := range usecase.Request.Fields {
				fieldDefinition := b.FieldToParamUsecaseRequest[field]
				if fieldDefinition == nil {
					continue
				}
				for _, validation := range fieldDefinition.Validations {
					if validation.Rule == coredomaindefinition.ValidationRuleMIMETypes {
						mimeTypes := ""
						switch validation.Value.(type) {
						case string:
							mimeTypes = validation.Value.(string)
						case []string:
							mimeTypes = "[]string{"
							for _, mimeType := range validation.Value.([]string) {
								mimeTypes += fmt.Sprintf(`"%s",`, mimeType)
							}
							mimeTypes = mimeTypes[:len(mimeTypes)-1]
							mimeTypes += "}"
						}
						strBuilder.Append(&stringbuilder.Chainable{
							Elem: &stringbuilder.If{
								Condition: fmt.Sprintf(`err := %s.%s.%s(ctx, %s, %s.%s, "%s"); err != nil`,
									strct,
									validator,
									validateMimeTypes,
									mimeTypes,
									f.Args[1].Name,
									field.Name,
									field.JsonName,
								),
								Then: &stringbuilder.Chainable{
									Elem: stringbuilder.String("return nil, err"),
								},
							},
						})
						// str += fmt.Sprintf(
						// 	"if err := %s.%s.%s(ctx, %s, %s, %s); err != nil {",
						// 	usecaseValidator.GetMethodName(),
						// 	usecaseValidator.Fields[0].Name,
						// 	usecaseValidator.Fields[0].Type.(*model.Interface).Methods[0].Name,
						// 	f.Args[1].Name,
						// 	validation.Value,
						// 	field.Name,
						// ) + consts.LN
						// str += fmt.Sprintf("return nil, %s.%s.%s", usecaseValidator.GetMethodName(), usecaseValidator.Fields[0].Name, usecaseValidator.Fields[0].Type.(*model.Interface).Methods[0].Name) + consts.LN
						// str += "}" + consts.LN
					}
				}
			}

			strBuilder.Append(stringbuilder.String(fmt.Sprintf(
				"return %s.%s.%s(ctx, %s)",
				strct, usecaseValidator.Fields[1].Name, f.Name, f.Args[1].Name,
			)))
			// str += fmt.Sprintf("return %s.%s.%s(ctx, %s)", usecaseValidator.GetMethodName(), usecaseValidator.Fields[1].Name, f.Name, f.Args[1].Name) + consts.LN
			// return str, nil
			return strBuilder.String(), nil
		}
		usecaseValidator.Methods = append(usecaseValidator.Methods, f)
	}

	b.UsecaseValidatorImpl = usecaseValidator

	return usecaseValidator
}
