package stringifier

import (
	"context"
	"fmt"

	"github.com/cleoGitHub/golem-common/pkg/merror"
	"github.com/cleoGitHub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
)

func StringifyTypeDefinitionUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, td *model.TypeDefinition) (string, error) {
	str, err := StringifyTypeUsecase(ctx, pkgManager, td.Type)
	if err != nil {
		return "", merror.Stack(err)
	}
	return fmt.Sprintf("type %s %s", td.Name, str), nil
}
