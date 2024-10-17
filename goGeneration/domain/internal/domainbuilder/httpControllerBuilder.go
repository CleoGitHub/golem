package domainbuilder

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const ROUTE_TMPL = `ctx := r.Context()

request := &{{ .UsecasePkg }}.{{ .UsecaseRequest }}{}

err := json.NewDecoder(r.Body).Decode(request)
if err != nil {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
	return
}

{{ .OptionalFieldExtraction }}

result, err := {{ .ControllerName }}.{{ .Usecases }}.{{ .UsecaseName }}(ctx, request)
if err != nil {
	if {{ .ControllerName }}.{{ .Validator }}.{{ .IsValidationErrorMethodName }}(ctx, err) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	} else if errors.Is({{ .RepositoryPkg }}.{{ .RepostiroyErrNotFound }}, err) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

w.WriteHeader(http.StatusOK)
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(result)
`

const (
	ROUTE_REGISTRATION_PARAM = "router"
)

type RouteTemplate struct {
	ControllerName              string
	UsecaseName                 string
	UsecasePkg                  string
	UsecaseRequest              string
	Validator                   string
	RepositoryPkg               string
	RepostiroyErrNotFound       string
	Usecases                    string
	IsValidationErrorMethodName string
	OptionalFieldExtraction     string
}

type HttpControllerBuilder struct {
	EmptyBuilder

	domainDefinition *coredomaindefinition.Domain
	domain           *model.Domain

	controller *model.Struct

	routeRegistration string
}

func NewHttpControllerBuilder(ctx context.Context, domainDefinition *coredomaindefinition.Domain, domain *model.Domain) *HttpControllerBuilder {
	return &HttpControllerBuilder{
		domainDefinition: domainDefinition,
		domain:           domain,
		controller: &model.Struct{
			Name: GetHttpControllerName(ctx, domainDefinition),
			Fields: []*model.Field{
				{
					Name: GetDomainUsecaseName(ctx, domainDefinition.Name),
					Type: &model.PkgReference{
						Pkg: domain.Architecture.UsecasePkg,
						Reference: &model.ExternalType{
							Type: GetDomainUsecaseName(ctx, domainDefinition.Name),
						},
					},
				},
				{
					Name: VALIDATOR_NAME,
					Type: &model.PkgReference{
						Pkg: domain.Architecture.UsecasePkg,
						Reference: &model.ExternalType{
							Type: VALIDATOR_NAME,
						},
					},
				},
			},
		},
	}
}

func (builder *HttpControllerBuilder) WithCRUD(ctx context.Context, definition *coredomaindefinition.CRUD) {
	if builder.err != nil {
		return
	}

	if definition.Get.Active {
		builder.addCRUDAction(ctx, GET, definition.On)
	}

	if definition.GetActive.Active {
		builder.addCRUDAction(ctx, GET_ACTIVE, definition.On)
	}

	if definition.List.Active {
		builder.addCRUDAction(ctx, LIST, definition.On)
	}

	if definition.ListActive.Active {
		builder.addCRUDAction(ctx, LIST_ACTIVE, definition.On)
	}

	if definition.Create.Active {
		builder.addCRUDAction(ctx, CREATE, definition.On)
	}

	if definition.Update.Active {
		builder.addCRUDAction(ctx, UPDATE, definition.On)
	}

	if definition.Delete.Active {
		builder.addCRUDAction(ctx, DELETE, definition.On)
	}

	for _, relationCRUD := range definition.RelationCRUDs {
		var from *coredomaindefinition.Model
		var to *coredomaindefinition.Model
		if relationCRUD.Relation.Source == definition.On {
			from = relationCRUD.Relation.Source
			to = relationCRUD.Relation.Target
		} else if relationCRUD.Relation.Target == definition.On {
			from = relationCRUD.Relation.Target
			to = relationCRUD.Relation.Source
		} else {
			builder.err = merror.Stack(ErrRelationDoesNotBelongToModel)
			return
		}

		if relationCRUD.Add.Active {
			builder.addRelationCRUDAction(ctx, ADD, from, to)
		}

		if relationCRUD.Remove.Active {
			builder.addRelationCRUDAction(ctx, REMOVE, from, to)
		}

		if relationCRUD.List.Active {
			builder.addRelationCRUDAction(ctx, LIST, from, to)
		}

	}
}

func (builder *HttpControllerBuilder) addCRUDAction(ctx context.Context, action string, on *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}
	method := GetCRUDMethodName(ctx, action, on)
	route := GetHttpRoute(ctx, method)
	switch action {
	case GET_ACTIVE:
		method = GetCRUDMethodName(ctx, GET, on)
	case LIST_ACTIVE:
		method = GetCRUDMethodName(ctx, LIST, on)
	}
	route.Content = func() (content string, requiredPkg []*model.GoPkg) {
		return builder.getRouteContent(ctx, GetUsecaseMethodName(ctx, method), GetUsecaseRequestName(ctx, method), "")
	}

	builder.controller.Methods = append(builder.controller.Methods, route)

	builder.routeRegistration += fmt.Sprintf(
		`%s.Post("%s", %s.%s)`,
		ROUTE_REGISTRATION_PARAM, GetHttpRouteName(ctx, builder.domainDefinition, method), builder.controller.GetMethodName(), method,
	) + consts.LN
}

func (builder *HttpControllerBuilder) addRelationCRUDAction(ctx context.Context, action string, from *coredomaindefinition.Model, to *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}
	method := GetCRUDRelationMethodName(ctx, action, from, to)
	route := GetHttpRoute(ctx, method)
	route.Content = func() (content string, requiredPkg []*model.GoPkg) {
		return builder.getRouteContent(ctx, GetUsecaseMethodName(ctx, method), GetUsecaseRequestName(ctx, method), "")
	}

	builder.controller.Methods = append(builder.controller.Methods, route)

	builder.routeRegistration += fmt.Sprintf(
		`%s.Post("%s", %s.%s)`,
		ROUTE_REGISTRATION_PARAM, GetHttpRouteName(ctx, builder.domainDefinition, method), builder.controller.GetMethodName(), method,
	) + consts.LN
}

func (builder *HttpControllerBuilder) getRouteContent(ctx context.Context, method string, request string, optionalFieldExtraction string) (string, []*model.GoPkg) {
	tmpl := RouteTemplate{
		ControllerName:              builder.controller.GetMethodName(),
		UsecaseName:                 method,
		UsecasePkg:                  builder.domain.Architecture.UsecasePkg.Alias,
		UsecaseRequest:              request,
		Validator:                   VALIDATOR_NAME,
		RepositoryPkg:               builder.domain.Architecture.RepositoryPkg.Alias,
		RepostiroyErrNotFound:       REPOSITORY_ERROR_NOT_FOUND.Name,
		Usecases:                    GetDomainUsecaseName(ctx, builder.domainDefinition.Name),
		IsValidationErrorMethodName: VALIDATOR_IS_VALIDATION_ERROR_METHOD_NAME,
		OptionalFieldExtraction:     optionalFieldExtraction,
	}

	buffer := bytes.NewBufferString("")
	err := template.Must(template.New("route").Parse(ROUTE_TMPL)).Execute(buffer, tmpl)
	if err != nil {
		panic(err)
	}
	return buffer.String(), []*model.GoPkg{
		builder.domain.Architecture.RepositoryPkg,
		consts.CommonPkgs["errors"],
		consts.CommonPkgs["json"],
	}
}

func (builder *HttpControllerBuilder) WithUsecase(ctx context.Context, definition *coredomaindefinition.Usecase) {
	if builder.err != nil {
		return
	}

	method := GetUsecaseMethodName(ctx, definition.Name)
	route := GetHttpRoute(ctx, method)
	route.Content = func() (content string, requiredPkg []*model.GoPkg) {
		fileIdx := 0
		optionalContent := ""
		for _, arg := range definition.Args {
			if arg.Type == coredomaindefinition.PrimitiveTypeFile {
				str := fmt.Sprintf(`file%d, _, err := r.FormFile("%s")`, fileIdx, arg.Name) + consts.LN
				str += "if err != nil {" + consts.LN
				str += "w.WriteHeader(http.StatusBadRequest)" + consts.LN
				str += "return" + consts.LN
				str += "}" + consts.LN
				str += fmt.Sprintf(`defer file%d.Close()`, fileIdx) + consts.LN

				str += fmt.Sprintf(`fileBytes%d, err := io.ReadAll(file%d)`, fileIdx, fileIdx) + consts.LN
				str += "if err != nil {" + consts.LN
				str += "w.WriteHeader(http.StatusBadRequest)" + consts.LN
				str += "return" + consts.LN
				str += "}" + consts.LN

				str += fmt.Sprintf("request.%s = fileBytes%d", GetFieldName(ctx, arg.Name), fileIdx) + consts.LN

				optionalContent += str
				fileIdx++
			}
		}

		if optionalContent != "" {
			optionalContent = "r.ParseMultipartForm(10 << 20)" + consts.LN + consts.LN + optionalContent
		}

		str, pkgs := builder.getRouteContent(ctx, GetUsecaseMethodName(ctx, method), GetUsecaseRequestName(ctx, method), optionalContent)

		return str, append(pkgs, consts.CommonPkgs["io"])
	}

	builder.controller.Methods = append(builder.controller.Methods, route)

	builder.routeRegistration += fmt.Sprintf(
		`%s.Post("%s", %s.%s)`,
		ROUTE_REGISTRATION_PARAM, GetHttpRouteName(ctx, builder.domainDefinition, method), builder.controller.GetMethodName(), method,
	) + consts.LN
}

func (builder *HttpControllerBuilder) addRoutesRegistration(ctx context.Context) {
	if builder.err != nil {
		return
	}

	builder.controller.Methods = append(builder.controller.Methods, &model.Function{
		Name: "RegisterRoutes",
		Args: []*model.Param{
			{
				Name: "router",
				Type: &model.ExternalType{
					Type: "Router",
				},
			},
		},
		Content: func() (content string, requiredPkg []*model.GoPkg) {
			return builder.routeRegistration, []*model.GoPkg{}
		},
	})
}

func (builder *HttpControllerBuilder) addTypeRegisterRoute(ctx context.Context) {
	if builder.err != nil {
		return
	}

	t := &model.Interface{
		Name: "Router",
		Methods: []*model.Function{
			{
				Name: "Post",
				Args: []*model.Param{
					{
						Name: "route",
						Type: model.PrimitiveTypeString,
					},
					{
						Name: "handler",
						Type: &model.Function{
							Args: []*model.Param{
								{
									Type: &model.PkgReference{
										Pkg: consts.CommonPkgs["http"],
										Reference: &model.ExternalType{
											Type: "ResponseWriter",
										},
									},
								},
								{
									Type: &model.PkgReference{
										Pkg: consts.CommonPkgs["http"],
										Reference: &model.ExternalType{
											Type: "Request",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	f := &model.File{
		Name:     "Router",
		Pkg:      builder.domain.Architecture.HttpControllerPkg,
		Elements: []interface{}{t},
	}

	builder.domain.Files = append(builder.domain.Files, f)
}

func (builder *HttpControllerBuilder) Build(ctx context.Context) error {
	if builder.err != nil {
		return builder.err
	}

	builder.addTypeRegisterRoute(ctx)
	builder.addRoutesRegistration(ctx)

	f := &model.File{
		Name:     GetHttpControllerName(ctx, builder.domainDefinition),
		Pkg:      builder.domain.Architecture.HttpControllerPkg,
		Elements: []interface{}{builder.controller},
	}
	builder.domain.Files = append(builder.domain.Files, f)

	return nil
}
