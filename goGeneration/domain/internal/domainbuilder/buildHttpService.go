package domainbuilder

import (
	"context"
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

			template := `
			body, err := json.Marshal(request)
			if err != nil {
				return nil, err
			}

			req, err := http.NewRequest("GET", "http://example.com", nil)
			{{.}}.Add("Content-Type", "application/json")
			{{.}}.Add("Accept", "application/json")
			return {{.}}.Post(ctx, url, body)
			{{range .}}
			`
			return "", nil
		}
	}

	b.Domain.HttpService = httpService

	return b
}
