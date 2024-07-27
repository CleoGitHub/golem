package usecase

import (
	"context"
	"os/exec"

	"github.com/cleoGitHub/golem/goGeneration/domain/model"
	"github.com/cleoGitHub/golem/pkg/merror"
)

func (g *GenerationUsecaseImpl) GoTidyDomainUsecase(ctx context.Context, path string, domain *model.Domain) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = path + "/" + domain.Name

	if err := cmd.Run(); err != nil {
		return merror.Stack(err)
	}

	return nil
}
