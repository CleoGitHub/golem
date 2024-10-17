package consts

import "github.com/cleogithub/golem/coredomaindefinition"

const ID = "Id"

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
}
