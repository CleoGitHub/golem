package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleoGitHub/golem-common/pkg/merror"
	"github.com/cleoGitHub/golem-common/pkg/stringtool"
	"github.com/cleoGitHub/golem/coredomaindefinition"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
)

func (b *domainBuilder) buildUsecase(ctx context.Context, usecaseDefinition *coredomaindefinition.Usecase) *domainBuilder {
	if b.Err != nil {
		return b
	}

	usecase := &model.Usecase{
		Function: &model.Function{
			Name: fmt.Sprintf("%sUsecase", usecaseDefinition.Name),
		},
		Request: &model.Struct{Name: usecaseDefinition.Name + "Request"},
		Result:  &model.Struct{Name: usecaseDefinition.Name + "Result"},
		Roles:   usecaseDefinition.Roles,
	}

	for _, param := range usecaseDefinition.Args {
		t, err := TypeDefinitionToType(ctx, param.Type)
		if err != nil {
			b.Err = merror.Stack(err)
			return b
		}
		f := &model.Field{
			Name:     stringtool.UpperFirstLetter(param.Name),
			Type:     t,
			JsonName: param.Name,
			Tags: []*model.Tag{{
				Name:   "json",
				Values: []string{param.Name},
			}},
		}
		validationsValues := []string{}
		for _, validation := range param.Validations {
			switch validation.Rule {
			case coredomaindefinition.ValidationRuleRequired:
				validationsValues = append(validationsValues, "required")
			case coredomaindefinition.ValidationRuleEmail:
				validationsValues = append(validationsValues, "email")
			case coredomaindefinition.ValidationRuleUUID:
				validationsValues = append(validationsValues, "uuid")
			case coredomaindefinition.ValidationRuleHexColor:
				validationsValues = append(validationsValues, "hexcolor")
			case coredomaindefinition.ValidationRuleGT:
				if s, ok := validation.Value.(string); ok {
					validationsValues = append(validationsValues, "gt:"+s)
				} else {
					b.Err = NewErrValidationValueExpectedType(string(validation.Rule), "string")
					return b
				}
			case coredomaindefinition.ValidationRuleGTE:
				if s, ok := validation.Value.(string); ok {
					validationsValues = append(validationsValues, "gte:"+s)
				} else {
					b.Err = NewErrValidationValueExpectedType(string(validation.Rule), "string")
					return b
				}
			case coredomaindefinition.ValidationRuleLT:
				if s, ok := validation.Value.(string); ok {
					validationsValues = append(validationsValues, "lt:"+s)
				} else {
					b.Err = NewErrValidationValueExpectedType(string(validation.Rule), "string")
					return b
				}
			case coredomaindefinition.ValidationRuleLTE:
				if s, ok := validation.Value.(string); ok {
					validationsValues = append(validationsValues, "lte:"+s)
				} else {
					b.Err = NewErrValidationValueExpectedType(string(validation.Rule), "string")
					return b
				}
			case coredomaindefinition.ValidationRuleUnique, coredomaindefinition.ValidationRuleUniqueIn:
				b.Err = NewErrUnexpectedValidationRule(string(validation.Rule))
				return b
			}
		}

		if len(validationsValues) > 0 {
			f.Tags = append(f.Tags, &model.Tag{
				Name:   "validate",
				Values: validationsValues,
			})
		}
		b.FieldToParamUsecaseRequest[f] = param
		usecase.Request.Fields = append(usecase.Request.Fields, f)
	}

	b.Domain.Usecases = append(b.Domain.Usecases, usecase)
	b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Request)
	b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Result)

	return b
}
