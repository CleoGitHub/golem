package stringifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func StringifyModelUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, m *model.Model) (string, error) {
	str := ""

	for _, c := range m.Struct.Consts {
		c, err := StringifyConstUsecase(ctx, pkgManager, c)
		if err != nil {
			return "", merror.Stack(err)
		}
		str += c + consts.LN
	}

	str += fmt.Sprintf("type %s struct {", m.Struct.Name) + consts.LN
	str += "// Attributes" + consts.LN
	for _, field := range m.Struct.Fields {
		// Never pass model pkg in GetType call, model are always in the same package
		st, err := StringifyFieldUsecase(ctx, pkgManager, field)
		if err != nil {
			return "", merror.Stack(err)
		}
		str += st + consts.LN
	}

	if m.Relations != nil && len(m.Relations) > 0 {
		// Add relation field
		str += consts.LN
		str += "// Relations" + consts.LN
		for _, relation := range m.Relations {
			if relation.Type == model.RelationSingleMandatory || relation.Type == model.RelationSingleOptionnal {
				s, err := StringifyFieldUsecase(ctx, pkgManager, &model.Field{
					Name: relation.On.Struct.Name + "Id",
					Type: model.PrimitiveTypeString,
					Tags: []*model.Tag{
						{
							Name:   "json",
							Values: []string{relation.On.JsonName + "Id"},
						},
					},
				})
				if err != nil {
					return "", merror.Stack(err)
				}
				str += s + consts.LN
				s, err = StringifyFieldUsecase(ctx, pkgManager, &model.Field{
					Name: relation.On.Struct.Name,
					Type: &model.PointerType{
						Type: relation.On,
					},
					Tags: []*model.Tag{
						{
							Name:   "json",
							Values: []string{relation.On.JsonName},
						},
					},
				})
				if err != nil {
					return "", merror.Stack(err)
				}
				str += s + consts.LN
			} else if relation.Type == model.RelationMultiple {
				relationName := relation.On.Struct.Name + "s"
				if strings.HasSuffix(relationName, "ys") {
					relationName = relationName[:len(relationName)-2] + "ies"
				}

				jsonRelationName := relation.On.JsonName + "s"
				if strings.HasSuffix(jsonRelationName, "ys") {
					jsonRelationName = jsonRelationName[:len(jsonRelationName)-1] + "ies"
				}

				s, err := StringifyFieldUsecase(ctx, pkgManager, &model.Field{
					Name: relationName,
					Type: &model.ArrayType{
						Type: &model.PointerType{
							Type: relation.On,
						},
					},
					Tags: []*model.Tag{
						{
							Name:   "json",
							Values: []string{jsonRelationName},
						},
					},
				})
				if err != nil {
					return "", merror.Stack(err)
				}
				str += s + consts.LN
			}
		}
	}

	str += "}" + consts.LN

	for _, method := range m.Struct.Methods {
		st, err := StringifyFunctionDefinitionUsecase(ctx, pkgManager, method)
		if err != nil {
			return "", merror.Stack(err)
		}
		str += fmt.Sprintf("func (%s *%s) %s {", m.Struct.GetMethodName(), m.Struct.Name, st) + consts.LN
		s, pkgs := method.Content()
		for _, pkg := range pkgs {
			if err := pkgManager.ImportPkg(pkg); err != nil {
				return "", merror.Stack(err)
			}
		}
		str += s + consts.LN
		str += "}" + consts.LN
	}

	return str, nil
}
