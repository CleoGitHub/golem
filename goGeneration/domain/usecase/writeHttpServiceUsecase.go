package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/internal/stringifier"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (g *GenerationUsecaseImpl) WriteHttpServiceUsecase(ctx context.Context, domain *model.Domain, httpService *model.Struct, path string) error {
	// if file path does not exist, create it
	filepath := path + "/" + domain.Architecture.SdkPkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	f, err := os.Create(filepath + "/" + stringtool.LowerFirstLetter(httpService.Name) + ".go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	pkgManager := &gopkgmanager.GoPkgManager{
		Pkg: domain.Architecture.SdkPkg.ShortName,
	}

	str := ""
	s, err := stringifier.StringifyStructUsecase(ctx, pkgManager, httpService)
	if err != nil {
		return merror.Stack(err)
	}
	str += s

	str += fmt.Sprintf("var _ %s = &%s{}", domain.Service.Name, httpService.Name)

	str = pkgManager.ToString() + consts.LN + str
	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	return nil
}
