package domainbuilder

import (
	"context"
	"fmt"
	"strings"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (b *domainBuilder) buildHttpService(ctx context.Context) *domainBuilder {
	if b.Err != nil {
		return b
	}

	httpService := &model.Struct{
		Name:    stringtool.UpperFirstLetter(b.Domain.Name) + "HttpService",
		Methods: []*model.Function{},
		Fields: []*model.Field{
			{
				Name: "",
				Type: &model.PkgReference{
					Pkg:       consts.CommonPkgs["httpclient"],
					Reference: &model.ExternalType{Type: "HttpClient"},
				},
			},
		},
	}

	for _, usecase := range b.Domain.Usecases {
		m := usecase.Function.Copy().(*model.Function)
		m.Name = strings.TrimSuffix(m.Name, "Usecase")
		m.Args[1].Type = &model.PointerType{
			Type: m.Args[1].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.Copy(),
		}
		m.Results[0].Type = &model.PointerType{
			Type: m.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.Copy(),
		}
		httpService.Methods = append(httpService.Methods, m)
		m.Content = func() (string, []*model.GoPkg) {
			pkgs := []*model.GoPkg{
				consts.CommonPkgs["json"],
				consts.CommonPkgs["http"],
			}

			str := "body, err := json.Marshal(request)" + consts.LN
			str += "if err != nil { return nil, err }" + consts.LN
			str += fmt.Sprintf(
				`req, err := %s.NewRequest(%s.MethodPost, "%s/%s", body, map[string]string{"Content-Type": "application/json"})`,
				httpService.GetMethodName(), consts.CommonPkgs["http"].Alias, b.Domain.Name, b.GetHttpRoute(ctx, usecase),
			) + consts.LN
			str += "if err != nil { return nil, err }" + consts.LN
			str += fmt.Sprintf(`resp, status, err := %s.Do(req)`, httpService.GetMethodName()) + consts.LN
			str += "if err != nil { return nil, err }" + consts.LN
			str += fmt.Sprintf("if status != 200 { return nil, %s.ErrUnexpectedStatus }", consts.CommonPkgs["httpclient"].Alias) + consts.LN
			str += fmt.Sprintf("result := &%s{}", m.Results[0].Type.(*model.PointerType).Type.GetType()) + consts.LN
			str += "if err = json.Unmarshal(resp, result); err != nil { return nil, err }" + consts.LN

			str += "return result, nil"

			return str, pkgs
		}
	}

	b.Domain.HttpService = httpService

	return b
}
