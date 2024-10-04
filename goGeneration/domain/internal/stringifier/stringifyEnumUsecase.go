package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func StringifyEnumUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, enum *model.Enum) (string, error) {
	t := enum.Type.GetType(model.InPkg(pkgManager.Pkg))

	str := "const (" + consts.LN
	for key, value := range enum.Values {
		if _, ok := value.(string); ok {
			value = fmt.Sprintf(`"%v"`, value)
		}
		str += fmt.Sprintf("%s %s = %s", key, t, value) + consts.LN
	}
	str += ")" + consts.LN
	return str, nil
}
