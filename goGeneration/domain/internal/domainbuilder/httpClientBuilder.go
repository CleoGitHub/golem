package domainbuilder

import (
	"bytes"
	"context"
	"text/template"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const HTTP_CLIENT_ROUTE_TEMPLATE = `	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := {{ .StructName }}.NewRequest(http.MethodPost, "{{ .Route }}", body, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	resp, status, err := {{ .StructName }}.Do(req)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, httpclient.ErrUnexpectedStatus
	}
	result := &{{ .UsecasePkg }}.{{ .ResultType }}{}
	if err = json.Unmarshal(resp, result); err != nil {
		return nil, err
	}
	return result, nil
`

type HttpClientRouteTemplate struct {
	Route      string
	ResultType string
	UsecasePkg string
	StructName string
}

type HttpClientBuilder struct {
	EmptyBuilder

	domainDefinition *coredomaindefinition.Domain
	domain           *model.Domain

	client *model.Struct

	err error
}

func NewHttpClientBuilder(ctx context.Context, domainDefinition *coredomaindefinition.Domain, domain *model.Domain) *HttpClientBuilder {
	return &HttpClientBuilder{
		domainDefinition: domainDefinition,
		domain:           domain,

		client: &model.Struct{
			Name: GetHttpClientName(ctx, domainDefinition),
			Fields: []*model.Field{
				{
					Type: &model.PkgReference{
						Pkg: consts.CommonPkgs["httpclient"],
						Reference: &model.ExternalType{
							Type: "HttpClient",
						},
					},
				},
			},
		},
	}
}

func (builder *HttpClientBuilder) WithCRUD(ctx context.Context, definition *coredomaindefinition.CRUD) {
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

func (builder *HttpClientBuilder) addCRUDAction(ctx context.Context, action string, on *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}
	method := GetCRUDMethodName(ctx, action, on)
	request := GetUsecaseRequestName(ctx, method)
	response := GetUsecaseResponseName(ctx, method)
	switch action {
	case GET_ACTIVE:
		request = GetUsecaseRequestName(ctx, GetCRUDMethodName(ctx, GET, on))
		response = GetUsecaseResponseName(ctx, GetCRUDMethodName(ctx, GET, on))
	case LIST_ACTIVE:
		request = GetUsecaseRequestName(ctx, GetCRUDMethodName(ctx, LIST, on))
		response = GetUsecaseResponseName(ctx, GetCRUDMethodName(ctx, LIST, on))
	}
	route := &model.Function{
		Name: method,
		Args: []*model.Param{
			CTX,
			{
				Name: "request",
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.domain.Architecture.UsecasePkg,
						Reference: &model.ExternalType{
							Type: request,
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.domain.Architecture.UsecasePkg,
						Reference: &model.ExternalType{
							Type: response,
						},
					},
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
	}
	route.Content = func() (content string, requiredPkg []*model.GoPkg) {
		return builder.getRouteContent(ctx, GetUsecaseMethodName(ctx, method), response)
	}

	builder.client.Methods = append(builder.client.Methods, route)
}

func (builder *HttpClientBuilder) addRelationCRUDAction(ctx context.Context, action string, from *coredomaindefinition.Model, to *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}
	method := GetCRUDRelationMethodName(ctx, action, from, to)
	route := &model.Function{
		Name: method,
		Args: []*model.Param{
			CTX,
			{
				Name: "request",
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.domain.Architecture.UsecasePkg,
						Reference: &model.ExternalType{
							Type: GetUsecaseRequestName(ctx, method),
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.domain.Architecture.UsecasePkg,
						Reference: &model.ExternalType{
							Type: GetUsecaseResponseName(ctx, method),
						},
					},
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
	}
	route.Content = func() (content string, requiredPkg []*model.GoPkg) {
		return builder.getRouteContent(ctx, GetUsecaseMethodName(ctx, method), GetUsecaseResponseName(ctx, method))
	}

	builder.client.Methods = append(builder.client.Methods, route)
}

func (builder *HttpClientBuilder) getRouteContent(ctx context.Context, method string, response string) (string, []*model.GoPkg) {
	tmpl := HttpClientRouteTemplate{
		Route:      GetHttpRouteName(ctx, builder.domainDefinition, method),
		ResultType: response,
		UsecasePkg: builder.domain.Architecture.UsecasePkg.Alias,
		StructName: builder.client.GetMethodName(),
	}

	buffer := bytes.NewBufferString("")
	err := template.Must(template.New("route").Parse(HTTP_CLIENT_ROUTE_TEMPLATE)).Execute(buffer, tmpl)
	if err != nil {
		panic(err)
	}
	return buffer.String(), []*model.GoPkg{
		consts.CommonPkgs["json"],
		consts.CommonPkgs["http"],
	}
}

func (builder *HttpClientBuilder) WithUsecase(ctx context.Context, definition *coredomaindefinition.Usecase) {
	if builder.err != nil {
		return
	}

	method := GetUsecaseMethodName(ctx, definition.Name)
	route := &model.Function{
		Name: method,
		Args: []*model.Param{
			CTX,
			{
				Name: "request",
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.domain.Architecture.UsecasePkg,
						Reference: &model.ExternalType{
							Type: GetUsecaseRequestName(ctx, method),
						},
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.domain.Architecture.UsecasePkg,
						Reference: &model.ExternalType{
							Type: GetUsecaseResponseName(ctx, method),
						},
					},
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
	}
	route.Content = func() (content string, requiredPkg []*model.GoPkg) {
		return builder.getRouteContent(ctx, GetUsecaseMethodName(ctx, method), GetUsecaseResponseName(ctx, method))
	}

	builder.client.Methods = append(builder.client.Methods, route)
}

func (builder *HttpClientBuilder) Build(ctx context.Context) error {
	if builder.err != nil {
		return builder.err
	}

	f := &model.File{
		Name:     GetHttpClientName(ctx, builder.domainDefinition),
		Pkg:      builder.domain.Architecture.SdkPkg,
		Elements: []interface{}{builder.client},
	}
	builder.domain.Files = append(builder.domain.Files, f)

	return nil
}
