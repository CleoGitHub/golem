package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func StringifyConstUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, c *model.Var) (string, error) {
	var err error
	t := ""
	if c.Type != nil {
		t, err = StringifyTypeUsecase(ctx, pkgManager, c.Type)
		if err != nil {
			return "", merror.Stack(err)
		}
	}

	var v interface{}
	switch tpe := c.Value.(type) {
	case model.Type:
		t, err := StringifyTypeUsecase(ctx, pkgManager, tpe)
		if err != nil {
			return "", merror.Stack(err)
		}
		v = t
	case string:
		v = fmt.Sprintf(`"%s"`, tpe)
	default:
		v = c.Value
	}

	variable := "var"
	if c.IsConst {
		variable = "const"
	}

	str := fmt.Sprintf("%s %s %s = %v", variable, c.Name, t, v)
	return str, nil
}
