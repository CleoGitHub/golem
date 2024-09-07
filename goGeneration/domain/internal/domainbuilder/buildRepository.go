package domainbuilder

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (b *domainBuilder) buildRepository(ctx context.Context, repositoryDefinition *coredomaindefinition.Repository) *domainBuilder {
	if b.Err != nil {
		return b
	}

	on, err := b.GetModel(ctx, repositoryDefinition.On)
	if err != nil {
		b.Err = merror.Stack(err)
		return b
	}

	repo := &model.Repository{
		On:        on,
		Name:      fmt.Sprintf("%sRepository", on.Struct.Name),
		TableName: repositoryDefinition.TableName,
		AllowedOrderBys: model.Consts{
			Name:   strings.ToUpper(stringtool.SnakeCase(repositoryDefinition.On.Name)) + "_ALLOWED_ORDER_BY",
			Values: []interface{}{},
		},
		AllowedWheres: model.Consts{
			Name:   strings.ToUpper(stringtool.SnakeCase(repositoryDefinition.On.Name)) + "_ALLOWED_WHERE",
			Values: []interface{}{},
		},
		FieldToColumn: model.Map{
			Name: strings.ToUpper(stringtool.SnakeCase(repositoryDefinition.On.Name)) + "_FIELD_TO_COLUMN",
			Type: model.MapType{
				Key:   model.PrimitiveTypeString,
				Value: model.PrimitiveTypeString,
			},
			Values: []model.MapValue{},
		},
	}

	transactionField := &model.Field{
		Name: "Transaction",
		Type: &model.PkgReference{
			Pkg:       b.Domain.Architecture.RepositoryPkg,
			Reference: b.GetTransation(ctx),
		},
	}

	if repo.DefaultOrderBy == "" {
		repo.DefaultOrderBy = "UpdatedAt"
	}

	for _, defaultField := range consts.DefaultModelFields {
		repo.FieldToColumn.Values = append(repo.FieldToColumn.Values, model.MapValue{
			Key:   GetFieldName(ctx, defaultField),
			Value: stringtool.SnakeCase(defaultField.Name),
		})
		repo.AllowedWheres.Values = append(repo.AllowedWheres.Values, GetFieldName(ctx, defaultField))
		repo.AllowedOrderBys.Values = append(repo.AllowedOrderBys.Values, GetFieldName(ctx, defaultField))
	}

	for _, field := range repositoryDefinition.On.Fields {
		repo.FieldToColumn.Values = append(repo.FieldToColumn.Values, model.MapValue{
			Key:   GetFieldName(ctx, field),
			Value: stringtool.SnakeCase(field.Name),
		})
		repo.AllowedWheres.Values = append(repo.AllowedWheres.Values, GetFieldName(ctx, field))
		repo.AllowedOrderBys.Values = append(repo.AllowedOrderBys.Values, GetFieldName(ctx, field))
	}

	for _, relation := range b.RelationDefinitionsToBuild {
		if relation.Source == repositoryDefinition.On {
			if slices.Contains(
				[]coredomaindefinition.RelationType{
					coredomaindefinition.RelationTypeBelongsTo,
					coredomaindefinition.RelationTypeManyToOne,
					coredomaindefinition.RelationTypeOneToOne,
				},
				relation.Type,
			) {
				k := GetSingleRelationName(ctx, relation.Target)
				repo.FieldToColumn.Values = append(repo.FieldToColumn.Values, model.MapValue{
					Key:   k,
					Value: stringtool.SnakeCase(k),
				})
				repo.AllowedWheres.Values = append(repo.AllowedWheres.Values, GetSingleRelationName(ctx, relation.Target))
				repo.AllowedOrderBys.Values = append(repo.AllowedOrderBys.Values, GetSingleRelationName(ctx, relation.Target))
			}
		}

		if relation.Target == repositoryDefinition.On {
			if slices.Contains(
				[]coredomaindefinition.RelationType{
					coredomaindefinition.RelationTypeOneToMany,
					coredomaindefinition.RelationTypeOneToOne,
				},
				relation.Type,
			) {
				k := GetSingleRelationName(ctx, relation.Source)
				repo.FieldToColumn.Values = append(repo.FieldToColumn.Values, model.MapValue{
					Key:   k,
					Value: stringtool.SnakeCase(k),
				})
				repo.AllowedWheres.Values = append(repo.AllowedWheres.Values, GetSingleRelationName(ctx, relation.Source))
				repo.AllowedOrderBys.Values = append(repo.AllowedOrderBys.Values, GetSingleRelationName(ctx, relation.Source))
			}
		}
	}

	repo.Functions = append(repo.Functions, consts.DefaultRepositoryMethods...)

	modelPkgReference := &model.PointerType{
		Type: &model.PkgReference{
			Pkg:       b.Domain.Architecture.ModelPkg,
			Reference: on,
		},
	}

	if repositoryDefinition.GetMethod {
		method := &model.Function{
			Name: fmt.Sprintf("Get%s", on.Struct.Name),
			Results: []*model.Param{
				{
					Type: modelPkgReference,
				},
				{
					Type: model.PrimitiveTypeError,
				},
			},
		}
		repoMethod := &model.RepositoryMethod{}
		// Create method Context
		ctxMethod := &model.Struct{
			Name:   fmt.Sprintf("Get%sContext", GetModelName(ctx, repositoryDefinition.On)),
			Fields: []*model.Field{},
		}

		// Add active field
		ctxMethod.Fields = append(ctxMethod.Fields, &model.Field{
			Name: "RetriveInactive",
			Type: model.PrimitiveTypeBool,
		})

		// Add by field
		ctxMethod.Fields = append(ctxMethod.Fields, &model.Field{
			Name: "By",
			Type: &model.MapType{
				Key:   model.PrimitiveTypeString,
				Value: model.PrimitiveTypeInterface,
			},
		})

		// Add not field
		ctxMethod.Fields = append(ctxMethod.Fields, &model.Field{
			Name: "Not",
			Type: &model.MapType{
				Key:   model.PrimitiveTypeString,
				Value: model.PrimitiveTypeInterface,
			},
		})
		ctxMethod.Fields = append(ctxMethod.Fields, transactionField)
		// ctxMethod.Fields = append(ctxMethod.Fields, withableFields...)

		repoMethod.Context = ctxMethod

		repoMethod.Opt = &model.TypeDefinition{
			Name: fmt.Sprintf("Get%sOpt", GetModelName(ctx, repositoryDefinition.On)),
			Type: &model.Function{
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Pkg:       b.Domain.Architecture.RepositoryPkg,
								Reference: ctxMethod,
							},
						},
					},
				},
			},
		}

		for _, field := range ctxMethod.Fields {
			opt := &model.Function{
				Name: fmt.Sprintf("%sWith%s", method.Name, field.Name),
				Args: []*model.Param{
					{
						Name: stringtool.LowerFirstLetter(field.Name),
						Type: field.Type,
					},
				},
				Results: []*model.Param{
					{
						Type: &model.PkgReference{
							Pkg:       b.Domain.Architecture.RepositoryPkg,
							Reference: repoMethod.Opt,
						},
					},
				},
			}
			opt.Content = func() (string, []*model.GoPkg) {
				str := fmt.Sprintf("return func(ctx *%s) {", ctxMethod.Name)
				str += fmt.Sprintf(" ctx.%s = %s", field.Name, stringtool.LowerFirstLetter(field.Name))
				str += " }"
				return str, nil
			}
			repoMethod.Opts = append(repoMethod.Opts, opt)
		}
		method.Args = append(method.Args, &model.Param{
			Name: "ctx",
			Type: &model.PkgReference{
				Pkg: consts.CommonPkgs["context"],
				Reference: &model.ExternalType{
					Type: "Context",
				},
			},
		})
		method.Args = append(method.Args, &model.Param{
			Name: "opts",
			Type: &model.VariaidicType{
				Type: &model.PkgReference{
					Pkg:       b.Domain.Architecture.RepositoryPkg,
					Reference: repoMethod.Opt,
				},
			},
		})
		// repoMethod.Function = method
		b.Domain.DomainRepository.Methods = append(b.Domain.DomainRepository.Methods, method)

		repo.Methods = append(repo.Methods, repoMethod)
	}

	if repositoryDefinition.ListMethod {
		method := &model.Function{
			Name: fmt.Sprintf("List%s", on.Struct.Name),
			Results: []*model.Param{
				{
					Type: &model.ArrayType{
						Type: modelPkgReference,
					},
				},
				{
					Type: model.PrimitiveTypeError,
				},
			},
		}
		repoMethod := &model.RepositoryMethod{}
		// Create method Context
		ctxMethod := &model.Struct{
			Name:   fmt.Sprintf("List%sContext", GetModelName(ctx, repositoryDefinition.On)),
			Fields: []*model.Field{},
		}

		// Add active field
		ctxMethod.Fields = append(ctxMethod.Fields, &model.Field{
			Name: "RetriveInactive",
			Type: model.PrimitiveTypeBool,
		})

		ctxMethod.Fields = append(ctxMethod.Fields, &model.Field{
			Name: "Pagination",
			Type: &model.PkgReference{
				Pkg:       b.Domain.Architecture.RepositoryPkg,
				Reference: b.GetPagination(ctx),
			},
		})

		// Add by field
		ctxMethod.Fields = append(ctxMethod.Fields, &model.Field{
			Name: "Ordering",
			Type: &model.PkgReference{
				Pkg:       b.Domain.Architecture.RepositoryPkg,
				Reference: b.GetOrdering(ctx),
			},
		})

		// Add by field
		ctxMethod.Fields = append(ctxMethod.Fields, &model.Field{
			Name: "By",
			Type: &model.MapType{
				Key:   model.PrimitiveTypeString,
				Value: model.PrimitiveTypeInterface,
			},
		})

		// Add not field
		ctxMethod.Fields = append(ctxMethod.Fields, &model.Field{
			Name: "Not",
			Type: &model.MapType{
				Key:   model.PrimitiveTypeString,
				Value: model.PrimitiveTypeInterface,
			},
		})

		// ctxMethod.Fields = append(ctxMethod.Fields, &model.Field{
		// 	Name: "OrderBy",
		// 	Type: model.PrimitiveTypeString,
		// })

		// ctxMethod.Fields = append(ctxMethod.Fields, &model.Field{
		// 	Name: "Order",
		// 	Type: model.PrimitiveTypeString,
		// })

		ctxMethod.Fields = append(ctxMethod.Fields, transactionField)
		// ctxMethod.Fields = append(ctxMethod.Fields, withableFields...)

		repoMethod.Context = ctxMethod

		repoMethod.Opt = &model.TypeDefinition{
			Name: fmt.Sprintf("List%sOpt", GetModelName(ctx, repositoryDefinition.On)),
			Type: &model.Function{
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Pkg: b.Domain.Architecture.RepositoryPkg,
								Reference: &model.ExternalType{
									Type: ctxMethod.Name,
								},
							},
						},
					},
				},
			},
		}

		for _, field := range ctxMethod.Fields {
			opt := &model.Function{
				Name: fmt.Sprintf("%sWith%s", method.Name, field.Name),
				Args: []*model.Param{
					{
						Name: stringtool.LowerFirstLetter(field.Name),
						Type: field.Type,
					},
				},
				Results: []*model.Param{
					{
						Type: &model.PkgReference{
							Pkg:       b.Domain.Architecture.RepositoryPkg,
							Reference: repoMethod.Opt,
						},
					},
				},
			}
			opt.Content = func() (string, []*model.GoPkg) {
				str := fmt.Sprintf("return func(ctx *%s) {", ctxMethod.Name)
				str += fmt.Sprintf(" ctx.%s = %s", field.Name, stringtool.LowerFirstLetter(field.Name))
				str += " }"
				return str, nil
			}
			// opt.Content = fmt.Sprintf("return func(ctx *%s) {", ctxMethod.Name)
			// opt.Content += fmt.Sprintf(" ctx.%s = %s", field.Name, stringtool.LowerFirstLetter(field.Name))
			// opt.Content += " }"
			repoMethod.Opts = append(repoMethod.Opts, opt)
		}
		method.Args = append(method.Args, &model.Param{
			Name: "ctx",
			Type: &model.PkgReference{
				Pkg: consts.CommonPkgs["context"],
				Reference: &model.ExternalType{
					Type: "Context",
				},
			},
		})
		method.Args = append(method.Args, &model.Param{
			Name: "opts",
			Type: &model.VariaidicType{
				Type: &model.PkgReference{
					Pkg:       b.Domain.Architecture.RepositoryPkg,
					Reference: repoMethod.Opt,
				},
			},
		})
		// repoMethod.Function = method
		b.Domain.DomainRepository.Methods = append(b.Domain.DomainRepository.Methods, method)

		repo.Methods = append(repo.Methods, repoMethod)
	}

	if repositoryDefinition.CreateMethod {
		method := &model.Function{
			Name: fmt.Sprintf("Create%s", on.Struct.Name),
			Results: []*model.Param{
				{
					Type: modelPkgReference,
				},
				{
					Type: model.PrimitiveTypeError,
				},
			},
		}
		repoMethod := &model.RepositoryMethod{}
		// Create method Context
		ctxMethod := &model.Struct{
			Name:   fmt.Sprintf("Create%sContext", GetModelName(ctx, repositoryDefinition.On)),
			Fields: []*model.Field{},
		}

		ctxMethod.Fields = append(ctxMethod.Fields, transactionField)

		repoMethod.Context = ctxMethod

		repoMethod.Opt = &model.TypeDefinition{
			Name: fmt.Sprintf("Create%sOpt", GetModelName(ctx, repositoryDefinition.On)),
			Type: &model.Function{
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Pkg: b.Domain.Architecture.RepositoryPkg,
								Reference: &model.ExternalType{
									Type: ctxMethod.Name,
								},
							},
						},
					},
				},
			},
		}

		for _, field := range ctxMethod.Fields {
			opt := &model.Function{
				Name: fmt.Sprintf("%sWith%s", method.Name, field.Name),
				Args: []*model.Param{
					{
						Name: stringtool.LowerFirstLetter(field.Name),
						Type: field.Type,
					},
				},
				Results: []*model.Param{
					{
						Type: &model.PkgReference{
							Pkg:       b.Domain.Architecture.RepositoryPkg,
							Reference: repoMethod.Opt,
						},
					},
				},
			}
			opt.Content = func() (string, []*model.GoPkg) {
				str := fmt.Sprintf("return func(ctx *%s) {", ctxMethod.Name)
				str += fmt.Sprintf(" ctx.%s = %s", field.Name, stringtool.LowerFirstLetter(field.Name))
				str += " }"
				return str, nil
			}
			repoMethod.Opts = append(repoMethod.Opts, opt)
		}
		method.Args = append(method.Args, &model.Param{
			Name: "ctx",
			Type: &model.PkgReference{
				Pkg: consts.CommonPkgs["context"],
				Reference: &model.ExternalType{
					Type: "Context",
				},
			},
		})
		// entity to create
		method.Args = append(method.Args, &model.Param{
			Name: "entity",
			Type: modelPkgReference,
		})

		method.Args = append(method.Args, &model.Param{
			Name: "opts",
			Type: &model.VariaidicType{
				Type: &model.PkgReference{
					Pkg:       b.Domain.Architecture.RepositoryPkg,
					Reference: repoMethod.Opt,
				},
			},
		})
		// repoMethod.Function = method
		b.Domain.DomainRepository.Methods = append(b.Domain.DomainRepository.Methods, method)

		repo.Methods = append(repo.Methods, repoMethod)
	}

	if repositoryDefinition.UpdateMethod {
		method := &model.Function{
			Name: fmt.Sprintf("Update%s", on.Struct.Name),
			Results: []*model.Param{
				{
					Type: modelPkgReference,
				},
				{
					Type: model.PrimitiveTypeError,
				},
			},
		}
		repoMethod := &model.RepositoryMethod{}
		// Update method Context
		ctxMethod := &model.Struct{
			Name:   fmt.Sprintf("Update%sContext", GetModelName(ctx, repositoryDefinition.On)),
			Fields: []*model.Field{},
		}

		ctxMethod.Fields = append(ctxMethod.Fields, transactionField)

		repoMethod.Context = ctxMethod

		repoMethod.Opt = &model.TypeDefinition{
			Name: fmt.Sprintf("Update%sOpt", GetModelName(ctx, repositoryDefinition.On)),
			Type: &model.Function{
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Pkg: b.Domain.Architecture.RepositoryPkg,
								Reference: &model.ExternalType{
									Type: ctxMethod.Name,
								},
							},
						},
					},
				},
			},
		}

		for _, field := range ctxMethod.Fields {
			opt := &model.Function{
				Name: fmt.Sprintf("%sWith%s", method.Name, field.Name),
				Args: []*model.Param{
					{
						Name: stringtool.LowerFirstLetter(field.Name),
						Type: field.Type,
					},
				},
				Results: []*model.Param{
					{
						Type: &model.PkgReference{
							Pkg:       b.Domain.Architecture.RepositoryPkg,
							Reference: repoMethod.Opt,
						},
					},
				},
			}
			opt.Content = func() (string, []*model.GoPkg) {
				str := fmt.Sprintf("return func(ctx *%s) {", ctxMethod.Name)
				str += fmt.Sprintf(" ctx.%s = %s", field.Name, stringtool.LowerFirstLetter(field.Name))
				str += " }"
				return str, nil
			}
			repoMethod.Opts = append(repoMethod.Opts, opt)
		}
		method.Args = append(method.Args, &model.Param{
			Name: "ctx",
			Type: &model.PkgReference{
				Pkg: consts.CommonPkgs["context"],
				Reference: &model.ExternalType{
					Type: "Context",
				},
			},
		})
		// entity to update
		method.Args = append(method.Args, &model.Param{
			Name: "entity",
			Type: modelPkgReference,
		})

		method.Args = append(method.Args, &model.Param{
			Name: "opts",
			Type: &model.VariaidicType{
				Type: &model.PkgReference{
					Pkg:       b.Domain.Architecture.RepositoryPkg,
					Reference: repoMethod.Opt,
				},
			},
		})
		// repoMethod.Function = method
		b.Domain.DomainRepository.Methods = append(b.Domain.DomainRepository.Methods, method)

		repo.Methods = append(repo.Methods, repoMethod)
	}

	if repositoryDefinition.DeleteMethod {
		method := &model.Function{
			Name: fmt.Sprintf("Delete%s", on.Struct.Name),
			Results: []*model.Param{
				{
					Type: model.PrimitiveTypeError,
				},
			},
		}
		repoMethod := &model.RepositoryMethod{}
		// Delete method Context
		ctxMethod := &model.Struct{
			Name:   fmt.Sprintf("Delete%sContext", GetModelName(ctx, repositoryDefinition.On)),
			Fields: []*model.Field{},
		}

		ctxMethod.Fields = append(ctxMethod.Fields, transactionField)

		repoMethod.Context = ctxMethod

		repoMethod.Opt = &model.TypeDefinition{
			Name: fmt.Sprintf("Delete%sOpt", GetModelName(ctx, repositoryDefinition.On)),
			Type: &model.Function{
				Args: []*model.Param{
					{
						Name: "ctx",
						Type: &model.PkgReference{
							Pkg: b.Domain.Architecture.RepositoryPkg,
							Reference: &model.PointerType{
								Type: &model.PkgReference{
									Pkg: b.Domain.Architecture.RepositoryPkg,
									Reference: &model.ExternalType{
										Type: ctxMethod.Name,
									},
								},
							},
						},
					},
				},
			},
		}

		for _, field := range ctxMethod.Fields {
			opt := &model.Function{
				Name: fmt.Sprintf("%sWith%s", method.Name, field.Name),
				Args: []*model.Param{
					{
						Name: stringtool.LowerFirstLetter(field.Name),
						Type: field.Type,
					},
				},
				Results: []*model.Param{
					{
						Type: &model.PkgReference{
							Pkg:       b.Domain.Architecture.RepositoryPkg,
							Reference: repoMethod.Opt,
						},
					},
				},
			}
			opt.Content = func() (string, []*model.GoPkg) {
				str := fmt.Sprintf("return func(ctx *%s) {", ctxMethod.Name)
				str += fmt.Sprintf(" ctx.%s = %s", field.Name, stringtool.LowerFirstLetter(field.Name))
				str += " }"
				return str, nil
			}
			repoMethod.Opts = append(repoMethod.Opts, opt)
		}
		method.Args = append(method.Args, &model.Param{
			Name: "ctx",
			Type: &model.PkgReference{
				Pkg: consts.CommonPkgs["context"],
				Reference: &model.ExternalType{
					Type: "Context",
				},
			},
		})
		// entity to Delete
		method.Args = append(method.Args, &model.Param{
			Name: "id",
			Type: model.PrimitiveTypeString,
		})

		method.Args = append(method.Args, &model.Param{
			Name: "opts",
			Type: &model.VariaidicType{
				Type: &model.PkgReference{
					Pkg:       b.Domain.Architecture.RepositoryPkg,
					Reference: repoMethod.Opt,
				},
			},
		})
		// repoMethod.Function = method
		b.Domain.DomainRepository.Methods = append(b.Domain.DomainRepository.Methods, method)

		repo.Methods = append(repo.Methods, repoMethod)
	}

	b.RepositoryDefinitionToRepository[repositoryDefinition] = repo
	b.ModelToRepository[on] = repo

	b.Repositories = append(b.Repositories, repo)

	b.Domain.Repositories = append(b.Domain.Repositories, repo)

	return b
}
