package stringifier

import (
	"context"
	"fmt"
	"reflect"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func StringifyConstsUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, c *model.Consts) (string, error) {
	t := ""
	str := ""
	for _, value := range c.Values {
		if t == "" {
			t = reflect.TypeOf(value).String()
		} else if t != "interface{}" {
			if t != reflect.TypeOf(value).String() {
				t = "interface{}"
			}
		}
		v := value
		if _, ok := v.(string); ok {
			v = fmt.Sprintf(`"%v"`, v)
		}
		str += fmt.Sprintf("%v,", v) + consts.LN
	}
	str = fmt.Sprintf("var %s = []%s{", c.Name, t) + consts.LN + str
	str += "}" + consts.LN
	return str, nil
}
