package usecase

import (
	"bytes"
	"context"
	"os"
	"text/template"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (g *GenerationUsecaseImpl) WriteTemplateUsecase(ctx context.Context, tmpl *model.Template, path string) (err error) {
	// if file path does not exist, create it
	filepath := path + "/" + tmpl.Pkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	tmpl.Data.SetPkg(tmpl.Pkg)
	temp, err := template.New(tmpl.Name).Parse(tmpl.StringTemplate)
	if err != nil {
		return merror.Stack(err)
	}

	buffer := bytes.NewBufferString("")
	err = temp.Execute(buffer, tmpl.Data)
	if err != nil {
		return merror.Stack(err)
	}

	f, err := os.Create(filepath + "/" + stringtool.CamelCase(tmpl.Name) + ".golem.go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	str := consts.HEADER + consts.LN + buffer.String()
	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	return nil
}
