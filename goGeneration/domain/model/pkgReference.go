package model

type PkgReference struct {
	Pkg       *GoPkg
	Reference Type
}

func (r *PkgReference) GetType(opts ...GetTypeOpt) string {
	ctx := &GetTypeContext{}
	for _, opt := range opts {
		opt(ctx)
	}

	if ctx.InPkg == r.Pkg.Alias {
		return r.Reference.GetType(opts...)
	}

	return r.Pkg.Alias + "." + r.Reference.GetType(opts...)
}

func (r *PkgReference) SubTypes() []Type {
	return []Type{r.Reference}
}

func (r *PkgReference) Copy() *PkgReference {
	return &PkgReference{
		Pkg:       r.Pkg,
		Reference: Copy(r.Reference),
	}
}
