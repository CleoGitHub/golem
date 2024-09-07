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
	// default model fields are added to all models
	DefaultModelFields []*coredomaindefinition.Field

	// INTERNAL

	// if error happened during build, error is returned on Build method
	Err error

	// Store definition to build in order to build in normal order but call it in disered order
	ModelDefinitionsToBuild      []*coredomaindefinition.Model
	RepositoryDefinitionsToBuild []*coredomaindefinition.Repository
	RelationDefinitionsToBuild   []*coredomaindefinition.Relation
	CRUDToBuild                  []*coredomaindefinition.CRUD
	UsecaseDefinitionsToBuild    []*coredomaindefinition.Usecase

	// Store definition to model
	ModelDefinitionToModel           map[*coredomaindefinition.Model]*model.Model
	ModelToModelDefinition           map[*model.Model]*coredomaindefinition.Model
	RepositoryDefinitionToRepository map[*coredomaindefinition.Repository]*model.Repository
	UsecaseDefinitionToUsecase       map[*coredomaindefinition.Usecase]*model.Usecase
	RelationToRelationDefinition     map[*model.Relation]*coredomaindefinition.Relation
	ModelToRepository                map[*model.Model]*model.Repository
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

	// Usecase elements
	DomainUsecase              *model.Interface
	UsecaseCRUDImpl            *model.Struct
	FieldToParamUsecaseRequest map[*model.Field]*coredomaindefinition.Param
	CRUDActionToUsecase        map[*coredomaindefinition.CRUDAction]*model.Usecase
	FieldToValidationRules     map[*model.Field][]*coredomaindefinition.Validation
	CRUDCreateRequest          map[*model.Model]*model.Struct
	CRUDUpdateRequest          map[*model.Model]*model.Struct
	CRUDDeleteRequest          map[*model.Model]*model.Struct
	CRUDGetRequest             map[*model.Model]*model.Struct
	CRUDListRequest            map[*model.Model]*model.Struct

	ModelUsecaseStruct map[*model.Model]*model.Struct

	// Validator of struct request
	Validator            *model.Interface
	UsecaseValidatorImpl *model.Struct
}

func NewDomainBuilder(
	d *coredomaindefinition.Domain,
	defaultModelFields []*coredomaindefinition.Field,
) *domainBuilder {
	return &domainBuilder{
		DefaultModelFields:               defaultModelFields,
		RepositoryDefinitionToRepository: map[*coredomaindefinition.Repository]*model.Repository{},
		RelationToRelationDefinition:     map[*model.Relation]*coredomaindefinition.Relation{},
		Models:                           []*model.Model{},
		Repositories:                     []*model.Repository{},
		ModelDefinitionToModel:           map[*coredomaindefinition.Model]*model.Model{},
		ModelToModelDefinition:           map[*model.Model]*coredomaindefinition.Model{},
		ModelUsecaseStruct:               map[*model.Model]*model.Struct{},
		FieldToParamUsecaseRequest:       map[*model.Field]*coredomaindefinition.Param{},
		CRUDActionToUsecase:              map[*coredomaindefinition.CRUDAction]*model.Usecase{},
		ModelToRepository:                map[*model.Model]*model.Repository{},
		// RepoToDomainRepoGetEntityRepoMethod: map[*model.Repository]*model.Function{},
		RepositoryErrors:       map[string]*model.Var{},
		FieldToValidationRules: map[*model.Field][]*coredomaindefinition.Validation{},
		ModelToGormModel:       map[*model.Model]*model.GormModel{},

		Domain: &model.Domain{
			Name: d.Name,
			Architecture: &model.Architecture{
				ModelPkg: &model.GoPkg{
					ShortName: "model",
					Alias:     "model",
					FullName: fmt.Sprintf(
						"%s/domain/model",
						d.Name,
					),
				},
				RepositoryPkg: &model.GoPkg{
					ShortName: "repository",
					Alias:     "repository",
					FullName: fmt.Sprintf(
						"%s/domain/port/repository",
						d.Name,
					),
				},
				UsecasePkg: &model.GoPkg{
					ShortName: "usecase",
					Alias:     "usecase",
					FullName: fmt.Sprintf(
						"%s/domain/usecase",
						d.Name,
					),
				},
				ControllerPkg: &model.GoPkg{
					ShortName: "controller",
					Alias:     "controller",
					FullName: fmt.Sprintf(
						"%s/adapter/controller",
						d.Name,
					),
				},
				GormAdapterPkg: &model.GoPkg{
					ShortName: "gormadapter",
					Alias:     "gormadapter",
					FullName: fmt.Sprintf(
						"%s/adapter/repository/gormadapter",
						d.Name,
					),
				},
				SdkPkg: &model.GoPkg{
					ShortName: "client",
					Alias:     "client",
					FullName: fmt.Sprintf(
						"%s/sdk/client",
						d.Name,
					),
				},
				JavascriptClient: fmt.Sprintf(
					"%s/sdk/JS",
					d.Name,
				),
			},
		},
	}
}

func (b *domainBuilder) GetModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) (*model.Model, error) {
	if m, ok := b.ModelDefinitionToModel[modelDefinition]; ok {
		return m, nil
	}
	return nil, NewErrModelNotFound(modelDefinition.Name)
}

func (b *domainBuilder) WithModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) *domainBuilder {
	if b.Err != nil {
		return b
	}

	b.ModelDefinitionsToBuild = append(b.ModelDefinitionsToBuild, modelDefinition)

	return b
}

func (b *domainBuilder) WithRepository(ctx context.Context, repositoryDefinition *coredomaindefinition.Repository) *domainBuilder {
	if b.Err != nil {
		return b
	}

	b.RepositoryDefinitionsToBuild = append(b.RepositoryDefinitionsToBuild, repositoryDefinition)

	return b
}

func (b *domainBuilder) WithRelation(ctx context.Context, relationDefinition *coredomaindefinition.Relation) *domainBuilder {
	if b.Err != nil {
		return b
	}

	b.RelationDefinitionsToBuild = append(b.RelationDefinitionsToBuild, relationDefinition)

	return b
}

func (b *domainBuilder) WithCRUD(ctx context.Context, crudDefinition *coredomaindefinition.CRUD) *domainBuilder {
	if b.Err != nil {
		return b
	}

	b.CRUDToBuild = append(b.CRUDToBuild, crudDefinition)

	return b
}

func (b *domainBuilder) WithUsecase(ctx context.Context, usecaseDefinition *coredomaindefinition.Usecase) *domainBuilder {
	if b.Err != nil {
		return b
	}

	b.UsecaseDefinitionsToBuild = append(b.UsecaseDefinitionsToBuild, usecaseDefinition)

	return b
}

func (b *domainBuilder) Build(ctx context.Context) (*model.Domain, error) {
	if b.Err != nil {
		return nil, b.Err
	}

	b.buildRepositoryErrors(ctx)

	for _, modelDefinition := range b.ModelDefinitionsToBuild {
		b.buildModel(ctx, modelDefinition)
	}
	b.buildDomainRepository(ctx)
	for _, repositoryDefinition := range b.RepositoryDefinitionsToBuild {
		b.buildRepository(ctx, repositoryDefinition)
	}
	for _, relationDefinition := range b.RelationDefinitionsToBuild {
		b.buildRelation(ctx, relationDefinition)
	}
	for _, crudDefinition := range b.CRUDToBuild {
		b.buildCRUD(ctx, crudDefinition)
	}

	for _, usecaseDefinition := range b.UsecaseDefinitionsToBuild {
		b.buildUsecase(ctx, usecaseDefinition)
	}

	for _, usecase := range b.Domain.Usecases {
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
					Pkg:       b.Domain.Architecture.UsecasePkg,
					Reference: usecase.Request,
				},
			},
		})
		usecase.Function.Results = append(usecase.Function.Results, &model.Param{
			Type: &model.PointerType{
				Type: &model.PkgReference{
					Pkg:       b.Domain.Architecture.UsecasePkg,
					Reference: usecase.Result,
				},
			},
		})
		usecase.Function.Results = append(usecase.Function.Results, &model.Param{
			Type: model.PrimitiveTypeError,
		})
		b.GetDomainUsecase(ctx).Methods = append(b.GetDomainUsecase(ctx).Methods, usecase.Function)
	}
	b.buildUsecaseCRUDImplementation(ctx)
	b.buildHttpController(ctx)
	b.buildGormAdpater(ctx)
	b.buildService(ctx)
	b.buildHttpService(ctx)

	if b.Err != nil {
		return nil, b.Err
	}

	return b.Domain, nil
}
