package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

// Builder for domain
// Store model definition to build and then build on Build method
// Only build on Build method as order is important
type domainBuilder struct {
	Definition *coredomaindefinition.Domain

	// default model fields are added to all models
	DefaultModelFields []*coredomaindefinition.Field

	// INTERNAL

	// if error happened during build, error is returned on Build method
	Err error

	ModelBuilders          []*ModelBuilder
	RepositoryBuilders     []*RepositoryBuilder
	GormRepositoryBuilders []*GormRepositoryBuilder

	// Store definition to build in order to build in normal order but call it in disered order
	ModelDefinitionsToBuild      []*coredomaindefinition.Model
	RepositoryDefinitionsToBuild []*coredomaindefinition.Repository
	RelationDefinitionsToBuild   []*coredomaindefinition.Relation
	CRUDToBuild                  []*coredomaindefinition.CRUD
	UsecaseDefinitionsToBuild    []*coredomaindefinition.Usecase

	// // Store definition to model
	ModelDefinitionToModel           map[*coredomaindefinition.Model]*model.Model
	ModelToModelDefinition           map[*model.Model]*coredomaindefinition.Model
	RepositoryDefinitionToRepository map[*coredomaindefinition.Repository]*model.Repository
	// UsecaseDefinitionToUsecase       map[*coredomaindefinition.Usecase]*model.Usecase
	RelationToRelationDefinition map[*model.Relation]*coredomaindefinition.Relation
	ModelToRepository            map[*model.Model]*model.Repository
	// RepoToDomainRepoGetEntityRepoMethod map[*model.Repository]*model.Function
	ModelToGormModel map[*model.Model]*model.GormModel
	Domain           *model.Domain
	Models           []*model.Model
	Repositories     []*model.Repository

	// Repository elements
	Ordering         *model.Struct
	Pagination       *model.Struct
	Transaction      *model.Interface
	RepositoryErrors map[string]*model.Var

	// // Usecase elements
	DomainUsecase *model.Interface
	// UsecaseCRUDImpl            *model.Struct
	FieldToParamUsecaseRequest map[*model.Field]*coredomaindefinition.Param
	CRUDActionToUsecase        map[*coredomaindefinition.CRUDAction]*model.Usecase
	FieldToValidationRules     map[*model.Field][]*coredomaindefinition.Validation
	// CRUDCreateRequest          map[*model.Model]*model.Struct
	// CRUDUpdateRequest          map[*model.Model]*model.Struct
	// CRUDDeleteRequest          map[*model.Model]*model.Struct
	// CRUDGetRequest             map[*model.Model]*model.Struct
	// CRUDListRequest            map[*model.Model]*model.Struct

	ModelUsecaseStruct map[*model.Model]*model.Struct

	// // Validator of struct request
	Validator            *model.Interface
	UsecaseValidatorImpl *model.Struct

	RelationGraph *RelationGraph

	GormDomainRepositoryBuilder *GormDomainRepositoryBuilder

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

		// RepositoryDefinitionToRepository: map[*coredomaindefinition.Repository]*model.Repository{},
		RelationToRelationDefinition: map[*model.Relation]*coredomaindefinition.Relation{},
		ModelDefinitionToModel:       map[*coredomaindefinition.Model]*model.Model{},
		ModelToModelDefinition:       map[*model.Model]*coredomaindefinition.Model{},
		ModelUsecaseStruct:           map[*model.Model]*model.Struct{},
		FieldToParamUsecaseRequest:   map[*model.Field]*coredomaindefinition.Param{},
		CRUDActionToUsecase:          map[*coredomaindefinition.CRUDAction]*model.Usecase{},
		// ModelToRepository:                map[*model.Model]*model.Repository{},
		// RepoToDomainRepoGetEntityRepoMethod: map[*model.Repository]*model.Function{},
		// RepositoryErrors:       map[string]*model.Var{},
		FieldToValidationRules: map[*model.Field][]*coredomaindefinition.Validation{},
		// ModelToGormModel:       map[*model.Model]*model.GormModel{},
		RelationGraph: &RelationGraph{},

		Domain: &model.Domain{
			Name: definition.Name,
		},
	}
	builder.setArchitecture(ctx)

	builder.builders = append(builder.builders, NewGormDomainRepositoryBuilder(ctx, builder, builder.Definition))

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

func (domainBuilder *domainBuilder) NewModelBuilder(ctx context.Context, modelDefinition *coredomaindefinition.Model) Builder {
	return NewModelBuilder(ctx, domainBuilder, modelDefinition, domainBuilder.DefaultModelFields)
}

func (domainBuilder *domainBuilder) NewRepositoryBuilder(ctx context.Context, repositoryDefinition *coredomaindefinition.Repository) Builder {
	return NewRepositoryBuilder(ctx, domainBuilder, repositoryDefinition)
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

func (domainBuilder *domainBuilder) GetSdkPackage() *model.GoPkg {
	return domainBuilder.Domain.Architecture.SdkPkg
}

func (builder *domainBuilder) GetModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) (*model.Model, error) {
	if m, ok := builder.ModelDefinitionToModel[modelDefinition]; ok {
		return m, nil
	}
	return nil, NewErrModelNotFound(modelDefinition.Name)
}

func (builder *domainBuilder) WithModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) *domainBuilder {
	if builder.Err != nil {
		return builder
	}

	builder.builders = append(builder.builders, builder.NewModelBuilder(ctx, modelDefinition))

	builder.ModelDefinitionsToBuild = append(builder.ModelDefinitionsToBuild, modelDefinition)

	return builder
}

func (builder *domainBuilder) WithRepository(ctx context.Context, repositoryDefinition *coredomaindefinition.Repository) *domainBuilder {
	if builder.Err != nil {
		return builder
	}

	builder.builders = append(builder.builders, builder.NewRepositoryBuilder(ctx, repositoryDefinition))
	builder.builders = append(builder.builders, builder.NewGormRepositoryBuilder(ctx, repositoryDefinition))

	builder.RepositoryDefinitionsToBuild = append(builder.RepositoryDefinitionsToBuild, repositoryDefinition)

	return builder
}

func (builder *domainBuilder) WithRelation(ctx context.Context, relationDefinition *coredomaindefinition.Relation) *domainBuilder {
	if builder.Err != nil {
		return builder
	}

	builder.RelationDefinitionsToBuild = append(builder.RelationDefinitionsToBuild, relationDefinition)
	builder.RelationGraph.addRelationToGraph(ctx, relationDefinition)

	return builder
}

func (builder *domainBuilder) WithCRUD(ctx context.Context, crudDefinition *coredomaindefinition.CRUD) *domainBuilder {
	if builder.Err != nil {
		return builder
	}

	builder.CRUDToBuild = append(builder.CRUDToBuild, crudDefinition)

	return builder
}

func (builder *domainBuilder) WithUsecase(ctx context.Context, usecaseDefinition *coredomaindefinition.Usecase) *domainBuilder {
	if builder.Err != nil {
		return builder
	}

	builder.UsecaseDefinitionsToBuild = append(builder.UsecaseDefinitionsToBuild, usecaseDefinition)

	return builder
}

func (builder *domainBuilder) Build(ctx context.Context) (*model.Domain, error) {
	if builder.Err != nil {
		return nil, builder.Err
	}

	builder.addRepositoryErrors(ctx)
	builder.addPagination(ctx)
	builder.addTransaction(ctx)
	builder.addOrdering(ctx)

	builder.builders = append(builder.builders, NewGormDomainRepositoryBuilder(ctx, builder, builder.Definition))

	for _, modelDefinition := range builder.ModelDefinitionsToBuild {
		for _, b := range builder.builders {
			b.WithModel(ctx, modelDefinition)
		}
	}
	for _, relationDefinition := range builder.RelationDefinitionsToBuild {
		for _, b := range builder.builders {
			b.WithRelation(ctx, relationDefinition)
		}
	}
	// for _, crudDefinition := range builder.CRUDToBuild {
	// 	builder.buildCRUD(ctx, crudDefinition)
	// }

	// for _, usecaseDefinition := range builder.UsecaseDefinitionsToBuild {
	// 	builder.buildUsecase(ctx, usecaseDefinition)
	// }

	for _, usecase := range builder.Domain.Usecases {
		usecase.Function.Args = append(usecase.Function.Args, &model.Param{
			Name: "ctx",
			Type: &model.PkgReference{
				Pkg: consts.CommonPkgs["context"],
				Reference: &model.ExternalType{
					Type: "Context",
				},
			},
		})
		usecase.Function.Args = append(usecase.Function.Args, &model.Param{
			Name: "request",
			Type: &model.PointerType{
				Type: &model.PkgReference{
					Pkg:       builder.Domain.Architecture.UsecasePkg,
					Reference: usecase.Request,
				},
			},
		})
		usecase.Function.Results = append(usecase.Function.Results, &model.Param{
			Type: &model.PointerType{
				Type: &model.PkgReference{
					Pkg:       builder.Domain.Architecture.UsecasePkg,
					Reference: usecase.Result,
				},
			},
		})
		usecase.Function.Results = append(usecase.Function.Results, &model.Param{
			Type: model.PrimitiveTypeError,
		})
		builder.GetDomainUsecase(ctx).Methods = append(builder.GetDomainUsecase(ctx).Methods, usecase.Function)
	}
	// builder.buildUsecaseCRUDImplementation(ctx)
	// builder.buildHttpController(ctx)
	// builder.buildGormAdpater(ctx)
	// builder.buildService(ctx)
	// builder.buildHttpService(ctx)

	// for _, modelBuilder := range builder.ModelBuilders {
	// 	m, err := modelBuilder.Build(ctx)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	builder.Domain.ModelsV2 = append(builder.Domain.ModelsV2, m)

	// }

	port, err := builder.buildRelationGraph(ctx)
	if err != nil {
		return nil, err
	}
	builder.Domain.Ports = append(builder.Domain.Ports, port)

	port, err = builder.buildDomainRepository(ctx)
	if err != nil {
		return nil, err
	}
	builder.Domain.Ports = append(builder.Domain.Ports, port)

	for _, b := range builder.builders {
		err := b.Build(ctx)
		if err != nil {
			return nil, err
		}
	}

	return builder.Domain, nil
}

// func (builder *domainBuilder) buildDomainRepository(ctx context.Context) {
// 	domainRepo := &model.Interface{
// 		Name:    GetDomainRepositoryName(ctx, builder.Definition),
// 		Methods: []*model.Function{},
// 	}

// 	for _, repositoryBuilder := range builder.RepositoryBuilders {
// 		domainRepo.Methods = append(domainRepo.Methods, repositoryBuilder.Methods...)
// 	}

// 	builder.Domain.Ports = append(builder.Domain.Ports, &model.File{
// 		Name: GetDomainRepositoryName(ctx, builder.Definition),
// 		Pkg:  builder.GetRepositoryPackage(),
// 		Elements: []interface{}{
// 			domainRepo,
// 		},
// 	})
// }
