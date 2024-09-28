package domainbuilder

import (
	"context"
	"slices"

	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

type RepositoryBuilder struct {
	DomainBuilder *domainBuilder
	Definition    *coredomaindefinition.Repository
	Repository    *model.Port
	Err           error

	FieldToColumn   *model.Map
	AllowedOrderBys *model.Consts
	AllowedWheres   *model.Consts

	RelationHandled []*coredomaindefinition.Relation
}

func NewRepositoryBuilder(
	ctx context.Context,
	domainBuilder *domainBuilder,
	definition *coredomaindefinition.Repository,
) *RepositoryBuilder {
	builder := &RepositoryBuilder{
		DomainBuilder: domainBuilder,
		Definition:    definition,
		Repository: &model.Port{
			Name: GetRepositoryName(ctx, definition),
			Pkg:  domainBuilder.Domain.Architecture.RepositoryPkg,
			// Methods:         []*model.RepositoryMethod{},
			// Functions:       []*model.Function{},
		},
		Err: nil,
	}

	elements := []interface{}{}

	elements = append(elements, &model.Var{
		Name:    GetRepositoryTableName(ctx, definition),
		Type:    model.PrimitiveTypeString,
		IsConst: true,
		Value:   definition.TableName,
	})

	defaultOrderBy := builder.DomainBuilder.Definition.Configuration.DefaultOrderBy
	if definition.DefaultOrderBy != "" {
		defaultOrderBy = definition.DefaultOrderBy
	}

	elements = append(elements, &model.Var{
		Name:    GetRepositoryDefaultOrderBy(ctx, definition),
		Type:    model.PrimitiveTypeString,
		IsConst: true,
		Value:   defaultOrderBy,
	})

	builder.FieldToColumn = &model.Map{
		Name: GetRepositoryFieldToColumnName(ctx, definition),
		Type: model.MapType{
			Key:   model.PrimitiveTypeString,
			Value: model.PrimitiveTypeString,
		},
	}
	builder.AllowedOrderBys = &model.Consts{
		Name:   GetRepositoryAllowedOrderBy(ctx, definition),
		Values: []interface{}{},
	}
	builder.AllowedWheres = &model.Consts{
		Name:   GetRepositoryAllowedWhere(ctx, definition),
		Values: []interface{}{},
	}
	for _, f := range builder.DomainBuilder.DefaultModelFields {
		builder.FieldToColumn.Values = append(builder.FieldToColumn.Values, model.MapValue{
			Key:   GetFieldName(ctx, f),
			Value: GetColumnName(ctx, f),
		})
		builder.AllowedOrderBys.Values = append(builder.AllowedOrderBys.Values, GetFieldName(ctx, f))
		builder.AllowedWheres.Values = append(builder.AllowedWheres.Values, GetFieldName(ctx, f))
	}

	for _, f := range definition.On.Fields {
		builder.FieldToColumn.Values = append(builder.FieldToColumn.Values, model.MapValue{
			Key:   GetFieldName(ctx, f),
			Value: GetColumnName(ctx, f),
		})
		builder.AllowedOrderBys.Values = append(builder.AllowedOrderBys.Values, GetFieldName(ctx, f))
		builder.AllowedWheres.Values = append(builder.AllowedWheres.Values, GetFieldName(ctx, f))
	}
	elements = append(elements, builder.FieldToColumn)
	elements = append(elements, builder.AllowedOrderBys)
	elements = append(elements, builder.AllowedWheres)

	builder.Repository.Elements = elements

	builder.addGetMethod(ctx)

	return builder
}

func (builder *RepositoryBuilder) WithRelation(ctx context.Context, relation *coredomaindefinition.Relation) *RepositoryBuilder {
	if builder.Err != nil {
		return builder
	}

	if relation.Source != builder.Definition.On && relation.Target != builder.Definition.On {
		return builder
	}

	if slices.Contains(builder.RelationHandled, relation) {
		return builder
	}
	builder.RelationHandled = append(builder.RelationHandled, relation)

	if !IsRelationMultiple(ctx, builder.Definition.On, relation) {
		var to *coredomaindefinition.Model
		if relation.Source == builder.Definition.On {
			to = relation.Target
		} else {
			to = relation.Source
		}
		builder.AllowedOrderBys.Values = append(builder.AllowedOrderBys.Values, GetSingleRelationIdName(ctx, to))
		builder.AllowedWheres.Values = append(builder.AllowedWheres.Values, GetSingleRelationIdName(ctx, to))
		builder.FieldToColumn.Values = append(builder.FieldToColumn.Values, model.MapValue{
			Key:   GetSingleRelationIdName(ctx, to),
			Value: GetSingleRelationColumn(ctx, to),
		})
	}

	return builder
}

func (builder *RepositoryBuilder) addGetMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodCtx := &model.Struct{
		Name:   GetMethodContextName(ctx, GetRepositoryGetMethod(ctx, builder.Definition)),
		Fields: []*model.Field{},
	}

	builder.addDefaultContextField(ctx, methodCtx)
	builder.addRetriveMethodDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)
}

func (builder *RepositoryBuilder) addDefaultContextField(ctx context.Context, methodContext *model.Struct) {
	if builder.Err != nil {
		return
	}

	methodContext.Fields = append(methodContext.Fields, &model.Field{
		Name: GetRepositoryMethodContextTransactionField(ctx),
		Type: &model.PointerType{
			Type: &model.PkgReference{
				Pkg: builder.DomainBuilder.GetRepositoryPackage(),
				Reference: &model.ExternalType{
					Type: REPOSITORY_TRANSATION,
				},
			},
		},
	})
}

func (builder *RepositoryBuilder) addRetriveMethodDefaultContextField(ctx context.Context, methodContext *model.Struct) {
	if builder.Err != nil {
		return
	}

	methodContext.Fields = append(methodContext.Fields, &model.Field{
		Name: GetRepositoryMethodContextWithInacctiveField(ctx),
		Type: model.PrimitiveTypeBool,
	})

	methodContext.Fields = append(methodContext.Fields, &model.Field{
		Name: "by",
		Type: &model.ArrayType{
			Type: &model.PointerType{
				Type: &model.PkgReference{
					Pkg: builder.DomainBuilder.GetRepositoryPackage(),
					Reference: &model.ExternalType{
						Type: REPOSITORY_WHERE,
					},
				},
			},
		},
	})
}

func (builder *RepositoryBuilder) Build(ctx context.Context) (*model.Port, error) {
	return builder.Repository, builder.Err
}
