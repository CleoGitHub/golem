package usecase

import (
	"context"
	"os"
	"os/exec"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (g *GenerationUsecaseImpl) InitDomainUsecase(ctx context.Context, domain *model.Domain, path string) error {
	// if go.mod file does not exist at root folder, create it
	//check path exist or create if
	if _, err := os.Stat(path + "/" + domain.Name); os.IsNotExist(err) {
		if err := os.MkdirAll(path+"/"+domain.Name, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	// init go.mod file if not exist
	if _, err := os.Stat(path + "/" + domain.Name + "/go.mod"); os.IsNotExist(err) {
		// create go.mod file
		cmd := exec.Command("go", "mod", "init", domain.Name)
		cmd.Dir = path + "/" + domain.Name
		if err := cmd.Run(); err != nil {
			return merror.Stack(err)
		}
	}

	return nil
}
