package domainbuilder

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/cleoGitHub/golem/coredomaindefinition"
	"github.com/cleoGitHub/golem/goGeneration/domain/consts"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
	"github.com/cleoGitHub/golem/pkg/merror"
	"github.com/cleoGitHub/golem/pkg/stringtool"
)

func (b *domainBuilder) isRelationUpTreeActivable(on *model.Model, activableModels map[*model.Model]bool, relationDone []*model.Relation, joins *string, pkgs *[]*model.GoPkg) bool {
	activable := on.Activable
	activableModels[on] = activable
	for _, relation := range on.Relations {
		// avoid infinite recursion
		if relationDone != nil && slices.Contains(relationDone, relation) {
			continue
		}
		relationDone = append(relationDone, relation)

		if relation.Definition.Type == coredomaindefinition.RelationTypeBelongsTo && b.ModelDefinitionToModel[relation.Definition.Source] == on {
			*pkgs = append(*pkgs, consts.CommonPkgs["fmt"])
			// If already visited relation, skip
			if isActivable, ok := activableModels[relation.On]; ok {
				*joins = fmt.Sprintf(
					`request.Joins("JOIN %s ON ? = ?", fmt.Sprintf("%s.%s", %s.%s["%sId"]), fmt.Sprintf("%s.%s", %s.%s["Id"]))`,
					b.ModelToRepository[relation.On].TableName,
					b.ModelToRepository[on].TableName,
					"%s",
					b.Domain.Architecture.RepositoryPkg.Alias,
					b.ModelToRepository[on].FieldToColumn.Name,
					relation.On.Struct.Name,
					b.ModelToRepository[relation.On].TableName,
					"%s",
					b.Domain.Architecture.RepositoryPkg.Alias,
					b.ModelToRepository[relation.On].FieldToColumn.Name,
				) + consts.LN + *joins
				activable = isActivable || activable
				continue
			}
			if b.isRelationUpTreeActivable(relation.On, activableModels, relationDone, joins, pkgs) {
				*joins = fmt.Sprintf(
					`request.Joins("JOIN %s ON ? = ?", fmt.Sprintf("%s.%s", %s.%s["%sId"]), fmt.Sprintf("%s.%s", %s.%s["Id"]))`,
					b.ModelToRepository[relation.On].TableName,
					b.ModelToRepository[on].TableName,
					"%s",
					b.Domain.Architecture.RepositoryPkg.Alias,
					b.ModelToRepository[on].FieldToColumn.Name,
					relation.On.Struct.Name,
					b.ModelToRepository[relation.On].TableName,
					"%s",
					b.Domain.Architecture.RepositoryPkg.Alias,
					b.ModelToRepository[relation.On].FieldToColumn.Name,
				) + consts.LN + *joins
				activable = true
				activableModels[relation.On] = true
			}
		}
	}
	return activable
}

func (b *domainBuilder) buildGormAdpater(ctx context.Context) *domainBuilder {
	if b.Err != nil {
		return b
	}

	tx := &model.Field{
		Name: "Tx",
		Type: &model.PointerType{
			Type: &model.PkgReference{
				Pkg: consts.CommonPkgs["gorm"],
				Reference: &model.ExternalType{
					Type: "DB",
				},
			},
		},
	}
	b.Domain.GormTransaction = &model.Struct{
		Name: "GormTransaction",
		Fields: []*model.Field{
			tx,
		},
		Methods: []*model.Function{
			{
				Name: "Get",
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PkgReference{
							Pkg: consts.CommonPkgs["context"],
							Reference: &model.ExternalType{
								Type: "Context",
							},
						},
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeInterface,
					},
				},
				Content: func() (content string, requiredPkg []*model.GoPkg) {
					return fmt.Sprintf(
						"return %s.%s", b.Domain.GormTransaction.GetMethodName(), tx.Name,
					), nil
				},
			},
			{
				Name: "Commit",
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PkgReference{
							Pkg: consts.CommonPkgs["context"],
							Reference: &model.ExternalType{
								Type: "Context",
							},
						},
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
				Content: func() (content string, requiredPkg []*model.GoPkg) {
					return fmt.Sprintf(
						"return %s.%s.Commit().Error", b.Domain.GormTransaction.GetMethodName(), tx.Name,
					), nil
				},
			},
			{
				Name: "Rollback",
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PkgReference{
							Pkg: consts.CommonPkgs["context"],
							Reference: &model.ExternalType{
								Type: "Context",
							},
						},
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
				Content: func() (content string, requiredPkg []*model.GoPkg) {
					return fmt.Sprintf(
						"return %s.%s.Rollback().Error", b.Domain.GormTransaction.GetMethodName(), tx.Name,
					), nil
				},
			},
		},
	}

	// First build gorm models
	for _, m := range b.Domain.Models {
		gormModel := &model.GormModel{
			Struct: &model.Struct{
				Name: m.Struct.Name,
			},
		}

		gormModel.FromModel = &model.Function{
			Name: m.Struct.Name + "FromModel",
			Args: []*model.Param{
				{
					Name: "m",
					Type: &model.PointerType{
						Type: &model.PkgReference{
							Pkg:       b.Domain.Architecture.ModelPkg,
							Reference: m,
						},
					},
				},
			},
			Results: []*model.Param{
				{
					Type: &model.PointerType{
						Type: gormModel.Struct,
					},
				},
			},
			Content: func() (content string, requiredPkg []*model.GoPkg) {
				return fmt.Sprintf(
					"return &%s{}", gormModel.Struct.Name,
				), nil
			},
		}

		gormModel.ToModel = &model.Function{
			Name: m.Struct.Name + "ToModel",
			Args: []*model.Param{
				{
					Name: "m",
					Type: &model.PointerType{
						Type: gormModel.Struct,
					},
				},
			},
			Results: []*model.Param{
				{
					Type: &model.PointerType{
						Type: &model.PkgReference{
							Pkg:       b.Domain.Architecture.ModelPkg,
							Reference: m,
						},
					},
				},
			},
			Content: func() (content string, requiredPkg []*model.GoPkg) {
				return fmt.Sprintf(
					"return &%s.%s{}", b.Domain.Architecture.ModelPkg.Alias, gormModel.Struct.Name,
				), nil
			},
		}

		gormModel.FromModels = &model.Function{
			Name: m.Struct.Name + "FromModels",
			Args: []*model.Param{
				{
					Name: "models",
					Type: &model.ArrayType{
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Pkg:       b.Domain.Architecture.ModelPkg,
								Reference: m,
							},
						},
					},
				},
			},
			Results: []*model.Param{
				{
					Type: &model.ArrayType{
						Type: &model.PointerType{
							Type: gormModel.Struct,
						},
					},
				},
			},
			Content: func() (content string, requiredPkg []*model.GoPkg) {
				content = fmt.Sprintf("entities := []*%s{}", gormModel.Struct.Name) + consts.LN
				content += fmt.Sprintf("for _, m := range %s {", gormModel.FromModels.Args[0].Name) + consts.LN
				content += fmt.Sprintf(
					"entities = append(entities, %s(m))",
					gormModel.FromModel.Name,
				) + consts.LN
				content += "}" + consts.LN
				content += "return entities"
				return content, nil
			},
		}

		gormModel.ToModels = &model.Function{
			Name: m.Struct.Name + "ToModels",
			Args: []*model.Param{
				{
					Name: "m",
					Type: &model.ArrayType{
						Type: &model.PointerType{
							Type: gormModel.Struct,
						},
					},
				},
			},
			Results: []*model.Param{
				{
					Type: &model.ArrayType{
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Pkg:       b.Domain.Architecture.ModelPkg,
								Reference: m,
							},
						},
					},
				},
			},
			Content: func() (content string, requiredPkg []*model.GoPkg) {
				content = fmt.Sprintf("entities := []*%s.%s{}", b.Domain.Architecture.ModelPkg.Alias, m.Struct.Name) + consts.LN
				content += fmt.Sprintf("for _, m := range %s {", gormModel.ToModels.Args[0].Name) + consts.LN
				content += fmt.Sprintf(
					"entities = append(entities, %s(m))",
					gormModel.ToModel.Name,
				) + consts.LN
				content += "}" + consts.LN
				content += "return entities"
				return content, nil
			},
		}

		for _, f := range m.Struct.Fields {
			f = f.Copy()
			f.Tags = nil
			if f.Name == "DeletedAt" {
				f.Type = &model.PkgReference{
					Pkg: consts.CommonPkgs["gorm"],
					Reference: &model.ExternalType{
						Type: "DeletedAt",
					},
				}
			}
			gormModel.Struct.Fields = append(gormModel.Struct.Fields, f)
		}

		b.Domain.GormModels = append(b.Domain.GormModels, gormModel)
		b.ModelToGormModel[m] = gormModel
	}

	for _, m := range b.Domain.Models {
		gormModel := b.ModelToGormModel[m]
		for _, r := range m.Relations {
			if r.Type == model.RelationSingleMandatory || r.Type == model.RelationSingleOptionnal {
				f := &model.Field{
					Name: r.On.Struct.Name + "Id",
					Type: model.PrimitiveTypeString,
				}
				if r.Type == model.RelationSingleOptionnal {
					f.Type = &model.PointerType{
						Type: f.Type,
					}
				}
				gormModel.Struct.Fields = append(gormModel.Struct.Fields, f)
				f = &model.Field{
					Name: r.On.Struct.Name,
					Type: &model.PointerType{
						Type: b.ModelToGormModel[r.On].Struct,
					},
				}
				gormModel.Struct.Fields = append(gormModel.Struct.Fields, f)
			} else {
				relationName := r.On.Struct.Name + "s"
				if strings.HasSuffix(relationName, "ys") {
					relationName = relationName[:len(relationName)-2] + "ies"
				}
				f := &model.Field{
					Name: relationName,
					Type: &model.ArrayType{
						Type: &model.PointerType{
							Type: b.ModelToGormModel[r.On].Struct,
						},
					},
				}
				gormModel.Struct.Fields = append(gormModel.Struct.Fields, f)
			}
		}

		gormModel.FromModel.Content = func() (content string, requiredPkg []*model.GoPkg) {
			pkgs := []*model.GoPkg{}
			arg := gormModel.FromModel.Args[0].Name
			content = fmt.Sprintf("if %s == nil {", gormModel.FromModel.Args[0].Name) + consts.LN
			content += "return nil" + consts.LN
			content += "}" + consts.LN
			content += fmt.Sprintf("entity := &%s{}", gormModel.Struct.Name) + consts.LN
			for _, f := range m.Struct.Fields {
				if f.Name == "DeletedAt" {
					pkgs = append(pkgs, consts.CommonPkgs["gorm"])
					content += fmt.Sprintf("entity.DeletedAt = %s.DeletedAt{", consts.CommonPkgs["gorm"].Alias) + consts.LN
					content += fmt.Sprintf("Time: %s.DeletedAt,", arg) + consts.LN
					content += fmt.Sprintf(`Valid:  %s.DeletedAt.String() != "0000-00-00 00:00:00",`, arg) + consts.LN
					content += "}" + consts.LN
				} else {
					content += fmt.Sprintf("entity.%s = %s.%s", f.Name, arg, f.Name) + consts.LN
				}
			}
			for _, r := range m.Relations {
				if r.Type == model.RelationSingleMandatory || r.Type == model.RelationSingleOptionnal {
					if r.Type == model.RelationSingleMandatory {
						content += fmt.Sprintf("entity.%sId = %s.%sId", r.On.Struct.Name, arg, r.On.Struct.Name) + consts.LN
					} else if r.Type == model.RelationSingleOptionnal {
						content += fmt.Sprintf("entity.%sId = &%s.%sId", r.On.Struct.Name, arg, r.On.Struct.Name) + consts.LN
					}
					content += fmt.Sprintf("entity.%s = %s(%s.%s)", r.On.Struct.Name, b.ModelToGormModel[r.On].FromModel.Name, arg, r.On.Struct.Name) + consts.LN
				} else {
					relationName := r.On.Struct.Name + "s"
					if strings.HasSuffix(relationName, "ys") {
						relationName = relationName[:len(relationName)-2] + "ies"
					}
					content += fmt.Sprintf("entity.%s = %s(%s.%s)", relationName, b.ModelToGormModel[r.On].FromModels.Name, arg, relationName) + consts.LN
				}
			}
			content += "return entity"
			return content, pkgs
		}
		gormModel.ToModel.Content = func() (content string, requiredPkg []*model.GoPkg) {
			pkgs := []*model.GoPkg{}
			arg := gormModel.ToModel.Args[0].Name
			content = fmt.Sprintf("if %s == nil {", gormModel.FromModel.Args[0].Name) + consts.LN
			content += "return nil" + consts.LN
			content += "}" + consts.LN
			content += fmt.Sprintf("entity := &%s.%s{}", b.Domain.Architecture.ModelPkg.Alias, m.Struct.Name) + consts.LN
			for _, f := range m.Struct.Fields {
				if f.Name == "DeletedAt" {
					pkgs = append(pkgs, consts.CommonPkgs["gorm"])
					content += fmt.Sprintf("entity.DeletedAt = %s.DeletedAt.Time", arg) + consts.LN
				} else {
					content += fmt.Sprintf("entity.%s = %s.%s", f.Name, arg, f.Name) + consts.LN
				}
			}
			for _, r := range m.Relations {
				if r.Type == model.RelationSingleMandatory || r.Type == model.RelationSingleOptionnal {
					if r.Type == model.RelationSingleMandatory {
						content += fmt.Sprintf("entity.%sId = %s.%sId", r.On.Struct.Name, arg, r.On.Struct.Name) + consts.LN
					} else if r.Type == model.RelationSingleOptionnal {
						content += fmt.Sprintf("entity.%sId = *%s.%sId", r.On.Struct.Name, arg, r.On.Struct.Name) + consts.LN
					}
					content += fmt.Sprintf("entity.%s = %s(%s.%s)", r.On.Struct.Name, b.ModelToGormModel[r.On].ToModel.Name, arg, r.On.Struct.Name) + consts.LN
				} else {
					relationName := r.On.Struct.Name + "s"
					if strings.HasSuffix(relationName, "ys") {
						relationName = relationName[:len(relationName)-2] + "ies"
					}
					content += fmt.Sprintf("entity.%s = %s(%s.%s)", relationName, b.ModelToGormModel[r.On].ToModels.Name, arg, relationName) + consts.LN
				}
			}
			content += "return entity"
			return content, pkgs
		}
	}

	db := &model.Field{
		Name: "DB",
		Type: &model.PointerType{
			Type: &model.PkgReference{
				Pkg: consts.CommonPkgs["gorm"],
				Reference: &model.ExternalType{
					Type: "DB",
				},
			},
		},
	}
	b.Domain.GormDomainRepository = &model.Struct{
		Name: stringtool.UpperFirstLetter(b.Domain.Name) + "DomainRepository",
		Fields: []*model.Field{
			db,
		},
		Methods: []*model.Function{
			{
				Name: "Migrate",
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PkgReference{
							Pkg: consts.CommonPkgs["context"],
							Reference: &model.ExternalType{
								Type: "Context",
							},
						},
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
				Content: func() (string, []*model.GoPkg) {
					pkgs := []*model.GoPkg{consts.CommonPkgs["gorm"], consts.CommonPkgs["context"]}
					content := fmt.Sprintf("return %s.%s.AutoMigrate(\nctx, ", b.Domain.GormDomainRepository.GetMethodName(), db.Name) + consts.LN
					for _, m := range b.Domain.GormModels {
						content += fmt.Sprintf("&%s{}, ", m.Struct.Name) + consts.LN
					}
					content += ")"

					return content, pkgs
				},
			},
			{
				Name: "BeginTransaction",
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PkgReference{
							Pkg: consts.CommonPkgs["context"],
							Reference: &model.ExternalType{
								Type: "Context",
							},
						},
					},
				},
				Results: []*model.Param{
					{
						Type: &model.PkgReference{
							Pkg:       b.Domain.Architecture.RepositoryPkg,
							Reference: b.Domain.RepositoryTransaction,
						},
					},
				},
				Content: func() (string, []*model.GoPkg) {
					pkgs := []*model.GoPkg{}
					content := fmt.Sprintf("return &%s{", b.Domain.GormTransaction.Name) + consts.LN
					content += fmt.Sprintf("%s: %s.%s.Begin(),", tx.Name, b.Domain.GormDomainRepository.GetMethodName(), db.Name) + consts.LN
					content += "}"
					return content, pkgs
				},
			},
		},
	}

	for _, repoDefinition := range b.RepositoryDefinitionsToBuild {
		on, err := b.GetModel(ctx, repoDefinition.On)
		if err != nil {
			b.Err = merror.Stack(err)
			return b
		}

		repo := b.RepositoryDefinitionToRepository[repoDefinition]
		gormModel := b.ModelToGormModel[on]

		if repoDefinition.GetMethod {
			var function *model.Function
			for _, f := range b.Domain.DomainRepository.Methods {
				if f.Name == fmt.Sprintf("Get%s", on.Struct.Name) {
					function = f
					break
				}
			}

			function = function.Copy().(*model.Function)
			function.Content = func() (string, []*model.GoPkg) {
				pkgs := []*model.GoPkg{
					b.Domain.Architecture.RepositoryPkg,
					consts.CommonPkgs["gorm"],
					consts.CommonPkgs["errors"],
					consts.CommonPkgs["merror"],
				}
				// Init context
				content := fmt.Sprintf("methodCtx := &%s.%sContext{}", b.Domain.Architecture.RepositoryPkg.Alias, function.Name) + consts.LN
				content += fmt.Sprintf("for _, opt := range %s {", function.Args[len(function.Args)-1].Name) + consts.LN
				content += "opt(methodCtx)" + consts.LN
				content += "}" + consts.LN

				content += fmt.Sprintf("var db *%s.DB", consts.CommonPkgs["gorm"].Alias) + consts.LN
				content += "if methodCtx.Transaction != nil {" + consts.LN
				content += fmt.Sprintf("tx, ok := methodCtx.Transaction.Get(ctx).(*%s.DB)", consts.CommonPkgs["gorm"].Alias) + consts.LN
				content += fmt.Sprintf(`if !ok { return nil, %s.New("expected transaction to be *gorm.DB") }`, consts.CommonPkgs["errors"].Alias) + consts.LN
				content += "db = tx" + consts.LN
				content += "} else { " + consts.LN
				content += fmt.Sprintf("db = %s.%s", b.Domain.GormDomainRepository.GetMethodName(), db.Name) + consts.LN
				content += "}" + consts.LN

				activeConditions := ""
				content += fmt.Sprintf(`request := db.Table("%s")`, repo.TableName) + consts.LN
				content += consts.LN

				joins := ""
				activableModels := map[*model.Model]bool{}
				if b.isRelationUpTreeActivable(on, activableModels, []*model.Relation{}, &joins, &pkgs) {
					for activableModel := range activableModels {
						if activableModel.Activable {
							activeConditions += fmt.Sprintf(`request.Where("%s.active = ?", true)`, b.ModelToRepository[activableModel].TableName) + consts.LN
						}
					}
				}

				if activeConditions != "" {
					content += "if !methodCtx.RetriveInactive { " + consts.LN
					content += joins
					content += activeConditions
					content += "}" + consts.LN
				}
				content += consts.LN

				content += fmt.Sprintf("entity := &%s{}", gormModel.Struct.Name) + consts.LN

				content += "if err := request.First(entity).Error; err != nil { " + consts.LN
				content += fmt.Sprintf(`if err == %s.ErrRecordNotFound { return nil, %s.%s }`, consts.CommonPkgs["gorm"].Alias, b.Domain.Architecture.RepositoryPkg.Alias, b.RepositoryErrors["notFound"].Name) + consts.LN
				content += fmt.Sprintf(`return nil, %s.Stack(err)`, consts.CommonPkgs["merror"].Alias) + consts.LN
				content += "}" + consts.LN
				content += consts.LN

				content += fmt.Sprintf("return %s(entity), nil", gormModel.ToModel.Name) + consts.LN
				return content, pkgs
			}
			b.Domain.GormDomainRepository.Methods = append(b.Domain.GormDomainRepository.Methods, function)
		}

		if repoDefinition.ListMethod {
			var function *model.Function
			for _, f := range b.Domain.DomainRepository.Methods {
				if f.Name == fmt.Sprintf("List%s", on.Struct.Name) {
					function = f
					break
				}
			}

			function = function.Copy().(*model.Function)
			function.Content = func() (string, []*model.GoPkg) {
				pkgs := []*model.GoPkg{
					b.Domain.Architecture.RepositoryPkg,
					consts.CommonPkgs["gorm"],
					consts.CommonPkgs["errors"],
					consts.CommonPkgs["merror"],
				}
				// Init context
				content := fmt.Sprintf("methodCtx := &%s.%sContext{}", b.Domain.Architecture.RepositoryPkg.Alias, function.Name) + consts.LN
				content += fmt.Sprintf("for _, opt := range %s {", function.Args[len(function.Args)-1].Name) + consts.LN
				content += "opt(methodCtx)" + consts.LN
				content += "}" + consts.LN

				content += fmt.Sprintf("var db *%s.DB", consts.CommonPkgs["gorm"].Alias) + consts.LN
				content += "if methodCtx.Transaction != nil {" + consts.LN
				content += fmt.Sprintf("tx, ok := methodCtx.Transaction.Get(ctx).(*%s.DB)", consts.CommonPkgs["gorm"].Alias) + consts.LN
				content += fmt.Sprintf(`if !ok { return nil, %s.New("expected transaction to be *gorm.DB") }`, consts.CommonPkgs["errors"].Alias) + consts.LN
				content += "db = tx" + consts.LN
				content += "} else { " + consts.LN
				content += fmt.Sprintf("db = %s.%s", b.Domain.GormDomainRepository.GetMethodName(), db.Name) + consts.LN
				content += "}" + consts.LN

				activeConditions := ""
				content += fmt.Sprintf(`request := db.Table("%s")`, repo.TableName) + consts.LN
				content += consts.LN

				joins := ""
				activableModels := map[*model.Model]bool{}
				if b.isRelationUpTreeActivable(on, activableModels, []*model.Relation{}, &joins, &pkgs) {
					for activableModel := range activableModels {
						if activableModel.Activable {
							activeConditions += fmt.Sprintf(`request.Where("%s.active = ?", true)`, b.ModelToRepository[activableModel].TableName) + consts.LN
						}
					}
				}

				if activeConditions != "" {
					content += "if !methodCtx.RetriveInactive { " + consts.LN
					content += joins
					content += activeConditions
					content += "}" + consts.LN
				}
				content += consts.LN

				content += "request.Limit(int(methodCtx.Pagination.GetItemsPerPage()))" + consts.LN
				content += "request.Offset(int((methodCtx.Pagination.GetPage() - 1) * methodCtx.Pagination.GetItemsPerPage()))" + consts.LN
				content += consts.LN

				content += fmt.Sprintf("entity := []*%s{}", gormModel.Struct.Name) + consts.LN

				content += "if err := request.Find(entity).Error; err != nil { " + consts.LN
				content += fmt.Sprintf(`if err == %s.ErrRecordNotFound { return nil, %s.%s }`, consts.CommonPkgs["gorm"].Alias, b.Domain.Architecture.RepositoryPkg.Alias, b.RepositoryErrors["notFound"].Name) + consts.LN
				content += fmt.Sprintf(`return nil, %s.Stack(err)`, consts.CommonPkgs["merror"].Alias) + consts.LN
				content += "}" + consts.LN
				content += consts.LN

				content += fmt.Sprintf("return %s(entity), nil", gormModel.ToModels.Name) + consts.LN
				return content, pkgs
			}
			b.Domain.GormDomainRepository.Methods = append(b.Domain.GormDomainRepository.Methods, function)
		}
		if repoDefinition.CreateMethod {
			var function *model.Function
			for _, f := range b.Domain.DomainRepository.Methods {
				if f.Name == fmt.Sprintf("Create%s", on.Struct.Name) {
					function = f
					break
				}
			}

			function = function.Copy().(*model.Function)
			function.Content = func() (string, []*model.GoPkg) {
				pkgs := []*model.GoPkg{
					b.Domain.Architecture.RepositoryPkg,
					consts.CommonPkgs["gorm"],
					consts.CommonPkgs["errors"],
					consts.CommonPkgs["merror"],
				}
				return "return nil, nil", pkgs
			}
			b.Domain.GormDomainRepository.Methods = append(b.Domain.GormDomainRepository.Methods, function)
		}
	}

	implementedFunctions := []string{}
	for _, f := range b.Domain.GormDomainRepository.Methods {
		implementedFunctions = append(implementedFunctions, f.Name)
	}
	for _, f := range b.Domain.DomainRepository.Methods {
		if !slices.Contains(implementedFunctions, f.Name) {
			f = f.Copy().(*model.Function)
			f.Content = func() (string, []*model.GoPkg) {
				return `panic("not implemented")`, []*model.GoPkg{}
			}
			b.Domain.GormDomainRepository.Methods = append(b.Domain.GormDomainRepository.Methods, f)
		}
	}

	return b
}
