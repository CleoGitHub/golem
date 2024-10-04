package gopkgmanager

import (
	"fmt"

	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

type GoPkgManager struct {
	Pkg     string
	Imports map[string]*model.GoPkg
}

func (g *GoPkgManager) ToString() string {
	str := fmt.Sprintf("package %s\n", g.Pkg)

	// Print all imports
	str += "import (\n"
	for _, goPkg := range g.Imports {
		alias := ""
		if goPkg.Alias != goPkg.ShortName {
			alias = goPkg.Alias
		}
		str += fmt.Sprintf(`%s "%s"`, alias, goPkg.FullName) + consts.LN
	}
	str += ")\n"

	return str
}

func (g *GoPkgManager) ImportPkg(goPkg *model.GoPkg) error {
	if g.Imports == nil {
		g.Imports = make(map[string]*model.GoPkg)
	}

	if goPkg.Alias == g.Pkg {
		return nil
	}
	// Check if goPkg is already imported
	if _, ok := g.Imports[goPkg.Alias]; !ok {
		g.Imports[goPkg.FullName] = goPkg
	}

	return nil
}
