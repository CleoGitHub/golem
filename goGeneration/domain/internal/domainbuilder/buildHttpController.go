package domainbuilder

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/stringtool"
)

func (b *domainBuilder) buildHttpController(ctx context.Context) *domainBuilder {
	if b.Err != nil {
		return b
	}

	controller := &model.Struct{
		Name: fmt.Sprintf("%sHttpController", stringtool.UpperFirstLetter(b.Domain.Name)),
		Fields: []*model.Field{
			{
				Name: "Usecases",
				Type: &model.PkgReference{
					Pkg:       b.Domain.Architecture.UsecasePkg,
					Reference: b.GetDomainUsecase(ctx),
				},
			},
			{
				Name: "Validator",
				Type: &model.PkgReference{
					Pkg:       b.Domain.Architecture.UsecasePkg,
					Reference: b.GetValidator(ctx),
				},
			},
		},
	}

	for _, usecase := range b.Domain.Usecases {
		writer := &model.Param{
			Name: "w",
			Type: &model.PkgReference{
				Pkg: consts.CommonPkgs["http"],
				Reference: &model.ExternalType{
					Type: "ResponseWriter",
				},
			},
		}
		request := &model.Param{
			Name: "r",
			Type: &model.PointerType{
				Type: &model.PkgReference{
					Pkg: consts.CommonPkgs["http"],
					Reference: &model.ExternalType{
						Type: "Request",
					},
				},
			},
		}
		f := &model.Function{
			Name: strings.TrimSuffix(usecase.Function.Name, "Usecase"),
			Args: []*model.Param{writer, request},
		}
		f.Content = func() (string, []*model.GoPkg) {
			pkgs := []*model.GoPkg{
				b.Domain.Architecture.UsecasePkg,
				b.Domain.Architecture.RepositoryPkg,
				consts.CommonPkgs["json"],
				consts.CommonPkgs["http"],
				consts.CommonPkgs["errors"],
			}

			s := struct {
				Request        string
				Writer         string
				UsecasePkg     string
				UsecaseStruct  string
				Usecases       string
				Usecase        string
				UsecaseRequest string
				UsecaseResult  string
				Validator      string
				Files          []string
				FileFields     []string
				JsonPkg        string
				HttpPkg        string
				WritePkg       string
				IOPkg          string
				ErrorsPkg      string
				RepositoryPkg  string
				ErrNotFound    string
				Controller     string
			}{
				Request:        request.Name,
				Writer:         writer.Name,
				UsecasePkg:     b.Domain.Architecture.UsecasePkg.Alias,
				UsecaseStruct:  usecase.Request.Name,
				Usecase:        usecase.Function.Name,
				Usecases:       controller.Fields[0].Name,
				Files:          []string{},
				FileFields:     []string{},
				UsecaseRequest: "request",
				UsecaseResult:  "result",
				Validator:      controller.Fields[1].Name,
				JsonPkg:        consts.CommonPkgs["json"].Alias,
				HttpPkg:        consts.CommonPkgs["http"].Alias,
				WritePkg:       consts.CommonPkgs["http"].Alias,
				IOPkg:          consts.CommonPkgs["io"].Alias,
				ErrorsPkg:      consts.CommonPkgs["errors"].Alias,
				RepositoryPkg:  b.Domain.Architecture.RepositoryPkg.Alias,
				ErrNotFound:    b.RepositoryErrors["notFound"].Name,
				Controller:     controller.GetMethodName(),
			}

			for _, field := range usecase.Request.Fields {
				if field.Type == model.PrimitiveTypeBytes && b.FieldToParamUsecaseRequest[field].Type == coredomaindefinition.PrimitiveTypeFile {
					pkgs = append(pkgs, consts.CommonPkgs["io"])
					s.Files = append(s.Files, field.JsonName)
					s.FileFields = append(s.FileFields, field.Name)
				}
			}

			stringWriter := bytes.NewBufferString("")
			Content := `
			ctx := {{ .Request }}.Context()
				{{ .UsecaseRequest }} := &{{.UsecasePkg}}.{{.UsecaseStruct}}{}
				err := {{ .JsonPkg }}.NewDecoder({{ .Request }}.Body).Decode({{ .UsecaseRequest }})
				if err != nil {
				{{ .Writer }}.WriteHeader({{ .HttpPkg }}.StatusBadRequest)
				{{ .Writer }}.Write([]byte(err.Error()))
				return
				}

				{{ if .Files}} {{ .Request }}.ParseMultipartForm(10 << 20) {{ end }}

				{{ range $idx, $file := .Files }}
				file{{$idx}}, _, err := {{ $.Request }}.FormFile("{{ $file }}")
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				defer file{{$idx}}.Close()

				fileBytes{{$idx}}, err := {{ $.IOPkg }}.ReadAll(file{{$idx}})
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				{{ $.UsecaseRequest }}.{{ index $.FileFields $idx }} = fileBytes{{$idx}}
				{{ end }}

				{{ .UsecaseResult }}, err := {{ .Controller }}.{{ .Usecases }}.{{ .Usecase }}(ctx, {{ .UsecaseRequest }})
				if err != nil {
					if {{ .Controller }}.{{ .Validator }}.IsValidationError(ctx, err) {
						{{ .Writer }}.WriteHeader({{ .HttpPkg }}.StatusUnprocessableEntity)
						{{ .Writer }}.Write([]byte(err.Error()))
						return
					} else if {{ .ErrorsPkg }}.Is({{ .RepositoryPkg }}.{{ .ErrNotFound }}, err) {
						{{ .Writer }}.WriteHeader({{ .HttpPkg }}.StatusNotFound)
						return
					} else {
						{{ .Writer }}.WriteHeader({{ .HttpPkg }}.StatusInternalServerError)
						return
					}
				}

				{{ .JsonPkg }}.NewEncoder({{ .Writer }}).Encode({{ .UsecaseResult }})
			`
			tpl := template.Must(template.New("tpl").Parse(Content))
			err := tpl.Execute(stringWriter, s)
			if err != nil {
				panic(err)
			}
			str := stringWriter.String()

			return str, pkgs
		}
		controller.Methods = append(controller.Methods, f)
	}

	router := &model.Param{
		Name: "rtr",
		Type: &model.PkgReference{
			Pkg: consts.CommonPkgs["router"],
			Reference: &model.ExternalType{
				Type: "Router",
			},
		},
	}
	f := &model.Function{
		Name: "RegisterRoutes",
		Args: []*model.Param{
			{
				Name: "ctx",
				Type: &model.PkgReference{
					Pkg: consts.CommonPkgs["context"],
					Reference: &model.ExternalType{
						Type: "Context",
					},
				},
			},
			router,
		},
	}
	f.Content = func() (string, []*model.GoPkg) {
		pkgs := []*model.GoPkg{
			consts.CommonPkgs["router"],
		}

		str := fmt.Sprintf(`routes := []%s.Route{`, consts.CommonPkgs["router"].Alias) + consts.LN
		for _, usecase := range b.Domain.Usecases {
			str += "{" + consts.LN
			str += fmt.Sprintf("Type: %s.Post,", consts.CommonPkgs["router"].Alias) + consts.LN
			str += fmt.Sprintf(`Pattern: "/%s/%s",`, b.Domain.Name, b.GetHttpRoute(ctx, usecase)) + consts.LN
			str += fmt.Sprintf(`Handler: %s.%s,`, controller.GetMethodName(), b.GetHttpRoute(ctx, usecase)) + consts.LN
			if usecase.Roles != nil && len(usecase.Roles) > 0 {
				str += fmt.Sprintf(`Roles: []string{"%s"},`, strings.Join(usecase.Roles, `","`)) + consts.LN
			}
			str += "}," + consts.LN
		}
		str += "}" + consts.LN

		str += fmt.Sprintf(`%s.AddRoutes(ctx, routes...)`, router.Name)

		return str, pkgs
	}
	controller.Methods = append(controller.Methods, f)

	b.Domain.Controllers = append(b.Domain.Controllers, controller)

	return b
}

func (b *domainBuilder) GetHttpRoute(ctx context.Context, usecase *model.Usecase) string {
	return strings.TrimSuffix(usecase.Function.Name, "Usecase")
}
