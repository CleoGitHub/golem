package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func StringifyConstsUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, c *model.Consts) (string, error) {
	str := fmt.Sprintf("var %s = []interface{}{", c.Name) + consts.LN
	for _, value := range c.Values {
		v := value
		if _, ok := v.(string); ok {
			v = fmt.Sprintf(`"%v"`, v)
		}
		str += fmt.Sprintf("%v,", v) + consts.LN
	}
	str += "}" + consts.LN
	return str, nil
}
