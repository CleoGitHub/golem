package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

type DomainRepositoryBuilder struct {
	*EmptyBuilder

	domainBuilder *domainBuilder

	repositoryBuilders []*RepositoryBuilder

	File                 *model.File
	repositoryDefinition *model.Interface

	Err error

	repositoryDefinitions []*coredomaindefinition.Repository
	relationDefinitions   []*coredomaindefinition.Relation
	modelsDefinitions     []*coredomaindefinition.Model
}

const (
	REPOSITORY_MIGRATE           = "Migrate"
	REPOSITORY_BEGIN_TRANSACTION = "BeginTransaction"
	REPOSITORY_COMMIT            = "Commit"
	REPOSITORY_ROLLBACK          = "Rollback"
)

func NewDomainRepositoryBuilder(ctx context.Context, domainBuilder *domainBuilder) Builder {
	builder := &DomainRepositoryBuilder{
		domainBuilder: domainBuilder,
	}

	repoName := GetDomainRepositoryName(ctx, builder.domainBuilder.Definition)
	builder.repositoryDefinition = &model.Interface{
		Name: repoName,
		Methods: []*model.Function{
			{
				Name: REPOSITORY_MIGRATE,
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
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
			{
				Name: REPOSITORY_BEGIN_TRANSACTION,
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
				},
				Results: []*model.Param{
					{
						Type: &model.PkgReference{
							Pkg: builder.domainBuilder.GetRepositoryPackage(),
							Reference: &model.ExternalType{
								Type: TRANSACTION_NAME,
							},
						},
					},
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
		},
	}

	// for _, repositoryBuilder := range builder.repositoryBuilders {
	// 	builder.repositoryDefinition.Methods = append(builder.repositoryDefinition.Methods, repositoryBuilder.Methods...)
	// }

	builder.File = &model.File{
		Name: repoName,
		Pkg:  builder.domainBuilder.GetRepositoryPackage(),
		Elements: []interface{}{
			builder.repositoryDefinition,
		},
	}
	return builder

}

func (builder *DomainRepositoryBuilder) WithModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) {
	if builder.Err != nil {
		return
	}

	builder.modelsDefinitions = append(builder.modelsDefinitions, modelDefinition)
}

func (builder *DomainRepositoryBuilder) WithRelation(ctx context.Context, relationDefinition *coredomaindefinition.Relation) {
	if builder.Err != nil {
		return
	}

	builder.relationDefinitions = append(builder.relationDefinitions, relationDefinition)
}

func (builder *DomainRepositoryBuilder) WithRepository(ctx context.Context, repositoryDefinition *coredomaindefinition.Repository) {
	if builder.Err != nil {
		return
	}

	builder.repositoryDefinitions = append(builder.repositoryDefinitions, repositoryDefinition)
	builder.repositoryBuilders = append(builder.repositoryBuilders, NewRepositoryBuilder(ctx, builder.domainBuilder, repositoryDefinition).(*RepositoryBuilder))
}

func (builder *DomainRepositoryBuilder) Build(ctx context.Context) error {
	if builder.Err != nil {
		return builder.Err
	}

	for _, modelDefinition := range builder.modelsDefinitions {
		for _, b := range builder.repositoryBuilders {
			b.WithModel(ctx, modelDefinition)
		}
	}

	for _, relationDefinition := range builder.relationDefinitions {
		for _, b := range builder.repositoryBuilders {
			b.WithRelation(ctx, relationDefinition)
		}
	}

	for _, repositoryBuilder := range builder.repositoryBuilders {
		if err := repositoryBuilder.Build(ctx); err != nil {
			return err
		}

		builder.repositoryDefinition.Methods = append(builder.repositoryDefinition.Methods, repositoryBuilder.Methods...)
	}

	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, builder.File)

	// builder.addOrdering(ctx)
	builder.addRepositoryErrors(ctx)
	builder.addWhere(ctx)
	builder.addTransaction(ctx)

	return builder.Err
}

var REPOSITORY_ERROR_NOT_FOUND = &model.Var{
	Name: "ErrNotFound",
	Value: &model.PkgReference{
		Pkg: consts.CommonPkgs["fmt"],
		Reference: &model.ExternalType{
			Type: `Errorf("not found")`,
		},
	},
}

func (b *DomainRepositoryBuilder) addRepositoryErrors(ctx context.Context) {
	if b.Err != nil {
		return
	}

	b.domainBuilder.Domain.Files = append(b.domainBuilder.Domain.Files, &model.File{
		Name: "errors",
		Pkg:  b.domainBuilder.GetRepositoryPackage(),
		Elements: []interface{}{
			REPOSITORY_ERROR_NOT_FOUND,
		},
	})
}

const (
	REPOSITORY_WHERE                    = "Where"
	REPOSITORY_WHERE_KEY                = "Key"
	REPOSITORY_WHERE_OPERATOR           = "Operator"
	REPOSITORY_WHERE_VALUE              = "Value"
	REPOSITORY_WHERE_OPERATOR_EQUAL     = "EQUAL"
	REPOSITORY_WHERE_OPERATOR_NOT_EQUAL = "NOT_EQUAL"
	REPOSITORY_WHERE_OPERATOR_IN        = "IN"
	REPOSITORY_WHERE_OPERATOR_NOT_IN    = "NOT_IN"
	REPOSITORY_WHERE_OPERATOR_TYPE      = "WHERE_OPERATOR"
)

var WHERE_OPERATOR_TYPE = &model.TypeDefinition{
	Name: REPOSITORY_WHERE_OPERATOR_TYPE,
	Type: model.PrimitiveTypeString,
}

var WHERE_OPERATOR = &model.Enum{
	Name: "Operator",
	Type: WHERE_OPERATOR_TYPE,
	Values: map[string]interface{}{
		REPOSITORY_WHERE_OPERATOR_EQUAL:     REPOSITORY_WHERE_OPERATOR_EQUAL,
		REPOSITORY_WHERE_OPERATOR_NOT_EQUAL: REPOSITORY_WHERE_OPERATOR_NOT_EQUAL,
		REPOSITORY_WHERE_OPERATOR_IN:        REPOSITORY_WHERE_OPERATOR_IN,
		REPOSITORY_WHERE_OPERATOR_NOT_IN:    REPOSITORY_WHERE_OPERATOR_NOT_IN,
	},
}

var WHERE = &model.Struct{
	Name:       REPOSITORY_WHERE,
	MethodName: stringtool.LowerFirstLetter(REPOSITORY_WHERE),
	Fields: []*model.Field{
		{
			Name: REPOSITORY_WHERE_KEY,
			Type: model.PrimitiveTypeString,
		},
		{

			Name: REPOSITORY_WHERE_OPERATOR,
			Type: WHERE_OPERATOR_TYPE,
		},
		{
			Name: REPOSITORY_WHERE_VALUE,
			Type: model.PrimitiveTypeInterface,
		},
	},
}

func (b *DomainRepositoryBuilder) addWhere(ctx context.Context) {
	if b.Err != nil {
		return
	}

	b.domainBuilder.Domain.Files = append(b.domainBuilder.Domain.Files, &model.File{
		Name: PAGINATION_NAME,
		Pkg:  b.domainBuilder.GetRepositoryPackage(),
		Elements: []interface{}{
			WHERE_OPERATOR_TYPE,
			WHERE,
			WHERE_OPERATOR,
		},
	})
}

const (
	TRANSACTION_NAME     = "Transaction"
	TRANSACTION_GET      = "Get"
	TRANSACTION_COMMIT   = "Commit"
	TRANSACTION_ROLLBACK = "Rollback"
)

var TRANSACTION = &model.Interface{
	Name: TRANSACTION_NAME,
	Methods: []*model.Function{
		{
			Name: TRANSACTION_GET,
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
			},
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeInterface,
				},
			},
		},
		{
			Name: TRANSACTION_COMMIT,
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
			},
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeError,
				},
			},
		},
		{
			Name: TRANSACTION_ROLLBACK,
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
			},
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeError,
				},
			},
		},
	},
}

func (b *DomainRepositoryBuilder) addTransaction(ctx context.Context) {
	if b.Err != nil {
		return
	}

	b.domainBuilder.Domain.Files = append(b.domainBuilder.Domain.Files, &model.File{
		Name: TRANSACTION_NAME,
		Pkg:  b.domainBuilder.GetRepositoryPackage(),
		Elements: []interface{}{
			TRANSACTION,
		},
	})
}

// const (
// 	ORDERING_NAME        = "Ordering"
// 	ORDERING_ORDER       = "Order"
// 	ORDERING_ORDERBY     = "OrderBy"
// 	ORDERING_GET_ORDER   = "GetOrder"
// 	ORDERING_GET_ORDERBY = "GetOrderBy"
// )

// var ASC = &model.Var{
// 	Name:    "ASC",
// 	Type:    model.PrimitiveTypeString,
// 	Value:   "ASC",
// 	IsConst: true,
// }
// var DESC = &model.Var{
// 	Name:    "DESC",
// 	Type:    model.PrimitiveTypeString,
// 	Value:   "DESC",
// 	IsConst: true,
// }

// var ORDERING = &model.Struct{
// 	Name:       ORDERING_NAME,
// 	MethodName: stringtool.LowerFirstLetter(ORDERING_NAME),
// 	Consts: []*model.Var{
// 		ASC,
// 		DESC,
// 	},
// 	Fields: []*model.Field{
// 		{
// 			Name: ORDERING_ORDER,
// 			Type: model.PrimitiveTypeString,
// 			Tags: []*model.Tag{
// 				{
// 					Name:   "json",
// 					Values: []string{stringtool.LowerFirstLetter(ORDERING_ORDER)},
// 				},
// 			},
// 		},
// 		{
// 			Name: ORDERING_ORDERBY,
// 			Type: model.PrimitiveTypeString,
// 			Tags: []*model.Tag{
// 				{
// 					Name:   "json",
// 					Values: []string{stringtool.LowerFirstLetter(ORDERING_ORDERBY)},
// 				},
// 			},
// 		},
// 	},
// 	Methods: []*model.Function{
// 		{
// 			Name: "GetOrder",
// 			Results: []*model.Param{
// 				{
// 					Type: model.PrimitiveTypeString,
// 				},
// 			},
// 			Content: func() (string, []*model.GoPkg) {
// 				str := ""
// 				str += fmt.Sprintf(
// 					`if strings.ToUpper(%s.%s) != ASC && strings.ToUpper(%s.%s) != DESC { return ASC }`,
// 					stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDER, stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDER,
// 				) + consts.LN
// 				str += fmt.Sprintf("return  strings.ToUpper(%s.%s)", stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDER) + consts.LN
// 				return str, []*model.GoPkg{consts.CommonPkgs["strings"]}
// 			},
// 		}, {
// 			Name: "GetOrderBy",
// 			Args: []*model.Param{
// 				{
// 					Name: "allowedOrderBys",
// 					Type: &model.ArrayType{
// 						Type: model.PrimitiveTypeString,
// 					},
// 				}, {
// 					Name: "defaultOrderBy",
// 					Type: model.PrimitiveTypeString,
// 				},
// 			},
// 			Results: []*model.Param{
// 				{
// 					Type: model.PrimitiveTypeString,
// 				},
// 			},
// 			Content: func() (string, []*model.GoPkg) {
// 				str := ""
// 				str += fmt.Sprintf(
// 					`if slices.Contains(allowedOrderBys, %s.%s)  { return  %s.%s }`,
// 					stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDERBY,
// 					stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDERBY,
// 				) + consts.LN
// 				str += "return  defaultOrderBy" + consts.LN
// 				return str, []*model.GoPkg{consts.CommonPkgs["slices"]}
// 			},
// 		},
// 	},
// }

// func (b *DomainRepositoryBuilder) addOrdering(ctx context.Context) {
// 	if b.Err != nil {
// 		return
// 	}

// 	b.domainBuilder.Domain.Files = append(b.domainBuilder.Domain.Files, &model.File{
// 		Name: ORDERING_NAME,
// 		Pkg:  b.domainBuilder.GetRepositoryPackage(),
// 		Elements: []interface{}{
// 			ORDERING,
// 		},
// 	})
// }
