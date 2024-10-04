package model

type Function struct {
	On      Type
	OnName  string
	Name    string
	Args    []*Param
	Results []*Param
	Content func() (content string, requiredPkg []*GoPkg)
}

func (t Function) GetType(opts ...GetTypeOpt) string {
	args := ""
	for i, arg := range t.Args {
		if i > 0 {
			args += ", "
		}
		args += arg.Name + " " + arg.Type.GetType(opts...)
	}

	results := ""
	for i, ret := range t.Results {
		if i > 0 {
			results += ", "
		}
		results += ret.Name + " " + ret.Type.GetType(opts...)
	}

	return "func(" + args + ") (" + results + ")"
}

func (t Function) SubTypes() []Type {
	ts := []Type{}
	for _, arg := range t.Args {
		ts = append(ts, arg.Type)
		ts = append(ts, arg.Type.SubTypes()...)
	}
	for _, ret := range t.Results {
		ts = append(ts, ret.Type)
		ts = append(ts, ret.Type.SubTypes()...)
	}

	return ts
}

// Add copy function
func (t *Function) Copy() Type {
	return &Function{
		Name:    t.Name,
		Args:    Params(t.Args).Copy(),
		Results: Params(t.Results).Copy(),
		Content: t.Content,
	}
}

type Functions []*Function

func (t Functions) Copy() Functions {
	ts := []*Function{}
	for _, arg := range t {
		ts = append(ts, arg.Copy().(*Function))
	}
	return ts
}
