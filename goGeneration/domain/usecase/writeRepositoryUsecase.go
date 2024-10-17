package usecase

// func (g *GenerationUsecaseImpl) WriteRepositoryUsecase(ctx context.Context, domain *model.Domain, repo *model.Repository, path string) error {
// 	// if file path does not exist, create it
// 	filepath := path + "/" + domain.Architecture.RepositoryPkg.FullName
// 	if _, err := os.Stat(filepath); os.IsNotExist(err) {
// 		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
// 			return merror.Stack(err)
// 		}
// 	}

// 	f, err := os.Create(filepath + "/" + stringtool.LowerFirstLetter(repo.Name) + ".go")
// 	if err != nil {
// 		return merror.Stack(err)
// 	}
// 	defer f.Close()

// 	pkgManager := &gopkgmanager.GoPkgManager{
// 		Pkg: domain.Architecture.RepositoryPkg.ShortName,
// 	}

// 	str, err := stringifier.StringifyRepositoryUsecase(ctx, pkgManager, repo)
// 	if err != nil {
// 		return merror.Stack(err)
// 	}

// 	str = pkgManager.ToString() + consts.LN + str
// 	_, err = f.WriteString(str)
// 	if err != nil {
// 		return merror.Stack(err)
// 	}

// 	return nil
// }
