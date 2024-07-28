package domainbuilder

import (
	"context"
	"strings"

	"github.com/cleoGitHub/golem/coredomaindefinition"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
	"github.com/cleoGitHub/golem/pkg/merror"
)

func (b *domainBuilder) buildCRUD(ctx context.Context, crudDefinition *coredomaindefinition.CRUD) *domainBuilder {
	if b.Err != nil {
		return b
	}

	on, err := b.GetModel(ctx, crudDefinition.On)
	if err != nil {
		b.Err = merror.Stack(err)
		return b
	}

	if crudDefinition.Create != nil && crudDefinition.Create.Active {
		usecase := &model.Usecase{
			Roles: crudDefinition.Create.Roles,
			Function: &model.Function{
				Name: "Create" + on.Struct.Name + "Usecase",
			},
			Request: &model.Struct{
				Name: "Create" + on.Struct.Name + "Request",
			},
			Result: &model.Struct{
				Name: "Create" + on.Struct.Name + "Result",
				Fields: []*model.Field{
					{
						Name: on.Struct.Name,
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Reference: on,
								Pkg:       b.Domain.Architecture.ModelPkg,
							},
						},
						Tags: []*model.Tag{
							{
								Name:   "json",
								Values: []string{on.JsonName},
							},
						},
					},
				},
			},
		}

		for _, field := range b.ModelUsecaseStruct[on].Fields {
			usecase.Request.Fields = append(usecase.Request.Fields, field.Copy())
		}

		for _, relation := range on.Relations {
			relationDefinition := b.RelationToRelationDefinition[relation]
			if relationDefinition.Type == coredomaindefinition.RelationTypeSubresourcesOf &&
				b.ModelDefinitionToModel[relationDefinition.Target] == on {
				subresources := b.ModelDefinitionToModel[relationDefinition.Source]
				subresourcesRequest := &model.Struct{
					Name: on.Struct.Name + subresources.Struct.Name + "Request",
				}

				for _, field := range b.ModelUsecaseStruct[subresources].Fields {
					subresourcesRequest.Fields = append(subresourcesRequest.Fields, field.Copy())
				}
				b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, subresourcesRequest)

				subresourcesName := subresources.Struct.Name + "s"
				if strings.HasSuffix(subresourcesName, "ys") {
					subresourcesName = subresourcesName[:len(subresourcesName)-2] + "ies"
				}
				subresourcesJsonName := subresources.JsonName + "s"
				if strings.HasSuffix(subresourcesJsonName, "ys") {
					subresourcesJsonName = subresourcesJsonName[:len(subresourcesJsonName)-2] + "ies"
				}
				usecase.Request.Fields = append(usecase.Request.Fields, &model.Field{
					Name: subresourcesName,
					Type: &model.ArrayType{
						Type: subresourcesRequest,
					},
					Tags: []*model.Tag{
						{
							Name:   "json",
							Values: []string{subresourcesJsonName},
						},
					},
				})
			}
		}

		b.Domain.Usecases = append(b.Domain.Usecases, usecase)
		b.CRUDActionToUsecase[crudDefinition.Create] = usecase
		b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Request)
		b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Result)
	}

	if crudDefinition.Update != nil && crudDefinition.Update.Active {
		usecase := &model.Usecase{
			Roles: crudDefinition.Update.Roles,
			Function: &model.Function{
				Name: "Update" + on.Struct.Name + "Usecase",
			},
			Request: &model.Struct{
				Name: "Update" + on.Struct.Name + "Request",
			},
			Result: &model.Struct{
				Name: "Update" + on.Struct.Name + "Result",
				Fields: []*model.Field{
					{
						Name: on.Struct.Name,
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Reference: on,
								Pkg:       b.Domain.Architecture.ModelPkg,
							},
						},
						Tags: []*model.Tag{
							{
								Name:   "json",
								Values: []string{on.JsonName},
							},
						},
					},
				},
			},
		}

		for _, field := range b.ModelUsecaseStruct[on].Fields {
			usecase.Request.Fields = append(usecase.Request.Fields, field.Copy())
		}

		usecase.Request.Fields = append([]*model.Field{
			{
				Name: "Id",
				Type: model.PrimitiveTypeString,
				Tags: []*model.Tag{
					{
						Name:   "json",
						Values: []string{"id"},
					},
					{
						Name:   "validate",
						Values: []string{"required", "uuid"},
					},
				},
			},
		}, usecase.Request.Fields...)

		b.Domain.Usecases = append(b.Domain.Usecases, usecase)
		b.CRUDActionToUsecase[crudDefinition.Update] = usecase
		b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Request)
		b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Result)
	}

	if crudDefinition.Get != nil && crudDefinition.Get.Active {
		usecase := &model.Usecase{
			Roles: crudDefinition.Get.Roles,
			Function: &model.Function{
				Name: "Get" + on.Struct.Name + "Usecase",
			},
			Request: &model.Struct{
				Name: "Get" + on.Struct.Name + "Request",
			},
			Result: &model.Struct{
				Name: "Get" + on.Struct.Name + "Result",
				Fields: []*model.Field{
					{
						Name: on.Struct.Name,
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Reference: on,
								Pkg:       b.Domain.Architecture.ModelPkg,
							},
						},
						Tags: []*model.Tag{
							{
								Name:   "json",
								Values: []string{on.JsonName},
							},
						},
					},
				},
			},
		}

		usecase.Request.Fields = append(usecase.Request.Fields, &model.Field{
			Name: "Id",
			Type: model.PrimitiveTypeString,
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{"id"},
				},
				{
					Name:   "validate",
					Values: []string{"required", "uuid"},
				},
			},
		})

		b.Domain.Usecases = append(b.Domain.Usecases, usecase)
		b.CRUDActionToUsecase[crudDefinition.Get] = usecase
		b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Request)
		b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Result)

		if on.Activable {
			usecase = usecase.Copy()
			usecase.Function.Name = strings.Replace(usecase.Function.Name, "Get", "GetActive", 1)
			usecase.Request.Name = strings.Replace(usecase.Request.Name, "Get", "GetActive", 1)
			usecase.Result.Name = strings.Replace(usecase.Result.Name, "Get", "GetActive", 1)
			usecase.Roles = crudDefinition.List.RolesForActive
			b.Domain.Usecases = append(b.Domain.Usecases, usecase)
			b.CRUDActionToUsecase[crudDefinition.List] = usecase
			b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Request)
			b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Result)
		}
	}

	if crudDefinition.List != nil && crudDefinition.List.Active {
		usecase := &model.Usecase{
			Roles: crudDefinition.List.Roles,
			Function: &model.Function{
				Name: "List" + on.Struct.Name + "Usecase",
			},
			Request: &model.Struct{
				Name: "List" + on.Struct.Name + "Request",
			},
			Result: &model.Struct{
				Name: "List" + on.Struct.Name + "Result",
				Fields: []*model.Field{
					{
						Name: on.Struct.Name + "s",
						Type: &model.ArrayType{
							Type: &model.PointerType{
								Type: &model.PkgReference{
									Reference: on,
									Pkg:       b.Domain.Architecture.ModelPkg,
								},
							},
						},
						Tags: []*model.Tag{
							{
								Name:   "json",
								Values: []string{on.JsonName + "s"},
							},
						},
					},
				},
			},
		}

		usecase.Request.Fields = append(usecase.Request.Fields, &model.Field{
			Name: "Ordering",
			Type: &model.PkgReference{
				Pkg:       b.Domain.Architecture.RepositoryPkg,
				Reference: b.GetOrdering(ctx),
			},
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{"ordering"},
				},
			},
		})

		usecase.Request.Fields = append(usecase.Request.Fields, &model.Field{
			Name: "Pagination",
			Type: &model.PkgReference{
				Pkg:       b.Domain.Architecture.RepositoryPkg,
				Reference: b.GetPagination(ctx),
			},
			Tags: []*model.Tag{
				{
					Name:   "json",
					Values: []string{"pagination"},
				},
			},
		})

		b.Domain.Usecases = append(b.Domain.Usecases, usecase)
		b.CRUDActionToUsecase[crudDefinition.List] = usecase
		b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Request)
		b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Result)

		if on.Activable {
			usecase = usecase.Copy()
			usecase.Function.Name = strings.Replace(usecase.Function.Name, "List", "ListActive", 1)
			usecase.Request.Name = strings.Replace(usecase.Request.Name, "List", "ListActive", 1)
			usecase.Result.Name = strings.Replace(usecase.Result.Name, "List", "ListActive", 1)
			usecase.Roles = crudDefinition.List.RolesForActive
			b.Domain.Usecases = append(b.Domain.Usecases, usecase)
			b.CRUDActionToUsecase[crudDefinition.List] = usecase
			b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Request)
			b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Result)
		}
	}

	if crudDefinition.Delete != nil && crudDefinition.Delete.Active {
		usecase := &model.Usecase{
			Roles: crudDefinition.Delete.Roles,
			Function: &model.Function{
				Name: "Delete" + on.Struct.Name + "Usecase",
			},
			Request: &model.Struct{
				Name: "Delete" + on.Struct.Name + "Request",
			},
			Result: &model.Struct{
				Name:   "Delete" + on.Struct.Name + "Result",
				Fields: []*model.Field{},
			},
		}

		usecase.Request.Fields = append([]*model.Field{
			{
				Name: "Id",
				Type: model.PrimitiveTypeString,
				Tags: []*model.Tag{
					{
						Name:   "json",
						Values: []string{"id"},
					},
					{
						Name:   "validate",
						Values: []string{"required", "uuid"},
					},
				},
			},
		}, usecase.Request.Fields...)

		b.Domain.Usecases = append(b.Domain.Usecases, usecase)
		b.CRUDActionToUsecase[crudDefinition.Delete] = usecase
		b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Request)
		b.Domain.UsecaseStructs = append(b.Domain.UsecaseStructs, usecase.Result)
	}

	return b
}
