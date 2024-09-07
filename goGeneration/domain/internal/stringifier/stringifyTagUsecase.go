package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/merror"
)

func StringifyTagUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, tag *model.Tag) (string, error) {
	str := ""
	for idx, value := range tag.Values {
		str += value
		if idx < len(tag.Values)-1 {
			str += ","
		}
	}
	str = fmt.Sprintf(`%s:"%s"`, tag.Name, str)
	return str, nil
}

func StringifyTagsUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, tags []*model.Tag) (string, error) {
	str := ""
	for idx, tag := range tags {
		s, err := StringifyTagUsecase(ctx, pkgManager, tag)
		if err != nil {
			return "", merror.Stack(err)
		}

		str += s
		if idx < len(tags)-1 {
			str += " "
		}
	}
	if str == "" {
		return "", nil
	}
	return fmt.Sprintf("`%s`", str), nil
}
