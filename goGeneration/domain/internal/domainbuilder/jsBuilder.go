package domainbuilder

import (
	"context"
	"fmt"
	"strings"

	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
)

type JSBuilder struct {
	EmptyBuilder

	domainBuilder    *domainBuilder
	domainDefinition *coredomaindefinition.Domain

	err error
}

func NewJSBuilder(ctx context.Context, domainDefinition *coredomaindefinition.Domain, domainBuilder *domainBuilder) Builder {
	domainBuilder.AddBuilder(ctx, NewJSServiceBuilder(domainBuilder))
	return &JSBuilder{
		domainBuilder:    domainBuilder,
		domainDefinition: domainDefinition,
	}
}

func (builder *JSBuilder) WithModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) {
	if builder.err != nil {
		return
	}

	builder.domainBuilder.AddBuilder(ctx, NewJSClassBuilder(ctx, builder.domainBuilder, modelDefinition))
}

func (builder *JSBuilder) WithUsecase(ctx context.Context, usecaseDefinition *coredomaindefinition.Usecase) {
	if builder.err != nil {
		return
	}

	typeImports := []string{}

	content := fmt.Sprintf("export class %s {", GetUsecaseRequestName(ctx, usecaseDefinition.Name)) + consts.LN
	for _, arg := range usecaseDefinition.Args {
		typeImports = append(typeImports, fmt.Sprintf("%s: string", GetFieldName(ctx, arg.Name)))
	}
	content += consts.TAB + fmt.Sprintf("constructor(%s) {", strings.Join(typeImports, ", ")) + consts.LN
	content += consts.TAB + consts.TAB + fmt.Sprintf("this.%s = %s", GetUsecaseRequestName(ctx, usecaseDefinition.Name), GetUsecaseRequestName(ctx, usecaseDefinition.Name)) + consts.LN
	content += consts.TAB + "}" + consts.LN
	content += consts.LN

	if builder.domainBuilder.Domain.JSFiles == nil {

	}
}

func (builder *JSBuilder) Build(ctx context.Context) error {
	if builder.err != nil {
		return builder.err
	}

	content := StructToClass(PAGINATION, "")
	content += consts.LN

	content += StructToClass(ORDERING, "")
	content += consts.LN

	if builder.domainBuilder.Domain.JSFiles == nil {
		builder.domainBuilder.Domain.JSFiles = map[string]string{}
	}
	builder.domainBuilder.Domain.JSFiles["utils"] = content

	return nil
}
