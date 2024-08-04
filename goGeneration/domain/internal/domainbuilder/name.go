package domainbuilder

import (
	"context"
	"strings"

	"github.com/cleoGitHub/golem-common/pkg/stringtool"
	"github.com/cleoGitHub/golem/coredomaindefinition"
)

func GetFieldName(ctx context.Context, fieldDefinition *coredomaindefinition.Field) string {
	return stringtool.UpperFirstLetter(fieldDefinition.Name)
}

func GetModelName(ctx context.Context, modelDefinition *coredomaindefinition.Model) string {
	return stringtool.UpperFirstLetter(modelDefinition.Name)
}

func GetSingleRelationName(ctx context.Context, m *coredomaindefinition.Model) string {
	return stringtool.UpperFirstLetter(m.Name) + "Id"
}

func GetMultiRelationName(ctx context.Context, m *coredomaindefinition.Model) string {
	if strings.HasSuffix(m.Name, "y") {
		return stringtool.UpperFirstLetter(m.Name[:len(m.Name)-1]) + "ies"
	}
	return stringtool.UpperFirstLetter(m.Name) + "s"
}
