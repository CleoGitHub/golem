package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/merror"
)

func StringifyFunctionUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, function *model.Function) (string, error) {
	def, err := StringifyFunctionDefinitionUsecase(ctx, pkgManager, function)
	if err != nil {
		return "", merror.Stack(err)
	}

	str := fmt.Sprintf("func %s {", def) + consts.LN
	s, pkgs := function.Content()
	for _, pkg := range pkgs {
		if err := pkgManager.ImportPkg(pkg); err != nil {
			return "", merror.Stack(err)
		}
	}
	str += s + consts.LN
	str += "}"
	return str, nil
}

func StringifyFunctionDefinitionUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, function *model.Function) (string, error) {
	str := ""
	argsStr, err := StringifyParamsUsecase(ctx, pkgManager, function.Args)
	if err != nil {
		return "", merror.Stack(err)
	}

	resultsStr, err := StringifyParamsUsecase(ctx, pkgManager, function.Results)
	if err != nil {
		return "", merror.Stack(err)
	}

	str += fmt.Sprintf("%s(%s)(%s)", function.Name, argsStr, resultsStr)
	return str, nil
}
