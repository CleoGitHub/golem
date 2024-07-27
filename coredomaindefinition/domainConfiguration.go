package coredomaindefinition

type DomainConfiguration struct {
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
}

func (domainConfiguration *DomainConfiguration) GetDomainFolder() string {
	if domainConfiguration.DomainFolder == "" {
		return DOMAIN_FOLDER
	}
	return domainConfiguration.DomainFolder
}

func (domainConfiguration *DomainConfiguration) GetModelFolder() string {
	if domainConfiguration.ModelFolder == "" {
		return MODEL_FOLDER
	}
	return domainConfiguration.ModelFolder
}

func (domainConfiguration *DomainConfiguration) GetUsecaseFolder() string {
	if domainConfiguration.UsecaseFolder == "" {
		return USECASE_FOLDER
	}
	return domainConfiguration.UsecaseFolder
}

func (domainConfiguration *DomainConfiguration) GetPortFolder() string {
	if domainConfiguration.PortFolder == "" {
		return PORT_FOLDER
	}
	return domainConfiguration.PortFolder
}

func (domainConfiguration *DomainConfiguration) GetAdapterFolder() string {
	if domainConfiguration.AdapterFolder == "" {
		return ADAPTER_FOLDER
	}
	return domainConfiguration.AdapterFolder
}

func (domainConfiguration *DomainConfiguration) GetControllerFolder() string {
	if domainConfiguration.ControllerFolder == "" {
		return CONTROLLER_FOLDER
	}
	return domainConfiguration.ControllerFolder
}

func (domainConfiguration *DomainConfiguration) GetRepositoryFolder() string {
	if domainConfiguration.RepositoryFolder == "" {
		return REPOSITORY_FOLDER
	}
	return domainConfiguration.RepositoryFolder
}

func (domainConfiguration *DomainConfiguration) GetDomainPath() string {
	if domainConfiguration.DomainPath == "" {
		return DOMAIN_PATH
	}
	return domainConfiguration.DomainPath
}

func (domainConfiguration *DomainConfiguration) GetUsecasePath() string {
	if domainConfiguration.UsecasePath == "" {
		return USECASE_PATH
	}
	return domainConfiguration.UsecasePath
}

func (domainConfiguration *DomainConfiguration) GetModelPath() string {
	if domainConfiguration.ModelPath == "" {
		return MODEL_PATH
	}
	return domainConfiguration.ModelPath
}

func (domainConfiguration *DomainConfiguration) GetPortPath() string {
	if domainConfiguration.PortPath == "" {
		return PORT_PATH
	}
	return domainConfiguration.PortPath
}

func (domainConfiguration *DomainConfiguration) GetControllerPath() string {
	if domainConfiguration.ControllerPath == "" {
		return CONTROLLER_PATH
	}
	return domainConfiguration.ControllerPath
}

func (domainConfiguration *DomainConfiguration) GetRepositoryPath() string {
	if domainConfiguration.RepositoryPath == "" {
		return REPOSITORY_PATH
	}
	return domainConfiguration.RepositoryPath
}

func (domainConfiguration *DomainConfiguration) GetAdapterPath() string {
	if domainConfiguration.AdapterPath == "" {
		return ADAPTER_PATH
	}
	return domainConfiguration.AdapterPath
}
