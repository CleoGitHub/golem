package domainwriter

import (
	"os"

	"github.com/cleogithub/golem/pkg/merror"
)

// DomainWriter is a struct to generate files and folder of hexagonal architecture.
type DomainWriter struct {
	DomainGenerationPath string

	// DomainFolder is the folder path of domain layer.
	DomainFolder string
	// ModelFolder is the folder path of model layer.
	ModelFolder string
	// UsecaseFolder is the folder path of use case layer.
	UsecaseFolder string
	// PortFolder is the folder path of port layer.
	PortFolder string
	// AdapterFolder is the folder path of adapter layer.
	AdapterFolder string
	// ControllerFolder is the folder path of controller layer.
	ControllerFolder string
	// RepositoryFolder is the folder path of repository layer.
	RepositoryFolder string

	// DomainPath is the path of domain layer.
	DomainPath string
	// UsecasePath is the path of use case layer.
	UsecasePath string
	// ModelPath is the path of model layer.
	ModelPath string
	// PortPath is the path of port layer.
	PortPath string
	// ControllerPath is the path of controller layer.
	ControllerPath string
	// RepositoryPath is the path of repository layer.
	RepositoryPath string
	// AdapterPath is the path of adapter layer.
	AdapterPath string

	// ModelFiles are the generated files of model layer.
	Files []*File
}

func NewDomainWriter() *DomainWriter {
	dw := &DomainWriter{
		DomainFolder:     DOMAIN_FOLDER,
		ModelFolder:      MODEL_FOLDER,
		UsecaseFolder:    USECASE_FOLDER,
		PortFolder:       PORT_FOLDER,
		AdapterFolder:    ADAPTER_FOLDER,
		ControllerFolder: CONTROLLER_FOLDER,
		RepositoryFolder: REPOSITORY_FOLDER,

		DomainPath:     DOMAIN_PATH,
		UsecasePath:    USECASE_PATH,
		ModelPath:      MODEL_PATH,
		PortPath:       PORT_PATH,
		ControllerPath: CONTROLLER_PATH,
		RepositoryPath: REPOSITORY_PATH,
		AdapterPath:    ADAPTER_PATH,
	}

	return dw
}

func (domainWriter *DomainWriter) GetModelFile() *File {
	file := &File{
		Path:     domainWriter.ModelPath,
		Contents: []Content{},
	}

	return file
}

func (domainWriter *DomainWriter) GetUseCaseFile() *File {
	file := &File{
		Path:     domainWriter.UsecasePath,
		Contents: []Content{},
	}

	return file
}

func (domainWriter *DomainWriter) GetPortFile() *File {
	file := &File{
		Path:     domainWriter.PortPath,
		Contents: []Content{},
	}

	return file
}

func (domainWriter *DomainWriter) GetRepositoryFile() *File {
	file := &File{
		Path:     domainWriter.RepositoryPath,
		Contents: []Content{},
	}

	return file
}

func (domainWriter *DomainWriter) Write() error {
	path := domainWriter.DomainGenerationPath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	for _, file := range domainWriter.Files {
		// Create folder if not exist
		if _, err := os.Stat(path + "/" + file.Path); os.IsNotExist(err) {
			if err := os.MkdirAll(path+"/"+file.Path, os.ModePerm); err != nil {
				return merror.Stack(err)
			}
		}

		// Create file
		f, err := os.Create(path + "/" + file.Path + "/" + file.Name)
		if err != nil {
			return merror.Stack(err)
		}
		defer f.Close()

		// Write data
		for _, content := range file.Contents {
			if _, err := f.WriteString(content.ToString()); err != nil {
				return merror.Stack(err)
			}
		}
	}

	return nil
}
