package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/cleogithub/golem-common/pkg/merror"
	"github.com/cleogithub/golem-common/pkg/stringtool"
	"github.com/cleogithub/golem/goGeneration/domain/consts"
	"github.com/cleogithub/golem/goGeneration/domain/model"
)

func getLevel(initialLevel string) (func(), func(), func() string, func()) {
	level := initialLevel

	add := func() {
		level += consts.TAB
	}

	sub := func() {
		level = strings.TrimSuffix(level, consts.TAB)
	}

	get := func() string {
		return level
	}

	reset := func() {
		level = initialLevel
	}

	return add, sub, get, reset
}

func (g *GenerationUsecaseImpl) GenerateJavascriptClientUsecase(ctx context.Context, domain *model.Domain, path string) error {
	// if file path does not exist, create it
	filepath := path + "/" + domain.Architecture.JavascriptClient
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			return merror.Stack(err)
		}
	}

	typesExport := ""

	// Generate service
	f, err := os.Create(filepath + "/types.js")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	addLevel, subLevel, level, resetLevel := getLevel("")
	addLevel()
	str := "import {" + consts.LN
	for _, m := range domain.Models {
		str += level() + fmt.Sprintf("%s,", m.Struct.Name) + consts.LN
	}
	subLevel()
	str += "} from './index.js'" + consts.LN
	str += consts.LN

	str += structToClass(domain.Pagination, level()) + consts.LN

	str += structToClass(domain.Ordering, level()) + consts.LN

	for i, usecase := range domain.Usecases {
		str += fmt.Sprintf("// %s", strings.TrimSuffix(usecase.Function.Name, "Usecase")) + consts.LN
		str += structToClass(usecase.Request, level())
		typesExport += consts.TAB + fmt.Sprintf("%s,", usecase.Request.Name) + consts.LN
		str += structToClass(usecase.Result, level())
		typesExport += consts.TAB + fmt.Sprintf("%s,", usecase.Result.Name) + consts.LN
		if i < len(domain.Usecases)-1 {
			str += consts.LN
		}

	}

	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	modelsExport := ""

	// Generate models
	f, err = os.Create(filepath + "/models.js")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	str = ""
	resetLevel()
	for i, m := range domain.Models {
		modelsExport += consts.TAB + fmt.Sprintf("%s,", m.Struct.Name) + consts.LN
		str += fmt.Sprintf("export class %s {", m.Struct.Name) + consts.LN
		constructor := consts.TAB + "constructor("
		constructorContent := ""
		hydratorContent := ""
		hydratorArg := "data"

		for idx, field := range m.Struct.Fields {
			str += consts.TAB + fmt.Sprintf("%s //%s", field.JsonName, field.Type.GetType()) + consts.LN

			constructor += field.JsonName
			constructorContent += consts.TAB + consts.TAB + fmt.Sprintf("this.%s = %s", field.JsonName, field.JsonName) + consts.LN
			if idx < len(m.Struct.Fields)-1 {
				constructor += ", "
			}
			hydratorContent += consts.TAB + consts.TAB + fmt.Sprintf("if (%s.%s != undefined) { this.%s = %s.%s }", hydratorArg, field.JsonName, field.JsonName, hydratorArg, field.JsonName) + consts.LN
		}
		for _, relation := range m.Relations {
			name := stringtool.LowerFirstLetter(relation.On.Struct.Name)
			t := relation.On.GetType()
			if relation.Type == model.RelationMultiple {
				if strings.HasSuffix(name, "y") {
					name = name[:len(name)-1] + "ies"
				} else {
					name = name + "s"
				}
				t = fmt.Sprintf("[]%s", t)
			} else {
				str += consts.TAB + fmt.Sprintf("%sId //string", name) + consts.LN
				if constructor != consts.TAB+"constructor(" {
					constructor += ", "
				}
				constructor += name + "Id"
				constructorContent += consts.TAB + consts.TAB + fmt.Sprintf("this.%sId = %sId", name, name) + consts.LN
				hydratorContent += consts.TAB + consts.TAB + fmt.Sprintf("if (%s.%sId != undefined) { this.%sId = %s.%sId }", hydratorArg, name, name, hydratorArg, name) + consts.LN
			}
			str += consts.TAB + fmt.Sprintf("%s //%s", name, t) + consts.LN
		}
		str += consts.LN + constructor + ") {" + consts.LN + constructorContent + consts.TAB + "}" + consts.LN
		str += consts.LN

		str += consts.TAB + "hydrate(data) {" + consts.LN + hydratorContent + consts.TAB + consts.TAB + "return this" + consts.LN + consts.TAB + "}" + consts.LN
		str += consts.LN

		str += consts.TAB + "static From(data) {" + consts.LN
		str += consts.TAB + consts.TAB + "return (new " + m.Struct.Name + "()).hydrate(data)" + consts.LN
		str += consts.TAB + "}" + consts.LN

		str += "}" + consts.LN

		if i < len(domain.Models)-1 {
			str += consts.LN
		}
	}

	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	// Generate service

	serviceExport := ""
	f, err = os.Create(filepath + "/service.js")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	resetLevel()
	str = fmt.Sprintf("export class %sService {", stringtool.UpperFirstLetter(domain.Name)) + consts.LN
	serviceExport += consts.TAB + fmt.Sprintf("%sService,", stringtool.UpperFirstLetter(domain.Name)) + consts.LN
	addLevel()

	str += level() + "host" + consts.LN
	str += level() + "port" + consts.LN
	str += level() + "httpClient" + consts.LN
	str += consts.LN

	str += level() + "constructor(host, port, httpClient) {" + consts.LN
	addLevel()
	str += level() + "this.host = host" + consts.LN
	str += level() + "this.port = port" + consts.LN
	str += level() + "this.httpClient = httpClient" + consts.LN
	subLevel()
	str += level() + "}" + consts.LN
	str += consts.LN

	for i, usecase := range domain.Usecases {
		endpoint := strings.TrimSuffix(usecase.Function.Name, "Usecase")
		str += level() + fmt.Sprintf("%s(data) {", endpoint) + consts.LN
		addLevel()
		str += level() + "return new Promise((resolve, reject) => {" + consts.LN
		addLevel()
		str += level() + "this.httpClient.post(" + consts.LN
		addLevel()
		str += level() + fmt.Sprintf("this.host+':'+this.port+'/%s/%s',", domain.Name, endpoint) + consts.LN
		str += level() + "JSON.stringify(data)," + consts.LN
		str += level() + `{ headers: { 'Content-Type': 'application/json', 'Accept': 'application/json' } }` + consts.LN
		subLevel()
		str += level() + ")" + consts.LN
		str += level() + ".then(resp => {" + consts.LN
		addLevel()
		str += level() + fmt.Sprintf("resolve(%s.From(JSON.parse(resp)))", usecase.Result.Name) + consts.LN
		subLevel()
		str += level() + "})" + consts.LN
		str += level() + ".catch(err => { reject(err) })" + consts.LN
		subLevel()
		str += level() + "})" + consts.LN
		subLevel()
		str += level() + "}" + consts.LN
		if i < len(domain.Usecases)-1 {
			str += consts.LN
		}
	}

	str += "}" + consts.LN

	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	// Generate service
	f, err = os.Create(filepath + "/index.js")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	resetLevel()

	str = "import {" + consts.LN
	str += typesExport
	str += "} from './types'" + consts.LN
	str += consts.LN

	str += "import {" + consts.LN
	str += modelsExport
	str += "} from './models'" + consts.LN
	str += consts.LN

	str += "import {" + consts.LN
	str += serviceExport
	str += "} from './service'" + consts.LN
	str += consts.LN

	str += "export {" + consts.LN
	str += typesExport
	str += modelsExport
	str += serviceExport
	str += "}" + consts.LN

	_, err = f.WriteString(str)
	if err != nil {
		return merror.Stack(err)
	}

	module := map[string]interface{}{
		"name": domain.Name,
		"main": "index.js",
		"type": "module",
	}

	f, err = os.Create(filepath + "/package.json")
	if err != nil {
		return merror.Stack(err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(module, "", "\t")
	if err != nil {
		return merror.Stack(err)
	}

	_, err = f.WriteString(string(b))
	if err != nil {
		return merror.Stack(err)
	}

	cmd := exec.Command("npm", "install")
	cmd.Dir = path + "/" + domain.Architecture.JavascriptClient

	if err := cmd.Run(); err != nil {
		return merror.Stack(err)
	}

	return nil
}

func structToClass(s *model.Struct, initialLevel string) string {
	str := ""
	addLevel, subLevel, level, _ := getLevel(initialLevel)
	str += level() + fmt.Sprintf("export class %s {", s.Name) + consts.LN
	addLevel()
	constructor := level() + "constructor("
	constructorContent := ""

	hydrate := level() + "hydrate(data) {" + consts.LN
	hydrateContent := ""

	for idx, field := range s.Fields {
		name := stringtool.LowerFirstLetter(field.Name)
		str += level() + name + consts.LN
		constructor += name
		constructorContent += level() + consts.TAB + fmt.Sprintf("this.%s = %s", name, name) + consts.LN
		if idx < len(s.Fields)-1 {
			constructor += ", "
		}

		if slices.Contains([]model.Type{
			model.PrimitiveTypeInt,
			model.PrimitiveTypeFloat,
			model.PrimitiveTypeBool,
			model.PrimitiveTypeString,
		}, field.Type) {
			hydrateContent += level() + consts.TAB + fmt.Sprintf("if ( data.%s ) { this.%s = data.%s }", name, name, name) + consts.LN
		} else if t, ok := field.Type.(*model.PointerType); ok {
			if pkgRef, ok := t.Type.(*model.PkgReference); ok {
				if mo, ok := pkgRef.Reference.(*model.Model); ok {
					addLevel()
					hydrateContent += level() + fmt.Sprintf("if ( data.%s ) { this.%s = %s.From(data.%s) }", name, name, mo.Struct.Name, name) + consts.LN
					subLevel()
				}
			}
		} else if arr, ok := field.Type.(*model.ArrayType); ok {
			if t, ok := arr.Type.(*model.PointerType); ok {
				if pkgRef, ok := t.Type.(*model.PkgReference); ok {
					if mo, ok := pkgRef.Reference.(*model.Model); ok {
						addLevel()
						hydrateContent += level() + fmt.Sprintf("if ( data.%s ) {", name) + consts.LN
						addLevel()
						hydrateContent += level() + fmt.Sprintf("this.%s = data.%s.Map(el => %s.From(el))", name, name, mo.Struct.Name) + consts.LN
						subLevel()
						hydrateContent += level() + "}" + consts.LN
						subLevel()
					}
				}
			}
		}
	}
	constructor += ") {" + consts.LN + constructorContent + level() + "}"
	str += consts.LN + constructor + consts.LN
	hydrate += hydrateContent + level() + "}" + consts.LN
	str += consts.LN + hydrate
	str += consts.LN

	str += level() + fmt.Sprintf("static From(data) { return new %s().hydrate(data) }", s.Name) + consts.LN

	subLevel()

	str += level() + "}" + consts.LN

	return str
}
