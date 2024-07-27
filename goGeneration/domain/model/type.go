package model

// ensure all types implement Type interface
var _ Type = PrimitiveType("")
var _ Type = &ArrayType{}
var _ Type = &MapType{}
var _ Type = &ExternalType{}
var _ Type = &PointerType{}
var _ Type = &VariaidicType{}
var _ Type = &GormModel{}
var _ Type = &Function{}
var _ Type = &PkgReference{}

type GetTypeContext struct {
	InPkg string
}

type GetTypeOpt func(*GetTypeContext)

func InPkg(pkg string) GetTypeOpt {
	return func(ctx *GetTypeContext) {
		ctx.InPkg = pkg
	}
}

type Type interface {
	// GetType return the string type version
	GetType(...GetTypeOpt) string
	// SubTypes return all types used to stringify the type
	SubTypes() []Type
	// Copy return a copy of the type
	Copy() Type
}

type PrimitiveType string

// add all Golang primitives and other recurrent types
const (
	PrimitiveTypeInt    PrimitiveType = "int64"
	PrimitiveTypeFloat  PrimitiveType = "float64"
	PrimitiveTypeString PrimitiveType = "string"
	PrimitiveTypeBool   PrimitiveType = "bool"
	PrimitiveTypeByte   PrimitiveType = "byte"
	PrimitiveTypeBytes  PrimitiveType = "[]byte"
	PrimitiveTypeError  PrimitiveType = "error"

	PrimitiveTypeInterface PrimitiveType = "interface{}"
)

func (t PrimitiveType) GetType(...GetTypeOpt) string {
	return string(t)
}

func (t PrimitiveType) SubTypes() []Type {
	return []Type{}
}

func (t PrimitiveType) Copy() Type {
	return t
}

type ArrayType struct {
	Type Type
}

func (t *ArrayType) GetType(opts ...GetTypeOpt) string {
	return "[]" + t.Type.GetType(opts...)
}

func (t *ArrayType) SubTypes() []Type {
	ts := []Type{
		t.Type,
	}
	return append(ts, t.Type.SubTypes()...)
}

func (t *ArrayType) Copy() Type {
	return &ArrayType{Type: t.Type.Copy()}
}

type MapType struct {
	Key   Type
	Value Type
}

func (t *MapType) GetType(opts ...GetTypeOpt) string {
	return "map[" + t.Key.GetType(opts...) + "]" + t.Value.GetType(opts...)
}

func (t *MapType) SubTypes() []Type {
	ts := []Type{t.Key, t.Value}

	ts = append(ts, t.Key.SubTypes()...)
	return append(ts, t.Value.SubTypes()...)
}

func (t *MapType) Copy() Type {
	return &MapType{
		Key:   t.Key.Copy(),
		Value: t.Value.Copy(),
	}
}

type PointerType struct {
	Type Type
}

func (t *PointerType) GetType(opts ...GetTypeOpt) string {
	return "*" + t.Type.GetType(opts...)
}

func (t *PointerType) SubTypes() []Type {
	ts := []Type{t.Type}

	return append(ts, t.Type.SubTypes()...)
}

func (t *PointerType) Copy() Type {
	return &PointerType{Type: t.Type.Copy()}
}

type ExternalType struct {
	Type string
}

func (t *ExternalType) GetType(...GetTypeOpt) string {
	return t.Type
}

func (t *ExternalType) SubTypes() []Type {
	return []Type{}
}

func (t *ExternalType) Copy() Type {
	return &ExternalType{Type: t.Type}
}

// Variadic type represent parameters that can be in a non deterministic number in a function
type VariaidicType struct {
	Type Type
}

func (t *VariaidicType) GetType(opts ...GetTypeOpt) string {
	return "..." + t.Type.GetType(opts...)
}

func (t *VariaidicType) SubTypes() []Type {
	ts := []Type{t.Type}

	return append(ts, t.Type.SubTypes()...)
}

func (t *VariaidicType) Copy() Type {
	return &VariaidicType{Type: t.Type.Copy()}
}
