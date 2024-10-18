package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	GORM_DOMAIN_REPO_METHOD_NAME = "repo"
)

type GormDomainRepositoryBuilder struct {
	*EmptyBuilder

	DomainBuilder    *domainBuilder
	DomainRepository *model.File
	Migrations       string

	Err error
}

func NewGormDomainRepositoryBuilder(
	ctx context.Context,
	domainBuilder *domainBuilder,
	definition *coredomaindefinition.Domain,
) Builder {
	builder := &GormDomainRepositoryBuilder{
		DomainBuilder: domainBuilder,
		DomainRepository: &model.File{
			Name: GetDomainRepositoryName(ctx, definition),
			Pkg:  domainBuilder.Domain.Architecture.RepositoryPkg,
		},
	}

	return builder
}

var _ Builder = (*GormDomainRepositoryBuilder)(nil)

func (builder *GormDomainRepositoryBuilder) WithModel(ctx context.Context, model *coredomaindefinition.Model) {
	if builder.Err != nil {
		return
	}

	builder.Migrations += fmt.Sprintf("&%s{},", GetModelName(ctx, model)) + "\n"
}

func (builder *GormDomainRepositoryBuilder) addTransaction(ctx context.Context) {
	if builder.Err != nil {
		return
	}
	builder.DomainBuilder.Domain.Files = append(builder.DomainBuilder.Domain.Files, &model.File{
		Name: TRANSACTION_NAME,
		Pkg:  builder.DomainBuilder.GetGormAdapterPackage(),
		Elements: []interface{}{
			&model.Struct{
				Name:       TRANSACTION_NAME,
				MethodName: stringtool.LowerFirstLetter(TRANSACTION_NAME),
				Fields: []*model.Field{
					{
						Name: "db",
						Type: &model.PointerType{
							Type: &model.PkgReference{
								Pkg: consts.CommonPkgs["gorm"],
								Reference: &model.ExternalType{
									Type: "DB",
								},
							},
						},
					},
				},
				Methods: []*model.Function{
					{
						Name: TRANSACTION_GET,
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
						Content: func() (string, []*model.GoPkg) {
							return fmt.Sprintf("return %s.db", stringtool.LowerFirstLetter(TRANSACTION_NAME)), nil
						},
					},
					{
						Name: TRANSACTION_COMMIT,
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
							return fmt.Sprintf("return %s.db.Commit().Error", stringtool.LowerFirstLetter(TRANSACTION_NAME)), nil
						},
					},
					{
						Name: TRANSACTION_ROLLBACK,
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
							return fmt.Sprintf("return %s.db.Rollback().Error", stringtool.LowerFirstLetter(TRANSACTION_NAME)), nil
						},
					},
				},
			},
		},
	})
}

func (builder *GormDomainRepositoryBuilder) Build(ctx context.Context) error {
	if builder.Err != nil {
		return builder.Err
	}

	builder.addTransaction(ctx)

	gormDomainRepo := &model.Struct{
		Name:       GetGormDomainRepositoryName(ctx, builder.DomainBuilder.Definition),
		MethodName: GORM_DOMAIN_REPO_METHOD_NAME,
		Fields: []*model.Field{
			{
				Name: GORM_DOMAIN_REPOSITORY_DB_FIELD_NAME,
				Type: &model.PointerType{
					Type: &model.PkgReference{
						Pkg: consts.CommonPkgs["gorm"],
						Reference: &model.ExternalType{
							Type: "DB",
						},
					},
				},
			},
		},
	}
	gormDomainRepo.Methods = append(gormDomainRepo.Methods, &model.Function{
		Name: REPOSITORY_MIGRATE,
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
			str := fmt.Sprintf("return %s.%s.AutoMigrate(", gormDomainRepo.GetMethodName(), GORM_DOMAIN_REPOSITORY_DB_FIELD_NAME) + consts.LN
			str += builder.Migrations
			str += ")" + consts.LN
			return str, []*model.GoPkg{
				consts.CommonPkgs["gorm"],
			}
		},
	})

	gormDomainRepo.Methods = append(gormDomainRepo.Methods, &model.Function{
		Name: REPOSITORY_BEGIN_TRANSACTION,
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
					Pkg: builder.DomainBuilder.GetRepositoryPackage(),
					Reference: &model.ExternalType{
						Type: TRANSACTION_NAME,
					},
				},
			},
			{
				Type: model.PrimitiveTypeError,
			},
		},
		Content: func() (string, []*model.GoPkg) {
			str := fmt.Sprintf(
				"return &%s{db: %s.%s.Begin()}, nil",
				TRANSACTION_NAME, gormDomainRepo.GetMethodName(), GORM_DOMAIN_REPOSITORY_DB_FIELD_NAME,
			) + consts.LN
			return str, []*model.GoPkg{}
		},
	})

	byOperatorToGormOperator := &model.Function{
		Name: OPERATOR_TO_GORM_OPERATOR,
		Args: []*model.Param{
			{
				Name: "operator",
				Type: &model.PkgReference{
					Pkg: builder.DomainBuilder.GetRepositoryPackage(),
					Reference: &model.ExternalType{
						Type: REPOSITORY_WHERE_OPERATOR_TYPE,
					},
				},
			},
		},
		Results: []*model.Param{
			{
				Type: model.PrimitiveTypeString,
			},
		},
		Content: func() (string, []*model.GoPkg) {
			str := "switch operator {" + consts.LN
			str += fmt.Sprintf("case %s.%s:", builder.DomainBuilder.GetRepositoryPackage().Alias, REPOSITORY_WHERE_OPERATOR_EQUAL) + consts.LN
			str += `return "="` + consts.LN
			str += fmt.Sprintf("case %s.%s:", builder.DomainBuilder.GetRepositoryPackage().Alias, REPOSITORY_WHERE_OPERATOR_NOT_EQUAL) + consts.LN
			str += `return "<>"` + consts.LN
			str += fmt.Sprintf("case %s.%s:", builder.DomainBuilder.GetRepositoryPackage().Alias, REPOSITORY_WHERE_OPERATOR_IN) + consts.LN
			str += `return "IN"` + consts.LN
			str += fmt.Sprintf("case %s.%s:", builder.DomainBuilder.GetRepositoryPackage().Alias, REPOSITORY_WHERE_OPERATOR_NOT_IN) + consts.LN
			str += `return "NOT IN"` + consts.LN

			str += "}" + consts.LN
			str += `return ""`
			return str, nil
		},
	}
	valueToGormValue := &model.Function{
		Name: VALUE_TO_GORM_VALUE,
		Args: []*model.Param{
			{
				Name: "value",
				Type: model.PrimitiveTypeInterface,
			},
		},
		Results: []*model.Param{
			{
				Type: model.PrimitiveTypeString,
			},
		},
		Content: func() (string, []*model.GoPkg) {
			str := "v := reflect.TypeOf(value).Kind()" + consts.LN
			str += "switch  v {" + consts.LN
			str += "case reflect.Array, reflect.Slice:" + consts.LN
			str += `str := "["` + consts.LN
			str += "for i := 0; i < reflect.ValueOf(value).Len(); i++ {" + consts.LN
			str += fmt.Sprintf("str += %s(reflect.ValueOf(value).Index(i).Interface())", VALUE_TO_GORM_VALUE) + consts.LN
			str += `str += ","` + consts.LN
			str += "}" + consts.LN
			str += `str = strings.TrimSuffix(str, ",") + "]"` + consts.LN
			str += `return str` + consts.LN
			str += "case reflect.String:" + consts.LN
			str += `return fmt.Sprintf("\"%+v\"", v)` + consts.LN
			str += "default:" + consts.LN
			str += `return fmt.Sprintf("%+v", v)` + consts.LN
			str += "}"
			return str, []*model.GoPkg{
				consts.CommonPkgs["fmt"],
				consts.CommonPkgs["strings"],
				consts.CommonPkgs["reflect"],
			}
		},
	}
	builder.DomainBuilder.Domain.Files = append(builder.DomainBuilder.Domain.Files, &model.File{
		Name: GetGormDomainRepositoryName(ctx, builder.DomainBuilder.Definition),
		Elements: []interface{}{
			byOperatorToGormOperator,
			valueToGormValue,
			gormDomainRepo,
		},
		Pkg: builder.DomainBuilder.GetGormAdapterPackage(),
	})

	if builder.Err != nil {
		return builder.Err
	}

	return nil
}
