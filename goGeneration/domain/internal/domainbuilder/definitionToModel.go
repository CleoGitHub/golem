package domainbuilder

import (
	"context"

	"github.com/cleoGitHub/golem/coredomaindefinition"
	"github.com/cleoGitHub/golem/goGeneration/domain/consts"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
)

func TypeDefinitionToType(ctx context.Context, typeDefinition coredomaindefinition.Type) (model.Type, error) {
	switch typeDefinition.GetType() {
	case coredomaindefinition.PrimitiveTypeBool.GetType():
		return model.PrimitiveTypeBool, nil
	case coredomaindefinition.PrimitiveTypeInt.GetType():
		return model.PrimitiveTypeInt, nil
	case coredomaindefinition.PrimitiveTypeFloat.GetType():
		return model.PrimitiveTypeFloat, nil
	case coredomaindefinition.PrimitiveTypeString.GetType():
		return model.PrimitiveTypeString, nil
	case coredomaindefinition.PrimitiveTypeByte.GetType():
		return model.PrimitiveTypeByte, nil
	case coredomaindefinition.PrimitiveTypeBytes.GetType(), coredomaindefinition.PrimitiveTypeFile.GetType():
		return model.PrimitiveTypeBytes, nil
	case coredomaindefinition.PrimitiveTypeDate.GetType(),
		coredomaindefinition.PrimitiveTypeDateTime.GetType(),
		coredomaindefinition.PrimitiveTypeTime.GetType():
		return &model.PkgReference{
			Pkg: consts.CommonPkgs["time"],
			Reference: &model.ExternalType{
				Type: "Time",
			},
		}, nil
	default:
		return nil, NewErrUnknownType(typeDefinition.GetType())
	}
}

func FieldDefinitionToField(ctx context.Context, fieldDefinition *coredomaindefinition.Field) (*model.Field, error) {
	t, err := TypeDefinitionToType(ctx, fieldDefinition.Type)
	if err != nil {
		return nil, err
	}

	return &model.Field{
		Name: GetFieldName(ctx, fieldDefinition),
		Type: t,
		Tags: []*model.Tag{
			{
				Name:   "json",
				Values: []string{fieldDefinition.Name},
			},
		},
	}, nil
}
