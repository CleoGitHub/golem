package domainbuilder

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	REPOSIOTY_METHOD_CONTEXT_OPTS_NAME = "opts"
)

func GetSingleRelationColumn(ctx context.Context, m *coredomaindefinition.Model) string {
	return stringtool.SnakeCase(GetSingleRelationIdName(ctx, m))
}

func GetManyToManyColumn(ctx context.Context, relation *coredomaindefinition.Relation) string {
	names := []string{relation.Source.Name, relation.Target.Name}
	slices.Sort(names)

	return fmt.Sprintf("%s_%s", stringtool.SnakeCase(names[0]), stringtool.SnakeCase(names[1]))
}

func GetRepositoryName(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return stringtool.UpperFirstLetter(definition.On.Name) + "Repository"
}

func GetRepositoryTableName(ctx context.Context, definition *coredomaindefinition.Repository) string {
	tableName := definition.TableName
	if tableName == "" {
		tableName = stringtool.SnakeCase(PluralizeName(ctx, definition.On.Name))
	}

	return tableName
}

func GetRepositoryConstTableName(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("%s_TABLE_NAME", strings.ToUpper(definition.On.Name))
}

func GetRepositoryFieldToColumnName(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("%s_FIELD_TO_COLUMN", strings.ToUpper(definition.On.Name))
}

func GetColumnName(ctx context.Context, fieldDefinition *coredomaindefinition.Field) string {
	return stringtool.SnakeCase(fieldDefinition.Name)
}

func GetColumnNameFromName(ctx context.Context, name string) string {
	return stringtool.SnakeCase(name)
}

func GetDomainRepositoryName(ctx context.Context, definition *coredomaindefinition.Domain) string {
	return stringtool.UpperFirstLetter(definition.Name) + "Repository"
}

func GetRepositoryAllowedOrderBy(ctx context.Context, definition *coredomaindefinition.Repository) string {
	return fmt.Sprintf("%s_ALLOWED_ORDER_BY", strings.ToUpper(definition.On.Name))
}

func GetRepositoryAllowedWhere(ctx context.Context, m *coredomaindefinition.Model) string {
	return fmt.Sprintf("%s_ALLOWED_WHERE", strings.ToUpper(m.Name))
}

func GetRepositoryDefaultOrderBy(ctx context.Context, m *coredomaindefinition.Model) string {
	return fmt.Sprintf("%s_DEFAULT_ORDER_BY", strings.ToUpper(m.Name))
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

func GetRepositoryAddRelationMethod(ctx context.Context, definition *coredomaindefinition.Repository, relation *coredomaindefinition.Relation) string {
	var to *coredomaindefinition.Model
	if relation.Source == definition.On {
		to = relation.Target
	} else {
		to = relation.Source
	}
	return fmt.Sprintf("Add%sTo%s", GetModelName(ctx, to), GetModelName(ctx, definition.On))
}

func GetRepositoryRemoveRelationMethod(ctx context.Context, definition *coredomaindefinition.Repository, relation *coredomaindefinition.Relation) string {
	var to *coredomaindefinition.Model
	if relation.Source == definition.On {
		to = relation.Target
	} else {
		to = relation.Source
	}
	return fmt.Sprintf("Remove%sTo%s", GetModelName(ctx, to), GetModelName(ctx, definition.On))
}

func GetRepositoryMethodContextTransactionField(ctx context.Context) string {
	return stringtool.UpperFirstLetter(TRANSACTION_NAME)
}

func GetRepositoryMethodOptionName(ctx context.Context, methodName string) string {
	return fmt.Sprintf("%sOpt", methodName)
}

func GetRepositoryGetSignature(ctx context.Context, repository *coredomaindefinition.Repository, repositoryPkg *model.GoPkg, modelPkg *model.GoPkg) *model.Function {
	methodName := GetRepositoryGetMethod(ctx, repository)
	return &model.Function{
		Name: methodName,
		Args: []*model.Param{
			{
				Name: "ctx",
				Type: &model.PkgReference{
					Pkg: consts.CommonPkgs["context"],
					Reference: &model.ExternalType{
						Type: "Context",
					},
				},
			},
			{
				Name: REPOSIOTY_METHOD_CONTEXT_OPTS_NAME,
				Type: &model.VariaidicType{
					Type: &model.PkgReference{
						Pkg: repositoryPkg,
						Reference: &model.ExternalType{
							Type: GetRepositoryMethodOptionName(ctx, methodName),
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: modelPkg,
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, repository.On),
						},
					},
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
	}
}

func GetRepositoryListSignature(ctx context.Context, repository *coredomaindefinition.Repository, repositoryPkg *model.GoPkg, modelPkg *model.GoPkg) *model.Function {
	methodName := GetRepositoryListMethod(ctx, repository)
	return &model.Function{
		Name: methodName,
		Args: []*model.Param{
			{
				Name: "ctx",
				Type: &model.PkgReference{
					Pkg: consts.CommonPkgs["context"],
					Reference: &model.ExternalType{
						Type: "Context",
					},
				},
			},
			{
				Name: REPOSIOTY_METHOD_CONTEXT_OPTS_NAME,
				Type: &model.VariaidicType{
					Type: &model.PkgReference{
						Pkg: repositoryPkg,
						Reference: &model.ExternalType{
							Type: GetRepositoryMethodOptionName(ctx, methodName),
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.ArrayType{
					Type: &model.PointerType{
						Type: &model.PkgReference{
							Pkg: modelPkg,
							Reference: &model.ExternalType{
								Type: GetModelName(ctx, repository.On),
							},
						},
					},
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
	}
}

func GetRepositoryCreateSignature(ctx context.Context, repository *coredomaindefinition.Repository, repositoryPkg *model.GoPkg, modelPkg *model.GoPkg) *model.Function {
	methodName := GetRepositoryCreateMethod(ctx, repository)
	return &model.Function{
		Name: methodName,
		Args: []*model.Param{
			{
				Name: "ctx",
				Type: &model.PkgReference{
					Pkg: consts.CommonPkgs["context"],
					Reference: &model.ExternalType{
						Type: "Context",
					},
				},
			},
			{
				Name: "entity",
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: modelPkg,
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, repository.On),
						},
					},
				},
			},
			{
				Name: REPOSIOTY_METHOD_CONTEXT_OPTS_NAME,
				Type: &model.VariaidicType{
					Type: &model.PkgReference{
						Pkg: repositoryPkg,
						Reference: &model.ExternalType{
							Type: GetRepositoryMethodOptionName(ctx, methodName),
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: modelPkg,
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, repository.On),
						},
					},
				},
			},
		},
	}
}

func GetRepositoryUpdateSignature(ctx context.Context, repository *coredomaindefinition.Repository, repositoryPkg *model.GoPkg, modelPkg *model.GoPkg) *model.Function {
	methodName := GetRepositoryUpdateMethod(ctx, repository)
	return &model.Function{
		Name: methodName,
		Args: []*model.Param{
			{
				Name: "ctx",
				Type: &model.PkgReference{
					Pkg: consts.CommonPkgs["context"],
					Reference: &model.ExternalType{
						Type: "Context",
					},
				},
			},
			{
				Name: "entity",
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: modelPkg,
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, repository.On),
						},
					},
				},
			},
			{
				Name: REPOSIOTY_METHOD_CONTEXT_OPTS_NAME,
				Type: &model.VariaidicType{
					Type: &model.PkgReference{
						Pkg: repositoryPkg,
						Reference: &model.ExternalType{
							Type: GetRepositoryMethodOptionName(ctx, methodName),
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: modelPkg,
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, repository.On),
						},
					},
				},
			},
		},
	}
}

func GetRepositoryDeleteSignature(ctx context.Context, repository *coredomaindefinition.Repository, repositoryPkg *model.GoPkg, modelPkg *model.GoPkg) *model.Function {
	methodName := GetRepositoryDeleteMethod(ctx, repository)
	return &model.Function{
		Name: methodName,
		Args: []*model.Param{
			{
				Name: "ctx",
				Type: &model.PkgReference{
					Pkg: consts.CommonPkgs["context"],
					Reference: &model.ExternalType{
						Type: "Context",
					},
				},
			},
			{
				Name: "id",
				Type: model.PrimitiveTypeString,
			},
			{
				Name: REPOSIOTY_METHOD_CONTEXT_OPTS_NAME,
				Type: &model.VariaidicType{
					Type: &model.PkgReference{
						Pkg: repositoryPkg,
						Reference: &model.ExternalType{
							Type: GetRepositoryMethodOptionName(ctx, methodName),
						},
					},
				},
			},
		},
	}
}

func GetRepositoryRelationNodeName(ctx context.Context, m *coredomaindefinition.Model) string {
	return GetModelName(ctx, m) + "RelationNode"
}

func GetPreloadName(ctx context.Context, m *coredomaindefinition.Model) string {
	return "Preload" + GetModelName(ctx, m)
}

func GetGormDomainRepositoryName(ctx context.Context, m *coredomaindefinition.Domain) string {
	return stringtool.UpperFirstLetter(m.Name) + "Repository"
}
