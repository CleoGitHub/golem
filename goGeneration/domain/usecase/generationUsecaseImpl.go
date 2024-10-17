package usecase

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/domainbuilder"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/internal/stringifier"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

type GenerationUsecaseImpl struct {
}

var _ GenerationUsecase = &GenerationUsecaseImpl{}

// GenerateDomainUsecase implements GenerationUsecase.
func (g *GenerationUsecaseImpl) GenerateDomainUsecase(ctx context.Context, domainDefinition coredomaindefinition.Domain, path string) error {
	domainBuilder := domainbuilder.NewDomainBuilder(
		ctx,
		&domainDefinition,
		consts.DefaultModelFields,
	)

	for _, m := range domainDefinition.Models {
		domainBuilder.WithModel(ctx, m)
	}

	for _, r := range domainDefinition.Relations {
		domainBuilder.WithRelation(ctx, r)
	}

	for _, r := range domainDefinition.Repositories {
		domainBuilder.WithRepository(ctx, r)
	}

	for _, r := range domainDefinition.CRUDs {
		domainBuilder.WithCRUD(ctx, r)
	}

	for _, r := range domainDefinition.Usecases {
		domainBuilder.WithUsecase(ctx, r)
	}

	domain, err := domainBuilder.Build(ctx)
	if err != nil {
		return merror.Stack(err)
	}
	_ = domain

	if err := g.removeGenerationsUsecase(ctx, stringtool.RemoveDuplicate(path+"/"+domain.Name, '/')); err != nil {
		return merror.Stack(err)
	}

	if err := g.initDomainUsecase(ctx, domain, path); err != nil {
		return merror.Stack(err)
	}

	for _, m := range domain.Models {
		if err := g.write(ctx, domain.Architecture.ModelPkg, m, path); err != nil {
			return merror.Stack(err)
		}
	}

	for _, port := range domain.Files {
		if err := g.write(ctx, port.Pkg, port, path); err != nil {
			return merror.Stack(err)
		}
	}

	if err := g.formatDomainUsecase(ctx, domain, path); err != nil {
		return merror.Stack(err)
	}

	if err := g.goTidyDomainUsecase(ctx, path, domain); err != nil {
		return merror.Stack(err)
	}

	if err := g.generateJavascriptClientUsecase(ctx, domain, path); err != nil {
		return merror.Stack(err)
	}

	return nil
}

func (g *GenerationUsecaseImpl) initDomainUsecase(ctx context.Context, domain *model.Domain, path string) error {
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

func (g *GenerationUsecaseImpl) goTidyDomainUsecase(ctx context.Context, path string, domain *model.Domain) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = path + "/" + domain.Name

	if err := cmd.Run(); err != nil {
		return merror.Stack(err)
	}

	return nil
}

func (u GenerationUsecaseImpl) formatDomainUsecase(ctx context.Context, domain *model.Domain, path string) error {
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

// GenerateDomainUsecase implements GenerationUsecase.
func (g *GenerationUsecaseImpl) removeGenerationsUsecase(ctx context.Context, path string) error {
	// go deep through the domain and remove all generated filed ending with golem.go

	// list all files in the path
	files, err := os.ReadDir(path)
	if err != nil {
		return merror.Stack(err)
	}

	for _, f := range files {
		if !f.IsDir() && (strings.Contains(f.Name(), ".golem.")) {
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
		if !f.IsDir() && (strings.Contains(f.Name(), ".golem.")) {
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

func (g *GenerationUsecaseImpl) write(ctx context.Context, inPkg *model.GoPkg, elem interface{}, path string) (err error) {
	// if file path does not exist, create it
	filepath := path + "/" + inPkg.FullName
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	pkgManager := &gopkgmanager.GoPkgManager{
		Pkg: inPkg.ShortName,
	}
	str := ""
	name := ""

	switch t := elem.(type) {
	case *model.Struct:
		str, err = stringifier.StringifyStructUsecase(ctx, pkgManager, t)
		if err != nil {
			return merror.Stack(err)
		}

		name = t.Name
	case *model.Interface:
		str, err = stringifier.StringifyInterfaceUsecase(ctx, pkgManager, t)
		if err != nil {
			return merror.Stack(err)
		}

		name = t.Name
	case *model.File:
		str, err = stringifier.StringifyFileUsecase(ctx, pkgManager, t)
		if err != nil {
			return merror.Stack(err)
		}

		name = t.Name
	default:
		return merror.Stack(ErrUnknowTypeToWrite)
	}

	f, err := os.Create(filepath + "/" + stringtool.LowerFirstLetter(name) + ".golem.go")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	str = consts.HEADER + consts.LN + pkgManager.ToString() + consts.LN + str
	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	return nil
}

func (g *GenerationUsecaseImpl) generateJavascriptClientUsecase(ctx context.Context, domain *model.Domain, path string) error {
	// if file path does not exist, create it
	filepath := path + "/" + domain.Architecture.JavascriptClient
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	fileImport := ""
	export := ""

	for name, content := range domain.JSFiles {
		filename := name + ".golem.js"
		fileImport += fmt.Sprintf("export * from './%s';", filename) + consts.LN
		export += consts.TAB + fmt.Sprintf("%s,", name) + consts.LN

		// Generate service
		f, err := os.Create(filepath + "/" + filename)
		if err != nil {
			return merror.Stack(err)
		}
		defer f.Close()

		_, err = f.WriteString(content)
		if err != nil {
			return merror.Stack(err)
		}
	}

	f, err := os.Create(filepath + "/index.js")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	_, err = f.WriteString(fileImport)
	if err != nil {
		return merror.Stack(err)
	}

	return nil
}
