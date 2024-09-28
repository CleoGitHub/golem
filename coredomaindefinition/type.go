package coredomaindefinition

type Type interface {
	GetType() string
}

type PrimitiveType string

// All meta type to define primitive type
const (
	PrimitiveTypeInt      PrimitiveType = "int"
	PrimitiveTypeFloat    PrimitiveType = "float"
	PrimitiveTypeString   PrimitiveType = "string"
	PrimitiveTypeBool     PrimitiveType = "bool"
	PrimitiveTypeByte     PrimitiveType = "byte"
	PrimitiveTypeBytes    PrimitiveType = "bytes"
	PrimitiveTypeDate     PrimitiveType = "date"
	PrimitiveTypeDateTime PrimitiveType = "datetime"
	PrimitiveTypeTime     PrimitiveType = "time"
	PrimitiveTypeFile     PrimitiveType = "file"
)

func (t PrimitiveType) GetType() string {
	return string(t)
}
