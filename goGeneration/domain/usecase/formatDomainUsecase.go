package usecase

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/cleoGitHub/golem/goGeneration/domain/model"
	"github.com/cleoGitHub/golem/pkg/merror"
)

func (u GenerationUsecaseImpl) FormatDomainUsecase(ctx context.Context, domain *model.Domain, path string) error {
	// use command gofmt to format go files in generation folder
	cmd := exec.Command("gofmt", "-w", "-s", path+"/"+domain.Name)
	errWriter := bytes.NewBufferString("")
	cmd.Stderr = errWriter
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", errWriter.String())
		return merror.Stack(err)
	}
	return nil
}
