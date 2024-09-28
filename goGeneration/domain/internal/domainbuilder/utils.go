package domainbuilder

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
)

func PluralizeName(ctx context.Context, name string) string {
	if strings.HasSuffix(name, "y") {
		return stringtool.UpperFirstLetter(name[:len(name)-1]) + "ies"
	}
	return stringtool.UpperFirstLetter(name) + "s"
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

func GetSingleRelationColumn(ctx context.Context, m *coredomaindefinition.Model) string {
	return stringtool.SnakeCase(GetSingleRelationIdName(ctx, m))
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
			coredomaindefinition.RelationTypeManyToMany, coredomaindefinition.RelationTypeOneToMany,
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
			coredomaindefinition.RelationTypeManyToOne,
		}, relation.Type)
	} else if relation.Target == m {
		return slices.Contains([]coredomaindefinition.RelationType{
			coredomaindefinition.RelationTypeManyToMany,
			coredomaindefinition.RelationTypeOneToMany,
			coredomaindefinition.RelationTypeSubresourcesOf,
			coredomaindefinition.RelationTypeBelongsTo,
		}, relation.Type)
	}
	return false
}

func GetRepositoryName(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return stringtool.UpperFirstLetter(definition.On.Name) + "Repository"
}

func GetRepositoryTableName(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("%s_TABLE_NAME", strings.ToUpper(definition.On.Name))
}

func GetRepositoryFieldToColumnName(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("%s_FIELD_TO_COLUMN", strings.ToUpper(definition.On.Name))
}

func GetColumnName(ctx context.Context, fieldDefinition *coredomaindefinition.Field) string {
	return stringtool.SnakeCase(fieldDefinition.Name)
}

func GetRepositoryAllowedOrderBy(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("%s_ALLOWED_ORDER_BY", strings.ToUpper(definition.On.Name))
}

func GetRepositoryAllowedWhere(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("%s_ALLOWED_WHERE", strings.ToUpper(definition.On.Name))
}

func GetRepositoryDefaultOrderBy(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("%s_DEFAULT_ORDER_BY", strings.ToUpper(definition.On.Name))
}

func GetRepositoryGetMethod(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("Get%s", stringtool.UpperFirstLetter(definition.On.Name))
}

func GetRepositoryListMethod(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("List%s", PluralizeName(ctx, definition.On.Name))
}

func GetRepositoryCreateMethod(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("Create%s", stringtool.UpperFirstLetter(definition.On.Name))
}

func GetRepositoryUpdateMethod(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("Update%s", stringtool.UpperFirstLetter(definition.On.Name))
}

func GetRepositoryDeleteMethod(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("Delete%s", stringtool.UpperFirstLetter(definition.On.Name))
}

func GetMethodContextName(ctx context.Context, methodName string) string {
	return fmt.Sprintf("%sContext", methodName)
}

func GetRepositoryMethodContextTransactionField(ctx context.Context) string {
	return stringtool.LowerFirstLetter(REPOSITORY_TRANSATION)
}

func GetRepositoryMethodContextWithInacctiveField(ctx context.Context) string {
	return "retriveInactive"
}
