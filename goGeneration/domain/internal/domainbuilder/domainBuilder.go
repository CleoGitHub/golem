package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

// Builder for domain
// Store model definition to build and then build on Build method
// Only build on Build method as order is important
type domainBuilder struct {
	*EmptyBuilder

	Definition *coredomaindefinition.Domain

	// default model fields are added to all models
	DefaultModelFields []*coredomaindefinition.Field

	// INTERNAL

	// if error happened during build, error is returned on Build method
	err error

	// Store definition to build in order to build in normal order but call it in disered order
	ModelDefinitionsToBuild      []*coredomaindefinition.Model
	RepositoryDefinitionsToBuild []*coredomaindefinition.Repository
	RelationDefinitionsToBuild   []*coredomaindefinition.Relation
	CRUDToBuild                  []*coredomaindefinition.CRUD
	UsecaseDefinitionsToBuild    []*coredomaindefinition.Usecase

	Domain *model.Domain

	RelationGraph *RelationGraph

	builders []Builder
}

func NewDomainBuilder(
	ctx context.Context,
	definition *coredomaindefinition.Domain,
	defaultModelFields []*coredomaindefinition.Field,
) *domainBuilder {
	builder := &domainBuilder{
		Definition:         definition,
		DefaultModelFields: defaultModelFields,
		RelationGraph:      &RelationGraph{},

		Domain: &model.Domain{
			Name: definition.Name,
		},
	}
	builder.setArchitecture(ctx)

	builder.builders = append(builder.builders, NewDomainRepositoryBuilder(ctx, builder))
	builder.builders = append(builder.builders, NewGormDomainRepositoryBuilder(ctx, builder, builder.Definition))
	builder.builders = append(builder.builders, NewDomainUsecaseBuilder(ctx, builder))
	builder.builders = append(builder.builders, NewGormDomainRepositoryBuilder(ctx, builder, builder.Definition))

	if definition.Controllers.Http {
		builder.builders = append(builder.builders, NewHttpControllerBuilder(ctx, definition, builder.Domain))
		builder.builders = append(builder.builders, NewHttpClientBuilder(ctx, definition, builder.Domain))
		builder.builders = append(builder.builders, NewJSBuilder(ctx, definition, builder))
	}

	return builder
}

func (builder *domainBuilder) setArchitecture(ctx context.Context) *domainBuilder {
	builder.Domain.Architecture = &model.Architecture{
		ModelPkg: &model.GoPkg{
			ShortName: "model",
			Alias:     "model",
			FullName: fmt.Sprintf(
				"%s/domain/model",
				builder.Definition.Name,
			),
		},
		RepositoryPkg: &model.GoPkg{
			ShortName: "repository",
			Alias:     "repository",
			FullName: fmt.Sprintf(
				"%s/domain/port/repository",
				builder.Definition.Name,
			),
		},
		UsecasePkg: &model.GoPkg{
			ShortName: "usecase",
			Alias:     "usecase",
			FullName: fmt.Sprintf(
				"%s/domain/usecase",
				builder.Definition.Name,
			),
		},
		ControllerPkg: &model.GoPkg{
			ShortName: "controller",
			Alias:     "controller",
			FullName: fmt.Sprintf(
				"%s/adapter/controller",
				builder.Definition.Name,
			),
		},
		HttpControllerPkg: &model.GoPkg{
			ShortName: "httpadapter",
			Alias:     "httpadapter",
			FullName: fmt.Sprintf(
				"%s/adapter/controller/httpadapter",
				builder.Definition.Name,
			),
		},
		GormAdapterPkg: &model.GoPkg{
			ShortName: "gormadapter",
			Alias:     "gormadapter",
			FullName: fmt.Sprintf(
				"%s/adapter/repository/gormadapter",
				builder.Definition.Name,
			),
		},
		SdkPkg: &model.GoPkg{
			ShortName: "client",
			Alias:     "client",
			FullName: fmt.Sprintf(
				"%s/sdk/client",
				builder.Definition.Name,
			),
		},
		JavascriptClient: fmt.Sprintf(
			"%s/sdk/JS",
			builder.Definition.Name,
		),
	}
	return builder
}

func (domainBuilder *domainBuilder) AddBuilder(ctx context.Context, builder Builder) {
	if domainBuilder.err != nil {
		return
	}
	for _, m := range domainBuilder.ModelDefinitionsToBuild {
		builder.WithModel(ctx, m)
	}

	for _, r := range domainBuilder.RepositoryDefinitionsToBuild {
		builder.WithRepository(ctx, r)
	}

	for _, r := range domainBuilder.RelationDefinitionsToBuild {
		builder.WithRelation(ctx, r)
	}

	for _, c := range domainBuilder.CRUDToBuild {
		builder.WithCRUD(ctx, c)
	}

	for _, c := range domainBuilder.UsecaseDefinitionsToBuild {
		builder.WithUsecase(ctx, c)
	}

	domainBuilder.builders = append(domainBuilder.builders, builder)
}

func (domainBuilder *domainBuilder) NewModelBuilder(ctx context.Context, modelDefinition *coredomaindefinition.Model) Builder {
	return NewModelBuilder(ctx, domainBuilder, modelDefinition, domainBuilder.DefaultModelFields)
}

func (domainBuilder *domainBuilder) NewGormRepositoryBuilder(ctx context.Context, repositoryDefinition *coredomaindefinition.Repository) Builder {
	return NewGormRepositoryBuilder(ctx, domainBuilder, repositoryDefinition)
}

func (domainBuilder *domainBuilder) GetModelPackage() *model.GoPkg {
	return domainBuilder.Domain.Architecture.ModelPkg
}

func (domainBuilder *domainBuilder) GetRepositoryPackage() *model.GoPkg {
	return domainBuilder.Domain.Architecture.RepositoryPkg
}

func (domainBuilder *domainBuilder) GetUsecasePackage() *model.GoPkg {
	return domainBuilder.Domain.Architecture.UsecasePkg
}

func (domainBuilder *domainBuilder) GetControllerPackage() *model.GoPkg {
	return domainBuilder.Domain.Architecture.ControllerPkg
}

func (domainBuilder *domainBuilder) GetGormAdapterPackage() *model.GoPkg {
	return domainBuilder.Domain.Architecture.GormAdapterPkg
}

func (domainBuilder *domainBuilder) GetHttpControllerPackage() *model.GoPkg {
	return domainBuilder.Domain.Architecture.HttpControllerPkg
}

func (domainBuilder *domainBuilder) GetSdkPackage() *model.GoPkg {
	return domainBuilder.Domain.Architecture.SdkPkg
}

func (builder *domainBuilder) WithModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) *domainBuilder {
	if builder.err != nil {
		return builder
	}

	builder.AddBuilder(ctx, builder.NewModelBuilder(ctx, modelDefinition))

	builder.ModelDefinitionsToBuild = append(builder.ModelDefinitionsToBuild, modelDefinition)

	for _, b := range builder.builders {
		b.WithModel(ctx, modelDefinition)
	}

	return builder
}

func (builder *domainBuilder) WithRepository(ctx context.Context, repositoryDefinition *coredomaindefinition.Repository) *domainBuilder {
	if builder.err != nil {
		return builder
	}

	builder.AddBuilder(ctx, builder.NewGormRepositoryBuilder(ctx, repositoryDefinition))

	builder.RepositoryDefinitionsToBuild = append(builder.RepositoryDefinitionsToBuild, repositoryDefinition)

	for _, b := range builder.builders {
		b.WithRepository(ctx, repositoryDefinition)
	}

	return builder
}

func (builder *domainBuilder) WithRelation(ctx context.Context, relationDefinition *coredomaindefinition.Relation) *domainBuilder {
	if builder.err != nil {
		return builder
	}

	builder.RelationDefinitionsToBuild = append(builder.RelationDefinitionsToBuild, relationDefinition)
	builder.RelationGraph.addRelationToGraph(ctx, relationDefinition)

	for _, b := range builder.builders {
		b.WithRelation(ctx, relationDefinition)
	}

	return builder
}

func (builder *domainBuilder) WithCRUD(ctx context.Context, crudDefinition *coredomaindefinition.CRUD) *domainBuilder {
	if builder.err != nil {
		return builder
	}

	builder.CRUDToBuild = append(builder.CRUDToBuild, crudDefinition)

	for _, b := range builder.builders {
		b.WithCRUD(ctx, crudDefinition)
	}

	return builder
}

func (builder *domainBuilder) WithUsecase(ctx context.Context, usecaseDefinition *coredomaindefinition.Usecase) *domainBuilder {
	if builder.err != nil {
		return builder
	}

	builder.UsecaseDefinitionsToBuild = append(builder.UsecaseDefinitionsToBuild, usecaseDefinition)

	for _, b := range builder.builders {
		b.WithUsecase(ctx, usecaseDefinition)
	}

	return builder
}

func (builder *domainBuilder) Build(ctx context.Context) (*model.Domain, error) {
	if builder.err != nil {
		return nil, builder.err
	}

	builder.addOrdering(ctx)
	builder.addPagination(ctx)

	port, err := builder.buildRelationGraph(ctx)
	if err != nil {
		return nil, err
	}
	builder.Domain.Files = append(builder.Domain.Files, port)

	for _, b := range builder.builders {
		err := b.Build(ctx)
		if err != nil {
			return nil, err
		}
	}

	return builder.Domain, nil
}

const (
	ORDERING_NAME        = "Ordering"
	ORDERING_ORDER       = "Order"
	ORDERING_ORDERBY     = "OrderBy"
	ORDERING_GET_ORDER   = "GetOrder"
	ORDERING_GET_ORDERBY = "GetOrderBy"
)

var ASC = &model.Var{
	Name:    "ASC",
	Type:    model.PrimitiveTypeString,
	Value:   "ASC",
	IsConst: true,
}
var DESC = &model.Var{
	Name:    "DESC",
	Type:    model.PrimitiveTypeString,
	Value:   "DESC",
	IsConst: true,
}

var ORDERING = &model.Struct{
	Name:       ORDERING_NAME,
	MethodName: stringtool.LowerFirstLetter(ORDERING_NAME),
	Consts: []*model.Var{
		ASC,
		DESC,
	},
	Fields: []*model.Field{
		{
			Name: ORDERING_ORDER,
			Type: model.PrimitiveTypeString,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{stringtool.LowerFirstLetter(ORDERING_ORDER)},
				},
			},
		},
		{
			Name: ORDERING_ORDERBY,
			Type: model.PrimitiveTypeString,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{stringtool.LowerFirstLetter(ORDERING_ORDERBY)},
				},
			},
		},
	},
	Methods: []*model.Function{
		{
			Name: "GetOrder",
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeString,
				},
			},
			Content: func() (string, []*model.GoPkg) {
				str := ""
				str += fmt.Sprintf(
					`if strings.ToUpper(%s.%s) != ASC && strings.ToUpper(%s.%s) != DESC { return ASC }`,
					stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDER, stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDER,
				) + consts.LN
				str += fmt.Sprintf("return  strings.ToUpper(%s.%s)", stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDER) + consts.LN
				return str, []*model.GoPkg{consts.CommonPkgs["strings"]}
			},
		}, {
			Name: "GetOrderBy",
			Args: []*model.Param{
				{
					Name: "allowedOrderBys",
					Type: &model.ArrayType{
						Type: model.PrimitiveTypeString,
					},
				}, {
					Name: "defaultOrderBy",
					Type: model.PrimitiveTypeString,
				},
			},
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeString,
				},
			},
			Content: func() (string, []*model.GoPkg) {
				str := ""
				str += fmt.Sprintf(
					`if slices.Contains(allowedOrderBys, %s.%s)  { return  %s.%s }`,
					stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDERBY,
					stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDERBY,
				) + consts.LN
				str += "return  defaultOrderBy" + consts.LN
				return str, []*model.GoPkg{consts.CommonPkgs["slices"]}
			},
		},
	},
}

func (b *domainBuilder) addOrdering(ctx context.Context) {
	if b.err != nil {
		return
	}

	b.Domain.Files = append(b.Domain.Files, &model.File{
		Name: ORDERING_NAME,
		Pkg:  b.GetModelPackage(),
		Elements: []interface{}{
			ORDERING,
		},
	})
}

const (
	PAGINATION_NAME                   = "Pagination"
	PAGINATION_Page                   = "Page"
	PAGINATION_ItemsPerPage           = "ItemsPerPage"
	PAGINATION_GetItemsPerPage        = "GetItemsPerPage"
	PAGINATION_GetPage                = "GetPage"
	PAGINATION_MIN_ITEMS_PER_PAGE     = "MIN_ITEMS_PER_PAGE"
	PAGINATION_MAX_ITEMS_PER_PAGE     = "MAX_ITEMS_PER_PAGE"
	PAGINATION_DEFAULT_ITEMS_PER_PAGE = "DEFAULT_ITEMS_PER_PAGE"
)

var PAGINATION = &model.Struct{
	Name:       PAGINATION_NAME,
	MethodName: stringtool.LowerFirstLetter(PAGINATION_NAME),
	Consts: []*model.Var{
		{
			Name:    PAGINATION_MIN_ITEMS_PER_PAGE,
			Type:    model.PrimitiveTypeInt,
			Value:   5,
			IsConst: true,
		},
		{
			Name:    PAGINATION_MAX_ITEMS_PER_PAGE,
			Type:    model.PrimitiveTypeInt,
			Value:   100,
			IsConst: true,
		},
		{
			Name:    PAGINATION_DEFAULT_ITEMS_PER_PAGE,
			Type:    model.PrimitiveTypeInt,
			Value:   30,
			IsConst: true,
		},
	},
	Fields: []*model.Field{
		{
			Name: PAGINATION_Page,
			Type: model.PrimitiveTypeInt,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{stringtool.LowerFirstLetter(PAGINATION_Page)},
				},
			},
		}, {
			Name: PAGINATION_ItemsPerPage,
			Type: model.PrimitiveTypeInt,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{stringtool.LowerFirstLetter(PAGINATION_GetItemsPerPage)},
				},
			},
		},
	},
	Methods: []*model.Function{
		{
			Name: PAGINATION_GetPage,
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeInt,
				},
			},
			Content: func() (string, []*model.GoPkg) {
				str := ""
				str += fmt.Sprintf("if %s.%s < 1 { return 1 }", stringtool.LowerFirstLetter(PAGINATION_NAME), PAGINATION_Page) + consts.LN
				str += fmt.Sprintf("return %s.%s", stringtool.LowerFirstLetter(PAGINATION_NAME), PAGINATION_Page) + consts.LN
				return str, nil
			},
		},
		{
			Name: "GetItemsPerPage",
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeInt,
				},
			},
			Content: func() (string, []*model.GoPkg) {
				str := ""
				str += fmt.Sprintf(
					"if %[1]s.%[2]s < %s || %[1]s.%[2]s > %[4]s { return %s }",
					stringtool.LowerFirstLetter(PAGINATION_NAME), PAGINATION_Page, PAGINATION_MIN_ITEMS_PER_PAGE, PAGINATION_MIN_ITEMS_PER_PAGE, PAGINATION_DEFAULT_ITEMS_PER_PAGE,
				) + consts.LN
				str += fmt.Sprintf("return %s.%s ", stringtool.LowerFirstLetter(PAGINATION_NAME), PAGINATION_Page) + consts.LN
				return str, nil
			},
		},
	},
}

func (builder *domainBuilder) addPagination(ctx context.Context) {
	if builder.err != nil {
		return
	}

	builder.Domain.Files = append(builder.Domain.Files, &model.File{
		Name: PAGINATION_NAME,
		Pkg:  builder.GetModelPackage(),
		Elements: []interface{}{
			WHERE_OPERATOR_TYPE,
			WHERE,
			WHERE_OPERATOR,
			PAGINATION,
		},
	})
}
