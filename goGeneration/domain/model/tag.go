package model

type Tag struct {
	Name   string
	Values []string
}

func (t *Tag) Copy() *Tag {
	return &Tag{
		Name:   t.Name,
		Values: t.Values,
	}
}

type Tags []*Tag

func (t Tags) Copy() Tags {
	tags := make(Tags, 0, len(t))
	for _, tag := range t {
		tags = append(tags, tag.Copy())
	}
	return tags
}
