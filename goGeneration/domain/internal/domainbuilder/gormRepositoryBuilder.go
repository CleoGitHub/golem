package domainbuilder

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	GORM_MODEL_METHOD_NAME               = "gormModel"
	GORM_DOMAIN_REPOSITORY_DB_FIELD_NAME = "db"
	GORM_METHOD_CONTEXT_NAME             = "methodCtx"
	GORM_REQUEST_NAME                    = "request"
	OPERATOR_TO_GORM_OPERATOR            = "RepositoryOperatorToGormOperator"
	VALUE_TO_GORM_VALUE                  = "ValueToGormValue"
)

func ModelToGormModel(ctx context.Context, on *coredomaindefinition.Model) string {
	return fmt.Sprintf("%sFromModel", GetModelName(ctx, on))
}

func GormModelToModel(ctx context.Context, on *coredomaindefinition.Model) string {
	return fmt.Sprintf("%sToModel", GetModelName(ctx, on))
}

func ModelsToGormModels(ctx context.Context, on *coredomaindefinition.Model) string {
	return fmt.Sprintf("%sFromModels", PluralizeName(ctx, GetModelName(ctx, on)))
}

func GormModelsToModels(ctx context.Context, on *coredomaindefinition.Model) string {
	return fmt.Sprintf("%sToModels", PluralizeName(ctx, GetModelName(ctx, on)))
}

type GormRepositoryBuilder struct {
	DomainBuilder *domainBuilder
	GormDomainRepositoryBuilder
	Definition *coredomaindefinition.Repository
	Repository *model.File
	GormModel  *model.File
	Model      *model.Struct
	Err        error

	ModelToGormModel   []func() string
	ModelsToGormModels []func() string

	GormModelToModel   []func() string
	GormModelsToModels []func() string
}

var _ Builder = (*GormRepositoryBuilder)(nil)

func NewGormRepositoryBuilder(
	ctx context.Context,
	domainBuilder *domainBuilder,
	// gormDomainRepositoryBuilder *GormDomainRepositoryBuilder,
	definition *coredomaindefinition.Repository,
) Builder {
	builder := &GormRepositoryBuilder{
		DomainBuilder: domainBuilder,
		Definition:    definition,
		Repository: &model.File{
			Name: GetRepositoryName(ctx, definition),
			Pkg:  domainBuilder.Domain.Architecture.GormAdapterPkg,
		},
		GormModel: &model.File{
			Name:     GetModelName(ctx, definition.On),
			Pkg:      domainBuilder.Domain.Architecture.GormAdapterPkg,
			Elements: []interface{}{},
		},
		Err: nil,
	}
	builder.Model = &model.Struct{
		Name: GetModelName(ctx, definition.On),
		Methods: []*model.Function{
			{
				Name: "TableName",
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeString,
					},
				},
				Content: func() (content string, requiredPkg []*model.GoPkg) {
					return fmt.Sprintf(
							`return %s.%s`,
							builder.DomainBuilder.GetRepositoryPackage().Alias, GetRepositoryConstTableName(ctx, definition),
						), []*model.GoPkg{
							builder.DomainBuilder.GetRepositoryPackage(),
						}
				},
			},
		},
	}
	builder.GormModel.Elements = append(builder.GormModel.Elements, builder.Model)

	modelFieldNames := []string{}
	for _, f := range builder.DomainBuilder.DefaultModelFields {
		modelFieldNames = append(modelFieldNames, f.Name)
		field, err := builder.DomainBuilder.FieldDefinitionToField(ctx, f)
		if err != nil {
			builder.Err = merror.Stack(err)
			return builder
		}
		field.Tags = append(field.Tags, &model.Tag{
			Name:   "gorm",
			Values: []string{"column:" + GetColumnNameFromName(ctx, field.Name)},
		})
		builder.Model.Fields = append(builder.Model.Fields, PrepareFieldFormGorm(ctx, field))

		builder.ModelToGormModel = append(builder.ModelToGormModel, func() string {
			return fmt.Sprintf(
				"%s: %s.%s",
				GetFieldName(ctx, f.Name), GORM_MODEL_METHOD_NAME, GetFieldName(ctx, f.Name),
			) + "," + consts.LN
		})

		builder.GormModelToModel = append(builder.GormModelToModel, func() string {
			return fmt.Sprintf(
				"%s: %s.%s",
				GetFieldName(ctx, f.Name), GORM_MODEL_METHOD_NAME, GetFieldName(ctx, f.Name),
			) + "," + consts.LN
		})
	}
	if definition.On.Archivable {
		modelFieldNames = append(modelFieldNames, "deleted")
		fieldName := "DeletedAt"
		field := &model.Field{
			Name: fieldName,
			Type: &model.PointerType{
				Type: &model.PkgReference{
					Pkg: consts.CommonPkgs["gorm"],
					Reference: &model.ExternalType{
						Type: fieldName,
					},
				},
			},
			Tags: []*model.Tag{
				{
					Name:   "gorm",
					Values: []string{"column:" + GetColumnNameFromName(ctx, fieldName), "index"},
				},
			},
		}
		builder.Model.Fields = append(builder.Model.Fields, PrepareFieldFormGorm(ctx, field))

		builder.ModelToGormModel = append(builder.ModelToGormModel, func() string {
			return fmt.Sprintf(
				`%s:  &%s.DeletedAt{ Time: %s.%s, Valid: %s.%s.String() != "0000-00-00 00:00:00"}`,
				fieldName,
				consts.CommonPkgs["gorm"].Alias,
				GORM_MODEL_METHOD_NAME,
				fieldName,
				GORM_MODEL_METHOD_NAME,
				fieldName,
			) + "," + consts.LN
		})

		builder.GormModelToModel = append(builder.GormModelToModel, func() string {
			return fmt.Sprintf(
				`%s:  %s.DeletedAt.Time`,
				fieldName,
				GORM_MODEL_METHOD_NAME,
			) + "," + consts.LN
		})
	}

	// Add activable field if model is activable
	if definition.On.Activable {
		modelFieldNames = append(modelFieldNames, "active")
		field, err := builder.DomainBuilder.FieldDefinitionToField(ctx, &coredomaindefinition.Field{
			Name: "active",
			Type: coredomaindefinition.PrimitiveTypeBool,
		})
		if err != nil {
			builder.Err = merror.Stack(err)
			return builder
		}
		field.Tags = append(field.Tags, &model.Tag{
			Name:   "gorm",
			Values: []string{"column:" + GetColumnNameFromName(ctx, field.Name)},
		})
		builder.Model.Fields = append(builder.Model.Fields, PrepareFieldFormGorm(ctx, field))

		builder.ModelToGormModel = append(builder.ModelToGormModel, func() string {
			return fmt.Sprintf(
				"%s: %s.%s", ACTIVE_FIELD_NAME, GORM_MODEL_METHOD_NAME, ACTIVE_FIELD_NAME,
			) + "," + consts.LN
		})

		builder.GormModelToModel = append(builder.GormModelToModel, func() string {
			return fmt.Sprintf(
				"%s: %s.%s", ACTIVE_FIELD_NAME, GORM_MODEL_METHOD_NAME, ACTIVE_FIELD_NAME,
			) + "," + consts.LN
		})
	}

	// Add default fields to definition
	for _, field := range definition.On.Fields {
		if slices.Contains(modelFieldNames, field.Name) {
			builder.Err = merror.Stack(NewErrDefaultFiedlRedefined(field.Name))
			return builder
		}
		f, err := builder.DomainBuilder.FieldDefinitionToField(ctx, field)
		if err != nil {
			builder.Err = merror.Stack(err)
			return builder
		}
		f.Tags = append(f.Tags, &model.Tag{
			Name:   "gorm",
			Values: []string{"column:" + GetColumnName(ctx, field)},
		})
		builder.Model.Fields = append(builder.Model.Fields, PrepareFieldFormGorm(ctx, f))

		builder.ModelToGormModel = append(builder.ModelToGormModel, func() string {
			return fmt.Sprintf(
				"%s: %s.%s",
				GetFieldName(ctx, field.Name), GORM_MODEL_METHOD_NAME, GetFieldName(ctx, field.Name),
			) + "," + consts.LN
		})
	}

	return builder
}

func (builder *GormRepositoryBuilder) WithRelation(ctx context.Context, relation *coredomaindefinition.Relation) {
	if builder.Err != nil {
		return
	}

	fmt.Printf("Relation: %s => %s", relation.Source.Name, relation.Target.Name)

	if relation.Source != builder.Definition.On && relation.Target != builder.Definition.On {
		return
	}

	fmt.Printf("Relation: %s => %s", relation.Source.Name, relation.Target.Name)
	var to *coredomaindefinition.Model
	if relation.Source == builder.Definition.On {
		to = relation.Target
	} else {
		if relation.IgnoreReverse {
			return
		}
		to = relation.Source
	}
	if to.Name == "user" {
		fmt.Printf("It is user")
	}

	if !IsRelationMultiple(ctx, builder.Definition.On, relation) {
		if to.Name == "user" {
			fmt.Printf("It not multiple relation")
		}
		field := &model.Field{
			Name: GetSingleRelationName(ctx, to),
			Type: &model.PointerType{
				Type: &model.PkgReference{
					Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
					Reference: &model.ExternalType{
						Type: GetModelName(ctx, to),
					},
				},
			},
		}
		builder.Model.Fields = append(builder.Model.Fields, PrepareFieldFormGorm(ctx, field))

		builder.GormModelToModel = append(builder.GormModelToModel, func() string {
			return fmt.Sprintf(
				"%s: %s(%s.%s)", GetSingleRelationName(ctx, to), GormModelToModel(ctx, to), GORM_MODEL_METHOD_NAME, GetSingleRelationName(ctx, to),
			) + "," + consts.LN
		})

		builder.ModelToGormModel = append(builder.ModelToGormModel, func() string {
			return fmt.Sprintf(
				"%s: %s(%s.%s)", GetSingleRelationName(ctx, to), ModelToGormModel(ctx, to), GORM_MODEL_METHOD_NAME, GetSingleRelationName(ctx, to),
			) + "," + consts.LN
		})

		optionnal, err := IsRelationOptionnal(ctx, builder.Definition.On, relation)
		if err != nil {
			builder.Err = merror.Stack(err)
			return
		}

		var t model.Type = model.PrimitiveTypeString
		dereference := ""
		reference := ""
		if optionnal {
			t = &model.PointerType{
				Type: t,
			}

			dereference = "*"
			reference = "&"
		}
		field = &model.Field{
			Name: GetSingleRelationIdName(ctx, to),
			Type: t,
			Tags: []*model.Tag{
				{
					Name:   "gorm",
					Values: []string{"column:" + GetSingleRelationColumn(ctx, to)},
				},
			},
		}
		builder.Model.Fields = append(builder.Model.Fields, PrepareFieldFormGorm(ctx, field))

		builder.GormModelToModel = append(builder.GormModelToModel, func() string {
			return fmt.Sprintf(
				"%s: %s%s.%s", GetSingleRelationIdName(ctx, to), dereference, GORM_MODEL_METHOD_NAME, GetSingleRelationIdName(ctx, to),
			) + "," + consts.LN
		})

		builder.ModelToGormModel = append(builder.ModelToGormModel, func() string {
			return fmt.Sprintf(
				"%s: %s%s.%s", GetSingleRelationIdName(ctx, to), reference, GORM_MODEL_METHOD_NAME, GetSingleRelationIdName(ctx, to),
			) + "," + consts.LN
		})
	} else {
		field := &model.Field{
			Name: GetMultipleRelationName(ctx, to),
			Type: &model.ArrayType{
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, to),
						},
					},
				},
			},
		}
		if relation.Type == coredomaindefinition.RelationTypeManyToMany {
			field.Tags = append(field.Tags, &model.Tag{
				Name:   "gorm",
				Values: []string{"many2many:" + GetManyToManyColumn(ctx, relation)},
			})
		}
		builder.Model.Fields = append(builder.Model.Fields, PrepareFieldFormGorm(ctx, field))
		builder.ModelToGormModel = append(builder.ModelToGormModel, func() string {
			return fmt.Sprintf(
				"%s: %s(%s.%s)", GetMultipleRelationName(ctx, to), ModelsToGormModels(ctx, to), GORM_MODEL_METHOD_NAME, GetMultipleRelationName(ctx, to),
			) + "," + consts.LN
		})
		builder.GormModelToModel = append(builder.GormModelToModel, func() string {
			return fmt.Sprintf(
				"%s: %s(%s.%s)", GetMultipleRelationName(ctx, to), GormModelsToModels(ctx, to), GORM_MODEL_METHOD_NAME, GetMultipleRelationName(ctx, to),
			) + "," + consts.LN
		})

		if relation.Type == coredomaindefinition.RelationTypeManyToMany {
			builder.addManyToManyMethods(ctx, relation)
		}
	}
}

func (builder *GormRepositoryBuilder) addGetMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodName := GetRepositoryGetMethod(ctx, builder.Definition.On)

	ctxName := GetMethodContextName(ctx, methodName)
	method := GetRepositoryGetSignature(ctx, builder.Definition, builder.DomainBuilder.GetRepositoryPackage(), builder.DomainBuilder.GetModelPackage())
	method.Content = func() (string, []*model.GoPkg) {
		str := ""
		pkg := []*model.GoPkg{}

		s, p := builder.getInitContext(ctx, ctxName)
		str += s
		pkg = append(pkg, p...)

		s, p = builder.getGormTransactionInitialisation(ctx, "nil")
		str += s
		pkg = append(pkg, p...)

		s, p = builder.getRequestModelWithDependencyTree(ctx)
		str += s
		pkg = append(pkg, p...)

		str += fmt.Sprintf("entity := &%s{}", GetModelName(ctx, builder.Definition.On)) + consts.LN
		str += fmt.Sprintf("err := %s.First(entity).Error", GORM_REQUEST_NAME) + consts.LN
		str += "if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {" + consts.LN
		str += fmt.Sprintf("return nil, %s.%s", builder.DomainBuilder.GetRepositoryPackage().Alias, REPOSITORY_ERROR_NOT_FOUND.Name) + consts.LN
		str += "} else if err != nil {" + consts.LN
		str += "return nil, err" + consts.LN
		str += "} " + consts.LN
		str += fmt.Sprintf("return %s(entity), nil", GormModelToModel(ctx, builder.Definition.On)) + consts.LN

		return str, pkg
	}
	method.On = &model.PointerType{
		Type: &model.PkgReference{
			Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
			Reference: &model.ExternalType{
				Type: GetGormDomainRepositoryName(ctx, builder.DomainBuilder.Definition),
			},
		},
	}
	method.OnName = GORM_DOMAIN_REPO_METHOD_NAME
	builder.Repository.Elements = append(builder.Repository.Elements, method)
}

func (builder *GormRepositoryBuilder) addListMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodName := GetRepositoryListMethod(ctx, builder.Definition.On)

	ctxName := GetMethodContextName(ctx, methodName)
	method := GetRepositoryListSignature(ctx, builder.Definition, builder.DomainBuilder.GetRepositoryPackage(), builder.DomainBuilder.GetModelPackage())
	method.Content = func() (string, []*model.GoPkg) {
		str := ""
		pkg := []*model.GoPkg{}

		s, p := builder.getInitContext(ctx, ctxName)
		str += s
		pkg = append(pkg, p...)

		s, p = builder.getGormTransactionInitialisation(ctx, "nil")
		str += s
		pkg = append(pkg, p...)

		s, p = builder.getRequestModelWithDependencyTree(ctx)
		str += s
		pkg = append(pkg, p...)

		str += fmt.Sprintf("if %s.%s != (%s.%s{}) {", GORM_METHOD_CONTEXT_NAME, PAGINATION_NAME, builder.DomainBuilder.GetModelPackage().Alias, PAGINATION_NAME) + consts.LN
		limit := fmt.Sprintf("%s.%s.%s()", GORM_METHOD_CONTEXT_NAME, PAGINATION_NAME, PAGINATION_GetItemsPerPage)
		offset := fmt.Sprintf("%s.%s.%s()", GORM_METHOD_CONTEXT_NAME, PAGINATION_NAME, PAGINATION_GetPage)
		str += fmt.Sprintf(
			"%s = %s.Offset(int(%s * (%s - 1))).Limit(int(%s))",
			GORM_REQUEST_NAME, GORM_REQUEST_NAME, limit, offset, limit,
		) + consts.LN
		str += "}" + consts.LN

		str += fmt.Sprintf("entities := []*%s{}", GetModelName(ctx, builder.Definition.On)) + consts.LN
		str += fmt.Sprintf("err := %s.Find(entities).Error", GORM_REQUEST_NAME) + consts.LN
		str += "if err != nil {" + consts.LN
		str += "return nil, err" + consts.LN
		str += "} " + consts.LN
		str += fmt.Sprintf("return %s(entities), nil", GormModelsToModels(ctx, builder.Definition.On)) + consts.LN

		return str, pkg
	}
	method.On = &model.PkgReference{
		Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
		Reference: &model.ExternalType{
			Type: GetGormDomainRepositoryName(ctx, builder.DomainBuilder.Definition),
		},
	}
	method.OnName = GORM_DOMAIN_REPO_METHOD_NAME
	builder.Repository.Elements = append(builder.Repository.Elements, method)
}

func (builder *GormRepositoryBuilder) addCreateMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodName := GetRepositoryCreateMethod(ctx, builder.Definition.On)

	ctxName := GetMethodContextName(ctx, methodName)
	method := GetRepositoryCreateSignature(ctx, builder.Definition, builder.DomainBuilder.GetRepositoryPackage(), builder.DomainBuilder.GetModelPackage())
	method.Content = func() (string, []*model.GoPkg) {
		str := ""
		pkg := []*model.GoPkg{}

		s, p := builder.getInitContext(ctx, ctxName)
		str += s
		pkg = append(pkg, p...)

		s, p = builder.getGormTransactionInitialisation(ctx, "nil")
		str += s
		pkg = append(pkg, p...)

		str += fmt.Sprintf("result := %s(%s)", ModelToGormModel(ctx, builder.Definition.On), REPOSITORY_ENTITY_PARAM_NAME) + consts.LN
		str += fmt.Sprintf("err := db.Model(&%s{}).Create(result).Error", GetModelName(ctx, builder.Definition.On)) + consts.LN
		str += "if err != nil {" + consts.LN
		str += "return nil, err" + consts.LN
		str += "} " + consts.LN
		str += fmt.Sprintf("return %s(result), nil", GormModelToModel(ctx, builder.Definition.On)) + consts.LN

		return str, pkg
	}
	method.On = &model.PkgReference{
		Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
		Reference: &model.ExternalType{
			Type: GetGormDomainRepositoryName(ctx, builder.DomainBuilder.Definition),
		},
	}
	method.OnName = GORM_DOMAIN_REPO_METHOD_NAME
	builder.Repository.Elements = append(builder.Repository.Elements, method)
}

func (builder *GormRepositoryBuilder) addUpdateMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodName := GetRepositoryUpdateMethod(ctx, builder.Definition.On)

	ctxName := GetMethodContextName(ctx, methodName)
	method := GetRepositoryUpdateSignature(ctx, builder.Definition, builder.DomainBuilder.GetRepositoryPackage(), builder.DomainBuilder.GetModelPackage())
	method.Content = func() (string, []*model.GoPkg) {
		str := ""
		pkg := []*model.GoPkg{
			consts.CommonPkgs["gorm/clause"],
		}

		s, p := builder.getInitContext(ctx, ctxName)
		str += s
		pkg = append(pkg, p...)

		s, p = builder.getGormTransactionInitialisation(ctx, "nil")
		str += s
		pkg = append(pkg, p...)

		str += fmt.Sprintf("result := %s(%s)", ModelToGormModel(ctx, builder.Definition.On), REPOSITORY_ENTITY_PARAM_NAME) + consts.LN
		str += fmt.Sprintf("err := db.Model(&%s{}).Clauses(clause.Returning{}).Updates(result).Error", GetModelName(ctx, builder.Definition.On)) + consts.LN
		str += "if err != nil {" + consts.LN
		str += "return nil, err" + consts.LN
		str += "} " + consts.LN
		str += fmt.Sprintf("return %s(result), nil", GormModelToModel(ctx, builder.Definition.On)) + consts.LN

		return str, pkg
	}
	method.On = &model.PkgReference{
		Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
		Reference: &model.ExternalType{
			Type: GetGormDomainRepositoryName(ctx, builder.DomainBuilder.Definition),
		},
	}
	method.OnName = GORM_DOMAIN_REPO_METHOD_NAME
	builder.Repository.Elements = append(builder.Repository.Elements, method)
}

func (builder *GormRepositoryBuilder) addDeleteMethod(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	methodName := GetRepositoryDeleteMethod(ctx, builder.Definition.On)

	ctxName := GetMethodContextName(ctx, methodName)
	method := GetRepositoryDeleteSignature(ctx, builder.Definition, builder.DomainBuilder.GetRepositoryPackage(), builder.DomainBuilder.GetModelPackage())
	method.Content = func() (string, []*model.GoPkg) {
		str := ""
		pkg := []*model.GoPkg{
			consts.CommonPkgs["gorm/clause"],
		}

		s, p := builder.getInitContext(ctx, ctxName)
		str += s
		pkg = append(pkg, p...)

		s, p = builder.getGormTransactionInitialisation(ctx, "")
		str += s
		pkg = append(pkg, p...)

		str += fmt.Sprintf("err := db.Model(&%s{}).Delete(&%s{Id: id}).Error", GetModelName(ctx, builder.Definition.On), GetModelName(ctx, builder.Definition.On)) + consts.LN
		str += "if err != nil {" + consts.LN
		str += "return err" + consts.LN
		str += "} " + consts.LN
		str += "return nil" + consts.LN

		return str, pkg
	}
	method.On = &model.PkgReference{
		Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
		Reference: &model.ExternalType{
			Type: GetGormDomainRepositoryName(ctx, builder.DomainBuilder.Definition),
		},
	}
	method.OnName = GORM_DOMAIN_REPO_METHOD_NAME
	builder.Repository.Elements = append(builder.Repository.Elements, method)
}

func (builder *GormRepositoryBuilder) addManyToManyMethods(ctx context.Context, relation *coredomaindefinition.Relation) {
	if builder.Err != nil {
		return
	}

	addMethodName := GetRepositoryAddRelationMethod(ctx, builder.Definition.On, relation)
	addCtxName := GetMethodContextName(ctx, addMethodName)

	var to *coredomaindefinition.Model
	if relation.Source == builder.Definition.On {
		to = relation.Target
	} else {
		to = relation.Source
	}

	method := &model.Function{
		Name: addMethodName,
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
							Type: GetRepositoryMethodOptionName(ctx, addMethodName),
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
	}
	method.Content = func() (string, []*model.GoPkg) {
		str := ""
		pkg := []*model.GoPkg{
			consts.CommonPkgs["gorm/clause"],
		}

		s, p := builder.getInitContext(ctx, addCtxName)
		str += s
		pkg = append(pkg, p...)

		s, p = builder.getGormTransactionInitialisation(ctx, "")
		str += s
		pkg = append(pkg, p...)

		str += fmt.Sprintf(
			`err := db.Model(&%s{Id: %s}).Association("%s").Append(&%s{Id: %s})`,
			GetModelName(ctx, builder.Definition.On),
			builder.Definition.On.Name+"Id",
			GetMultipleRelationName(ctx, to),
			GetModelName(ctx, to),
			to.Name+"Id",
		) + consts.LN
		str += "if err != nil {" + consts.LN
		str += "return err" + consts.LN
		str += "} " + consts.LN
		str += "return nil" + consts.LN

		return str, pkg
	}
	method.On = &model.PkgReference{
		Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
		Reference: &model.ExternalType{
			Type: GetGormDomainRepositoryName(ctx, builder.DomainBuilder.Definition),
		},
	}
	method.OnName = GORM_DOMAIN_REPO_METHOD_NAME
	builder.Repository.Elements = append(builder.Repository.Elements, method)

	removeMethodName := GetRepositoryRemoveRelationMethod(ctx, builder.Definition.On, relation)
	removeCtxName := GetMethodContextName(ctx, removeMethodName)

	method = &model.Function{
		Name: removeMethodName,
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
							Type: GetRepositoryMethodOptionName(ctx, removeMethodName),
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
	}
	method.Content = func() (string, []*model.GoPkg) {
		str := ""
		pkg := []*model.GoPkg{
			consts.CommonPkgs["gorm/clause"],
		}

		s, p := builder.getInitContext(ctx, removeCtxName)
		str += s
		pkg = append(pkg, p...)

		s, p = builder.getGormTransactionInitialisation(ctx, "")
		str += s
		pkg = append(pkg, p...)

		str += fmt.Sprintf(
			`err := db.Model(&%s{Id: %s}).Association("%s").Delete(&%s{Id: %s})`,
			GetModelName(ctx, builder.Definition.On),
			builder.Definition.On.Name+"Id",
			GetMultipleRelationName(ctx, to),
			GetModelName(ctx, to),
			to.Name+"Id",
		) + consts.LN
		str += "if err != nil {" + consts.LN
		str += "return err" + consts.LN
		str += "} " + consts.LN
		str += "return nil" + consts.LN

		return str, pkg
	}
	method.On = &model.PkgReference{
		Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
		Reference: &model.ExternalType{
			Type: GetGormDomainRepositoryName(ctx, builder.DomainBuilder.Definition),
		},
	}
	method.OnName = GORM_DOMAIN_REPO_METHOD_NAME
	builder.Repository.Elements = append(builder.Repository.Elements, method)
}

func (builder *GormRepositoryBuilder) addMethods(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	builder.addGetMethod(ctx)
	builder.addListMethod(ctx)
	builder.addCreateMethod(ctx)
	builder.addUpdateMethod(ctx)
	builder.addDeleteMethod(ctx)
}

func (builder *GormRepositoryBuilder) addGormModelToModel(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	builder.GormModel.Elements = append(builder.GormModel.Elements, &model.Function{
		Name: GormModelToModel(ctx, builder.Definition.On),
		Args: []*model.Param{
			{
				Name: GORM_MODEL_METHOD_NAME,
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, builder.Definition.On),
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{

				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.DomainBuilder.GetModelPackage(),
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, builder.Definition.On),
						},
					},
				},
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			str := fmt.Sprintf("if %s == nil { return nil }", GORM_MODEL_METHOD_NAME) + consts.LN
			str += fmt.Sprintf("return &%s.%s{", builder.DomainBuilder.GetModelPackage().Alias, GetModelName(ctx, builder.Definition.On)) + consts.LN
			for _, gormModelToModel := range builder.GormModelToModel {
				str += gormModelToModel()
			}
			str += "}"
			return str, nil
		},
	})
}

func (builder *GormRepositoryBuilder) addModelToGormModel(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	builder.GormModel.Elements = append(builder.GormModel.Elements, &model.Function{
		Name: ModelToGormModel(ctx, builder.Definition.On),
		Args: []*model.Param{
			{
				Name: GORM_MODEL_METHOD_NAME,
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.DomainBuilder.GetModelPackage(),
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, builder.Definition.On),
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{

				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, builder.Definition.On),
						},
					},
				},
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			str := fmt.Sprintf("if %s == nil { return nil }", GORM_MODEL_METHOD_NAME) + consts.LN
			str += fmt.Sprintf("return &%s{", GetModelName(ctx, builder.Definition.On)) + consts.LN
			for _, modelToGormModel := range builder.ModelToGormModel {
				str += modelToGormModel()
			}
			str += "}"
			return str, nil
		},
	})
}

func (builder *GormRepositoryBuilder) getInitContext(ctx context.Context, contextName string) (string, []*model.GoPkg) {
	// Init context
	str := fmt.Sprintf("%s := &%s.%s{}", GORM_METHOD_CONTEXT_NAME, builder.DomainBuilder.GetRepositoryPackage().Alias, contextName) + consts.LN
	str += fmt.Sprintf("for _, opt := range %s {", REPOSIOTY_METHOD_CONTEXT_OPTS_NAME) + consts.LN
	str += fmt.Sprintf("opt(%s)", GORM_METHOD_CONTEXT_NAME) + consts.LN
	str += "}" + consts.LN

	return str, []*model.GoPkg{
		builder.DomainBuilder.GetRepositoryPackage(),
	}
}

func (builder *GormRepositoryBuilder) getGormTransactionInitialisation(ctx context.Context, extraReturns string) (string, []*model.GoPkg) {
	str := fmt.Sprintf("var db *%s.DB", consts.CommonPkgs["gorm"].Alias) + consts.LN
	str += fmt.Sprintf("if %s.Transaction != nil {", GORM_METHOD_CONTEXT_NAME) + consts.LN
	str += fmt.Sprintf("tx, ok := %s.Transaction.Get(ctx).(*%s.DB)", GORM_METHOD_CONTEXT_NAME, consts.CommonPkgs["gorm"].Alias) + consts.LN
	if extraReturns != "" {
		extraReturns = fmt.Sprintf("%s, ", extraReturns)
	}
	str += fmt.Sprintf(`if !ok { return %s %s.New("expected transaction to be *gorm.DB") }`, extraReturns, consts.CommonPkgs["errors"].Alias) + consts.LN
	str += fmt.Sprintf("%s = tx", GORM_DOMAIN_REPOSITORY_DB_FIELD_NAME) + consts.LN
	str += "} else { " + consts.LN
	str += fmt.Sprintf(
		"%s = %s.%s",
		GORM_DOMAIN_REPOSITORY_DB_FIELD_NAME,
		GORM_DOMAIN_REPO_METHOD_NAME,
		GORM_DOMAIN_REPOSITORY_DB_FIELD_NAME,
	) + consts.LN
	str += "}" + consts.LN

	return str, []*model.GoPkg{
		builder.DomainBuilder.GetRepositoryPackage(),
		consts.CommonPkgs["gorm"],
		consts.CommonPkgs["errors"],
	}
}

func (builder *GormRepositoryBuilder) getRequestModelWithDependencyTree(ctx context.Context) (string, []*model.GoPkg) {
	str := fmt.Sprintf("%s := %s.Model(&%s{})", GORM_REQUEST_NAME, GORM_DOMAIN_REPOSITORY_DB_FIELD_NAME, builder.Model.Name) + consts.LN
	if builder.Definition.On.Activable {
		str += fmt.Sprintf("if !%s.%s {", GORM_METHOD_CONTEXT_NAME, REPOSITORY_RETRIEVE_INACTIVE) + consts.LN
		str += fmt.Sprintf("%s.Where(&%s{%s: true})", GORM_REQUEST_NAME, GetModelName(ctx, builder.Definition.On), ACTIVE_FIELD_NAME) + consts.LN
		str += "}" + consts.LN
	}

	node := builder.DomainBuilder.RelationGraph.GetNode(builder.Definition.On)
	path := ""
	if node != nil {
		i := 0
		for i < len(node.Links) {
			link := node.Links[i]
			if link.Type == RelationNodeLinkType_DEPEND {
				node = link.To
				i = 0
				path = fmt.Sprintf("%s.%s", path, GetSingleRelationName(ctx, node.Model))
				path = strings.TrimPrefix(path, ".")
				if node.RequireRetriveInactive() && node.Model.Activable {
					str += fmt.Sprintf("if %s.%s {", GORM_METHOD_CONTEXT_NAME, REPOSITORY_RETRIEVE_INACTIVE) + consts.LN
					str += fmt.Sprintf(`%s.Joins("%s")`, GORM_REQUEST_NAME, path) + consts.LN
					str += "} else {" + consts.LN
					str += fmt.Sprintf(`%s.Joins("%s", %s.Where(&%s{%s: true}))`, GORM_REQUEST_NAME, path, GORM_REQUEST_NAME, GetModelName(ctx, node.Model), ACTIVE_FIELD_NAME) + consts.LN
					str += "}" + consts.LN
				}
			} else {
				i++
			}
		}
	}
	str += fmt.Sprintf("if %s.%s != nil {", GORM_METHOD_CONTEXT_NAME, REPOSITORY_BY) + consts.LN
	str += fmt.Sprintf("for _, where := range %s.%s {", GORM_METHOD_CONTEXT_NAME, REPOSITORY_BY) + consts.LN
	str += fmt.Sprintf("if slices.Contains(%s.%s, where.Key){", builder.DomainBuilder.GetRepositoryPackage().Alias, GetRepositoryAllowedWhere(ctx, builder.Definition.On)) + consts.LN
	str += fmt.Sprintf(
		`%s = %s.Where(fmt.Sprintf("%s %s ?", %s.%s[where.Key], %s(where.Operator)), %s(where.Value))`,
		GORM_REQUEST_NAME, GORM_REQUEST_NAME, "%s", "%s", builder.DomainBuilder.GetRepositoryPackage().Alias, GetRepositoryFieldToColumnName(ctx, builder.Definition), OPERATOR_TO_GORM_OPERATOR, VALUE_TO_GORM_VALUE,
	) + consts.LN
	str += "}" + consts.LN
	str += "}" + consts.LN
	str += "}" + consts.LN

	return str, []*model.GoPkg{
		consts.CommonPkgs["slices"],
		consts.CommonPkgs["fmt"],
		builder.DomainBuilder.GetRepositoryPackage(),
	}
}

func (builder *GormRepositoryBuilder) addGormModelsToModels(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	builder.GormModel.Elements = append(builder.GormModel.Elements, &model.Function{
		Name: GormModelsToModels(ctx, builder.Definition.On),
		Args: []*model.Param{
			{
				Name: PluralizeName(ctx, GORM_MODEL_METHOD_NAME),
				Type: &model.ArrayType{
					Type: &model.PointerType{
						Type: &model.PkgReference{
							Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
							Reference: &model.ExternalType{
								Type: GetModelName(ctx, builder.Definition.On),
							},
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
							Pkg: builder.DomainBuilder.GetModelPackage(),
							Reference: &model.ExternalType{
								Type: GetModelName(ctx, builder.Definition.On),
							},
						},
					},
				},
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			str := fmt.Sprintf("entities := []*%s.%s{}", builder.DomainBuilder.GetModelPackage().Alias, GetModelName(ctx, builder.Definition.On)) + consts.LN
			str += fmt.Sprintf("for _, %s := range %s {", GORM_MODEL_METHOD_NAME, PluralizeName(ctx, GORM_MODEL_METHOD_NAME)) + consts.LN
			str += fmt.Sprintf("entities = append(entities,%s(%s))", GormModelToModel(ctx, builder.Definition.On), GORM_MODEL_METHOD_NAME) + consts.LN
			str += "}" + consts.LN
			str += "return entities"
			return str, nil
		},
	})
}

func (builder *GormRepositoryBuilder) addModelsToGormModels(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	builder.GormModel.Elements = append(builder.GormModel.Elements, &model.Function{
		Name: ModelsToGormModels(ctx, builder.Definition.On),
		Args: []*model.Param{
			{
				Name: PluralizeName(ctx, GORM_MODEL_METHOD_NAME),
				Type: &model.ArrayType{
					Type: &model.PointerType{
						Type: &model.PkgReference{
							Pkg: builder.DomainBuilder.GetModelPackage(),
							Reference: &model.ExternalType{
								Type: GetModelName(ctx, builder.Definition.On),
							},
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
							Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
							Reference: &model.ExternalType{
								Type: GetModelName(ctx, builder.Definition.On),
							},
						},
					},
				},
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			str := fmt.Sprintf("entities := []*%s{}", GetModelName(ctx, builder.Definition.On)) + consts.LN
			str += fmt.Sprintf("for _, %s := range %s {", GORM_MODEL_METHOD_NAME, PluralizeName(ctx, GORM_MODEL_METHOD_NAME)) + consts.LN
			str += fmt.Sprintf("entities = append(entities,%s(%s))", ModelToGormModel(ctx, builder.Definition.On), GORM_MODEL_METHOD_NAME) + consts.LN
			str += "}" + consts.LN
			str += "return entities"
			return str, nil
		},
	})
}

func (builder *GormRepositoryBuilder) Build(ctx context.Context) (err error) {
	if builder.Err != nil {
		return builder.Err
	}

	builder.addGormModelToModel(ctx)
	builder.addModelToGormModel(ctx)
	builder.addGormModelsToModels(ctx)
	builder.addModelsToGormModels(ctx)
	builder.addMethods(ctx)

	// builder.addMethods(ctx, builder.GormDomainRepositoryBuilder)
	builder.DomainBuilder.Domain.Files = append(builder.DomainBuilder.Domain.Files, builder.GormModel)
	builder.DomainBuilder.Domain.Files = append(builder.DomainBuilder.Domain.Files, builder.Repository)

	// builder.DomainBuilder.Domain.Files = append(builder.DomainBuilder.Domain.Files, &model.File{
	// 	Name: builder.Definition.On.Name,
	// 	Pkg:  builder.DomainBuilder.GetGormAdapterPackage(),
	// 	Elements: []interface{}{
	// 		builder.GormModel,
	// 	},
	// })

	return builder.Err
}

func PrepareFieldFormGorm(ctx context.Context, field *model.Field) *model.Field {
	gormField := field.Copy()
	gormField.Tags = []*model.Tag{}
	gormTag := &model.Tag{
		Name: "gorm",
	}
	for _, tag := range field.Tags {
		if tag.Name == "gorm" {
			gormTag = tag
			break
		}
	}

	if len(gormTag.Values) != 0 {
		gormField.Tags = append(gormField.Tags, gormTag)
	}
	return gormField
}
