package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func StringifyStructUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, s *model.Struct) (string, error) {
	str := ""

	if s.Consts != nil {
		for _, c := range s.Consts {
			c, err := StringifyConstUsecase(ctx, pkgManager, c)
			if err != nil {
				return "", merror.Stack(err)
			}
			str += c + consts.LN
		}
	}

	str += fmt.Sprintf("type %s struct {", s.Name) + consts.LN
	for _, field := range s.Fields {
		// Never pass model pkg in GetType call, model are always in the same package
		st, err := StringifyFieldUsecase(ctx, pkgManager, field)
		if err != nil {
			return "", merror.Stack(err)
		}
		str += st + consts.LN
	}
	str += "}" + consts.LN

	for _, method := range s.Methods {
		st, err := StringifyFunctionDefinitionUsecase(ctx, pkgManager, method)
		if err != nil {
			return "", merror.Stack(err)
		}
		str += fmt.Sprintf("func (%s *%s) %s {", s.GetMethodName(), s.Name, st) + consts.LN
		if method.Content == nil {
			fmt.Printf("%s\n", method.Name)
		}
		s, pkgs := method.Content()
		for _, pkg := range pkgs {
			if err := pkgManager.ImportPkg(pkg); err != nil {
				return "", merror.Stack(err)
			}
		}
		str += s
		str += "}" + consts.LN
		str += consts.LN
	}

	return str, nil
}
