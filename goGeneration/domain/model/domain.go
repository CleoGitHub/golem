package model

type Domain struct {
	Name         string
	Architecture *Architecture
	Models       []*Struct
	Files        []*File
	JSFiles      map[string]string
}
