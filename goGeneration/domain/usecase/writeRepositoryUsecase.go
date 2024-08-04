package usecase

import (
	"context"
	"os"

	"github.com/cleoGitHub/golem-common/pkg/merror"
	"github.com/cleoGitHub/golem-common/pkg/stringtool"
	"github.com/cleoGitHub/golem/goGeneration/domain/consts"
	"github.com/cleoGitHub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleoGitHub/golem/goGeneration/domain/internal/stringifier"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
)

func (g *GenerationUsecaseImpl) WriteRepositoryUsecase(ctx context.Context, domain *model.Domain, repo *model.Repository, path string) error {
	// if file path does not exist, create it
	filepath := path + "/" + domain.Architecture.RepositoryPkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	f, err := os.Create(filepath + "/" + stringtool.LowerFirstLetter(repo.Name) + ".go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	pkgManager := &gopkgmanager.GoPkgManager{
		Pkg: domain.Architecture.RepositoryPkg.ShortName,
	}

	str, err := stringifier.StringifyRepositoryUsecase(ctx, pkgManager, repo)
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
