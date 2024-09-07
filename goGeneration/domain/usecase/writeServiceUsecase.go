package usecase

import (
	"context"
	"os"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/internal/stringifier"
	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/merror"
	"github.com/cleogithub/golem/pkg/stringtool"
)

func (g *GenerationUsecaseImpl) WriteServiceUsecase(ctx context.Context, domain *model.Domain, service *model.Interface, path string) error {
	// if file path does not exist, create it
	filepath := path + "/" + domain.Architecture.SdkPkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	f, err := os.Create(filepath + "/" + stringtool.LowerFirstLetter(service.Name) + ".go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	pkgManager := &gopkgmanager.GoPkgManager{
		Pkg: domain.Architecture.SdkPkg.ShortName,
	}

	str := ""
	s, err := stringifier.StringifyInterfaceUsecase(ctx, pkgManager, service)
	if err != nil {
		return merror.Stack(err)
	}
	str += s

	str = pkgManager.ToString() + consts.LN + str
	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	// if file path does not exist, create it
	filepath = path + "/" + domain.Architecture.SdkPkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	f, err = os.Create(filepath + "/structs.go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	pkgManager = &gopkgmanager.GoPkgManager{
		Pkg: domain.Architecture.SdkPkg.ShortName,
	}

	str = ""
	for _, method := range service.Methods {
		s, err := stringifier.StringifyStructUsecase(ctx, pkgManager, method.Args[1].Type.(*model.PointerType).Type.(*model.Struct))
		if err != nil {
			return merror.Stack(err)
		}
		str += s

		s, err = stringifier.StringifyStructUsecase(ctx, pkgManager, method.Results[0].Type.(*model.PointerType).Type.(*model.Struct))
		if err != nil {
			return merror.Stack(err)
		}
		str += s
	}

	str = pkgManager.ToString() + consts.LN + str
	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	return nil
}
