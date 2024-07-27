package model

type Field struct {
	Name     string
	Type     Type
	Tags     []*Tag
	JsonName string
}

func (f *Field) Copy() *Field {
	return &Field{
		Name: f.Name,
		Type: f.Type,
		Tags: Tags(f.Tags).Copy(),
	}
}

type Fields []*Field

func (f Fields) Copy() Fields {
	fields := make(Fields, 0, len(f))
	for _, field := range f {
		fields = append(fields, field.Copy())
	}
	return fields
}
