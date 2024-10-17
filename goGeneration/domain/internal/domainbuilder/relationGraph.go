package domainbuilder

import (
	"context"
	"fmt"
	"slices"

	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

type RelationNode struct {
	Model *coredomaindefinition.Model
	Links []*RelationNodeLink
}

func (n *RelationNode) RequireRetriveInactive() bool {
	if n.Model.Activable {
		return true
	}

	for _, link := range n.Links {
		if link.Type == RelationNodeLinkType_DEPEND && link.To.RequireRetriveInactive() {
			return true
		}
	}

	return false
}

type RelationNodeLink struct {
	To   *RelationNode
	Type RelationNodeLinkType
	Name string
}

type RelationNodeLinkType int

const (
	RelationNodeLinkType_DEPEND RelationNodeLinkType = iota
	RelationNodeLinkType_MANY
	RelationNodeLinkType_ONE
)

type RelationGraph []*RelationNode

func (g *RelationGraph) addRelationToGraph(ctx context.Context, relationDefinition *coredomaindefinition.Relation) {
	graph := *g

	var sourceNode *RelationNode
	var targetNode *RelationNode

	for _, n := range graph {
		if n.Model == relationDefinition.Source {
			sourceNode = n
		}
		if n.Model == relationDefinition.Target {
			targetNode = n
		}
		if sourceNode != nil && targetNode != nil {
			break
		}
	}
	if sourceNode == nil {
		sourceNode = &RelationNode{
			Model: relationDefinition.Source,
		}
	}
	if targetNode == nil {
		targetNode = &RelationNode{
			Model: relationDefinition.Target,
		}
	}

	if slices.Contains(
		[]coredomaindefinition.RelationType{coredomaindefinition.RelationTypeBelongsTo, coredomaindefinition.RelationTypeSubresourcesOf},
		relationDefinition.Type,
	) {
		sourceNode.Links = append(sourceNode.Links, &RelationNodeLink{
			To:   targetNode,
			Type: RelationNodeLinkType_DEPEND,
		})
		targetNode.Links = append(targetNode.Links, &RelationNodeLink{
			To:   sourceNode,
			Type: RelationNodeLinkType_MANY,
		})
	} else {
		if IsRelationMultiple(ctx, sourceNode.Model, relationDefinition) {
			sourceNode.Links = append(sourceNode.Links, &RelationNodeLink{
				To:   targetNode,
				Type: RelationNodeLinkType_MANY,
			})
		} else {
			sourceNode.Links = append(sourceNode.Links, &RelationNodeLink{
				To:   targetNode,
				Type: RelationNodeLinkType_ONE,
			})
		}
		if IsRelationMultiple(ctx, targetNode.Model, relationDefinition) {
			targetNode.Links = append(targetNode.Links, &RelationNodeLink{
				To:   sourceNode,
				Type: RelationNodeLinkType_MANY,
			})
		} else {
			targetNode.Links = append(targetNode.Links, &RelationNodeLink{
				To:   sourceNode,
				Type: RelationNodeLinkType_ONE,
			})
		}
	}
	if !slices.Contains(graph, sourceNode) {
		graph = append(graph, sourceNode)
	}
	if !slices.Contains(graph, targetNode) {
		graph = append(graph, targetNode)
	}

	*g = graph
}

func (g *RelationGraph) GetNode(model *coredomaindefinition.Model) *RelationNode {
	for _, n := range *g {
		if n.Model == model {
			return n
		}
	}
	return nil
}

func (builder *domainBuilder) buildRelationGraph(ctx context.Context) (*model.File, error) {
	if builder.err != nil {
		return nil, builder.err
	}

	file := &model.File{
		Name:     "relationGraph",
		Pkg:      builder.GetRepositoryPackage(),
		Elements: []interface{}{},
	}

	for _, node := range *builder.RelationGraph {
		relation := &model.Struct{
			Name: GetRepositoryRelationNodeName(ctx, node.Model),
		}

		for _, link := range node.Links {
			relation.Methods = append(relation.Methods, &model.Function{
				Name: GetPreloadName(ctx, link.To.Model),
				Results: []*model.Param{
					{
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Pkg: builder.GetRepositoryPackage(),
								Reference: &model.ExternalType{
									Type: GetRepositoryRelationNodeName(ctx, link.To.Model),
								},
							},
						},
					},
				},
				Content: func() (string, []*model.GoPkg) {
					return fmt.Sprintf("return &%s{}", GetRepositoryRelationNodeName(ctx, link.To.Model)), []*model.GoPkg{
						builder.GetRepositoryPackage(),
					}
				},
			})
		}
		file.Elements = append(file.Elements, relation)
	}
	return file, nil
}
