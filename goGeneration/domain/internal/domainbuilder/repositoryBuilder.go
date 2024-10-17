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
	*EmptyBuilder

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
			Key:   GetFieldName(ctx, f.Name),
			Value: GetColumnName(ctx, f),
		})
		builder.AllowedOrderBys.Values = append(builder.AllowedOrderBys.Values, GetFieldName(ctx, f.Name))
		builder.AllowedWheres.Values = append(builder.AllowedWheres.Values, GetFieldName(ctx, f.Name))
	}

	for _, f := range definition.On.Fields {
		builder.FieldToColumn.Values = append(builder.FieldToColumn.Values, model.MapValue{
			Key:   GetFieldName(ctx, f.Name),
			Value: GetColumnName(ctx, f),
		})
		builder.AllowedOrderBys.Values = append(builder.AllowedOrderBys.Values, GetFieldName(ctx, f.Name))
		builder.AllowedWheres.Values = append(builder.AllowedWheres.Values, GetFieldName(ctx, f.Name))
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

	builder.adCustomMethods(ctx)

	return builder
}

func (builder *RepositoryBuilder) adCustomMethods(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	for _, method := range builder.Definition.Methods {
		builder.adCustomMethod(ctx, method)
	}
}

func (builder *RepositoryBuilder) adCustomMethod(ctx context.Context, method *coredomaindefinition.RepositoryMethod) {
	if builder.Err != nil {
		return
	}

	f, err := GetRepositoryMethodSignature(ctx, method, builder.DomainBuilder.GetRepositoryPackage(), builder.DomainBuilder.TypeDefinitionToType)
	if err != nil {
		builder.Err = err
		return
	}

	methodCtx := &model.Struct{
		Name:   GetMethodContextName(ctx, f.Name),
		Fields: []*model.Field{},
	}

	builder.addDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)

	builder.addContextFieldOpt(ctx, methodCtx, f.Name)

	builder.Methods = append(builder.Methods, f)
}

func (builder *RepositoryBuilder) WithRelation(ctx context.Context, relation *coredomaindefinition.Relation) {
	if builder.Err != nil {
		return
	}

	if relation.Source != builder.Definition.On && relation.Target != builder.Definition.On {
		return
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
}

func (builder *RepositoryBuilder) addGetMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodName := GetRepositoryGetMethod(ctx, builder.Definition.On)

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

	methodName := GetRepositoryListMethod(ctx, builder.Definition.On)

	methodCtx := &model.Struct{
		Name: GetMethodContextName(ctx, methodName),
		Fields: []*model.Field{
			{
				Name: PAGINATION_NAME,
				Type: &model.PkgReference{
					Pkg: builder.DomainBuilder.GetModelPackage(),
					Reference: &model.ExternalType{
						Type: PAGINATION_NAME,
					},
				},
			},
			{
				Name: ORDERING_NAME,
				Type: &model.PkgReference{
					Pkg: builder.DomainBuilder.GetModelPackage(),
					Reference: &model.ExternalType{
						Type: ORDERING_NAME,
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
		Name:   GetMethodContextName(ctx, GetRepositoryCreateMethod(ctx, builder.Definition.On)),
		Fields: []*model.Field{},
	}

	builder.addDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)

	builder.addContextFieldOpt(ctx, methodCtx, GetRepositoryCreateMethod(ctx, builder.Definition.On))

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

	methodName := GetRepositoryUpdateMethod(ctx, builder.Definition.On)

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
		Name:   GetMethodContextName(ctx, GetRepositoryDeleteMethod(ctx, builder.Definition.On)),
		Fields: []*model.Field{},
	}

	builder.addDefaultContextField(ctx, methodCtx)

	builder.Repository.Elements = append(builder.Repository.Elements, methodCtx)

	builder.addContextFieldOpt(ctx, methodCtx, GetRepositoryDeleteMethod(ctx, builder.Definition.On))

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

	methodName := GetRepositoryAddRelationMethod(ctx, builder.Definition.On, relation)

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

	methodName = GetRepositoryRemoveRelationMethod(ctx, builder.Definition.On, relation)

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
						Type: REPOSITORY_WHERE,
					},
				},
			},
		},
	})
}

func (builder *RepositoryBuilder) addContextFieldOpt(ctx context.Context, methodContext *model.Struct, methodName string) {
	if builder.Err != nil {
		return
	}

	optGetter := &model.TypeDefinition{
		Name: methodName + "OptGettter",
		Type: model.PrimitiveTypeInt,
	}
	builder.Repository.Elements = append(builder.Repository.Elements, optGetter)

	builder.Repository.Elements = append(builder.Repository.Elements, &model.Var{
		Name:    methodName,
		Type:    &model.ExternalType{Type: methodName + "OptGettter"},
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
			Name: fmt.Sprintf("%s%s", methodName, GetOptName(ctx, field.Name)),
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
			Name: GetOptName(ctx, field.Name),
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
}

func (builder *RepositoryBuilder) addToHttpService(ctx context.Context, method *model.Function, request *model.Struct, response *model.Struct) {

}

func (builder *RepositoryBuilder) Build(ctx context.Context) error {
	if builder.Err != nil {
		return builder.Err
	}

	builder.DomainBuilder.Domain.Files = append(builder.DomainBuilder.Domain.Files, builder.Repository)

	return nil
}
