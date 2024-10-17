package domainbuilder

import (
	"context"
	"fmt"
	"strings"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	ERR_NOT_ALLOWED_NAME  = "ErrNotAllowed"
	ERR_NOT_ALLOWED       = "Not allowed"
	ERR_UNAUTHORIZED_NAME = "ErrForbiden"
	ERR_FORBIDEN          = "Forbiden"
)

const (
	VALIDATOR_USECASE_FIELD_NAME = "Usecase"
	CRUD_IMPL_REPO_NAME          = "DomainRepository"
	CRUD_IMPL_STUCT_NAME         = "crud"
)

type DomainUsecaseBuilder struct {
	EmptyBuilder

	domainBuilder *domainBuilder

	Err error

	structs []*model.Struct

	domainUsecase            *model.Interface
	domainUsecaseStructsFile *model.File

	validator *model.Struct
	security  *model.Struct

	repositoryDefinitions []*coredomaindefinition.Repository
	relationDefinitions   []*coredomaindefinition.Relation
	modelsDefinitions     []*coredomaindefinition.Model

	domainUsecaseImpl *model.Struct

	crudBuilders []*CRUDBuilder

	GOHttpService    *model.Struct
	GoHttpController *model.Struct
}

func NewDomainUsecaseBuilder(
	ctx context.Context,
	domainBuilder *domainBuilder,
) *DomainUsecaseBuilder {
	builder := &DomainUsecaseBuilder{
		domainBuilder: domainBuilder,
		domainUsecaseStructsFile: &model.File{
			Name: "structs",
			Pkg:  domainBuilder.GetUsecasePackage(),
		},
		domainUsecaseImpl: &model.Struct{
			Name:       GetDomainUsecaseName(ctx, domainBuilder.Domain.Name) + "CRUD",
			MethodName: CRUD_IMPL_STUCT_NAME,
			Fields: []*model.Field{
				{
					Name: CRUD_IMPL_REPO_NAME,
					Type: &model.PkgReference{
						Pkg: domainBuilder.GetRepositoryPackage(),
						Reference: &model.ExternalType{
							Type: GetDomainRepositoryName(ctx, domainBuilder.Definition),
						},
					},
				},
				{
					Name: VALIDATOR_NAME,
					Type: &model.ExternalType{
						Type: VALIDATOR_NAME,
					},
				},
			},
		},
		domainUsecase: &model.Interface{
			Name: GetDomainUsecaseName(ctx, domainBuilder.Domain.Name),
		},
	}

	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, builder.domainUsecaseStructsFile)
	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, &model.File{
		Name: domainBuilder.Domain.Name + "Usecase",
		Pkg:  domainBuilder.GetUsecasePackage(),
		Elements: []interface{}{
			builder.domainUsecase,
		},
	})
	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, &model.File{
		Name: stringtool.LowerFirstLetter(builder.domainUsecaseImpl.Name),
		Pkg:  domainBuilder.GetUsecasePackage(),
		Elements: []interface{}{
			builder.domainUsecaseImpl,
		},
	})

	builder.validator = &model.Struct{
		Name: builder.domainUsecase.Name + "Validator",
		Fields: []*model.Field{
			{
				Name: VALIDATOR_NAME,
				Type: &model.ExternalType{
					Type: VALIDATOR_NAME,
				},
			},
			{
				Name: VALIDATOR_USECASE_FIELD_NAME,
				Type: &model.ExternalType{
					Type: GetDomainUsecaseName(ctx, builder.domainBuilder.Domain.Name),
				},
			},
		},
	}
	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, &model.File{
		Name: builder.validator.Name,
		Pkg:  builder.domainBuilder.GetUsecasePackage(),
		Elements: []interface{}{
			builder.validator,
		},
	})
	builder.security = &model.Struct{
		Name: builder.domainUsecase.Name + "Security",
		Fields: []*model.Field{
			{
				Name: SECURITY_NAME,
				Type: &model.ExternalType{
					Type: SECURITY_NAME,
				},
			},
			{
				Name: VALIDATOR_USECASE_FIELD_NAME,
				Type: &model.ExternalType{
					Type: GetDomainUsecaseName(ctx, builder.domainBuilder.Domain.Name),
				},
			},
		},
	}
	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, &model.File{
		Name: builder.security.Name,
		Pkg:  builder.domainBuilder.GetUsecasePackage(),
		Elements: []interface{}{
			builder.security,
		},
	})

	return builder
}

// WithRelation implements Builder.
func (builder *DomainUsecaseBuilder) WithRelation(ctx context.Context, definition *coredomaindefinition.Relation) {
	if builder.Err != nil {
		return
	}
	builder.relationDefinitions = append(builder.relationDefinitions, definition)
}

// WithRepository implements Builder.
func (builder *DomainUsecaseBuilder) WithRepository(ctx context.Context, definition *coredomaindefinition.Repository) {
	if builder.Err != nil {
		return
	}

	builder.repositoryDefinitions = append(builder.repositoryDefinitions, definition)
}

func (builder *DomainUsecaseBuilder) WithCRUD(ctx context.Context, definition *coredomaindefinition.CRUD) {
	if builder.Err != nil {
		return
	}

	builder.crudBuilders = append(builder.crudBuilders, NewCRUDBuilder(ctx, builder.domainBuilder, definition).(*CRUDBuilder))

	if definition.Get.Active {
		method := GetCRUDMethodName(ctx, GET, definition.On)
		request := GetUsecaseRequestName(ctx, method)
		response := GetUsecaseResponseName(ctx, method)
		builder.addRolesCheck(ctx,
			method,
			request,
			response,
			definition.Get.Roles,
		)
	}
	if definition.GetActive.Active {
		method := GetCRUDMethodName(ctx, GET_ACTIVE, definition.On)
		request := GetUsecaseRequestName(ctx, GetCRUDMethodName(ctx, GET, definition.On))
		response := GetUsecaseResponseName(ctx, GetCRUDMethodName(ctx, GET, definition.On))
		builder.addRolesCheck(
			ctx,
			method,
			request,
			response,
			definition.GetActive.Roles,
		)
	}
	if definition.List.Active {
		method := GetCRUDMethodName(ctx, LIST, definition.On)
		request := GetUsecaseRequestName(ctx, method)
		response := GetUsecaseResponseName(ctx, method)
		builder.addRolesCheck(ctx,
			method,
			request,
			response,
			definition.Get.Roles,
		)
	}
	if definition.ListActive.Active {
		method := GetCRUDMethodName(ctx, LIST_ACTIVE, definition.On)
		request := GetUsecaseRequestName(ctx, GetCRUDMethodName(ctx, LIST, definition.On))
		response := GetUsecaseResponseName(ctx, GetCRUDMethodName(ctx, LIST, definition.On))
		builder.addRolesCheck(ctx,
			method,
			request,
			response,
			definition.Get.Roles,
		)
	}
	if definition.Create.Active {
		method := GetCRUDMethodName(ctx, CREATE, definition.On)
		request := GetUsecaseRequestName(ctx, method)
		response := GetUsecaseResponseName(ctx, method)
		builder.addRolesCheck(ctx,
			method,
			request,
			response,
			definition.Get.Roles,
		)
	}
	if definition.Update.Active {
		method := GetCRUDMethodName(ctx, UPDATE, definition.On)
		request := GetUsecaseRequestName(ctx, method)
		response := GetUsecaseResponseName(ctx, method)
		builder.addRolesCheck(ctx,
			method,
			request,
			response,
			definition.Get.Roles,
		)
	}
	if definition.Delete.Active {
		method := GetCRUDMethodName(ctx, DELETE, definition.On)
		request := GetUsecaseRequestName(ctx, method)
		response := GetUsecaseResponseName(ctx, method)
		builder.addRolesCheck(ctx,
			method,
			request,
			response,
			definition.Get.Roles,
		)
	}
	for _, relationCRUD := range definition.RelationCRUDs {
		var from *coredomaindefinition.Model
		var to *coredomaindefinition.Model
		if definition.On == relationCRUD.Relation.Source {
			from = relationCRUD.Relation.Source
			to = relationCRUD.Relation.Target
		} else if definition.On == relationCRUD.Relation.Target {
			from = relationCRUD.Relation.Target
			to = relationCRUD.Relation.Source
		} else {
			builder.Err = ErrRelationDoesNotBelongToModel
		}
		if relationCRUD.Add.Active {
			method := GetCRUDRelationMethodName(ctx, ADD, from, to)
			request := GetUsecaseRequestName(ctx, method)
			response := GetUsecaseResponseName(ctx, method)
			builder.addRolesCheck(
				ctx,
				method,
				request,
				response,
				relationCRUD.Add.Roles,
			)
		}
		if relationCRUD.Remove.Active {
			method := GetCRUDRelationMethodName(ctx, REMOVE, from, to)
			request := GetUsecaseRequestName(ctx, method)
			response := GetUsecaseResponseName(ctx, method)
			builder.addRolesCheck(
				ctx,
				method,
				request,
				response,
				relationCRUD.Remove.Roles,
			)
		}
		if relationCRUD.List.Active {
			method := GetCRUDRelationMethodName(ctx, LIST, from, to)
			request := GetUsecaseRequestName(ctx, method)
			response := GetUsecaseResponseName(ctx, method)
			builder.addRolesCheck(
				ctx,
				method,
				request,
				response,
				relationCRUD.List.Roles,
			)
		}
	}
}

func (builder *DomainUsecaseBuilder) addRolesCheck(ctx context.Context, method string, request string, response string, roles []string) {
	if builder.Err != nil {
		return
	}

	m := &model.Function{
		Name: method,
		Args: []*model.Param{
			CTX,
			{
				Name: REQUEST_PARAM_NAME,
				Type: &model.PointerType{
					Type: &model.ExternalType{
						Type: request,
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: &model.ExternalType{
						Type: response,
					},
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
	}
	m.Content = func() (string, []*model.GoPkg) {
		str := ""
		if len(roles) > 0 {
			str += `token := ctx.Value("token").(string)` + consts.LN
			str += `if token == "" {` + consts.LN
			str += fmt.Sprintf(`return nil, %s`, ERR_UNAUTHORIZED_NAME) + consts.LN
			str += "}" + consts.LN
			str += fmt.Sprintf(`allowed, err := %s.%s.%s(ctx, token, []string{"%s"})`, builder.security.GetMethodName(), SECURITY_NAME, SECURITY_IS_ALLOWED_METHOD_NAME, strings.Join(roles, `", "`)) + consts.LN
			str += "if err != nil {" + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN

			str += "if !allowed {" + consts.LN
			str += fmt.Sprintf("return nil, %s", ERR_NOT_ALLOWED_NAME) + consts.LN
			str += "}" + consts.LN
		}

		str += fmt.Sprintf("return %s.%s.%s(ctx, %s)",
			builder.security.GetMethodName(), VALIDATOR_USECASE_FIELD_NAME, GetUsecaseMethodName(ctx, method), REQUEST_PARAM_NAME,
		)

		return str, []*model.GoPkg{}
	}
	builder.security.Methods = append(builder.security.Methods, m)
}

func (builder *DomainUsecaseBuilder) WithUsecase(ctx context.Context, definition *coredomaindefinition.Usecase) {
	if builder.Err != nil {
		return
	}

	request := &model.Struct{
		Name:   GetUsecaseRequestName(ctx, definition.Name),
		Fields: []*model.Field{},
	}
	builder.structs = append(builder.structs, request)
	mimeTypeValidation := ""

	for _, arg := range definition.Args {
		t, err := builder.domainBuilder.TypeDefinitionToType(ctx, arg.Type)
		if err != nil {
			builder.Err = err
			return
		}

		f := &model.Field{
			Name: GetFieldName(ctx, arg.Name),
			Type: t,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{arg.Name},
				},
			},
			JsonName: arg.Name,
		}
		request.Fields = append(request.Fields, f)

		validationTags, err := GetValidationTags(ctx, arg.Validations)
		if err != nil {
			builder.Err = err
			return
		}
		if len(validationTags) > 0 {
			f.Tags = append(f.Tags, &model.Tag{
				Name:   "validate",
				Values: validationTags,
			})
		}

		if arg.Type == coredomaindefinition.PrimitiveTypeFile {
			mimeTypeValidation += fmt.Sprintf(
				`if err := %s.%s.%s(ctx, []string{"image/png", "image/jpeg", "image/webp"}, %s.%s, "%s"); err != nil {`,
				builder.validator.GetMethodName(), VALIDATOR_NAME, VALIDATOR_VALIDATE_MIME_TYPES_METHOD_NAME, REQUEST_PARAM_NAME, GetFieldName(ctx, arg.Name), arg.Name,
			) + consts.LN
			mimeTypeValidation += "return nil, err" + consts.LN
			mimeTypeValidation += "}" + consts.LN
		}
	}

	response := &model.Struct{
		Name:   GetUsecaseResponseName(ctx, definition.Name),
		Fields: []*model.Field{},
	}
	builder.structs = append(builder.structs, response)

	for _, res := range definition.Results {
		t, err := builder.domainBuilder.TypeDefinitionToType(ctx, res.Type)
		if err != nil {
			builder.Err = err
			return
		}

		f := &model.Field{
			Name: GetFieldName(ctx, res.Name),
			Type: t,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{res.Name},
				},
			},
			JsonName: res.Name,
		}
		response.Fields = append(response.Fields, f)
	}

	m := &model.Function{
		Name: GetUsecaseMethodName(ctx, definition.Name),
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
				Name: "request",
				Type: &model.PointerType{
					Type: &model.ExternalType{
						Type: GetUsecaseRequestName(ctx, definition.Name),
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: &model.ExternalType{
						Type: GetUsecaseResponseName(ctx, definition.Name),
					},
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (string, []*model.GoPkg) {
			str := fmt.Sprintf("if err := %s.%s.%s(ctx, %s); err != nil {", builder.validator.GetMethodName(), VALIDATOR_NAME, VALIDATOR_VALIDATE_METHOD_NAME, REQUEST_PARAM_NAME) + consts.LN
			str += "return nil, err" + consts.LN
			str += "}" + consts.LN

			str += mimeTypeValidation

			str += fmt.Sprintf("return %s.%s.%s(ctx, %s)",
				builder.validator.GetMethodName(), VALIDATOR_USECASE_FIELD_NAME, GetUsecaseMethodName(ctx, definition.Name), REQUEST_PARAM_NAME,
			)

			return str, []*model.GoPkg{}
		},
	}

	builder.domainUsecase.Methods = append(builder.domainUsecase.Methods, m)
	builder.validator.Methods = append(builder.validator.Methods, m)

	builder.addRolesCheck(ctx, GetUsecaseMethodName(ctx, definition.Name), GetUsecaseRequestName(ctx, definition.Name), GetUsecaseResponseName(ctx, definition.Name), definition.Roles)

	// m = m.Copy()
	// m.Content = func() (string, []*model.GoPkg) {
	// 	str := ""
	// 	if len(definition.Roles) > 0 {
	// 		str += `token := ctx.Value("token").(string)` + consts.LN
	// 		str += fmt.Sprintf(`allowed, err := %s.%s.%s(ctx, "", []string{%s})`, builder.security.GetMethodName(), SECURITY_NAME, SECURITY_IS_ALLOWED_METHOD_NAME, strings.Join(definition.Roles, ", ")) + consts.LN
	// 		str += "if err != nil {" + consts.LN
	// 		str += "return nil, err" + consts.LN
	// 		str += "}" + consts.LN

	// 		str += "if !allowed {" + consts.LN
	// 		str += fmt.Sprintf("return nil, %s", ERR_NOT_ALLOWED_NAME) + consts.LN
	// 		str += "}" + consts.LN
	// 	}

	// 	str += fmt.Sprintf("return %s.%s.%s(ctx, %s)",
	// 		builder.security.GetMethodName(), VALIDATOR_USECASE_FIELD_NAME, GetUsecaseMethodName(ctx, definition.Name), REQUEST_PARAM_NAME,
	// 	)

	// 	return str, []*model.GoPkg{}
	// }
	// builder.security.Methods = append(builder.security.Methods, m)
}

func (builder *DomainUsecaseBuilder) addSDK(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	f := &model.File{
		Name:     "structs",
		Pkg:      builder.domainBuilder.GetSdkPackage(),
		Elements: []interface{}{},
	}

	for _, s := range builder.structs {
		f.Elements = append(f.Elements, s.Copy())
	}
	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, f)

	service := &model.Interface{
		Name: GetServiceName(ctx, builder.domainBuilder.Definition),
	}
	f = &model.File{
		Name:     stringtool.LowerFirstLetter(GetServiceName(ctx, builder.domainBuilder.Definition)),
		Pkg:      builder.domainBuilder.GetSdkPackage(),
		Elements: []interface{}{service},
	}
	for _, m := range builder.domainUsecase.Methods {
		service.Methods = append(service.Methods, m.Copy())
	}
	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, f)
}

func (builder *DomainUsecaseBuilder) addHttpService(ctx context.Context) {
	if builder.Err != nil {
		return
	}
}

func (builder *DomainUsecaseBuilder) Build(ctx context.Context) error {
	if builder.Err != nil {
		return builder.Err
	}

	for _, m := range builder.modelsDefinitions {
		for _, b := range builder.crudBuilders {
			b.WithModel(ctx, m)
		}
	}

	for _, r := range builder.repositoryDefinitions {
		for _, b := range builder.crudBuilders {
			b.WithRepository(ctx, r)
		}
	}

	for _, r := range builder.relationDefinitions {
		for _, b := range builder.crudBuilders {
			b.WithRelation(ctx, r)
		}
	}

	for _, b := range builder.crudBuilders {
		err := b.Build(ctx)
		if err != nil {
			return merror.Stack(err)
		}

		for _, method := range b.Methods {
			m := method.Copy()
			m.Content = func() (string, []*model.GoPkg) {
				str := fmt.Sprintf("if err := %s.%s.%s(ctx, %s); err != nil {", builder.validator.GetMethodName(), VALIDATOR_NAME, VALIDATOR_VALIDATE_METHOD_NAME, REQUEST_PARAM_NAME) + consts.LN
				str += "return nil, err" + consts.LN
				str += "}" + consts.LN
				str += fmt.Sprintf("return %s.%s.%s(ctx, %s)",
					builder.validator.GetMethodName(), VALIDATOR_USECASE_FIELD_NAME, method.Name, REQUEST_PARAM_NAME,
				)

				return str, []*model.GoPkg{}
			}
			builder.validator.Methods = append(builder.validator.Methods, m)
		}

		builder.structs = append(builder.structs, b.Structs...)
		builder.domainUsecase.Methods = append(builder.domainUsecase.Methods, b.Methods...)
		builder.domainUsecaseImpl.Methods = append(builder.domainUsecaseImpl.Methods, b.Methods...)
	}

	for _, s := range builder.structs {
		builder.domainUsecaseStructsFile.Elements = append(builder.domainUsecaseStructsFile.Elements, s)
	}

	builder.addSDK(ctx)

	if builder.domainBuilder.Definition.Controllers.Http {
		builder.addHttpService(ctx)
	}

	builder.addValidator(ctx)
	builder.addSecurity(ctx)
	builder.addErrors(ctx)

	return nil
}

const (
	SECURITY_NAME                   = "SecurityValidator"
	SECURITY_IS_ALLOWED_METHOD_NAME = "IsAllowed"
)

const (
	VALIDATOR_NAME                            = "Validator"
	VALIDATOR_IS_VALIDATION_ERROR_METHOD_NAME = "IsValidationError"
	VALIDATOR_NEW_REFERENCE_ERROR_METHOD_NAME = "NewReferenceError"
	VALIDATOR_NEW_UNIQUE_ERROR_METHOD_NAME    = "NewUniqueError"
	VALIDATOR_VALIDATE_METHOD_NAME            = "Validate"
	VALIDATOR_VALIDATE_MIME_TYPES_METHOD_NAME = "ValidateMimeTypes"
)

func (builder *DomainUsecaseBuilder) addSecurity(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	security := &model.Interface{
		Name: SECURITY_NAME,
		Methods: []*model.Function{
			{
				Name: SECURITY_IS_ALLOWED_METHOD_NAME,
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
						Name: "token",
						Type: model.PrimitiveTypeString,
					},
					{
						Name: "roles",
						Type: &model.ArrayType{
							Type: model.PrimitiveTypeString,
						},
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeBool,
					},
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
		},
	}

	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, &model.File{
		Name: SECURITY_NAME,
		Pkg:  builder.domainBuilder.GetUsecasePackage(),
		Elements: []interface{}{
			security,
		},
	})
}

func (builder *DomainUsecaseBuilder) addValidator(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	validator := &model.Interface{
		Name: VALIDATOR_NAME,
		Methods: []*model.Function{
			{
				Name: VALIDATOR_VALIDATE_METHOD_NAME,
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
						Name: "request",
						Type: model.PrimitiveTypeInterface,
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
			{
				Name: VALIDATOR_IS_VALIDATION_ERROR_METHOD_NAME,
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
						Name: "err",
						Type: model.PrimitiveTypeError,
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeBool,
					},
				},
			},
			{
				Name: VALIDATOR_NEW_REFERENCE_ERROR_METHOD_NAME,
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
						Name: "reference",
						Type: model.PrimitiveTypeString,
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
			{
				Name: VALIDATOR_NEW_UNIQUE_ERROR_METHOD_NAME,
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
						Name: "field",
						Type: model.PrimitiveTypeString,
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
			{
				Name: VALIDATOR_VALIDATE_MIME_TYPES_METHOD_NAME,
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
						Name: "mimeTypes",
						Type: &model.ArrayType{
							Type: model.PrimitiveTypeString,
						},
					},
					{
						Name: "bytes",
						Type: model.PrimitiveTypeBytes,
					},
					{
						Name: "field",
						Type: model.PrimitiveTypeString,
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

	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, &model.File{
		Name: VALIDATOR_NAME,
		Pkg:  builder.domainBuilder.GetUsecasePackage(),
		Elements: []interface{}{
			validator,
		},
	})
}

func (builder *DomainUsecaseBuilder) addErrors(ctx context.Context) {
	if builder.Err != nil {
		return
	}

	errNotAllowed := &model.Var{
		Name: ERR_NOT_ALLOWED_NAME,
		Type: model.PrimitiveTypeError,
		Value: &model.PkgReference{
			Pkg: consts.CommonPkgs["errors"], Reference: &model.ExternalType{Type: fmt.Sprintf(`New("%s")`, ERR_NOT_ALLOWED)},
		},
	}

	errForbiden := &model.Var{
		Name: ERR_UNAUTHORIZED_NAME,
		Type: model.PrimitiveTypeError,
		Value: &model.PkgReference{
			Pkg: consts.CommonPkgs["errors"], Reference: &model.ExternalType{Type: fmt.Sprintf(`New("%s")`, ERR_FORBIDEN)},
		},
	}

	builder.domainBuilder.Domain.Files = append(builder.domainBuilder.Domain.Files, &model.File{
		Name: "errors",
		Pkg:  builder.domainBuilder.GetUsecasePackage(),
		Elements: []interface{}{
			errNotAllowed, errForbiden,
		},
	})
}
