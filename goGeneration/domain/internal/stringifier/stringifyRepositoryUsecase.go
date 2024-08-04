package stringifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/cleoGitHub/golem-common/pkg/merror"
	"github.com/cleoGitHub/golem/goGeneration/domain/consts"
	"github.com/cleoGitHub/golem/goGeneration/domain/internal/gopkgmanager"
	"github.com/cleoGitHub/golem/goGeneration/domain/model"
)

func StringifyRepositoryUsecase(ctx context.Context, pkgManager *gopkgmanager.GoPkgManager, repo *model.Repository) (string, error) {
	str := ""
	str += fmt.Sprintf(`const %s_TABLE_NAME = "%s"`, strings.ToUpper(repo.On.Struct.Name), repo.TableName) + consts.LN
	s, err := StringifyMapUsecase(ctx, pkgManager, &repo.FieldToColumn)
	if err != nil {
		return "", merror.Stack(err)
	}
	str += s

	s, err = StringifyConstsUsecase(ctx, pkgManager, &repo.AllowedOrderBys)
	if err != nil {
		return "", merror.Stack(err)
	}
	str += s

	s, err = StringifyConstsUsecase(ctx, pkgManager, &repo.AllowedWheres)
	if err != nil {
		return "", merror.Stack(err)
	}
	str += s

	for _, method := range repo.Methods {
		s, err = StringifyStructUsecase(ctx, pkgManager, method.Context)
		if err != nil {
			return "", merror.Stack(err)
		}
		str += s + consts.LN

		s, err = StringifyTypeDefinitionUsecase(ctx, pkgManager, method.Opt)
		if err != nil {
			return "", merror.Stack(err)
		}
		str += s + consts.LN

		for _, opt := range method.Opts {
			s, err = StringifyFunctionUsecase(ctx, pkgManager, opt)
			if err != nil {
				return "", merror.Stack(err)
			}
			str += s + consts.LN
		}
	}
	str += consts.LN

	// str += "type " + repo.Name + " interface {" + consts.LN
	// for _, f := range repo.Functions {
	// 	s, err := StringifyFunctionDefinitionUsecase(ctx, pkgManager, f)
	// 	if err != nil {
	// 		return "", merror.Stack(err)
	// 	}
	// 	str += s + consts.LN
	// }
	// for _, method := range repo.Methods {
	// 	s, err := StringifyFunctionDefinitionUsecase(ctx, pkgManager, method.Function)
	// 	if err != nil {
	// 		return "", merror.Stack(err)
	// 	}
	// 	str += s + consts.LN
	// }
	// str += "}" + consts.LN
	return str, nil
}
