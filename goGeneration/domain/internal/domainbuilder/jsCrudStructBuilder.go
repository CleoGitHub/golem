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

type JSCRUDStructBuilder struct {
	EmptyBuilder

	domainBuilder *domainBuilder
	definition    *coredomaindefinition.CRUD
	err           error

	GetImports        []string
	GetRequestFields  map[string]string
	GetResponseFields map[string]string

	ListImports        []string
	ListRequestFields  map[string]string
	ListResponseFields map[string]string

	CreateImports        []string
	CreateRequestFields  map[string]string
	CreateResponseFields map[string]string

	UpdateImports        []string
	UpdateRequestFields  map[string]string
	UpdateResponseFields map[string]string

	DeleteImports        []string
	DeleteRequestFields  map[string]string
	DeleteResponseFields map[string]string
}

func NewJSCRUDStructBuilder(ctx context.Context, domainBuilder *domainBuilder, definition *coredomaindefinition.CRUD) Builder {
	builder := &JSCRUDStructBuilder{
		EmptyBuilder:  EmptyBuilder{},
		domainBuilder: domainBuilder,
		definition:    definition,
	}

	if definition.Get.Active {
		builder.GetRequestFields = map[string]string{
			"id": "id",
		}

		builder.GetResponseFields = map[string]string{}
		builder.GetResponseFields[builder.definition.On.Name] = fmt.Sprintf("%s.from(%s)", GetModelName(ctx, builder.definition.On), HYDRATOR_PARAM_NAME)
		if !slices.Contains(builder.GetImports, GetModelName(ctx, builder.definition.On)) {
			builder.GetImports = append(builder.GetImports, GetModelName(ctx, builder.definition.On))
		}
	}

	if definition.List.Active {
		builder.ListRequestFields = map[string]string{
			"ordering":   fmt.Sprintf("Ordering.from(%s)", HYDRATOR_PARAM_NAME),
			"pagination": fmt.Sprintf("Pagination.from(%s)", HYDRATOR_PARAM_NAME),
		}

		builder.ListResponseFields = map[string]string{}
		builder.ListResponseFields[PluralizeName(ctx, builder.definition.On.Name)] = fmt.Sprintf("%s.map((elem) =>  %s.from(elem))", HYDRATOR_PARAM_NAME, GetModelName(ctx, builder.definition.On))
		if !slices.Contains(builder.ListImports, GetModelName(ctx, builder.definition.On)) {
			builder.ListImports = append(builder.ListImports, GetModelName(ctx, builder.definition.On))
		}
		builder.ListImports = append(builder.ListImports, "Ordering")
		builder.ListImports = append(builder.ListImports, "Pagination")
	}

	if definition.Create.Active {
		builder.CreateRequestFields = map[string]string{}
		builder.CreateRequestFields = builder.addModelFields(ctx, builder.CreateRequestFields)

		builder.CreateResponseFields = map[string]string{}
		builder.CreateResponseFields[builder.definition.On.Name] = fmt.Sprintf("%s.from(%s)", GetModelName(ctx, builder.definition.On), HYDRATOR_PARAM_NAME)
		if !slices.Contains(builder.CreateImports, GetModelName(ctx, builder.definition.On)) {
			builder.CreateImports = append(builder.CreateImports, GetModelName(ctx, builder.definition.On))
		}
	}

	if definition.Update.Active {
		builder.UpdateRequestFields = map[string]string{}
		builder.UpdateRequestFields[stringtool.LowerFirstLetter(consts.ID)] = stringtool.LowerFirstLetter(consts.ID)
		builder.UpdateRequestFields = builder.addModelFields(ctx, builder.UpdateRequestFields)

		builder.UpdateResponseFields = map[string]string{}
		builder.UpdateResponseFields[builder.definition.On.Name] = fmt.Sprintf("%s.from(%s)", GetModelName(ctx, builder.definition.On), HYDRATOR_PARAM_NAME)
		if !slices.Contains(builder.CreateImports, GetModelName(ctx, builder.definition.On)) {
			builder.CreateImports = append(builder.CreateImports, GetModelName(ctx, builder.definition.On))
		}
	}

	if definition.Delete.Active {
		builder.GetRequestFields = map[string]string{
			"id": "id",
		}
	}

	for _, relationCRUD := range definition.RelationCRUDs {
		builder.addRelationCRUD(ctx, relationCRUD)
	}

	return builder
}

func (builder *JSCRUDStructBuilder) addRelationCRUD(ctx context.Context, definition *coredomaindefinition.RelationCRUD) {
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

func (builder *JSCRUDStructBuilder) addRelationCRUDAdd(ctx context.Context, relation *coredomaindefinition.Relation, from *coredomaindefinition.Model, to *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}

	action := GetCRUDRelationMethodName(ctx, ADD, from, to)
	content := JSGetClassFromSimpleFields(
		GetUsecaseRequestName(ctx, action),
		[]string{
			stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, from)),
			stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, to)),
		},
	)
	content += consts.LN

	content += JSGetClassFromSimpleFields(
		GetUsecaseResponseName(ctx, action),
		[]string{},
	)

	if builder.domainBuilder.Domain.JSFiles == nil {
		builder.domainBuilder.Domain.JSFiles = map[string]string{}
	}

	builder.domainBuilder.Domain.JSFiles[stringtool.LowerFirstLetter(action)] = content
}

func (builder *JSCRUDStructBuilder) addRelationCRUDRemove(ctx context.Context, relation *coredomaindefinition.Relation, from *coredomaindefinition.Model, to *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}

	action := GetCRUDRelationMethodName(ctx, REMOVE, from, to)
	content := JSGetClassFromSimpleFields(
		GetUsecaseRequestName(ctx, action),
		[]string{
			stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, from)),
			stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, to)),
		},
	)
	content += consts.LN

	content += JSGetClassFromSimpleFields(
		GetUsecaseResponseName(ctx, action),
		[]string{},
	)

	if builder.domainBuilder.Domain.JSFiles == nil {
		builder.domainBuilder.Domain.JSFiles = map[string]string{}
	}

	builder.domainBuilder.Domain.JSFiles[stringtool.LowerFirstLetter(action)] = content
}

func (builder *JSCRUDStructBuilder) addRelationCRUDList(ctx context.Context, relation *coredomaindefinition.Relation, from *coredomaindefinition.Model, to *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}

	action := GetCRUDRelationMethodName(ctx, LIST, from, to)
	content := fmt.Sprintf("import %s from './';", GetModelName(ctx, to)) + consts.LN
	content += consts.LN
	content += JSGetClassFromSimpleFields(
		GetUsecaseRequestName(ctx, action),
		[]string{
			stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, from)),
		},
	)
	content += consts.LN

	content += JSGetClassFromTransformationFields(
		GetUsecaseResponseName(ctx, action),
		map[string]string{
			PluralizeName(ctx, to.Name): fmt.Sprintf("%s.from(data)", GetModelName(ctx, to)),
		},
	)

	if builder.domainBuilder.Domain.JSFiles == nil {
		builder.domainBuilder.Domain.JSFiles = map[string]string{}
	}

	builder.domainBuilder.Domain.JSFiles[stringtool.LowerFirstLetter(action)] = content
}

func (builder *JSCRUDStructBuilder) WithRelation(ctx context.Context, definition *coredomaindefinition.Relation) {
	if builder.err != nil {
		return
	}

	if definition.Source != builder.definition.On && definition.Target != builder.definition.On {
		return
	}

	var to *coredomaindefinition.Model
	if definition.Source == builder.definition.On {
		to = definition.Target
	} else if definition.Target == builder.definition.On {
		to = definition.Source
	} else {
		builder.err = merror.Stack(ErrRelationDoesNotBelongToModel)
	}

	if !IsRelationMultiple(ctx, builder.definition.On, definition) {
		fieldName := stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, to))
		if builder.definition.Create.Active {
			builder.CreateRequestFields[fieldName] = fmt.Sprintf("%s.%s", HYDRATOR_PARAM_NAME, fieldName)
		}
		if builder.definition.Update.Active {
			builder.UpdateRequestFields[fieldName] = fmt.Sprintf("%s.%s", HYDRATOR_PARAM_NAME, fieldName)
		}
	}
}

func (builder *JSCRUDStructBuilder) addFile(ctx context.Context, action string, imports []string, requestFields map[string]string, responseFields map[string]string) {
	if builder.err != nil {
		return
	}

	action = stringtool.LowerFirstLetter(action)
	if builder.domainBuilder.Domain.JSFiles == nil {
		builder.domainBuilder.Domain.JSFiles = map[string]string{}
	}
	if builder.domainBuilder.Domain.JSFiles[action] == "" {
		builder.domainBuilder.Domain.JSFiles[action] = ""
	}

	str := ""
	if len(imports) > 0 {
		str += fmt.Sprintf("import {\n\t%s\n} from './';", strings.Join(imports, ",\n\t")) + consts.LN
		str += consts.LN
	}

	str += JSGetClassFromTransformationFields(action, requestFields)
	str += consts.LN

	str += JSGetClassFromTransformationFields(action, responseFields)
	str += consts.LN

	builder.domainBuilder.Domain.JSFiles[action] = str
}

func (builder *JSCRUDStructBuilder) addModelFields(ctx context.Context, fields map[string]string) map[string]string {
	if builder.err != nil {
		return fields
	}

	for _, field := range builder.definition.On.Fields {
		fields[field.Name] = field.Name
	}

	if builder.definition.On.Activable {
		fields[stringtool.LowerFirstLetter(ACTIVE_FIELD_NAME)] = stringtool.LowerFirstLetter(ACTIVE_FIELD_NAME)
	}

	return fields
}

func (builder *JSCRUDStructBuilder) Build(ctx context.Context) error {
	if builder.err != nil {
		return builder.err
	}

	if builder.definition.Get.Active {
		builder.addFile(ctx, GetCRUDMethodName(ctx, GET, builder.definition.On), builder.GetImports, builder.GetRequestFields, builder.GetResponseFields)
	}

	if builder.definition.GetActive.Active {
		builder.addFile(ctx, GetCRUDMethodName(ctx, GET_ACTIVE, builder.definition.On), builder.GetImports, builder.GetRequestFields, builder.GetResponseFields)
	}

	if builder.definition.List.Active {
		builder.addFile(ctx, GetCRUDMethodName(ctx, LIST, builder.definition.On), builder.ListImports, builder.ListRequestFields, builder.ListResponseFields)
	}

	if builder.definition.ListActive.Active {
		builder.addFile(ctx, GetCRUDMethodName(ctx, LIST_ACTIVE, builder.definition.On), builder.ListImports, builder.ListRequestFields, builder.ListResponseFields)
	}

	if builder.definition.Create.Active {
		builder.addFile(ctx, GetCRUDMethodName(ctx, CREATE, builder.definition.On), builder.CreateImports, builder.CreateRequestFields, builder.CreateResponseFields)
	}

	if builder.definition.Update.Active {
		builder.addFile(ctx, GetCRUDMethodName(ctx, UPDATE, builder.definition.On), builder.UpdateImports, builder.UpdateRequestFields, builder.UpdateResponseFields)
	}

	if builder.definition.Delete.Active {
		builder.addFile(ctx, GetCRUDMethodName(ctx, DELETE, builder.definition.On), builder.DeleteImports, builder.DeleteRequestFields, builder.DeleteResponseFields)
	}

	return nil
}
