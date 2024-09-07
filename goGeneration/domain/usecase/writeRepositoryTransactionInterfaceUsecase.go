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

func (g *GenerationUsecaseImpl) WriteRepositoryTransactionInterfaceUsecase(ctx context.Context, domain *model.Domain, transaction *model.Interface, path string) error {
	// if file path does not exist, create it
	filepath := path + "/" + domain.Architecture.RepositoryPkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	f, err := os.Create(filepath + "/" + stringtool.LowerFirstLetter(transaction.Name) + ".go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	pkgManager := &gopkgmanager.GoPkgManager{
		Pkg: domain.Architecture.RepositoryPkg.ShortName,
	}

	str, err := stringifier.StringifyInterfaceUsecase(ctx, pkgManager, transaction)
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
