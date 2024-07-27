package usecase

import (
	"context"
	"os"

	"github.com/cleoGitHub/golem/goGeneration/domain/consts"
	"github.com/cleoGitHub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleoGitHub/golem/goGeneration/domain/internal/stringifier"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
	"github.com/cleoGitHub/golem/pkg/merror"
)

func (g *GenerationUsecaseImpl) WriteRepositoryErrorsUsecase(ctx context.Context, domain *model.Domain, path string) error {
	// if file path does not exist, create it
	filepath := path + "/" + domain.Architecture.RepositoryPkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	f, err := os.Create(filepath + "/errors.go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	pkgManager := &gopkgmanager.GoPkgManager{
		Pkg: domain.Architecture.RepositoryPkg.ShortName,
	}

	str := ""

	for _, err := range domain.RepositoryErrors {
		s, err := stringifier.StringifyConstUsecase(ctx, pkgManager, err)
		if err != nil {
			return merror.Stack(err)
		}

		str = s + consts.LN
	}

	str = pkgManager.ToString() + consts.LN + str
	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	return nil
}
