package model

type GormModel struct {
	Struct     *Struct
	FromModel  *Function
	ToModel    *Function
	FromModels *Function
	ToModels   *Function
}

func (m *GormModel) GetType(typeOpts ...GetTypeOpt) string {
	return m.Struct.Name
}

func (m *GormModel) SubTypes() []Type {
	return []Type{}
}

func (m *GormModel) Copy() Type {
	return &GormModel{
		Struct:     m.Struct.Copy().(*Struct),
		FromModel:  m.FromModel.Copy().(*Function),
		ToModel:    m.ToModel.Copy().(*Function),
		FromModels: m.FromModels.Copy().(*Function),
		ToModels:   m.ToModels.Copy().(*Function),
	}
}
