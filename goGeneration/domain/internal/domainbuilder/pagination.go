package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (b *domainBuilder) GetPagination(ctx context.Context) *model.Struct {
	if b.Err != nil {
		return nil
	}
	if b.Pagination != nil {
		return b.Pagination
	}
	pageField := &model.Field{
		Name: "Page",
		Type: model.PrimitiveTypeInt,
		Tags: []*model.Tag{
			{
				Name:   "json",
				Values: []string{"page"},
			},
		},
	}
	itemsPerPageField := &model.Field{
		Name: "ItemsPerPage",
		Type: model.PrimitiveTypeInt,
		Tags: []*model.Tag{
			{
				Name:   "json",
				Values: []string{"itemsPerPage"},
			},
		},
	}

	minItemPerPageConst := &model.Var{
		Name:    "MIN_ITEMS_PER_PAGE",
		Type:    model.PrimitiveTypeInt,
		Value:   5,
		IsConst: true,
	}
	maxItemPerPageConst := &model.Var{
		Name:    "MAX_ITEMS_PER_PAGE",
		Type:    model.PrimitiveTypeInt,
		Value:   100,
		IsConst: true,
	}
	defaultItemPerPageConst := &model.Var{
		Name:    "DEFAULT_ITEMS_PER_PAGE",
		Type:    model.PrimitiveTypeInt,
		Value:   30,
		IsConst: true,
	}
	pagination := &model.Struct{
		Name:   "Pagination",
		Consts: []*model.Var{minItemPerPageConst, maxItemPerPageConst, defaultItemPerPageConst},
		Fields: []*model.Field{pageField, itemsPerPageField},
	}

	pagination.Methods = append(pagination.Methods, &model.Function{
		Name: "GetPage",
		Results: []*model.Param{
			{
				Type: model.PrimitiveTypeInt,
			},
		},
		Content: func() (string, []*model.GoPkg) {
			str := ""
			str += fmt.Sprintf("if %s.%s < 1 { return 1 }", pagination.GetMethodName(), pageField.Name) + consts.LN
			str += fmt.Sprintf("return %s.%s", pagination.GetMethodName(), pageField.Name) + consts.LN
			return str, nil
		},
	})

	pagination.Methods = append(pagination.Methods, &model.Function{
		Name: "GetItemsPerPage",
		Results: []*model.Param{
			{
				Type: model.PrimitiveTypeInt,
			},
		},
		Content: func() (string, []*model.GoPkg) {
			str := ""
			str += fmt.Sprintf(
				"if %[1]s.%[2]s < %s || %[1]s.%[2]s > %[4]s { return %s }",
				pagination.GetMethodName(), itemsPerPageField.Name, minItemPerPageConst.Name, maxItemPerPageConst.Name, defaultItemPerPageConst.Name,
			) + consts.LN
			str += fmt.Sprintf("return %s.%s ", pagination.GetMethodName(), itemsPerPageField.Name) + consts.LN
			return str, nil
		},
	})

	b.Pagination = pagination
	b.Domain.Pagination = pagination

	return b.Pagination
}
