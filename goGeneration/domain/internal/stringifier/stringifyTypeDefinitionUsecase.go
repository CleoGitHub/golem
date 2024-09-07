package stringifier

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/merror"
)

func StringifyTypeDefinitionUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, td *model.TypeDefinition) (string, error) {
	str, err := StringifyTypeUsecase(ctx, pkgManager, td.Type)
	if err != nil {
		return "", merror.Stack(err)
	}
	return fmt.Sprintf("type %s %s", td.Name, str), nil
}
