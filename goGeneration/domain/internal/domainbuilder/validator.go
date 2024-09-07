package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func (b *domainBuilder) GetValidator(ctx context.Context) *model.Interface {
	if b.Err != nil {
		return nil
	}
	if b.Validator != nil {
		return b.Validator
	}

	validatorIntf := &model.Interface{
		Name: "Validator",
		Methods: []*model.Function{
			{
				Name: "Validate",
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
					{
						Name: "request",
						Type: model.PrimitiveTypeInterface,
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
			{
				Name: "IsValidationError",
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
					{
						Name: "err",
						Type: model.PrimitiveTypeError,
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeBool,
					},
				},
			},
			{
				Name: "NewReferenceError",
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
					{
						Name: "reference",
						Type: model.PrimitiveTypeString,
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
			{
				Name: "NewUniqueError",
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
					{
						Name: "field",
						Type: model.PrimitiveTypeString,
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
			{
				Name: "ValidateMimeTypes",
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
					{
						Name: "mimeTypes",
						Type: &model.ArrayType{
							Type: model.PrimitiveTypeString,
						},
					},
					{
						Name: "bytes",
						Type: model.PrimitiveTypeBytes,
					},
					{
						Name: "field",
						Type: model.PrimitiveTypeString,
					},
				},
				Results: []*model.Param{
					{
						Type: model.PrimitiveTypeError,
					},
				},
			},
		},
	}

	b.Validator = validatorIntf

	return validatorIntf
}
