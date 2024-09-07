package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/merror"
)

func StringifyParamUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, p *model.Param) (string, error) {
	t, err := StringifyTypeUsecase(ctx, pkgManager, p.Type)
	if err != nil {
		return "", merror.Stack(err)
	}
	return fmt.Sprintf("%s %s", p.Name, t), nil
}

func StringifyParamsUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, params []*model.Param) (string, error) {
	str := ""
	for idx, param := range params {
		s, err := StringifyParamUsecase(ctx, pkgManager, param)
		if err != nil {
			return "", merror.Stack(err)
		}
		str += s
		if idx < len(params)-1 {
			str += ", "
		}
	}
	return str, nil
}
