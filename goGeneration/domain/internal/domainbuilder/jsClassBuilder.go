package domainbuilder

import (
	"context"
	"fmt"
	"strings"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
)

const HYDRATOR_PARAM_NAME = "data"
const JS_FROM_METHOD_NAME = "from"
const JS_HYDRATE_METHOD_NAME = "hydrate"

type JSClassBuilder struct {
	*EmptyBuilder

	domainBuilder *domainBuilder

	err error

	definition *coredomaindefinition.Model

	fields            string
	constructorParams string
	constructor       string
	hydrator          string
}

func NewJSClassBuilder(ctx context.Context, domainBuilder *domainBuilder, m *coredomaindefinition.Model) Builder {
	builder := &JSClassBuilder{
		definition:    m,
		domainBuilder: domainBuilder,
	}
	builder.addModelFields(ctx)

	return builder
}

func (builder *JSClassBuilder) addModelFields(ctx context.Context) {
	if builder.err != nil {
		return
	}

	for _, f := range builder.domainBuilder.DefaultModelFields {
		builder.addModelField(ctx, f)
	}
	if builder.definition.Archivable {
		builder.addModelField(ctx, &coredomaindefinition.Field{
			Name: "deletedAt",
			Type: coredomaindefinition.PrimitiveTypeDateTime,
		})
	}
	// Add activable field if model is activable
	if builder.definition.Activable {
		builder.addModelField(ctx, &coredomaindefinition.Field{
			Name: "active",
			Type: coredomaindefinition.PrimitiveTypeBool,
		})
	}

	for _, f := range builder.definition.Fields {
		builder.addModelField(ctx, f)
	}
}

func (builder *JSClassBuilder) addModelField(ctx context.Context, f *coredomaindefinition.Field) {
	if builder.err != nil {
		return
	}

	t, err := builder.domainBuilder.TypeDefinitionToType(ctx, f.Type)
	if err != nil {
		builder.err = err
		return
	}

	builder.fields += consts.TAB + fmt.Sprintf("%s // %s", f.Name, t.GetType()) + consts.LN
	builder.constructorParams += fmt.Sprintf("%s,", f.Name)
	builder.constructor += consts.TAB + consts.TAB + fmt.Sprintf("this.%s = %s", f.Name, f.Name) + consts.LN
	builder.hydrator += consts.TAB + consts.TAB + fmt.Sprintf("if (%s.%s) this.%s = %s.%s;", HYDRATOR_PARAM_NAME, f.Name, f.Name, HYDRATOR_PARAM_NAME, f.Name) + consts.LN
}

func (builder *JSClassBuilder) WithRelation(ctx context.Context, definition *coredomaindefinition.Relation) {
	if builder.err != nil {
		return
	}

	if definition.Source != builder.definition && definition.Target != builder.definition {
		return
	}

	var to *coredomaindefinition.Model
	if definition.Source == builder.definition {
		to = definition.Target
	} else {
		if definition.IgnoreReverse {
			return
		}
		to = definition.Source
	}

	if IsRelationMultiple(ctx, builder.definition, definition) {
		builder.fields += consts.TAB + fmt.Sprintf("%s // []%s", PluralizeName(ctx, to.Name), GetModelName(ctx, to)) + consts.LN

		builder.hydrator += consts.TAB + consts.TAB + fmt.Sprintf(
			"if (%s.%s && Array.isArray(%s.%s)) {",
			HYDRATOR_PARAM_NAME, PluralizeName(ctx, to.Name), HYDRATOR_PARAM_NAME, PluralizeName(ctx, to.Name),
		) + consts.LN
		builder.hydrator += consts.TAB + consts.TAB + consts.TAB + fmt.Sprintf(
			"this.%s = %s.%s.map((elem) =>  %s.%s(elem));",
			PluralizeName(ctx, to.Name), HYDRATOR_PARAM_NAME, PluralizeName(ctx, to.Name), GetModelName(ctx, to), JS_FROM_METHOD_NAME,
		) + consts.LN
		builder.hydrator += consts.TAB + consts.TAB + "}" + consts.LN
	} else {
		builder.fields += consts.TAB + fmt.Sprintf("%s // %s", to.Name, GetModelName(ctx, to)) + consts.LN
		idRelationFieldName := stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, to))
		builder.fields += consts.TAB + fmt.Sprintf("%s // string", idRelationFieldName) + consts.LN
		optional, err := IsRelationOptionnal(ctx, builder.definition, definition)
		if err != nil {
			builder.err = err
			return
		}

		if !optional {
			builder.constructorParams += fmt.Sprintf("%s,", idRelationFieldName)
			builder.constructor += consts.TAB + consts.TAB + fmt.Sprintf("this.%s = %s", idRelationFieldName, idRelationFieldName) + consts.LN
		}

		builder.hydrator += consts.TAB + consts.TAB + fmt.Sprintf("if (%s.%s) this.%s = %s.%s;", HYDRATOR_PARAM_NAME, idRelationFieldName, idRelationFieldName, HYDRATOR_PARAM_NAME, idRelationFieldName) + consts.LN
	}
}

func (builder *JSClassBuilder) Build(ctx context.Context) error {
	if builder.err != nil {
		return builder.err
	}

	content := fmt.Sprintf("export class %s {", GetModelName(ctx, builder.definition)) + consts.LN

	// start class content
	content += builder.fields + consts.LN
	content += consts.TAB + fmt.Sprintf("constructor(%s) {", strings.TrimSuffix(builder.constructorParams, ",")) + consts.LN
	content += builder.constructor
	content += consts.TAB + "}" + consts.LN
	content += consts.LN

	content += consts.TAB + fmt.Sprintf("%s(%s) {", JS_HYDRATE_METHOD_NAME, HYDRATOR_PARAM_NAME) + consts.LN
	content += builder.hydrator
	content += consts.TAB + consts.TAB + "return this" + consts.LN
	content += consts.TAB + "}" + consts.LN
	content += consts.LN

	content += consts.TAB + fmt.Sprintf("static %s(data) {", JS_FROM_METHOD_NAME) + consts.LN
	content += consts.TAB + consts.TAB + fmt.Sprintf("return new %s().%s(data)", GetModelName(ctx, builder.definition), JS_HYDRATE_METHOD_NAME) + consts.LN
	content += consts.TAB + "}" + consts.LN

	// end class content
	content += "}" + consts.LN
	content += consts.LN

	if builder.domainBuilder.Domain.JSFiles == nil {
		builder.domainBuilder.Domain.JSFiles = map[string]string{}
	}
	if builder.domainBuilder.Domain.JSFiles["entities"] == "" {
		builder.domainBuilder.Domain.JSFiles["entities"] = ""
	}
	builder.domainBuilder.Domain.JSFiles["entities"] += content + consts.LN

	return nil
}
