package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (b *domainBuilder) TypeDefinitionToType(ctx context.Context, typeDefinition coredomaindefinition.Type) (model.Type, error) {
	switch ty := typeDefinition.(type) {
	case coredomaindefinition.PrimitiveType:
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

	case *coredomaindefinition.Array:
		t, err := b.TypeDefinitionToType(ctx, ty.Type)
		if err != nil {
			return nil, err
		}
		return &model.ArrayType{
			Type: t,
		}, nil
	case *coredomaindefinition.Model:
		return &model.PointerType{
			Type: &model.PkgReference{
				Pkg:       b.Domain.Architecture.ModelPkg,
				Reference: &model.ExternalType{Type: stringtool.UpperFirstLetter(ty.Name)},
			},
		}, nil
	default:
		return nil, NewErrUnknownType(typeDefinition.GetType())
	}
}

func (b *domainBuilder) FieldDefinitionToField(ctx context.Context, fieldDefinition *coredomaindefinition.Field) (*model.Field, error) {
	t, err := b.TypeDefinitionToType(ctx, fieldDefinition.Type)
	if err != nil {
		return nil, err
	}

	return &model.Field{
		Name: GetFieldName(ctx, fieldDefinition.Name),
		Type: t,
		Tags: []*model.Tag{
			{
				Name:   "json",
				Values: []string{fieldDefinition.Name},
			},
		},
		JsonName: fieldDefinition.Name,
	}, nil
}
