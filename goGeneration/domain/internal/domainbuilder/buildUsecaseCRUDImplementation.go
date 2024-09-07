package domainbuilder

import (
	"context"
	"fmt"
	"strings"

	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
	"github.com/cleogithub/golem/pkg/merror"
	"github.com/cleogithub/golem/pkg/stringtool"
)

func (b *domainBuilder) buildUsecaseCRUDImplementation(ctx context.Context) *domainBuilder {
	if b.Err != nil {
		return nil
	}

	domainRepository := &model.Field{
		Name: "DomainRepository",
		Type: &model.PkgReference{
			Pkg:       b.Domain.Architecture.RepositoryPkg,
			Reference: b.Domain.DomainRepository,
		},
	}

	validator := &model.Field{
		Name: "Validator",
		Type: b.GetValidator(ctx),
	}

	usecaseCRUDImpl := &model.Struct{
		Name:   stringtool.UpperFirstLetter(b.Domain.Name) + "UsecasesCRUDImpl",
		Fields: []*model.Field{domainRepository, validator},
	}

	structName := usecaseCRUDImpl.GetMethodName()

	for _, crud := range b.CRUDToBuild {
		on, err := b.GetModel(ctx, crud.On)
		if err != nil {
			b.Err = merror.Stack(err)
			return nil
		}

		createOrUpdateImpl := func(usecase *model.Usecase, action string) {
			f := usecase.Function.Copy().(*model.Function)
			request := usecase.Function.Args[1].Name
			f.Content = func() (content string, requiredPkg []*model.GoPkg) {
				pkgs := []*model.GoPkg{}
				for _, relation := range on.Relations {
					if relation.Type == model.RelationSingleMandatory || relation.Type == model.RelationSingleOptionnal {
						content += "// Check relation exist" + consts.LN
						for _, m := range b.Domain.DomainRepository.Methods {
							if strings.Contains(m.Name, fmt.Sprintf("Get%s", relation.On.Struct.Name)) {
								content += fmt.Sprintf(
									`if _,err := %s.%s.Get%s(%sctx, %s%s.Get%sWithBy(map[string]interface{}{"Id": %s.%s}),%s%s.Get%sWithRetriveInactive(true),%s); err != nil {`,
									structName,
									domainRepository.Name,
									relation.On.Struct.Name,
									consts.LN,
									consts.LN,
									b.Domain.Architecture.RepositoryPkg.Alias,
									relation.On.Struct.Name,
									request,
									relation.On.Struct.Name+"Id",
									consts.LN,
									b.Domain.Architecture.RepositoryPkg.Alias,
									relation.On.Struct.Name,
									consts.LN,
								) + consts.LN

								content += fmt.Sprintf("if errors.Is(err, %s.%s) {",
									b.Domain.Architecture.RepositoryPkg.Alias,
									b.RepositoryErrors["notFound"].Name,
								) + consts.LN

								content += fmt.Sprintf(
									"%s.%s.%s(ctx, %s.%s)",
									structName,
									validator.Name,
									b.GetValidator(ctx).Methods[2].Name,
									request,
									relation.On.Struct.Name+"Id",
								) + consts.LN
								content += "} else {" + consts.LN
								content += "return nil, err" + consts.LN
								content += "}" + consts.LN
								content += "}" + consts.LN

								pkgs = append(pkgs, b.Domain.Architecture.RepositoryPkg)
								pkgs = append(pkgs, consts.CommonPkgs["errors"])
							}
						}
					}
				}

				for _, field := range on.Struct.Fields {
					for _, validation := range b.FieldToValidationRules[field] {
						switch validation.Rule {
						case coredomaindefinition.ValidationRuleUnique:
							content += "// Check uniquiness of " + field.Name + consts.LN
							content += fmt.Sprintf(
								`_,err := %s.%s.Get%s(%sctx, %s%s.Get%sWithBy(map[string]interface{}{"%s": %s.%s}),%s`,
								structName,
								domainRepository.Name,
								on.Struct.Name,
								consts.LN,
								consts.LN,
								b.Domain.Architecture.RepositoryPkg.Alias,
								on.Struct.Name,
								field.Name,
								request,
								field.Name,
								consts.LN,
							)
							content += fmt.Sprintf("%s.Get%sWithRetriveInactive(true),%s",
								b.Domain.Architecture.RepositoryPkg.Alias,
								on.Struct.Name,
								consts.LN,
							)
							if action == "Update" {
								content += fmt.Sprintf(
									`%s.Get%sWithNot(map[string]interface{}{"Id": %s.Id}),%s`,
									b.Domain.Architecture.RepositoryPkg.Alias,
									on.Struct.Name,
									request,
									consts.LN,
								)
							}
							content += ")" + consts.LN
							content += fmt.Sprintf("if err != nil  && !errors.Is(err, %s.%s) {",
								b.Domain.Architecture.RepositoryPkg.Alias,
								b.RepositoryErrors["notFound"].Name,
							) + consts.LN
							content += fmt.Sprintf("return nil, %s.Stack(err)", consts.CommonPkgs["merror"].Alias) + consts.LN
							content += "} else if err == nil {" + consts.LN
							content += fmt.Sprintf(
								"%s.%s.%s(ctx, %s.%s)",
								structName,
								validator.Name,
								b.GetValidator(ctx).Methods[3].Name,
								request,
								field.Name,
							) + consts.LN
							content += "}" + consts.LN
						case coredomaindefinition.ValidationRuleUniqueIn:
							m, ok := validation.Value.(*coredomaindefinition.Model)
							if !ok {
								continue
							}
							uniqueIn, err := b.GetModel(ctx, m)
							if err != nil {
								continue
							}
							content += fmt.Sprintf("// Check uniquiness of %s %s in %s", on.Struct.Name, field.Name, uniqueIn.Struct.Name) + consts.LN
							content += fmt.Sprintf(
								`_,err := %s.%s.Get%s(%sctx, %s%s.Get%sWithBy(map[string]interface{}{"%s": %s.%s, "%sId": %s.%sId}),%s`,
								structName,
								domainRepository.Name,
								on.Struct.Name,
								consts.LN,
								consts.LN,
								b.Domain.Architecture.RepositoryPkg.Alias,
								on.Struct.Name,
								field.Name,
								request,
								field.Name,
								uniqueIn.Struct.Name,
								request,
								uniqueIn.Struct.Name,
								consts.LN,
							)
							if action == "Update" {
								content += fmt.Sprintf(
									`%s.Get%sWithNot(map[string]interface{}{"Id": %s.Id}),%s`,
									b.Domain.Architecture.RepositoryPkg.Alias,
									on.Struct.Name,
									request,
									consts.LN,
								)
							}
							content += ")" + consts.LN
							content += fmt.Sprintf("if err != nil  && !errors.Is(err, %s.%s) {",
								b.Domain.Architecture.RepositoryPkg.Alias,
								b.RepositoryErrors["notFound"].Name,
							) + consts.LN
							content += fmt.Sprintf("return nil, %s.Stack(err)", consts.CommonPkgs["merror"].Alias) + consts.LN
							content += "} else if err == nil {" + consts.LN
							content += fmt.Sprintf(
								"%s.%s.%s(ctx, %s.%s)",
								structName,
								validator.Name,
								b.GetValidator(ctx).Methods[3].Name,
								request,
								field.Name,
							) + consts.LN
							content += "}" + consts.LN
						}
					}
				}

				content += fmt.Sprintf(
					"entity := &%s.%s{",
					b.Domain.Architecture.ModelPkg.Alias,
					on.Struct.Name,
				) + consts.LN

				for _, field := range usecase.Function.Args[1].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct).Fields {
					content += fmt.Sprintf(
						"%s: %s.%s,",
						field.Name,
						usecase.Function.Args[1].Name,
						field.Name,
					) + consts.LN
				}
				pkgs = append(pkgs, b.Domain.Architecture.ModelPkg)

				content += "}" + consts.LN

				content += fmt.Sprintf(
					"%s, err := %s.%s.%s%s(ctx, entity)",
					on.JsonName,
					structName,
					domainRepository.Name,
					action,
					on.Struct.Name,
				) + consts.LN

				content += consts.IF_ERR
				content += fmt.Sprintf("return nil, %s.Stack(err)", consts.CommonPkgs["merror"].Alias) + consts.LN
				content += "}" + consts.LN
				pkgs = append(pkgs, consts.CommonPkgs["merror"])

				content += fmt.Sprintf(
					"return &%s{ %s: %s }, nil",
					usecase.Function.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct).Name,
					usecase.Function.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct).Fields[0].Name,
					on.JsonName,
				) + consts.LN

				return content, pkgs
			}

			usecaseCRUDImpl.Methods = append(usecaseCRUDImpl.Methods, f)
		}

		if crud.Create != nil && crud.Create.Active {
			usecase, ok := b.CRUDActionToUsecase[crud.Create]
			if !ok {
				b.Err = merror.Stack(fmt.Errorf("no usecase found for CRUDAction"))
				return nil
			}
			createOrUpdateImpl(usecase, "Create")
		}

		if crud.Update != nil && crud.Update.Active {
			usecase, ok := b.CRUDActionToUsecase[crud.Update]
			if !ok {
				b.Err = merror.Stack(fmt.Errorf("no usecase found for CRUDAction"))
				return nil
			}
			createOrUpdateImpl(usecase, "Update")
		}

		getImpl := func(usecase *model.Usecase, retriveInactive bool) {
			f := usecase.Function.Copy().(*model.Function)
			if !retriveInactive {
				f.Name = strings.Replace(f.Name, "Get", "GetActive", 1)
				s := f.Args[1].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct)
				s.Name = strings.Replace(s.Name, "Get", "GetActive", 1)
				s = f.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct)
				s.Name = strings.Replace(s.Name, "Get", "GetActive", 1)
			}
			request := usecase.Function.Args[1].Name
			f.Content = func() (string, []*model.GoPkg) {
				pkgs := make([]*model.GoPkg, 0)
				content := ""
				// getRepoMethod, ok := b.RepoToDomainRepoGetEntityRepoMethod[b.ModelToRepository[on]]
				// if !ok {
				// 	return "", nil
				// }
				content += fmt.Sprintf(
					`%s, err := %s.%s.Get%s(%sctx, %s%s.Get%sWithBy(map[string]interface{}{"Id": %s.Id}),`,
					on.JsonName,
					structName,
					domainRepository.Name,
					on.Struct.Name,
					consts.LN,
					consts.LN,
					b.Domain.Architecture.RepositoryPkg.Alias,
					on.Struct.Name,
					request,
				) + consts.LN

				if retriveInactive {
					content += fmt.Sprintf(
						`%s.Get%sWithRetriveInactive(true),`,
						b.Domain.Architecture.RepositoryPkg.Alias,
						on.Struct.Name,
					) + consts.LN
				}

				content += ")" + consts.LN

				content += consts.IF_ERR
				content += fmt.Sprintf("return nil, %s.Stack(err)", consts.CommonPkgs["merror"].Alias) + consts.LN
				content += "}" + consts.LN
				pkgs = append(pkgs, consts.CommonPkgs["merror"])

				content += fmt.Sprintf(
					"return &%s{ %s: %s }, nil",
					f.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct).Name,
					usecase.Function.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct).Fields[0].Name,
					on.JsonName,
				) + consts.LN
				return content, pkgs
			}
			usecaseCRUDImpl.Methods = append(usecaseCRUDImpl.Methods, f)
		}
		if crud.Get != nil && crud.Get.Active {
			usecase, ok := b.CRUDActionToUsecase[crud.Get]
			if !ok {
				b.Err = merror.Stack(fmt.Errorf("no usecase found for CRUDAction"))
				return nil
			}
			getImpl(usecase, true)
			if crud.On.Activable {
				getImpl(usecase, false)
			}
		}

		listImpl := func(usecase *model.Usecase, retriveInactive bool) {
			f := usecase.Function.Copy().(*model.Function)
			fmt.Printf("listImpl: %s\n", f.Name)
			if !retriveInactive {
				f.Name = strings.Replace(f.Name, "List", "ListActive", 1)
				s := f.Args[1].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct)
				s.Name = strings.Replace(s.Name, "List", "ListActive", 1)
				s = f.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct)
				s.Name = strings.Replace(s.Name, "List", "ListActive", 1)
			}
			request := usecase.Function.Args[1].Name
			f.Content = func() (string, []*model.GoPkg) {
				pkgs := make([]*model.GoPkg, 0)
				content := ""

				content += fmt.Sprintf(
					`%ss, err := %s.%s.List%s(%sctx,`,
					on.JsonName,
					structName,
					domainRepository.Name,
					on.Struct.Name,
					consts.LN,
				) + consts.LN
				content += fmt.Sprintf(
					"%s.List%sWithPagination(%s.%s),",
					b.Domain.Architecture.RepositoryPkg.Alias,
					on.Struct.Name,
					request,
					b.GetPagination(ctx).Name,
				) + consts.LN
				content += fmt.Sprintf(
					"%s.List%sWithOrdering(%s.%s),",
					b.Domain.Architecture.RepositoryPkg.Alias,
					on.Struct.Name,
					request,
					b.GetOrdering(ctx).Name,
				) + consts.LN
				if retriveInactive {
					content += fmt.Sprintf(
						`%s.List%sWithRetriveInactive(true),`,
						b.Domain.Architecture.RepositoryPkg.Alias,
						on.Struct.Name,
					) + consts.LN
				}
				content += ")" + consts.LN

				content += consts.IF_ERR
				content += fmt.Sprintf("return nil, %s.Stack(err)", consts.CommonPkgs["merror"].Alias) + consts.LN
				content += "}" + consts.LN
				pkgs = append(pkgs, consts.CommonPkgs["merror"])

				content += fmt.Sprintf(
					"return &%s{ %s: %ss }, nil",
					f.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct).Name,
					usecase.Function.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct).Fields[0].Name,
					on.JsonName,
				) + consts.LN
				return content, pkgs
			}
			usecaseCRUDImpl.Methods = append(usecaseCRUDImpl.Methods, f)
		}

		if crud.List != nil && crud.List.Active {
			usecase, ok := b.CRUDActionToUsecase[crud.List]
			if !ok {
				b.Err = merror.Stack(fmt.Errorf("no usecase found for CRUDAction"))
				return nil
			}
			listImpl(usecase, true)
			if crud.On.Activable {
				listImpl(usecase, false)
			}
		}

		if crud.Delete != nil && crud.Delete.Active {
			usecase, ok := b.CRUDActionToUsecase[crud.Delete]
			if !ok {
				b.Err = merror.Stack(fmt.Errorf("no usecase found for CRUDAction"))
				return nil
			}
			f := usecase.Function.Copy().(*model.Function)
			request := usecase.Function.Args[1].Name
			f.Content = func() (string, []*model.GoPkg) {
				pkgs := make([]*model.GoPkg, 0)
				content := ""
				// getRepoMethod, ok := b.RepoToDomainRepoGetEntityRepoMethod[b.ModelToRepository[on]]
				// if !ok {
				// 	return "", nil
				// }
				content += fmt.Sprintf(
					"err := %s.%s.Delete%s(ctx, %s.Id)",
					structName,
					domainRepository.Name,
					on.Struct.Name,
					request,
				) + consts.LN

				content += consts.IF_ERR
				content += fmt.Sprintf("return nil, %s.Stack(err)", consts.CommonPkgs["merror"].Alias) + consts.LN
				content += "}" + consts.LN
				pkgs = append(pkgs, consts.CommonPkgs["merror"])

				content += fmt.Sprintf("return &%s{}, nil", usecase.Function.Results[0].Type.(*model.PointerType).Type.(*model.PkgReference).Reference.(*model.Struct).Name) + consts.LN
				return content, pkgs
			}
			usecaseCRUDImpl.Methods = append(usecaseCRUDImpl.Methods, f)
		}
	}

	b.UsecaseCRUDImpl = usecaseCRUDImpl

	b.Domain.UsecasesCRUDImpl = usecaseCRUDImpl

	return b
}
