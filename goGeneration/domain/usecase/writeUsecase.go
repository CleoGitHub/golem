package usecase

import (
	"context"
	"os"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/internal/stringifier"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (g *GenerationUsecaseImpl) Write(ctx context.Context, inPkg *model.GoPkg, elem interface{}, path string) (err error) {
	// if file path does not exist, create it
	filepath := path + "/" + inPkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	pkgManager := &gopkgmanager.GoPkgManager{
		Pkg: inPkg.ShortName,
	}
	str := ""
	name := ""

	switch t := elem.(type) {
	case *model.Struct:
		str, err = stringifier.StringifyStructUsecase(ctx, pkgManager, t)
		if err != nil {
			return merror.Stack(err)
		}

		name = t.Name
	case *model.Interface:
		str, err = stringifier.StringifyInterfaceUsecase(ctx, pkgManager, t)
		if err != nil {
			return merror.Stack(err)
		}

		name = t.Name
	case *model.File:
		str, err = stringifier.StringifyPortUsecase(ctx, pkgManager, t)
		if err != nil {
			return merror.Stack(err)
		}

		name = t.Name
	default:
		return merror.Stack(ErrUnknowTypeToWrite)
	}

	f, err := os.Create(filepath + "/" + stringtool.LowerFirstLetter(name) + ".golem.go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	str = consts.HEADER + consts.LN + pkgManager.ToString() + consts.LN + str
	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	return nil
}
