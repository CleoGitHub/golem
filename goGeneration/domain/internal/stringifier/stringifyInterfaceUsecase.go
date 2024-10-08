package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func StringifyInterfaceUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, itf *model.Interface) (string, error) {
	str := fmt.Sprintf("type %s interface {", itf.Name) + consts.LN
	for _, method := range itf.Methods {
		s, err := StringifyFunctionDefinitionUsecase(ctx, pkgManager, method)
		if err != nil {
			return "", merror.Stack(err)
		}
		str += s + consts.LN
	}
	str += "}"

	return str, nil
}
