package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func StringifyFieldUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, field *model.Field) (string, error) {
	tagsStr, err := StringifyTagsUsecase(ctx, pkgManager, field.Tags)
	if err != nil {
		return "", merror.Stack(err)
	}
	typeStr, err := StringifyTypeUsecase(ctx, pkgManager, field.Type)
	if err != nil {
		return "", merror.Stack(err)
	}

	return fmt.Sprintf("%s %s %s", field.Name, typeStr, tagsStr), nil
}
