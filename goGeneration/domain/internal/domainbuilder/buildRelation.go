package domainbuilder

import (
	"context"
	"slices"

	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/merror"
)

func (b *domainBuilder) buildRelation(ctx context.Context, relationDefinition *coredomaindefinition.Relation) *domainBuilder {
	if b.Err != nil {
		return b
	}

	source, err := b.GetModel(ctx, relationDefinition.Source)
	if err != nil {
		b.Err = merror.Stack(err)
		return b
	}

	target, err := b.GetModel(ctx, relationDefinition.Target)
	if err != nil {
		b.Err = merror.Stack(err)
		return b
	}

	var sourceToTargetType model.RelationType
	switch relationDefinition.Type {
	case coredomaindefinition.RelationTypeBelongsTo, coredomaindefinition.RelationTypeSubresourcesOf:
		sourceToTargetType = model.RelationSingleMandatory
	case coredomaindefinition.RelationTypeManyToOne, coredomaindefinition.RelationTypeOneToOne:
		sourceToTargetType = model.RelationSingleOptionnal
	case coredomaindefinition.RelationTypeManyToMany, coredomaindefinition.RelationTypeOneToMany:
		sourceToTargetType = model.RelationMultiple
	}

	var targetToSourceType model.RelationType
	if !relationDefinition.IgnoreReverse {
		switch relationDefinition.Type {
		case coredomaindefinition.RelationTypeOneToMany, coredomaindefinition.RelationTypeOneToOne:
			targetToSourceType = model.RelationSingleOptionnal
		case coredomaindefinition.RelationTypeManyToMany,
			coredomaindefinition.RelationTypeBelongsTo,
			coredomaindefinition.RelationTypeManyToOne,
			coredomaindefinition.RelationTypeSubresourcesOf:
			targetToSourceType = model.RelationMultiple
		}
	}

	if slices.Contains([]model.RelationType{model.RelationSingleMandatory, model.RelationSingleOptionnal}, sourceToTargetType) {
		f := &model.Field{
			Name: target.Struct.Name + "Id",
			Type: model.PrimitiveTypeString,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{target.JsonName + "Id"},
				},
			},
		}

		if sourceToTargetType == model.RelationSingleMandatory {
			f.Tags = append(f.Tags, &model.Tag{
				Name:   "validate",
				Values: []string{"required"},
			})
		}

		b.ModelUsecaseStruct[source].Fields = append(b.ModelUsecaseStruct[source].Fields, f)
	}

	r := &model.Relation{
		On:         target,
		Type:       sourceToTargetType,
		Definition: relationDefinition,
	}
	b.RelationToRelationDefinition[r] = relationDefinition
	source.Relations = append(source.Relations, r)

	if slices.Contains([]model.RelationType{model.RelationSingleMandatory, model.RelationSingleOptionnal}, targetToSourceType) {
		f := &model.Field{
			Name: source.Struct.Name + "Id",
			Type: model.PrimitiveTypeString,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{source.JsonName + "Id"},
				},
			},
		}

		if targetToSourceType == model.RelationSingleMandatory {
			f.Tags = append(f.Tags, &model.Tag{
				Name:   "validate",
				Values: []string{"required"},
			})
		}
		b.ModelUsecaseStruct[target].Fields = append(b.ModelUsecaseStruct[target].Fields, f)
	}
	if !relationDefinition.IgnoreReverse {
		r := &model.Relation{
			On:         source,
			Type:       targetToSourceType,
			Definition: relationDefinition,
		}
		b.RelationToRelationDefinition[r] = relationDefinition
		target.Relations = append(target.Relations, r)
	}

	return b
}
