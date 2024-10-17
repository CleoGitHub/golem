package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func GetHttpClientName(ctx context.Context, domain *coredomaindefinition.Domain) string {
	return fmt.Sprintf("%sHttpClient", stringtool.UpperFirstLetter(domain.Name))
}

func GetHttpControllerName(ctx context.Context, domain *coredomaindefinition.Domain) string {
	return fmt.Sprintf("%sHttpController", stringtool.UpperFirstLetter(domain.Name))
}

func GetHttpRouteName(ctx context.Context, domain *coredomaindefinition.Domain, method string) string {
	return fmt.Sprintf("/%s/%s", domain.Name, stringtool.DashCase(method))
}

func GetHttpRoute(ctx context.Context, routeName string) *model.Function {
	return &model.Function{
		Name: routeName,
		Args: []*model.Param{
			{
				Name: "w",
				Type: &model.PkgReference{
					Pkg: consts.CommonPkgs["http"],
					Reference: &model.ExternalType{
						Type: "ResponseWriter",
					},
				},
			},
			{
				Name: "r",
				Type: &model.PkgReference{
					Pkg: consts.CommonPkgs["http"],
					Reference: &model.ExternalType{
						Type: "Request",
					},
				},
			},
		},
		Results: []*model.Param{},
	}
}
