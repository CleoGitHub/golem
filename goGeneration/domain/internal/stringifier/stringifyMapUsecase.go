package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func StringifyMapUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, m *model.Map) (string, error) {
	if m == nil || len(m.Values) == 0 {
		return "", nil
	}
	keyTypeStr, err := StringifyTypeUsecase(ctx, pkgManager, m.Type.Key)
	if err != nil {
		return "", merror.Stack(err)
	}
	valueTypeStr, err := StringifyTypeUsecase(ctx, pkgManager, m.Type.Value)
	if err != nil {
		return "", merror.Stack(err)
	}
	str := fmt.Sprintf("var %s = map[%s]%s{", m.Name, keyTypeStr, valueTypeStr) + consts.LN
	for _, value := range m.Values {
		k := value.Key
		v := value.Value
		if m.Type.Key == model.PrimitiveTypeString {
			k = fmt.Sprintf(`"%v"`, k)
		}
		if m.Type.Value == model.PrimitiveTypeString {
			v = fmt.Sprintf(`"%v"`, v)
		}
		str += fmt.Sprintf("%v: %v,", k, v) + consts.LN
	}
	str += "}" + consts.LN

	return str, nil
}
