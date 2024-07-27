package domainbuilder

import (
	"context"
	"fmt"
	"strings"

	"github.com/cleoGitHub/golem/goGeneration/domain/consts"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
	"github.com/cleoGitHub/golem/pkg/stringtool"
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
				`req, err := %s.NewRequest(%s.MethodPost, %s, body, map[string]string{"Content-Type": "application/json"})`,
				httpService.GetMethodName(), consts.CommonPkgs["http"].Alias, httpService.GetMethodName(),
			) + consts.LN
			str += "if err != nil { return nil, err }" + consts.LN
			str += fmt.Sprintf(`resp, err := %s.Do(req)`, httpService.GetMethodName()) + consts.LN
			str += "if err != nil { return nil, err }" + consts.LN

			return str, pkgs
		}
	}

	b.Domain.HttpService = httpService

	return b
}
