package model

type Template struct {
	Pkg            *GoPkg
	Name           string
	Data           TemplatePackageable
	StringTemplate string
}

type TemplatePackageable interface {
	SetPkg(*GoPkg)
}
