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

func (g *GenerationUsecaseImpl) WriteUsecaseValidatorInterfaceUsecase(ctx context.Context, domain *model.Domain, validator *model.Interface, path string) error {
	// if file path does not exist, create it
	filepath := path + "/" + domain.Architecture.UsecasePkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	f, err := os.Create(filepath + "/" + stringtool.LowerFirstLetter(validator.Name) + ".go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	pkgManager := &gopkgmanager.GoPkgManager{
		Pkg: domain.Architecture.UsecasePkg.ShortName,
	}

	str, err := stringifier.StringifyInterfaceUsecase(ctx, pkgManager, validator)
	if err != nil {
		return merror.Stack(err)
	}

	str = pkgManager.ToString() + consts.LN + str
	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	return nil
}
