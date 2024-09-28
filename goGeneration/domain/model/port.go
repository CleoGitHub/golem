package model

type Port struct {
	Name string

	Elements []interface{}
	Pkg      *GoPkg
}
