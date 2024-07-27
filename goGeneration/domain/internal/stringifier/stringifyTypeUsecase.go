package stringifier

import (
	"context"

	"github.com/cleoGitHub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
)

func StringifyTypeUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, t model.Type) (string, error) {
	str := t.GetType(model.InPkg(pkgManager.Pkg))
	types := []model.Type{t}
	types = append(types, t.SubTypes()...)
	for _, subType := range types {
		if pkgReference, ok := subType.(*model.PkgReference); ok && pkgManager.Pkg != pkgReference.Pkg.Alias {
			pkgManager.ImportPkg(pkgReference.Pkg)
		}
	}
	return str, nil
}
