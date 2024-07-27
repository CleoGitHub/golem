package model

type Map struct {
	Name   string
	Type   MapType
	Values []MapValue
}

type MapValue struct {
	Key   interface{}
	Value interface{}
}
