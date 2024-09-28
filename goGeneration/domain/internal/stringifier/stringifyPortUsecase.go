package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func StringifyPortUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, port *model.Port) (string, error) {
	str := ""
	for _, port := range port.Elements {
		switch t := port.(type) {
		case *model.Consts:
			s, err := StringifyConstsUsecase(ctx, pkgManager, t)
			if err != nil {
				return "", merror.Stack(err)
			}
			str += s + consts.LN
		case *model.Var:
			s, err := StringifyConstUsecase(ctx, pkgManager, t)
			if err != nil {
				return "", merror.Stack(err)
			}
			str += s + consts.LN
		case *model.Function:
			s, err := StringifyFunctionUsecase(ctx, pkgManager, t)
			if err != nil {
				return "", merror.Stack(err)
			}
			str += s + consts.LN
		case *model.Interface:
			s, err := StringifyInterfaceUsecase(ctx, pkgManager, t)
			if err != nil {
				return "", merror.Stack(err)
			}
			str += s + consts.LN
		case *model.Map:
			s, err := StringifyMapUsecase(ctx, pkgManager, t)
			if err != nil {
				return "", merror.Stack(err)
			}
			str += s + consts.LN
		case *model.Struct:
			s, err := StringifyStructUsecase(ctx, pkgManager, t)
			if err != nil {
				return "", merror.Stack(err)
			}
			str += s + consts.LN
		case *model.TypeDefinition:
			s, err := StringifyTypeDefinitionUsecase(ctx, pkgManager, t)
			if err != nil {
				return "", merror.Stack(err)
			}
			str += s + consts.LN
		default:
			return "", merror.Stack(fmt.Errorf("unexpected type %T in Port", t))
		}
	}

	return str, nil
}
