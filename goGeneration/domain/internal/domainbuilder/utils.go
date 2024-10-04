package domainbuilder

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func PluralizeName(ctx context.Context, name string) string {
	if strings.HasSuffix(name, "y") {
		return stringtool.UpperFirstLetter(name[:len(name)-1]) + "ies"
	}
	return stringtool.UpperFirstLetter(name) + "s"
}

func GetMethodName(ctx context.Context, s *model.Struct) string {
	return stringtool.LowerFirstLetter(s.Name)
}

func GetFieldName(ctx context.Context, fieldDefinition *coredomaindefinition.Field) string {
	return stringtool.UpperFirstLetter(fieldDefinition.Name)
}

func GetModelName(ctx context.Context, modelDefinition *coredomaindefinition.Model) string {
	return stringtool.UpperFirstLetter(modelDefinition.Name)
}

func GetSingleRelationName(ctx context.Context, m *coredomaindefinition.Model) string {
	return stringtool.UpperFirstLetter(m.Name)
}

func GetSingleRelationIdName(ctx context.Context, m *coredomaindefinition.Model) string {
	return stringtool.UpperFirstLetter(m.Name) + "Id"
}

func GetMultipleRelationName(ctx context.Context, m *coredomaindefinition.Model) string {
	return PluralizeName(ctx, m.Name)
}

func GetMultipleRelationIdsName(ctx context.Context, m *coredomaindefinition.Model) string {
	return stringtool.UpperFirstLetter(m.Name) + "Ids"
}

func IsRelationOptionnal(ctx context.Context, m *coredomaindefinition.Model, relation *coredomaindefinition.Relation) (bool, error) {
	if relation.Source == m {
		return slices.Contains([]coredomaindefinition.RelationType{
			coredomaindefinition.RelationTypeManyToMany, coredomaindefinition.RelationTypeManyToOne,
		}, relation.Type), nil
	} else if relation.Target == m {
		return slices.Contains([]coredomaindefinition.RelationType{
			coredomaindefinition.RelationTypeManyToMany, coredomaindefinition.RelationTypeOneToMany,
		}, relation.Type), nil
	} else {
		return false, ErrRelationDoesNotBelongToModel
	}
}

func IsRelationMultiple(ctx context.Context, m *coredomaindefinition.Model, relation *coredomaindefinition.Relation) bool {
	if relation.Source == m {
		return slices.Contains([]coredomaindefinition.RelationType{
			coredomaindefinition.RelationTypeManyToMany,
			coredomaindefinition.RelationTypeOneToMany,
		}, relation.Type)
	} else if relation.Target == m {
		return slices.Contains([]coredomaindefinition.RelationType{
			coredomaindefinition.RelationTypeManyToMany,
			coredomaindefinition.RelationTypeManyToOne,
			coredomaindefinition.RelationTypeSubresourcesOf,
			coredomaindefinition.RelationTypeBelongsTo,
		}, relation.Type)
	}
	return false
}

func GetMethodContextName(ctx context.Context, methodName string) string {
	return fmt.Sprintf("%sContext", methodName)
}
