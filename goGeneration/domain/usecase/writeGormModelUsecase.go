package usecase

// func (g *GenerationUsecaseImpl) WriteGormModelUsecase(ctx context.Context, domain *model.Domain, m *model.GormModel, path string) error {
// 	// if file path does not exist, create it
// 	filepath := path + "/" + domain.Architecture.GormAdapterPkg.FullName
// 	if _, err := os.Stat(filepath); os.IsNotExist(err) {
// 		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
// 			return merror.Stack(err)
// 		}
// 	}

// 	f, err := os.Create(filepath + "/" + stringtool.LowerFirstLetter(m.Struct.Name) + ".go")
// 	if err != nil {
// 		return merror.Stack(err)
// 	}
// 	defer f.Close()

// 	pkgManager := &gopkgmanager.GoPkgManager{
// 		Pkg: domain.Architecture.GormAdapterPkg.ShortName,
// 	}

// 	str, err := stringifier.StringifyStructUsecase(ctx, pkgManager, m.Struct)
// 	if err != nil {
// 		return merror.Stack(err)
// 	}

// 	s, err := stringifier.StringifyFunctionUsecase(ctx, pkgManager, m.FromModel)
// 	if err != nil {
// 		return merror.Stack(err)
// 	}
// 	str += s + consts.LN

// 	s, err = stringifier.StringifyFunctionUsecase(ctx, pkgManager, m.ToModel)
// 	if err != nil {
// 		return merror.Stack(err)
// 	}
// 	str += s + consts.LN

// 	s, err = stringifier.StringifyFunctionUsecase(ctx, pkgManager, m.FromModels)
// 	if err != nil {
// 		return merror.Stack(err)
// 	}
// 	str += s + consts.LN

// 	s, err = stringifier.StringifyFunctionUsecase(ctx, pkgManager, m.ToModels)
// 	if err != nil {
// 		return merror.Stack(err)
// 	}
// 	str += s + consts.LN

// 	str = pkgManager.ToString() + consts.LN + str
// 	_, err = f.WriteString(str)
// 	if err != nil {
// 		return merror.Stack(err)
// 	}

// 	return nil
// }
