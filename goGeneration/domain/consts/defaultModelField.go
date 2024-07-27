package consts

import "github.com/cleoGitHub/golem/coredomaindefinition"

var DefaultModelFields = []*coredomaindefinition.Field{
	{
		Name: "id",
		Type: coredomaindefinition.PrimitiveTypeString,
	},
	{
		Name: "createdAt",
		Type: coredomaindefinition.PrimitiveTypeDateTime,
	},
	{
		Name: "updatedAt",
		Type: coredomaindefinition.PrimitiveTypeDateTime,
	},
	{
		Name: "deletedAt",
		Type: coredomaindefinition.PrimitiveTypeDateTime,
	},
}
