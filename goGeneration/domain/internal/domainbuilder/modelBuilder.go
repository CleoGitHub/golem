package domainbuilder

import (
	"context"
	"slices"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

type ModelBuilder struct {
	DomainBuilder *domainBuilder

	Definition *coredomaindefinition.Model

	Model *model.Struct
	Err   error
}

func NewModelBuilder(
	ctx context.Context,
	domainBuilder *domainBuilder,
	definition *coredomaindefinition.Model,
	defaultFields []*coredomaindefinition.Field,
) Builder {
	builder := &ModelBuilder{
		Definition:    definition,
		DomainBuilder: domainBuilder,
		Model: &model.Struct{
			Name:   GetModelName(ctx, definition),
			Fields: []*model.Field{},
		},
		Err: nil,
	}

	modelFieldNames := []string{}
	for _, f := range defaultFields {
		modelFieldNames = append(modelFieldNames, f.Name)
		field, err := builder.DomainBuilder.FieldDefinitionToField(ctx, f)
		if err != nil {
			builder.Err = merror.Stack(err)
			return builder
		}
		builder.Model.Fields = append(builder.Model.Fields, field)
	}
	if definition.Archivable {
		modelFieldNames = append(modelFieldNames, "deleted")
		field, err := builder.DomainBuilder.FieldDefinitionToField(ctx, &coredomaindefinition.Field{
			Name: "deletedAt",
			Type: coredomaindefinition.PrimitiveTypeDateTime,
		})
		if err != nil {
			builder.Err = merror.Stack(err)
			return builder
		}
		builder.Model.Fields = append(builder.Model.Fields, field)
	}

	// Add activable field if model is activable
	if definition.Activable {
		modelFieldNames = append(modelFieldNames, "active")
		field, err := builder.DomainBuilder.FieldDefinitionToField(ctx, &coredomaindefinition.Field{
			Name: "active",
			Type: coredomaindefinition.PrimitiveTypeBool,
		})
		if err != nil {
			builder.Err = merror.Stack(err)
			return builder
		}
		builder.Model.Fields = append(builder.Model.Fields, field)
	}

	// Add default fields to definition
	for _, field := range definition.Fields {
		if slices.Contains(modelFieldNames, field.Name) {
			builder.Err = merror.Stack(NewErrDefaultFiedlRedefined(field.Name))
			return builder
		}
		f, err := builder.DomainBuilder.FieldDefinitionToField(ctx, field)
		if err != nil {
			builder.Err = merror.Stack(err)
			return builder
		}
		builder.Model.Fields = append(builder.Model.Fields, f)
	}

	return builder
}

func (builder *ModelBuilder) WithRelation(ctx context.Context, relationDefinition *coredomaindefinition.Relation) Builder {
	if builder.Err != nil {
		return builder
	}

	var to *coredomaindefinition.Model
	if relationDefinition.Source == builder.Definition {
		to = relationDefinition.Target
	} else if relationDefinition.Target == builder.Definition {
		// reverse not needed
		if relationDefinition.IgnoreReverse {
			return builder
		}

		to = relationDefinition.Source
	} else {
		return builder
	}

	if IsRelationMultiple(ctx, builder.Definition, relationDefinition) {
		builder.Model.Fields = append(builder.Model.Fields, &model.Field{
			Name: GetMultipleRelationName(ctx, to),
			Type: &model.ArrayType{
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: builder.DomainBuilder.GetModelPackage(),
						Reference: &model.ExternalType{
							Type: GetModelName(ctx, to),
						},
					},
				},
			},
			Tags: []*model.Tag{{Name: "json", Values: []string{
				stringtool.LowerFirstLetter(GetMultipleRelationName(ctx, to)),
			}}},
		})
	} else {
		builder.Model.Fields = append(builder.Model.Fields, &model.Field{
			Name: GetSingleRelationName(ctx, to),
			Type: &model.PointerType{
				Type: &model.PkgReference{
					Pkg: builder.DomainBuilder.GetModelPackage(),
					Reference: &model.ExternalType{
						Type: GetModelName(ctx, to),
					},
				},
			},
			Tags: []*model.Tag{{Name: "json", Values: []string{
				stringtool.LowerFirstLetter(GetSingleRelationName(ctx, to)),
			}}},
		})
		builder.Model.Fields = append(builder.Model.Fields, &model.Field{
			Name: GetSingleRelationIdName(ctx, to),
			Type: model.PrimitiveTypeString,
			Tags: []*model.Tag{{Name: "json", Values: []string{
				stringtool.LowerFirstLetter(GetSingleRelationIdName(ctx, to)),
			}}},
		})
	}

	return builder
}

// WithModel implements Builder.
func (builder *ModelBuilder) WithModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) Builder {
	return builder
}

// WithRepository implements Builder.
func (builder *ModelBuilder) WithRepository(ctx context.Context, repositoryDefinition *coredomaindefinition.Repository) Builder {
	return builder
}

func (builder *ModelBuilder) Build(ctx context.Context) error {
	if builder.Err != nil {
		return builder.Err
	}

	builder.DomainBuilder.Domain.ModelsV2 = append(builder.DomainBuilder.Domain.ModelsV2, builder.Model)

	return nil
}
