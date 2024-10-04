package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	PAGINATION_NAME                     = "Pagination"
	PAGINATION_WHERE                    = "Where"
	PAGINATION_Page                     = "Page"
	PAGINATION_ItemsPerPage             = "ItemsPerPage"
	PAGINATION_GetItemsPerPage          = "GetItemsPerPage"
	PAGINATION_GetPage                  = "GetPage"
	PAGINATION_MIN_ITEMS_PER_PAGE       = "MIN_ITEMS_PER_PAGE"
	PAGINATION_MAX_ITEMS_PER_PAGE       = "MAX_ITEMS_PER_PAGE"
	PAGINATION_DEFAULT_ITEMS_PER_PAGE   = "DEFAULT_ITEMS_PER_PAGE"
	REPOSITORY_WHERE_OPERATOR_EQUAL     = "EQUAL"
	REPOSITORY_WHERE_OPERATOR_NOT_EQUAL = "NOT_EQUAL"
	REPOSITORY_WHERE_OPERATOR_IN        = "IN"
	REPOSITORY_WHERE_OPERATOR_NOT_IN    = "NOT_IN"
	REPOSITORY_WHERE_OPERATOR_TYPE      = "WHERE_OPERATOR"
)

var PAGINATION = &model.Struct{
	Name:       PAGINATION_NAME,
	MethodName: stringtool.LowerFirstLetter(PAGINATION_NAME),
	Consts: []*model.Var{
		{
			Name:    PAGINATION_MIN_ITEMS_PER_PAGE,
			Type:    model.PrimitiveTypeInt,
			Value:   5,
			IsConst: true,
		},
		{
			Name:    PAGINATION_MAX_ITEMS_PER_PAGE,
			Type:    model.PrimitiveTypeInt,
			Value:   100,
			IsConst: true,
		},
		{
			Name:    PAGINATION_DEFAULT_ITEMS_PER_PAGE,
			Type:    model.PrimitiveTypeInt,
			Value:   30,
			IsConst: true,
		},
	},
	Fields: []*model.Field{
		{
			Name: PAGINATION_Page,
			Type: model.PrimitiveTypeInt,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{stringtool.LowerFirstLetter(PAGINATION_Page)},
				},
			},
		}, {
			Name: PAGINATION_ItemsPerPage,
			Type: model.PrimitiveTypeInt,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{stringtool.LowerFirstLetter(PAGINATION_GetItemsPerPage)},
				},
			},
		},
	},
	Methods: []*model.Function{
		{
			Name: PAGINATION_GetPage,
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeInt,
				},
			},
			Content: func() (string, []*model.GoPkg) {
				str := ""
				str += fmt.Sprintf("if %s.%s < 1 { return 1 }", stringtool.LowerFirstLetter(PAGINATION_NAME), PAGINATION_Page) + consts.LN
				str += fmt.Sprintf("return %s.%s", stringtool.LowerFirstLetter(PAGINATION_NAME), PAGINATION_Page) + consts.LN
				return str, nil
			},
		},
		{
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
					stringtool.LowerFirstLetter(PAGINATION_NAME), PAGINATION_Page, PAGINATION_MIN_ITEMS_PER_PAGE, PAGINATION_MIN_ITEMS_PER_PAGE, PAGINATION_DEFAULT_ITEMS_PER_PAGE,
				) + consts.LN
				str += fmt.Sprintf("return %s.%s ", stringtool.LowerFirstLetter(PAGINATION_NAME), PAGINATION_Page) + consts.LN
				return str, nil
			},
		},
	},
}

var WHERE_OPERATOR_TYPE = &model.TypeDefinition{
	Name: REPOSITORY_WHERE_OPERATOR_TYPE,
	Type: model.PrimitiveTypeString,
}

var WHERE_OPERATOR = &model.Enum{
	Name: "Operator",
	Type: WHERE_OPERATOR_TYPE,
	Values: map[string]interface{}{
		REPOSITORY_WHERE_OPERATOR_EQUAL:     REPOSITORY_WHERE_OPERATOR_EQUAL,
		REPOSITORY_WHERE_OPERATOR_NOT_EQUAL: REPOSITORY_WHERE_OPERATOR_NOT_EQUAL,
		REPOSITORY_WHERE_OPERATOR_IN:        REPOSITORY_WHERE_OPERATOR_IN,
		REPOSITORY_WHERE_OPERATOR_NOT_IN:    REPOSITORY_WHERE_OPERATOR_NOT_IN,
	},
}

var WHERE = &model.Struct{
	Name:       PAGINATION_WHERE,
	MethodName: stringtool.LowerFirstLetter(PAGINATION_WHERE),
	Fields: []*model.Field{
		{
			Name: "Key",
			Type: model.PrimitiveTypeString,
		},
		{

			Name: "Operator",
			Type: WHERE_OPERATOR_TYPE,
		},
		{
			Name: "Value",
			Type: model.PrimitiveTypeInterface,
		},
	},
}

func (b *domainBuilder) addPagination(ctx context.Context) *domainBuilder {
	if b.Err != nil {
		return b
	}

	b.Domain.Ports = append(b.Domain.Ports, &model.File{
		Name: PAGINATION_NAME,
		Pkg:  b.GetRepositoryPackage(),
		Elements: []interface{}{
			WHERE_OPERATOR_TYPE,
			WHERE,
			WHERE_OPERATOR,
			PAGINATION,
		},
	})

	return b
}
