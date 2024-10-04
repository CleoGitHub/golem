package usecase

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cleogithub/golem-common/pkg/merror"
)

// GenerateDomainUsecase implements GenerationUsecase.
func (g *GenerationUsecaseImpl) RemoveGenerationsUsecase(ctx context.Context, path string) error {
	// go deep through the domain and remove all generated filed ending with golem.go

	// list all files in the path
	files, err := os.ReadDir(path)
	if err != nil {
		return merror.Stack(err)
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".golem.go") {
			if err := os.Remove(path + "/" + f.Name()); err != nil {
				return merror.Stack(err)
			}
		} else if f.IsDir() {
			if err := removeGenerationInFolder(ctx, path+"/"+f.Name()); err != nil {
				return merror.Stack(err)
			}
		}
	}

	return nil
}

func removeGenerationInFolder(ctx context.Context, path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return merror.Stack(err)
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".golem.go") {
			fmt.Printf("Removing generated file: %s\n", path+"/"+f.Name())
			if err := os.Remove(path + "/" + f.Name()); err != nil {
				return merror.Stack(err)
			}
		} else if f.IsDir() {
			if err := removeGenerationInFolder(ctx, path+"/"+f.Name()); err != nil {
				return merror.Stack(err)
			}
		}
	}
	return nil
}
