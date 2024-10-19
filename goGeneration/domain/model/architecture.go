package model

type Architecture struct {
	ModelPkg          *GoPkg
	RepositoryPkg     *GoPkg
	UsecasePkg        *GoPkg
	GormAdapterPkg    *GoPkg
	ControllerPkg     *GoPkg
	SdkPkg            *GoPkg
	HttpControllerPkg *GoPkg
	ConstsPkg         *GoPkg
	JavascriptClient  string
}
