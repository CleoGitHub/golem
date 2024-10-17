package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
)

func GetDomainUsecaseName(ctx context.Context, name string) string {
	return stringtool.UpperFirstLetter(name) + "Usecase"
}

func GetCRUDMethodName(ctx context.Context, action string, model *coredomaindefinition.Model) string {
	return GetUsecaseMethodName(ctx, fmt.Sprintf("%s%s", stringtool.UpperFirstLetter(action), stringtool.UpperFirstLetter(model.Name)))
}

func GetUsecaseMethodName(ctx context.Context, action string) string {
	return stringtool.UpperFirstLetter(action)
}

func GetUsecaseRequestName(ctx context.Context, action string) string {
	return fmt.Sprintf("%sRequest", stringtool.UpperFirstLetter(action))
}

func GetUsecaseResponseName(ctx context.Context, action string) string {
	return fmt.Sprintf("%sResponse", stringtool.UpperFirstLetter(action))
}

func GetCRUDRelationMethodName(ctx context.Context, action string, from *coredomaindefinition.Model, to *coredomaindefinition.Model) string {
	switch action {
	case ADD:
		return fmt.Sprintf("Add%sTo%s", GetModelName(ctx, to), GetModelName(ctx, from))
	case REMOVE:
		return fmt.Sprintf("Remove%sFrom%s", GetModelName(ctx, to), GetModelName(ctx, from))
	case LIST:
		return fmt.Sprintf("List%sOf%s", PluralizeName(ctx, GetModelName(ctx, to)), GetModelName(ctx, from))
	default:
		return fmt.Sprintf("%s%sTo%s", stringtool.UpperFirstLetter(action), GetModelName(ctx, to), GetModelName(ctx, from))
	}
}

func GetValidationTags(ctx context.Context, validations []*coredomaindefinition.Validation) ([]string, error) {
	tags := make([]string, 0)
	for _, validation := range validations {
		switch validation.Rule {
		case coredomaindefinition.ValidationRuleRequired:
			tags = append(tags, "required")
		case coredomaindefinition.ValidationRuleEmail:
			tags = append(tags, "email")
		case coredomaindefinition.ValidationRuleUUID:
			tags = append(tags, "uuid")
		case coredomaindefinition.ValidationRuleHexColor:
			tags = append(tags, "hexcolor")
		case coredomaindefinition.ValidationRuleGT:
			if s, ok := validation.Value.(string); ok {
				tags = append(tags, "gt:"+s)
			} else {
				return nil, NewErrValidationValueExpectedType(string(validation.Rule), "string")
			}
		case coredomaindefinition.ValidationRuleGTE:
			if s, ok := validation.Value.(string); ok {
				tags = append(tags, "gte:"+s)
			} else {
				return nil, NewErrValidationValueExpectedType(string(validation.Rule), "string")
			}
		case coredomaindefinition.ValidationRuleLT:
			if s, ok := validation.Value.(string); ok {
				tags = append(tags, "lt:"+s)
			} else {
				return nil, NewErrValidationValueExpectedType(string(validation.Rule), "string")
			}
		case coredomaindefinition.ValidationRuleLTE:
			if s, ok := validation.Value.(string); ok {
				tags = append(tags, "lte:"+s)
			} else {
				return nil, NewErrValidationValueExpectedType(string(validation.Rule), "string")
			}
		}
	}
	return tags, nil
}
