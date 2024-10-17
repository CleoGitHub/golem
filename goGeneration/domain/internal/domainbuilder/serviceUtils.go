package domainbuilder

import (
	"context"
	"fmt"

	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/coredomaindefinition"
)

func GetServiceName(ctx context.Context, on *coredomaindefinition.Domain) string {
	return fmt.Sprintf("%sService", stringtool.UpperFirstLetter(on.Name))
}
