package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	REPOSITORY_RETRIEVE_INACTIVE = "RetriveInactive"
	REPOSITORY_BY                = "By"
)

type RepositoryBuilder struct {
	DomainBuilder *domainBuilder
	Definition    *coredomaindefinition.Repository
	Repository    *model.File
	Err           error

	FieldToColumn   *model.Map
	AllowedOrderBys *model.Consts
	AllowedWheres   *model.Consts

	Methods []*model.Function
}

func NewRepositoryBuilder(
	ctx context.Context,
	domainBuilder *domainBuilder,
	definition *coredomaindefinition.Repository,
) Builder {
	builder := &RepositoryBuilder{
		DomainBuilder: domainBuilder,
		Definition:    definition,
		Repository: &model.File{
			Name: GetRepositoryName(ctx, definition),
			Pkg:  domainBuilder.Domain.Architecture.RepositoryPkg,
		},
		Err: nil,
	}

	elements := []interface{}{}

	elements = append(elements, &model.Var{
		Name:    GetRepositoryConstTableName(ctx, definition),
		Type:    model.PrimitiveTypeString,
		IsConst: true,
		Value:   GetRepositoryTableName(ctx, definition),
	})

	defaultOrderBy := builder.DomainBuilder.Definition.Configuration.DefaultOrderBy
	if definition.DefaultOrderBy != "" {
		defaultOrderBy = definition.DefaultOrderBy
	}

	elements = append(elements, &model.Var{
		Name:    GetRepositoryDefaultOrderBy(ctx, definition.On),
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
		Name:   GetRepositoryAllowedWhere(ctx, definition.On),
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
	if definition.On.Activable {
		builder.FieldToColumn.Values = append(builder.FieldToColumn.Values, model.MapValue{
			Key:   ACTIVE_FIELD_NAME,
			Value: stringtool.SnakeCase(ACTIVE_FIELD_NAME),
		})
		builder.AllowedWheres.Values = append(builder.AllowedWheres.Values, ACTIVE_FIELD_NAME)
		builder.AllowedOrderBys.Values = append(builder.AllowedOrderBys.Values, ACTIVE_FIELD_NAME)
	}
	elements = append(elements, builder.FieldToColumn)
	elements = append(elements, builder.AllowedOrderBys)
	elements = append(elements, builder.AllowedWheres)

	builder.Repository.Elements = elements

	builder.addGetMethod(ctx)
	builder.addListMethod(ctx)
	builder.addCreateMethod(ctx)
	builder.addUpdateMethod(ctx)
	builder.addDeleteMethod(ctx)

	return builder
}

func (builder *RepositoryBuilder) WithRelation(ctx context.Context, relation *coredomaindefinition.Relation) Builder {
	if builder.Err != nil {
		return builder
	}

	if relation.Source != builder.Definition.On && relation.Target != builder.Definition.On {
		return builder
	}

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

	} else if relation.Type == coredomaindefinition.RelationTypeManyToMany {
		builder.addManyToManyMethods(ctx, relation)
	}

	return builder
}

// WithModel implements Builder.
func (builder *RepositoryBuilder) WithModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) Builder {
	return builder
}

// WithRepository implements Builder.
func (builder *RepositoryBuilder) WithRepository(ctx context.Context, repositoryDefinition *coredomaindefinition.Repository) Builder {
	return builder
}

func (builder *RepositoryBuilder) addGetMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodName := GetRepositoryGetMethod(ctx, builder.Definition)

	methodCtx := &model.Struct{
		Name:   GetMethodContextName(ctx, methodName),
		Fields: []*model.Field{},
	}

	builder.addDefaultContextField(ctx, methodCtx)
	builder.addRetriveMethodDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)

	builder.addContextFieldOpt(ctx, methodCtx, methodName)

	builder.Methods = append(builder.Methods, GetRepositoryGetSignature(
		ctx,
		builder.Definition,
		builder.DomainBuilder.GetRepositoryPackage(),
		builder.DomainBuilder.GetModelPackage(),
	))
}

func (builder *RepositoryBuilder) addListMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodName := GetRepositoryListMethod(ctx, builder.Definition)

	methodCtx := &model.Struct{
		Name: GetMethodContextName(ctx, methodName),
		Fields: []*model.Field{
			{
				Name: PAGINATION_NAME,
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.DomainBuilder.GetRepositoryPackage(),
						Reference: &model.ExternalType{
							Type: PAGINATION_NAME,
						},
					},
				},
			},
		},
	}

	builder.addDefaultContextField(ctx, methodCtx)
	builder.addRetriveMethodDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)

	builder.addContextFieldOpt(ctx, methodCtx, methodName)

	builder.Methods = append(builder.Methods, GetRepositoryListSignature(
		ctx,
		builder.Definition,
		builder.DomainBuilder.GetRepositoryPackage(),
		builder.DomainBuilder.GetModelPackage(),
	))
}

func (builder *RepositoryBuilder) addCreateMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodCtx := &model.Struct{
		Name:   GetMethodContextName(ctx, GetRepositoryCreateMethod(ctx, builder.Definition)),
		Fields: []*model.Field{},
	}

	builder.addDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)

	builder.addContextFieldOpt(ctx, methodCtx, GetRepositoryCreateMethod(ctx, builder.Definition))

	builder.Methods = append(builder.Methods, GetRepositoryCreateSignature(
		ctx,
		builder.Definition,
		builder.DomainBuilder.GetRepositoryPackage(),
		builder.DomainBuilder.GetModelPackage(),
	))
}

func (builder *RepositoryBuilder) addUpdateMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodName := GetRepositoryUpdateMethod(ctx, builder.Definition)

	methodCtx := &model.Struct{
		Name:   GetMethodContextName(ctx, methodName),
		Fields: []*model.Field{},
	}

	builder.addDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)

	builder.addContextFieldOpt(ctx, methodCtx, methodName)

	builder.Methods = append(builder.Methods, GetRepositoryUpdateSignature(
		ctx,
		builder.Definition,
		builder.DomainBuilder.GetRepositoryPackage(),
		builder.DomainBuilder.GetModelPackage(),
	))
}

func (builder *RepositoryBuilder) addDeleteMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodCtx := &model.Struct{
		Name:   GetMethodContextName(ctx, GetRepositoryDeleteMethod(ctx, builder.Definition)),
		Fields: []*model.Field{},
	}

	builder.addDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)

	builder.addContextFieldOpt(ctx, methodCtx, GetRepositoryDeleteMethod(ctx, builder.Definition))

	builder.Methods = append(builder.Methods, GetRepositoryDeleteSignature(
		ctx,
		builder.Definition,
		builder.DomainBuilder.GetRepositoryPackage(),
		builder.DomainBuilder.GetModelPackage(),
	))
}

func (builder *RepositoryBuilder) addManyToManyMethods(ctx context.Context, relation *coredomaindefinition.Relation) {
	if builder.Err != nil {
		return
	}

	methodName := GetRepositoryAddRelationMethod(ctx, builder.Definition, relation)

	methodCtx := &model.Struct{
		Name:   GetMethodContextName(ctx, methodName),
		Fields: []*model.Field{},
	}

	builder.addDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)

	builder.addContextFieldOpt(ctx, methodCtx, methodName)

	var to *coredomaindefinition.Model
	if relation.Source == builder.Definition.On {
		to = relation.Target
	} else {
		to = relation.Source
	}

	builder.Methods = append(builder.Methods, &model.Function{
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
				Name: builder.Definition.On.Name + "Id",
				Type: model.PrimitiveTypeString,
			},
			{
				Name: to.Name + "Id",
				Type: model.PrimitiveTypeString,
			},
			{
				Name: "opts",
				Type: &model.VariaidicType{
					Type: &model.PkgReference{
						Pkg: builder.DomainBuilder.GetRepositoryPackage(),
						Reference: &model.ExternalType{
							Type: GetRepositoryMethodOptionName(ctx, methodName),
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: model.PrimitiveTypeError,
			},
		},
	})

	methodName = GetRepositoryRemoveRelationMethod(ctx, builder.Definition, relation)

	methodCtx = &model.Struct{
		Name:   GetMethodContextName(ctx, methodName),
		Fields: []*model.Field{},
	}

	builder.addDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)

	builder.addContextFieldOpt(ctx, methodCtx, methodName)

	builder.Methods = append(builder.Methods, &model.Function{
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
				Name: builder.Definition.On.Name + "Id",
				Type: model.PrimitiveTypeString,
			},
			{
				Name: to.Name + "Id",
				Type: model.PrimitiveTypeString,
			},
			{
				Name: "opts",
				Type: &model.VariaidicType{
					Type: &model.PkgReference{
						Pkg: builder.DomainBuilder.GetRepositoryPackage(),
						Reference: &model.ExternalType{
							Type: GetRepositoryMethodOptionName(ctx, methodName),
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: model.PrimitiveTypeError,
			},
		},
	})
}

func (builder *RepositoryBuilder) addDefaultContextField(ctx context.Context, methodContext *model.Struct) {
	if builder.Err != nil {
		return
	}

	methodContext.Fields = append(methodContext.Fields, &model.Field{
		Name: GetRepositoryMethodContextTransactionField(ctx),
		Type: &model.PkgReference{
			Pkg: builder.DomainBuilder.GetRepositoryPackage(),
			Reference: &model.ExternalType{
				Type: TRANSACTION_NAME,
			},
		},
	})
}

func (builder *RepositoryBuilder) addRetriveMethodDefaultContextField(ctx context.Context, methodContext *model.Struct) {
	if builder.Err != nil {
		return
	}

	if builder.DomainBuilder.RelationGraph.GetNode(builder.Definition.On).RequireRetriveInactive() {
		methodContext.Fields = append(methodContext.Fields, &model.Field{
			Name: REPOSITORY_RETRIEVE_INACTIVE,
			Type: model.PrimitiveTypeBool,
		})
	}

	methodContext.Fields = append(methodContext.Fields, &model.Field{
		Name: REPOSITORY_BY,
		Type: &model.ArrayType{
			Type: &model.PointerType{
				Type: &model.PkgReference{
					Pkg: builder.DomainBuilder.GetRepositoryPackage(),
					Reference: &model.ExternalType{
						Type: PAGINATION_WHERE,
					},
				},
			},
		},
	})
}

func (builder *RepositoryBuilder) addContextFieldOpt(ctx context.Context, methodContext *model.Struct, methodName string) *RepositoryBuilder {
	if builder.Err != nil {
		return builder
	}

	optGetter := &model.TypeDefinition{
		Name: methodName + "OptGettter",
		Type: model.PrimitiveTypeInt,
	}
	builder.Repository.Elements = append(builder.Repository.Elements, optGetter)

	builder.Repository.Elements = append(builder.Repository.Elements, &model.Var{
		Name:    methodName,
		Type:    model.PrimitiveTypeInt,
		Value:   0,
		IsConst: true,
	})

	optType := &model.TypeDefinition{
		Name: GetRepositoryMethodOptionName(ctx, methodName),
		Type: &model.Function{
			Args: []*model.Param{
				{
					Name: "ctx",
					Type: &model.PointerType{
						Type: &model.PkgReference{
							Pkg:       builder.DomainBuilder.GetRepositoryPackage(),
							Reference: methodContext,
						},
					},
				},
			},
		},
	}
	builder.Repository.Elements = append(builder.Repository.Elements, optType)

	for _, field := range methodContext.Fields {
		opt := &model.Function{
			Name: fmt.Sprintf("%sWith%s", methodName, field.Name),
			Args: []*model.Param{
				{
					Name: stringtool.LowerFirstLetter(field.Name),
					Type: field.Type,
				},
			},
			Results: []*model.Param{
				{
					Type: &model.PkgReference{
						Pkg:       builder.DomainBuilder.GetRepositoryPackage(),
						Reference: optType,
					},
				},
			},
		}
		opt.Content = func() (string, []*model.GoPkg) {
			str := fmt.Sprintf("return func(ctx *%s) {", methodContext.Name)
			str += fmt.Sprintf(" ctx.%s = %s", field.Name, stringtool.LowerFirstLetter(field.Name))
			str += " }"
			return str, nil
		}
		builder.Repository.Elements = append(builder.Repository.Elements, opt)

		getOpt := &model.Function{
			Name: fmt.Sprintf("With%s", field.Name),
			On:   optGetter,
			Args: []*model.Param{
				{
					Name: stringtool.LowerFirstLetter(field.Name),
					Type: field.Type,
				},
			},
			Results: []*model.Param{
				{
					Type: &model.PkgReference{
						Pkg:       builder.DomainBuilder.GetRepositoryPackage(),
						Reference: optType,
					},
				},
			},
			Content: func() (string, []*model.GoPkg) {
				str := fmt.Sprintf("return %s(", opt.Name)
				for _, arg := range opt.Args {
					str += arg.Name
				}
				str += ")"
				return str, nil
			},
		}
		builder.Repository.Elements = append(builder.Repository.Elements, getOpt)
	}

	return builder
}

func (builder *RepositoryBuilder) Build(ctx context.Context) error {
	if builder.Err != nil {
		return builder.Err
	}

	builder.DomainBuilder.Domain.Ports = append(builder.DomainBuilder.Domain.Ports, builder.Repository)

	return nil
}
