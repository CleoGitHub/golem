package domainbuilder

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
)

const (
	JS_CLIENT_HOST        = "host"
	JS_CLIENT_PORT        = "port"
	JS_CLIENT_HTTP_CLIENT = "httpclient"
)

type JSServiceBuilder struct {
	EmptyBuilder

	domainBuilder *domainBuilder

	methods string

	// types        string
	typesImports []string

	err error
}

func NewJSServiceBuilder(domainBuilder *domainBuilder) Builder {
	builder := &JSServiceBuilder{
		domainBuilder: domainBuilder,
	}

	return builder
}

func (builder *JSServiceBuilder) WithCRUD(ctx context.Context, definition *coredomaindefinition.CRUD) {
	if builder.err != nil {
		return
	}

	builder.domainBuilder.AddBuilder(ctx, NewJSCRUDStructBuilder(ctx, builder.domainBuilder, definition))

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

func (builder *JSServiceBuilder) addCRUDAction(ctx context.Context, action string, on *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}
	method := GetCRUDMethodName(ctx, action, on)
	request := GetUsecaseRequestName(ctx, method)
	switch action {
	case GET_ACTIVE:
		request = GetUsecaseRequestName(ctx, GetCRUDMethodName(ctx, GET, on))
	case LIST_ACTIVE:
		request = GetUsecaseRequestName(ctx, GetCRUDMethodName(ctx, LIST, on))
	}
	builder.addMethod(ctx, method, request)
}

func (builder *JSServiceBuilder) addRelationCRUDAction(ctx context.Context, action string, from *coredomaindefinition.Model, to *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}
	method := GetCRUDRelationMethodName(ctx, action, from, to)
	request := GetUsecaseRequestName(ctx, method)
	builder.addMethod(ctx, method, request)
}

func (builder *JSServiceBuilder) WithUsecase(ctx context.Context, definition *coredomaindefinition.Usecase) {
	if builder.err != nil {
		return
	}

	method := GetUsecaseMethodName(ctx, definition.Name)
	request := GetUsecaseRequestName(ctx, method)
	builder.addMethod(ctx, method, request)

	requestFields := []string{}
	for _, field := range definition.Args {
		requestFields = append(requestFields, field.Name)
	}

	imports := []string{}
	responseFields := map[string]string{}
	for _, field := range definition.Results {
		if m, ok := field.Type.(*coredomaindefinition.Model); ok {
			if !slices.Contains(imports, GetModelName(ctx, m)) {
				imports = append(imports, GetModelName(ctx, m))
			}
			responseFields[field.Name] = fmt.Sprintf("%s.from(data)", GetModelName(ctx, m))
		} else if a, ok := field.Type.(*coredomaindefinition.Array); ok {
			if m, ok := a.Type.(*coredomaindefinition.Model); ok {
				if !slices.Contains(imports, GetModelName(ctx, m)) {
					imports = append(imports, GetModelName(ctx, m))
				}
				responseFields[field.Name] = fmt.Sprintf("data.map((el) => %s.from(el))", GetModelName(ctx, m))
			} else {
				responseFields[field.Name] = fmt.Sprintf("data.%s", field.Name)
			}
		} else {
			responseFields[field.Name] = fmt.Sprintf("data.%s", field.Name)
		}
	}

	content := ""
	if len(imports) > 0 {
		content += fmt.Sprintf("import {\n\t%s\n} from './';", strings.Join(imports, ",\n\t")) + consts.LN
		content += consts.LN
	}
	content += JSGetClassFromSimpleFields(request, requestFields) + consts.LN
	content += JSGetClassFromTransformationFields(GetUsecaseResponseName(ctx, method), responseFields)

	if builder.domainBuilder.Domain.JSFiles == nil {
		builder.domainBuilder.Domain.JSFiles = map[string]string{}
	}
	builder.domainBuilder.Domain.JSFiles[stringtool.LowerFirstLetter(method)] = content
}

func (builder *JSServiceBuilder) addMethod(ctx context.Context, method string, request string) {
	if builder.err != nil {
		return
	}

	builder.typesImports = append(builder.typesImports, GetUsecaseResponseName(ctx, method))

	builder.methods += consts.TAB + fmt.Sprintf("%s(%s) {", stringtool.LowerFirstLetter(method), stringtool.LowerFirstLetter(request)) + consts.LN
	builder.methods += consts.TAB + consts.TAB + "return new Promise((resolve, reject) => {" + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + fmt.Sprintf("this.%s.post(", JS_CLIENT_HTTP_CLIENT) + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + consts.TAB + fmt.Sprintf(
		"this.%s + ':' + this.%s + '%s',",
		JS_CLIENT_HOST, JS_CLIENT_PORT, GetHttpRouteName(ctx, builder.domainBuilder.Definition, stringtool.LowerFirstLetter(method)),
	) + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + consts.TAB + fmt.Sprintf("JSON.stringify(%s),", stringtool.LowerFirstLetter(request)) + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + consts.TAB + "{ headers: { 'Content-Type': 'application/json', 'Accept': 'application/json' } }" + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + ")" + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + ".then(response => {" + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + consts.TAB + fmt.Sprintf("resolve(%s.from(response.data))", GetUsecaseResponseName(ctx, method)) + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + "})" + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + ".catch(error => {" + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + consts.TAB + "reject(error)" + consts.LN
	builder.methods += consts.TAB + consts.TAB + consts.TAB + "})" + consts.LN
	builder.methods += consts.TAB + consts.TAB + "})" + consts.LN
	builder.methods += consts.TAB + "}" + consts.LN
	builder.methods += consts.LN
}

func (builder *JSServiceBuilder) Build(ctx context.Context) error {
	if builder.err != nil {
		return builder.err
	}

	content := fmt.Sprintf("import {\n\t%s\n} from './';", strings.Join(builder.typesImports, ",\n\t")) + consts.LN
	content += consts.LN

	content += fmt.Sprintf("export class %sService {", stringtool.UpperFirstLetter(builder.domainBuilder.Definition.Name)) + consts.LN
	content += consts.TAB + JS_CLIENT_HOST + consts.LN
	content += consts.TAB + JS_CLIENT_PORT + consts.LN
	content += consts.TAB + JS_CLIENT_HTTP_CLIENT + consts.LN
	content += consts.LN

	content += consts.TAB + fmt.Sprintf("constructor(%s, %s, %s) {", JS_CLIENT_HOST, JS_CLIENT_PORT, JS_CLIENT_HTTP_CLIENT) + consts.LN
	content += consts.TAB + consts.TAB + fmt.Sprintf("this.%s = %s", JS_CLIENT_HOST, JS_CLIENT_HOST) + consts.LN
	content += consts.TAB + consts.TAB + fmt.Sprintf("this.%s = %s", JS_CLIENT_PORT, JS_CLIENT_PORT) + consts.LN
	content += consts.TAB + consts.TAB + fmt.Sprintf("this.%s = %s", JS_CLIENT_HTTP_CLIENT, JS_CLIENT_HTTP_CLIENT) + consts.LN
	content += consts.TAB + "}" + consts.LN
	content += consts.LN

	content += builder.methods

	// end class content
	content += "}" + consts.LN
	content += consts.LN

	if builder.domainBuilder.Domain.JSFiles == nil {
		builder.domainBuilder.Domain.JSFiles = map[string]string{}
	}
	if builder.domainBuilder.Domain.JSFiles[fmt.Sprintf("%sService", builder.domainBuilder.Definition.Name)] == "" {
		builder.domainBuilder.Domain.JSFiles[fmt.Sprintf("%sService", builder.domainBuilder.Definition.Name)] = ""
	}
	builder.domainBuilder.Domain.JSFiles[fmt.Sprintf("%sService", builder.domainBuilder.Definition.Name)] += content + consts.LN

	return nil
}
