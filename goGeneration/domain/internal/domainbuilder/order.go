package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	ORDERING_NAME        = "Ordering"
	ORDERING_ORDER       = "Order"
	ORDERING_ORDERBY     = "OrderBy"
	ORDERING_GET_ORDER   = "GetOrder"
	ORDERING_GET_ORDERBY = "GetOrderBy"
)

var ASC = &model.Var{
	Name:    "ASC",
	Type:    model.PrimitiveTypeString,
	Value:   "ASC",
	IsConst: true,
}
var DESC = &model.Var{
	Name:    "DESC",
	Type:    model.PrimitiveTypeString,
	Value:   "DESC",
	IsConst: true,
}

var ORDERING = &model.Struct{
	Name:       ORDERING_NAME,
	MethodName: stringtool.LowerFirstLetter(ORDERING_NAME),
	Fields: []*model.Field{
		{
			Name: ORDERING_ORDER,
			Type: model.PrimitiveTypeString,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{stringtool.LowerFirstLetter(ORDERING_ORDER)},
				},
			},
		},
		{
			Name: ORDERING_ORDERBY,
			Type: model.PrimitiveTypeString,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{stringtool.LowerFirstLetter(ORDERING_ORDERBY)},
				},
			},
		},
	},
	Methods: []*model.Function{
		{
			Name: "GetOrder",
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeString,
				},
			},
			Content: func() (string, []*model.GoPkg) {
				str := ""
				str += fmt.Sprintf(
					`if strings.ToUpper(%s.%s) != ASC && strings.ToUpper(%s.%s) != DESC { return ASC }`,
					stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDER, stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDER,
				) + consts.LN
				str += fmt.Sprintf("return  strings.ToUpper(%s.%s)", stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDER) + consts.LN
				return str, []*model.GoPkg{consts.CommonPkgs["strings"]}
			},
		}, {
			Name: "GetOrderBy",
			Args: []*model.Param{
				{
					Name: "allowedOrderBys",
					Type: &model.ArrayType{
						Type: model.PrimitiveTypeString,
					},
				}, {
					Name: "defaultOrderBy",
					Type: model.PrimitiveTypeString,
				},
			},
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeString,
				},
			},
			Content: func() (string, []*model.GoPkg) {
				str := ""
				str += fmt.Sprintf(
					`if slices.Contains(allowedOrderBys, %s.%s)  { return  %s.%s }`,
					stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDERBY,
					stringtool.LowerFirstLetter(ORDERING_NAME), ORDERING_ORDERBY,
				) + consts.LN
				str += "return  defaultOrderBy" + consts.LN
				return str, []*model.GoPkg{consts.CommonPkgs["slices"]}
			},
		},
	},
}

func (b *domainBuilder) addOrdering(ctx context.Context) *domainBuilder {
	if b.Err != nil {
		return b
	}

	b.Domain.Ports = append(b.Domain.Ports, &model.File{
		Name: ORDERING_NAME,
		Pkg:  b.GetRepositoryPackage(),
		Elements: []interface{}{
			ASC,
			DESC,
			ORDERING,
		},
	})

	return b
}
