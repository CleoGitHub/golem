package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	REQUEST_PARAM_NAME = "request"
	GET                = "Get"
	GET_ACTIVE         = "GetActive"
	LIST               = "List"
	LIST_ACTIVE        = "ListActive"
	CREATE             = "Create"
	UPDATE             = "Update"
	DELETE             = "Delete"
	ADD                = "Add"
	REMOVE             = "Remove"
)

type CRUDBuilder struct {
	EmptyBuilder

	domainBuilder *domainBuilder
	definition    *coredomaindefinition.CRUD

	err error

	getRequest  *model.Struct
	getResponse *model.Struct
	get         *model.Function
	getActive   *model.Function

	listRequest  *model.Struct
	listResponse *model.Struct
	list         *model.Function
	listActive   *model.Function

	createRequest  *model.Struct
	createResponse *model.Struct
	create         *model.Function

	updateRequest  *model.Struct
	updateResponse *model.Struct
	update         *model.Function

	deleteRequest  *model.Struct
	deleteResponse *model.Struct
	delete         *model.Function

	createValidation string
	updateValidation string

	requestFieldToField string

	Methods []*model.Function
	Structs []*model.Struct
}

func NewCRUDBuilder(ctx context.Context, domainBuilder *domainBuilder, definition *coredomaindefinition.CRUD) Builder {
	builder := &CRUDBuilder{
		domainBuilder: domainBuilder,
		definition:    definition,
	}

	builder.addValidationChecks(ctx)
	builder.addModelRequestFieldToField(ctx)

	if definition.Get.Active {
		builder.addGet(ctx)
	}
	if definition.GetActive.Active {
		builder.addGetActive(ctx)
	}
	if definition.List.Active {
		builder.addList(ctx)
	}
	if definition.ListActive.Active {
		builder.addListActive(ctx)
	}
	if definition.Create.Active {
		builder.addCreate(ctx)
	}
	if definition.Update.Active {
		builder.addUpdate(ctx)
	}
	if definition.Delete.Active {
		builder.addDelete(ctx)
	}

	for _, relationCRUD := range definition.RelationCRUDs {
		builder.addRelationCRUD(ctx, relationCRUD)
	}

	return builder
}

func (builder *CRUDBuilder) addModelRequestFieldToField(ctx context.Context) {
	if builder.err != nil {
		return
	}

	if builder.definition.On.Activable {
		builder.requestFieldToField += fmt.Sprintf("%s: %s.%s,", ACTIVE_FIELD_NAME, REQUEST_PARAM_NAME, ACTIVE_FIELD_NAME) + consts.LN
	}

	for _, field := range builder.definition.On.Fields {
		builder.requestFieldToField += fmt.Sprintf("%s: %s.%s,", GetFieldName(ctx, field.Name), REQUEST_PARAM_NAME, GetFieldName(ctx, field.Name)) + consts.LN
	}
}

func (builder *CRUDBuilder) addValidationChecks(ctx context.Context) {
	if builder.err != nil {
		return
	}

	for _, field := range builder.definition.On.Fields {
		for _, validation := range field.Validations {
			switch validation.Rule {
			case coredomaindefinition.ValidationRuleUnique:
				builder.addUniqueCheck(ctx, field)
			case coredomaindefinition.ValidationRuleUniqueIn:
				value := validation.Value
				if in, ok := value.(*coredomaindefinition.Model); ok {
					builder.addUniqueInCheck(ctx, field, in)
				} else {
					builder.err = NewErrValidationValueExpectedType(string(validation.Rule), "*coredomaindefinition.Model")
				}
			}
		}
	}

}

func (builder *CRUDBuilder) addUniqueCheck(ctx context.Context, field *coredomaindefinition.Field) {
	if builder.err != nil {
		return
	}

	repoAlias := builder.domainBuilder.GetRepositoryPackage().Alias

	str := fmt.Sprintf("// validate uniqueness %s", field.Name) + consts.LN
	str += fmt.Sprintf("if _, err := %s.%s.%s(", CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryGetMethod(ctx, builder.definition.On)) + consts.LN
	str += "ctx," + consts.LN
	str += fmt.Sprintf("%s.%s.%s([]*%s.%s{", repoAlias, GetRepositoryGetMethod(ctx, builder.definition.On), GetOptName(ctx, REPOSITORY_BY), repoAlias, REPOSITORY_WHERE) + consts.LN
	str += "{" + consts.LN
	str += fmt.Sprintf(`%s: "%s",`, REPOSITORY_WHERE_KEY, GetFieldName(ctx, field.Name)) + consts.LN
	str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_OPERATOR, repoAlias, REPOSITORY_WHERE_OPERATOR_EQUAL) + consts.LN
	str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_VALUE, REQUEST_PARAM_NAME, GetFieldName(ctx, field.Name)) + consts.LN
	str += "}," + consts.LN

	builder.createValidation += str

	// Add where not id of request update element
	str += "{" + consts.LN
	str += fmt.Sprintf(`%s: "%s",`, REPOSITORY_WHERE_KEY, consts.ID) + consts.LN
	str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_OPERATOR, repoAlias, REPOSITORY_WHERE_OPERATOR_NOT_EQUAL) + consts.LN
	str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_VALUE, REQUEST_PARAM_NAME, consts.ID) + consts.LN
	str += "}," + consts.LN
	builder.updateValidation += str

	str = ""
	str += "})," + consts.LN
	str += fmt.Sprintf("); err != nil && err != %s.%s {", repoAlias, REPOSITORY_ERROR_NOT_FOUND.Name) + consts.LN
	str += "return nil, err" + consts.LN
	str += fmt.Sprintf("} else if err != %s.%s {", repoAlias, REPOSITORY_ERROR_NOT_FOUND.Name) + consts.LN
	str += "// should get ErrNotFound" + consts.LN
	str += fmt.Sprintf(`return nil, %s.%s.%s(ctx, "%s")`, CRUD_IMPL_STUCT_NAME, VALIDATOR_NAME, VALIDATOR_NEW_UNIQUE_ERROR_METHOD_NAME, GetFieldName(ctx, field.Name)) + consts.LN
	str += "}" + consts.LN

	builder.createValidation += str
	builder.updateValidation += str
}

func (builder *CRUDBuilder) addUniqueInCheck(ctx context.Context, field *coredomaindefinition.Field, in *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}

	repoAlias := builder.domainBuilder.GetRepositoryPackage().Alias

	str := fmt.Sprintf("// validate uniqueness %s in %s", field.Name, in.Name) + consts.LN
	str += fmt.Sprintf("if _, err := %s.%s.%s(", CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryGetMethod(ctx, builder.definition.On)) + consts.LN
	str += "ctx," + consts.LN
	str += fmt.Sprintf("%s.%s.%s([]*%s.%s{", repoAlias, GetRepositoryGetMethod(ctx, builder.definition.On), GetOptName(ctx, REPOSITORY_BY), repoAlias, REPOSITORY_WHERE) + consts.LN
	str += "{" + consts.LN
	str += fmt.Sprintf(`%s: "%s",`, REPOSITORY_WHERE_KEY, GetFieldName(ctx, field.Name)) + consts.LN
	str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_OPERATOR, repoAlias, REPOSITORY_WHERE_OPERATOR_EQUAL) + consts.LN
	str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_VALUE, REQUEST_PARAM_NAME, GetFieldName(ctx, field.Name)) + consts.LN
	str += "}," + consts.LN
	str += "{" + consts.LN
	str += fmt.Sprintf(`%s: "%s",`, REPOSITORY_WHERE_KEY, GetSingleRelationIdName(ctx, in)) + consts.LN
	str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_OPERATOR, repoAlias, REPOSITORY_WHERE_OPERATOR_EQUAL) + consts.LN
	str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_VALUE, REQUEST_PARAM_NAME, GetSingleRelationIdName(ctx, in)) + consts.LN
	str += "}," + consts.LN

	builder.createValidation += str

	// Add where not id of request update element
	str += "{" + consts.LN
	str += fmt.Sprintf(`%s: "%s",`, REPOSITORY_WHERE_KEY, consts.ID) + consts.LN
	str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_OPERATOR, repoAlias, REPOSITORY_WHERE_OPERATOR_NOT_EQUAL) + consts.LN
	str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_VALUE, REQUEST_PARAM_NAME, consts.ID) + consts.LN
	str += "}," + consts.LN
	builder.updateValidation += str

	str = ""
	str += "})," + consts.LN
	str += fmt.Sprintf("); err != nil && err != %s.%s {", repoAlias, REPOSITORY_ERROR_NOT_FOUND.Name) + consts.LN
	str += "return nil, err" + consts.LN
	str += fmt.Sprintf("} else if err != %s.%s {", repoAlias, REPOSITORY_ERROR_NOT_FOUND.Name) + consts.LN
	str += "// should get ErrNotFound" + consts.LN
	str += fmt.Sprintf(`return nil, %s.%s.%s(ctx, "%s")`, CRUD_IMPL_STUCT_NAME, VALIDATOR_NAME, VALIDATOR_NEW_UNIQUE_ERROR_METHOD_NAME, GetFieldName(ctx, field.Name)) + consts.LN
	str += "}" + consts.LN

	builder.createValidation += str
	builder.updateValidation += str
}

func (builder *CRUDBuilder) WithRelation(ctx context.Context, definition *coredomaindefinition.Relation) {
	if builder.err != nil {
		return
	}

	if definition.Source != builder.definition.On && definition.Target != builder.definition.On {
		return
	}

	var to *coredomaindefinition.Model
	if definition.Source == builder.definition.On {
		to = definition.Target
	} else {
		if definition.IgnoreReverse {
			return
		}
		to = definition.Source
	}

	if !IsRelationMultiple(ctx, builder.definition.On, definition) {
		tags := []string{"uuid"}
		optionnal, err := IsRelationOptionnal(ctx, builder.definition.On, definition)
		if err != nil {
			builder.err = err
			return
		}
		if !optionnal {
			tags = append(tags, "required")
		}
		field := &model.Field{
			Name: GetSingleRelationIdName(ctx, to),
			Type: model.PrimitiveTypeString,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, to))},
				},
				{
					Name:   "validate",
					Values: tags,
				},
			},
		}
		if builder.createRequest != nil {
			builder.createRequest.Fields = append(builder.createRequest.Fields, field)
		}
		if builder.updateRequest != nil {
			builder.updateRequest.Fields = append(builder.updateRequest.Fields, field)
		}

		repoAlias := builder.domainBuilder.GetRepositoryPackage().Alias
		str := fmt.Sprintf("// validate relation %s", GetModelName(ctx, to)) + consts.LN
		if optionnal {
			str += fmt.Sprintf(`if %s.%s != "" {`, REQUEST_PARAM_NAME, GetSingleRelationIdName(ctx, to)) + consts.LN
		}
		str += fmt.Sprintf("if _, err := %s.%s.%s(", CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryGetMethod(ctx, to)) + consts.LN
		str += "ctx," + consts.LN
		str += fmt.Sprintf("%s.%s.%s([]*%s.%s{", repoAlias, GetRepositoryGetMethod(ctx, to), GetOptName(ctx, REPOSITORY_BY), repoAlias, REPOSITORY_WHERE) + consts.LN
		str += "{" + consts.LN
		str += fmt.Sprintf(`%s: "%s",`, REPOSITORY_WHERE_KEY, consts.ID) + consts.LN
		str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_OPERATOR, repoAlias, REPOSITORY_WHERE_OPERATOR_EQUAL) + consts.LN
		str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_VALUE, REQUEST_PARAM_NAME, GetSingleRelationIdName(ctx, to)) + consts.LN
		str += "}," + consts.LN
		str += "})," + consts.LN
		str += ");  err != nil {" + consts.LN
		// str += "" + consts.LN
		str += fmt.Sprintf("if err == %s.%s {", repoAlias, REPOSITORY_ERROR_NOT_FOUND.Name) + consts.LN
		str += fmt.Sprintf(`return nil, %s.%s.%s(ctx, "%s")`, CRUD_IMPL_STUCT_NAME, VALIDATOR_NAME, VALIDATOR_NEW_REFERENCE_ERROR_METHOD_NAME, GetSingleRelationIdName(ctx, to)) + consts.LN
		str += "}" + consts.LN
		str += "return nil, err" + consts.LN
		str += "}" + consts.LN
		if optionnal {
			str += "}" + consts.LN
		}
		builder.createValidation += str
		builder.updateValidation += str

		builder.requestFieldToField += fmt.Sprintf("%s: %s.%s,", GetSingleRelationIdName(ctx, to), REQUEST_PARAM_NAME, GetSingleRelationIdName(ctx, to)) + consts.LN
	}
}

func (builder *CRUDBuilder) WithRepository(ctx context.Context, definition *coredomaindefinition.Repository) {
	if builder.err != nil {
		return
	}
}

func (builder *CRUDBuilder) WithUsecase(ctx context.Context, definition *coredomaindefinition.Usecase) {
	if builder.err != nil {
		return
	}
}

func (builder *CRUDBuilder) addGet(ctx context.Context) {
	if builder.err != nil {
		return
	}
	action := GetCRUDMethodName(ctx, GET, builder.definition.On)

	builder.buildGetStructs(ctx)

	builder.get = &model.Function{
		Name: action,
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
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: builder.getRequest,
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: builder.getResponse,
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			repoAlias := builder.domainBuilder.GetRepositoryPackage().Alias
			str := fmt.Sprintf("entity, err := %s.%s.%s(", CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryGetMethod(ctx, builder.definition.On)) + consts.LN
			str += "ctx," + consts.LN
			str += fmt.Sprintf("%s.%s.%s([]*%s.%s{", repoAlias, GetRepositoryGetMethod(ctx, builder.definition.On), GetOptName(ctx, REPOSITORY_BY), repoAlias, REPOSITORY_WHERE) + consts.LN
			str += "{" + consts.LN
			str += fmt.Sprintf(`%s: "%s",`, REPOSITORY_WHERE_KEY, consts.ID) + consts.LN
			str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_OPERATOR, repoAlias, REPOSITORY_WHERE_OPERATOR_EQUAL) + consts.LN
			str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_VALUE, REQUEST_PARAM_NAME, consts.ID) + consts.LN
			str += "}," + consts.LN
			str += "})," + consts.LN
			if node := builder.domainBuilder.RelationGraph.GetNode(builder.definition.On); node != nil && node.RequireRetriveInactive() {
				str += fmt.Sprintf(`%s.%s.%s(true),`, repoAlias, GetRepositoryGetMethod(ctx, builder.definition.On), GetOptName(ctx, REPOSITORY_RETRIEVE_INACTIVE)) + consts.LN
			}
			str += ")" + consts.LN
			str += "if err != nil {" + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN

			str += fmt.Sprintf("return &%s{%s: entity}, nil", GetUsecaseResponseName(ctx, action), GetModelName(ctx, builder.definition.On))
			return str, []*model.GoPkg{}
		},
	}
	builder.Methods = append(builder.Methods, builder.get)
}

func (builder *CRUDBuilder) addGetActive(ctx context.Context) {
	if builder.err != nil {
		return
	}
	if !builder.domainBuilder.RelationGraph.GetNode(builder.definition.On).RequireRetriveInactive() {
		builder.err = merror.Stack(ErrModelNotActivable)
		return
	}

	action := GetCRUDMethodName(ctx, GET_ACTIVE, builder.definition.On)

	builder.buildGetStructs(ctx)

	builder.getActive = &model.Function{
		Name: action,
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
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: builder.getRequest,
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: builder.getResponse,
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			repoAlias := builder.domainBuilder.GetRepositoryPackage().Alias
			str := fmt.Sprintf("entity, err := %s.%s.%s(", CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryGetMethod(ctx, builder.definition.On)) + consts.LN
			str += "ctx," + consts.LN
			str += fmt.Sprintf("%s.%s.%s([]*%s.%s{", repoAlias, GetRepositoryGetMethod(ctx, builder.definition.On), GetOptName(ctx, REPOSITORY_BY), repoAlias, REPOSITORY_WHERE) + consts.LN
			str += "{" + consts.LN
			str += fmt.Sprintf(`%s: "%s",`, REPOSITORY_WHERE_KEY, consts.ID) + consts.LN
			str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_OPERATOR, repoAlias, REPOSITORY_WHERE_OPERATOR_EQUAL) + consts.LN
			str += fmt.Sprintf(`%s: %s.%s,`, REPOSITORY_WHERE_VALUE, REQUEST_PARAM_NAME, consts.ID) + consts.LN
			str += "}," + consts.LN
			str += "})," + consts.LN
			str += ")" + consts.LN
			str += "if err != nil {" + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN

			str += fmt.Sprintf("return &%s{%s: entity}, nil", GetUsecaseResponseName(ctx, GetCRUDMethodName(ctx, GET, builder.definition.On)), GetModelName(ctx, builder.definition.On))
			return str, []*model.GoPkg{}
		},
	}
	builder.Methods = append(builder.Methods, builder.getActive)
}

func (builder *CRUDBuilder) addList(ctx context.Context) {
	if builder.err != nil {
		return
	}
	action := GetCRUDMethodName(ctx, LIST, builder.definition.On)

	builder.buildListStructs(ctx)

	builder.list = &model.Function{
		Name: action,
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
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: builder.listRequest,
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: builder.listResponse,
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			repoAlias := builder.domainBuilder.GetRepositoryPackage().Alias
			str := fmt.Sprintf("entity, err := %s.%s.%s(", CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryListMethod(ctx, builder.definition.On)) + consts.LN
			str += "ctx," + consts.LN
			str += fmt.Sprintf("%s.%s.%s([]*%s.%s{", repoAlias, GetRepositoryListMethod(ctx, builder.definition.On), GetOptName(ctx, REPOSITORY_BY), repoAlias, REPOSITORY_WHERE) + consts.LN
			str += "})," + consts.LN
			if node := builder.domainBuilder.RelationGraph.GetNode(builder.definition.On); node != nil && node.RequireRetriveInactive() {
				str += fmt.Sprintf(`%s.%s.%s(true),`, repoAlias, GetRepositoryListMethod(ctx, builder.definition.On), GetOptName(ctx, REPOSITORY_RETRIEVE_INACTIVE)) + consts.LN
			}
			str += fmt.Sprintf(`%s.%s.%s(%s.%s),`, repoAlias, GetRepositoryListMethod(ctx, builder.definition.On), GetOptName(ctx, PAGINATION_NAME), REQUEST_PARAM_NAME, PAGINATION_NAME) + consts.LN
			str += fmt.Sprintf(`%s.%s.%s(%s.%s),`, repoAlias, GetRepositoryListMethod(ctx, builder.definition.On), GetOptName(ctx, ORDERING_NAME), REQUEST_PARAM_NAME, ORDERING_NAME) + consts.LN
			str += ")" + consts.LN
			str += "if err != nil {" + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN

			str += fmt.Sprintf("return &%s{%s: entity}, nil", GetUsecaseResponseName(ctx, action), PluralizeName(ctx, GetModelName(ctx, builder.definition.On)))
			return str, []*model.GoPkg{}
		},
	}
	builder.Methods = append(builder.Methods, builder.list)
}

func (builder *CRUDBuilder) addListActive(ctx context.Context) {
	if builder.err != nil {
		return
	}
	if !builder.domainBuilder.RelationGraph.GetNode(builder.definition.On).RequireRetriveInactive() {
		builder.err = merror.Stack(NewErrRelationNotActivable(builder.definition.On.Name))
		return
	}

	action := GetCRUDMethodName(ctx, LIST_ACTIVE, builder.definition.On)

	builder.buildListStructs(ctx)

	builder.listActive = &model.Function{
		Name: action,
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
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: builder.listRequest,
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: builder.listResponse,
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			return "return nil, nil", []*model.GoPkg{}
		},
	}
	builder.Methods = append(builder.Methods, builder.listActive)
}

func (builder *CRUDBuilder) addCreate(ctx context.Context) {
	if builder.err != nil {
		return
	}
	action := GetCRUDMethodName(ctx, CREATE, builder.definition.On)
	builder.createRequest = &model.Struct{
		Name:   GetUsecaseRequestName(ctx, action),
		Fields: []*model.Field{},
	}
	builder.addDefaultFieldsToModificationStruct(ctx, builder.createRequest)
	builder.createResponse = &model.Struct{
		Name: GetUsecaseResponseName(ctx, action),
		Fields: []*model.Field{
			{
				Name: GetModelName(ctx, builder.definition.On),
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.domainBuilder.GetModelPackage(),
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, builder.definition.On),
						},
					},
				},
			},
		},
	}
	builder.Structs = append(builder.Structs, builder.createRequest, builder.createResponse)

	builder.create = &model.Function{
		Name: action,
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
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: builder.createRequest,
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: builder.createResponse,
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			str := builder.createValidation

			str += fmt.Sprintf("uuid := %s.NewString()", consts.CommonPkgs["uuid"].Alias) + consts.LN

			str += fmt.Sprintf("entity := &%s.%s{", builder.domainBuilder.GetModelPackage().Alias, GetModelName(ctx, builder.definition.On)) + consts.LN
			str += fmt.Sprintf("%s: uuid,", consts.ID) + consts.LN
			str += builder.requestFieldToField
			str += "}" + consts.LN

			str += fmt.Sprintf("entity, err := %s.%s.%s(ctx, entity)", CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryCreateMethod(ctx, builder.definition.On)) + consts.LN
			str += "if err != nil {" + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN

			str += fmt.Sprintf("return &%s{%s: entity}, nil", GetUsecaseResponseName(ctx, action), GetModelName(ctx, builder.definition.On)) + consts.LN
			return str, []*model.GoPkg{
				consts.CommonPkgs["uuid"],
				builder.domainBuilder.GetModelPackage(),
			}
		},
	}
	builder.Methods = append(builder.Methods, builder.create)
}

func (builder *CRUDBuilder) addUpdate(ctx context.Context) {
	if builder.err != nil {
		return
	}
	action := GetCRUDMethodName(ctx, UPDATE, builder.definition.On)
	builder.updateRequest = &model.Struct{
		Name:   GetUsecaseRequestName(ctx, action),
		Fields: []*model.Field{},
	}
	builder.addIdFieldToStruct(ctx, builder.updateRequest)
	builder.addDefaultFieldsToModificationStruct(ctx, builder.updateRequest)
	builder.Structs = append(builder.Structs, builder.updateRequest)

	builder.updateResponse = &model.Struct{
		Name: GetUsecaseResponseName(ctx, action),
		Fields: []*model.Field{
			{
				Name: GetModelName(ctx, builder.definition.On),
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.domainBuilder.GetModelPackage(),
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, builder.definition.On),
						},
					},
				},
			},
		},
	}
	builder.Structs = append(builder.Structs, builder.updateResponse)

	builder.update = &model.Function{
		Name: action,
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
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: builder.updateRequest,
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: builder.updateResponse,
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			str := builder.updateValidation

			str += fmt.Sprintf("entity := &%s.%s{", builder.domainBuilder.GetModelPackage().Alias, GetModelName(ctx, builder.definition.On)) + consts.LN
			str += fmt.Sprintf("%s: %s.%s,", consts.ID, REQUEST_PARAM_NAME, consts.ID) + consts.LN
			str += builder.requestFieldToField
			str += "}" + consts.LN

			str += fmt.Sprintf("entity, err := %s.%s.%s(ctx, entity)", CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryCreateMethod(ctx, builder.definition.On)) + consts.LN
			str += "if err != nil {" + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN

			str += fmt.Sprintf("return &%s{%s: entity}, nil", GetUsecaseResponseName(ctx, action), GetModelName(ctx, builder.definition.On)) + consts.LN
			return str, []*model.GoPkg{
				consts.CommonPkgs["uuid"],
				builder.domainBuilder.GetModelPackage(),
			}
		},
	}
	builder.Methods = append(builder.Methods, builder.update)
}

func (builder *CRUDBuilder) addDelete(ctx context.Context) {
	if builder.err != nil {
		return
	}
	action := GetCRUDMethodName(ctx, DELETE, builder.definition.On)
	builder.deleteRequest = &model.Struct{
		Name: GetUsecaseRequestName(ctx, action),
		Fields: []*model.Field{
			{
				Name: "Id",
				Type: model.PrimitiveTypeString,
				Tags: []*model.Tag{
					{
						Name:   "json",
						Values: []string{"id"},
					},
					{
						Name:   "validate",
						Values: []string{"required", "uuid"},
					},
				},
			},
		},
	}
	builder.Structs = append(builder.Structs, builder.deleteRequest)
	builder.deleteResponse = &model.Struct{
		Name:   GetUsecaseResponseName(ctx, action),
		Fields: []*model.Field{},
	}
	builder.Structs = append(builder.Structs, builder.deleteResponse)

	builder.delete = &model.Function{
		Name: action,
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
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: builder.deleteRequest,
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: builder.deleteResponse,
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			str := fmt.Sprintf("err := %s.%s.%s(ctx, %s.%s)", CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryDeleteMethod(ctx, builder.definition.On), REQUEST_PARAM_NAME, consts.ID) + consts.LN
			str += "if err != nil {" + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN
			str += fmt.Sprintf("return &%s{}, nil", GetUsecaseResponseName(ctx, action))
			return str, []*model.GoPkg{}
		},
	}
	builder.Methods = append(builder.Methods, builder.delete)
}

func (builder *CRUDBuilder) addRelationCRUD(ctx context.Context, definition *coredomaindefinition.RelationCRUD) {
	if builder.err != nil {
		return
	}
	var from *coredomaindefinition.Model
	var to *coredomaindefinition.Model
	if definition.Relation.Source == builder.definition.On {
		from = definition.Relation.Source
		to = definition.Relation.Target
	} else if definition.Relation.Target == builder.definition.On {
		from = definition.Relation.Target
		to = definition.Relation.Source
	} else {
		builder.err = merror.Stack(ErrRelationDoesNotBelongToModel)
	}

	if definition.Add.Active {
		builder.addRelationCRUDAdd(ctx, definition.Relation, from, to)
	}

	if definition.Remove.Active {
		builder.addRelationCRUDRemove(ctx, definition.Relation, from, to)
	}

	if definition.List.Active {
		builder.addRelationCRUDList(ctx, definition.Relation, from, to)
	}
}

func (builder *CRUDBuilder) addRelationCRUDAdd(ctx context.Context, relation *coredomaindefinition.Relation, from *coredomaindefinition.Model, to *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}

	action := GetCRUDRelationMethodName(ctx, ADD, from, to)
	request := &model.Struct{
		Name: GetUsecaseRequestName(ctx, action),
		Fields: []*model.Field{
			{
				Name: GetSingleRelationIdName(ctx, from),
				Type: model.PrimitiveTypeString,
				Tags: []*model.Tag{
					{
						Name:   "json",
						Values: []string{stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, from))},
					},
					{
						Name:   "validate",
						Values: []string{"required", "uuid"},
					},
				},
			},
			{
				Name: GetSingleRelationIdName(ctx, to),
				Type: model.PrimitiveTypeString,
				Tags: []*model.Tag{
					{
						Name:   "json",
						Values: []string{stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, to))},
					},
					{
						Name:   "validate",
						Values: []string{"required", "uuid"},
					},
				},
			},
		},
	}
	builder.Structs = append(builder.Structs, request)

	response := &model.Struct{
		Name:   GetUsecaseResponseName(ctx, action),
		Fields: []*model.Field{},
	}
	builder.Structs = append(builder.Structs, response)

	method := &model.Function{
		Name: action,
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
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: request,
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: response,
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			r := REQUEST_PARAM_NAME
			str := fmt.Sprintf(
				"err := %s.%s.%s(ctx, %s.%s, %s.%s)",
				CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryAddRelationMethod(ctx, from, relation), r, GetSingleRelationIdName(ctx, from), r, GetSingleRelationIdName(ctx, to),
			) + consts.LN
			str += "if err != nil {" + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN
			str += fmt.Sprintf("return &%s{}, nil", GetUsecaseResponseName(ctx, action))
			return str, []*model.GoPkg{}
		},
	}
	builder.Methods = append(builder.Methods, method)
}

func (builder *CRUDBuilder) addRelationCRUDRemove(ctx context.Context, relation *coredomaindefinition.Relation, from *coredomaindefinition.Model, to *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}

	action := GetCRUDRelationMethodName(ctx, REMOVE, from, to)
	request := &model.Struct{
		Name: GetUsecaseRequestName(ctx, action),
		Fields: []*model.Field{
			{
				Name: GetSingleRelationIdName(ctx, from),
				Type: model.PrimitiveTypeString,
				Tags: []*model.Tag{
					{
						Name:   "json",
						Values: []string{stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, from))},
					},
					{
						Name:   "validate",
						Values: []string{"required", "uuid"},
					},
				},
			},
			{
				Name: GetSingleRelationIdName(ctx, to),
				Type: model.PrimitiveTypeString,
				Tags: []*model.Tag{
					{
						Name:   "json",
						Values: []string{stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, to))},
					},
					{
						Name:   "validate",
						Values: []string{"required", "uuid"},
					},
				},
			},
		},
	}
	builder.Structs = append(builder.Structs, request)

	response := &model.Struct{
		Name:   GetUsecaseResponseName(ctx, action),
		Fields: []*model.Field{},
	}
	builder.Structs = append(builder.Structs, response)

	method := &model.Function{
		Name: action,
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
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: request,
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: response,
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			r := REQUEST_PARAM_NAME
			str := fmt.Sprintf(
				"err := %s.%s.%s(ctx, %s.%s, %s.%s)",
				CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryRemoveRelationMethod(ctx, from, relation), r, GetSingleRelationIdName(ctx, from), r, GetSingleRelationIdName(ctx, to),
			) + consts.LN
			str += "if err != nil {" + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN
			str += fmt.Sprintf("return &%s{}, nil", GetUsecaseResponseName(ctx, action))
			return str, []*model.GoPkg{}
		},
	}
	builder.Methods = append(builder.Methods, method)
}

func (builder *CRUDBuilder) addRelationCRUDList(ctx context.Context, relation *coredomaindefinition.Relation, from *coredomaindefinition.Model, to *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}

	action := GetCRUDRelationMethodName(ctx, LIST, from, to)

	request := &model.Struct{
		Name: GetUsecaseRequestName(ctx, action),
		Fields: []*model.Field{
			{
				Name: GetSingleRelationIdName(ctx, from),
				Type: model.PrimitiveTypeString,
				Tags: []*model.Tag{
					{
						Name:   "json",
						Values: []string{stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, from))},
					},
					{
						Name:   "validate",
						Values: []string{"required", "uuid"},
					},
				},
			},
		},
	}
	builder.Structs = append(builder.Structs, request)

	response := &model.Struct{
		Name: GetUsecaseResponseName(ctx, action),
		Fields: []*model.Field{
			{
				Name: PluralizeName(ctx, GetModelName(ctx, to)),
				Type: &model.ArrayType{
					Type: &model.PointerType{
						Type: &model.PkgReference{
							Pkg: builder.domainBuilder.GetModelPackage(),
							Reference: &model.ExternalType{
								Type: GetModelName(ctx, to),
							},
						},
					},
				},
			},
		},
	}
	builder.Structs = append(builder.Structs, response)

	method := &model.Function{
		Name: action,
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
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: request,
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: response,
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			// r := REQUEST_PARAM_NAME
			v := PluralizeName(ctx, to.Name)
			str := fmt.Sprintf(
				"%s, err := %s.%s.%s(ctx,",
				v, CRUD_IMPL_STUCT_NAME, CRUD_IMPL_REPO_NAME, GetRepositoryListMethod(ctx, to),
			) + consts.LN

			repoPkg := builder.domainBuilder.GetRepositoryPackage().Alias
			str += fmt.Sprintf(
				"%s.%s.%s([]*%s.%s{",
				repoPkg, GetRepositoryListMethod(ctx, to), GetOptName(ctx, REPOSITORY_BY), repoPkg, REPOSITORY_WHERE,
			) + consts.LN
			str += "{" + consts.LN
			str += fmt.Sprintf(`%s: "%s",`, REPOSITORY_WHERE_KEY, GetSingleRelationIdName(ctx, from)) + consts.LN
			str += fmt.Sprintf("%s: %s.%s,", REPOSITORY_WHERE_OPERATOR, repoPkg, REPOSITORY_WHERE_OPERATOR_EQUAL) + consts.LN
			str += fmt.Sprintf("%s: %s.%s,", REPOSITORY_WHERE_VALUE, REQUEST_PARAM_NAME, GetSingleRelationIdName(ctx, from)) + consts.LN
			str += "}," + consts.LN
			str += "}," + consts.LN
			str += ")," + consts.LN
			str += ")" + consts.LN
			str += "if err != nil {" + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN

			str += fmt.Sprintf("return &%s{%s: %s}, nil", GetUsecaseResponseName(ctx, action), PluralizeName(ctx, GetModelName(ctx, to)), v) + consts.LN
			return str, []*model.GoPkg{}
		},
	}
	builder.Methods = append(builder.Methods, method)
}

func (builder *CRUDBuilder) Build(ctx context.Context) error {
	if builder.err != nil {
		return builder.err
	}

	return nil
}

func (builder *CRUDBuilder) addIdFieldToStruct(ctx context.Context, request *model.Struct) {
	if builder.err != nil {
		return
	}
	request.Fields = append(request.Fields, &model.Field{
		Name: "Id",
		Type: model.PrimitiveTypeString,
		Tags: []*model.Tag{
			{
				Name:   "json",
				Values: []string{"id"},
			},
			{
				Name:   "validate",
				Values: []string{"required", "uuid"},
			},
		},
	})
}

func (builder *CRUDBuilder) addDefaultFieldsToModificationStruct(ctx context.Context, request *model.Struct) {
	if builder.err != nil {
		return
	}
	if builder.definition.On.Activable {
		request.Fields = append(request.Fields, &model.Field{
			Name: "Active",
			Type: model.PrimitiveTypeBool,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{"active"},
				},
			},
		})
	}
	for _, f := range builder.definition.On.Fields {
		field, err := builder.domainBuilder.FieldDefinitionToField(ctx, f)
		if err != nil {
			builder.err = merror.Stack(err)
			return
		}
		validationTags, err := GetValidationTags(ctx, f.Validations)
		if err != nil {
			builder.err = merror.Stack(err)
			return
		}
		if len(validationTags) > 0 {
			field.Tags = append(field.Tags, &model.Tag{
				Name:   "validate",
				Values: validationTags,
			})
		}
		request.Fields = append(request.Fields, field)
	}
}

func (builder *CRUDBuilder) buildGetStructs(ctx context.Context) {
	if builder.err != nil {
		return
	}
	action := GetCRUDMethodName(ctx, GET, builder.definition.On)
	if builder.getRequest == nil {
		builder.getRequest = &model.Struct{
			Name: GetUsecaseRequestName(ctx, action),
			Fields: []*model.Field{
				{
					Name: "Id",
					Type: model.PrimitiveTypeString,
					Tags: []*model.Tag{
						{
							Name:   "json",
							Values: []string{"id"},
						},
						{
							Name:   "validate",
							Values: []string{"required", "uuid"},
						},
					},
				},
			},
		}
		builder.Structs = append(builder.Structs, builder.getRequest)
	}
	if builder.getResponse == nil {
		builder.getResponse = &model.Struct{
			Name: GetUsecaseResponseName(ctx, action),
			Fields: []*model.Field{
				{
					Name: GetModelName(ctx, builder.definition.On),
					Type: &model.PointerType{
						Type: &model.PkgReference{
							Pkg: builder.domainBuilder.GetModelPackage(),
							Reference: &model.ExternalType{
								Type: GetModelName(ctx, builder.definition.On),
							},
						},
					},
				},
			},
		}
		builder.Structs = append(builder.Structs, builder.getResponse)
	}
}

func (builder *CRUDBuilder) buildListStructs(ctx context.Context) {
	if builder.err != nil {
		return
	}
	action := GetCRUDMethodName(ctx, LIST, builder.definition.On)
	if builder.listRequest == nil {
		builder.listRequest = &model.Struct{
			Name: GetUsecaseRequestName(ctx, action),
			Fields: []*model.Field{
				{
					Name: PAGINATION_NAME,
					Type: &model.PkgReference{
						Pkg: builder.domainBuilder.GetModelPackage(),
						Reference: &model.ExternalType{
							Type: PAGINATION_NAME,
						},
					},
				},
				{
					Name: ORDERING_NAME,
					Type: &model.PkgReference{
						Pkg: builder.domainBuilder.GetModelPackage(),
						Reference: &model.ExternalType{
							Type: ORDERING_NAME,
						},
					},
				},
			},
		}
		builder.Structs = append(builder.Structs, builder.listRequest)
	}
	if builder.listResponse == nil {
		builder.listResponse = &model.Struct{
			Name: GetUsecaseResponseName(ctx, action),
			Fields: []*model.Field{
				{
					Name: PluralizeName(ctx, GetModelName(ctx, builder.definition.On)),
					Type: &model.ArrayType{
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Pkg: builder.domainBuilder.GetModelPackage(),
								Reference: &model.ExternalType{
									Type: GetModelName(ctx, builder.definition.On),
								},
							},
						},
					},
				},
			},
		}
		builder.Structs = append(builder.Structs, builder.listResponse)
	}
}
