package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/internal/stringifier"
	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/merror"
	"github.com/cleogithub/golem/pkg/stringtool"
)

func (g *GenerationUsecaseImpl) WriteGormDomainRepositoryUsecase(ctx context.Context, domain *model.Domain, domainRepository *model.Struct, path string) error {
	// if file path does not exist, create it
	filepath := path + "/" + domain.Architecture.GormAdapterPkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	f, err := os.Create(filepath + "/" + stringtool.LowerFirstLetter(domainRepository.Name) + ".go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	pkgManager := &gopkgmanager.GoPkgManager{
		Pkg: domain.Architecture.GormAdapterPkg.ShortName,
	}

	str, err := stringifier.StringifyStructUsecase(ctx, pkgManager, domainRepository)
	if err != nil {
		return merror.Stack(err)
	}

	str += fmt.Sprintf("var _ %s.%s = &%s{}", domain.Architecture.RepositoryPkg.Alias, domain.DomainRepository.Name, domainRepository.Name)
	pkgManager.ImportPkg(domain.Architecture.RepositoryPkg)

	str = pkgManager.ToString() + consts.LN + str
	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	return nil
}
