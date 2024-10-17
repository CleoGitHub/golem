package domainbuilder

import (
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

const (
	ACTIVE_FIELD_NAME = "Active"
)

var CTX = &model.Param{
	Name: "ctx",
	Type: &model.PkgReference{
		Pkg: consts.CommonPkgs["context"],
		Reference: &model.ExternalType{
			Type: "Context",
		},
	},
}
