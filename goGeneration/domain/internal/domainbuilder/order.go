package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleoGitHub/golem/goGeneration/domain/consts"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
)

func (b *domainBuilder) GetOrdering(ctx context.Context) *model.Struct {
	if b.Err != nil {
		return nil
	}
	if b.Ordering != nil {
		return b.Ordering
	}
	orderField := &model.Field{
		Name: "Order",
		Type: model.PrimitiveTypeString,
		Tags: []*model.Tag{
			{
				Name:   "json",
				Values: []string{"order"},
			},
		},
	}
	orderByField := &model.Field{
		Name: "OrderBy",
		Type: model.PrimitiveTypeString,
		Tags: []*model.Tag{
			{
				Name:   "json",
				Values: []string{"orderBy"},
			},
		},
	}

	ASC := &model.Var{
		Name:    "ASC",
		Type:    model.PrimitiveTypeString,
		Value:   "ASC",
		IsConst: true,
	}
	DESC := &model.Var{
		Name:    "DESC",
		Type:    model.PrimitiveTypeString,
		Value:   "DESC",
		IsConst: true,
	}
	ordering := &model.Struct{
		Name:   "Ordering",
		Consts: []*model.Var{ASC, DESC},
		Fields: []*model.Field{orderField, orderByField},
	}

	ordering.Methods = append(ordering.Methods, &model.Function{
		Name: "GetOrder",
		Results: []*model.Param{
			{
				Type: model.PrimitiveTypeString,
			},
		},
		Content: func() (string, []*model.GoPkg) {
			str := ""
			str += fmt.Sprintf(
				`if strings.ToUpper(%s.%s) != %s && strings.ToUpper(%s.%s) != %s { return %s }`,
				ordering.GetMethodName(), orderField.Name, ASC.Name, ordering.GetMethodName(), orderField.Name, DESC.Name, ASC.Name,
			) + consts.LN
			str += fmt.Sprintf("return  strings.ToUpper(%s.%s)", ordering.GetMethodName(), orderField.Name) + consts.LN
			return str, []*model.GoPkg{consts.CommonPkgs["strings"]}
		},
	})

	allowedOrderBy := &model.Param{
		Name: "allowedOrderBys",
		Type: &model.ArrayType{
			Type: model.PrimitiveTypeString,
		},
	}
	defaultOrderBy := &model.Param{
		Name: "defaultOrderBy",
		Type: model.PrimitiveTypeString,
	}
	ordering.Methods = append(ordering.Methods, &model.Function{
		Name: "GetOrderBy",
		Args: []*model.Param{allowedOrderBy, defaultOrderBy},
		Results: []*model.Param{
			{
				Type: model.PrimitiveTypeString,
			},
		},
		Content: func() (string, []*model.GoPkg) {
			str := ""
			str += fmt.Sprintf(
				`if slices.Contains(%s, %s.%s)  { return  %s.%s }`,
				allowedOrderBy.Name, ordering.GetMethodName(), orderByField.Name, ordering.GetMethodName(), orderByField.Name,
			) + consts.LN
			str += fmt.Sprintf("return  %s", defaultOrderBy.Name) + consts.LN
			return str, []*model.GoPkg{consts.CommonPkgs["slices"]}
		},
	})

	b.Ordering = ordering

	return b.Ordering
}
