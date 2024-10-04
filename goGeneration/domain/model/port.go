package model

type File struct {
	Name string

	Elements []interface{}
	Pkg      *GoPkg
}
