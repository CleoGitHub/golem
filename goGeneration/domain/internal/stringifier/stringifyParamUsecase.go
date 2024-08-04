package stringifier

import (
	"context"
	"fmt"

	"github.com/cleoGitHub/golem-common/pkg/merror"
	"github.com/cleoGitHub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
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
