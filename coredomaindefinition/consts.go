package coredomaindefinition

const (
	DOMAIN_FOLDER     = "domain"
	MODEL_FOLDER      = "model"
	USECASE_FOLDER    = "usecase"
	PORT_FOLDER       = "port"
	ADAPTER_FOLDER    = "adapter"
	CONTROLLER_FOLDER = "controller"
	REPOSITORY_FOLDER = "repository"

	DOMAIN_PATH  = ""
	USECASE_PATH = DOMAIN_PATH + "/" + DOMAIN_FOLDER
	MODEL_PATH   = DOMAIN_PATH + "/" + DOMAIN_FOLDER
	PORT_PATH    = DOMAIN_PATH + "/" + DOMAIN_FOLDER

	CONTROLLER_PATH = PORT_PATH + "/" + PORT_FOLDER
	REPOSITORY_PATH = PORT_PATH + "/" + PORT_FOLDER

	ADAPTER_PATH = ""
)
