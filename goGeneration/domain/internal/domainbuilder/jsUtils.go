package domainbuilder

import (
	"fmt"
	"strings"

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

func StructToClass(s *model.Struct, initialLevel string) string {
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

		hydrateContent += level() + consts.TAB + fmt.Sprintf("if ( data.%s ) { this.%s = data.%s }", name, name, name) + consts.LN
	}
	constructor += ") {" + consts.LN + constructorContent + level() + "}"
	str += consts.LN + constructor + consts.LN
	hydrate += hydrateContent + level() + consts.TAB + "return this" + consts.LN + level() + "}" + consts.LN
	str += consts.LN + hydrate
	str += consts.LN

	str += level() + fmt.Sprintf("static from(data) { return new %s().hydrate(data) }", s.Name) + consts.LN

	subLevel()

	str += level() + "}" + consts.LN

	return str
}

func JSGetClassFromSimpleFields(name string, fields []string) string {
	str := fmt.Sprintf("export class %s {", name) + consts.LN
	str += consts.TAB + strings.Join(fields, "\n\t") + consts.LN
	str += consts.LN

	str += consts.TAB + "constructor("
	for _, field := range fields {
		str += field + ","
	}
	str = strings.TrimSuffix(str, ",")
	str += ") {" + consts.LN
	for _, field := range fields {
		str += consts.TAB + consts.TAB + fmt.Sprintf("this.%s = %s", field, field) + consts.LN
	}
	str += consts.TAB + "}" + consts.LN
	str += consts.LN

	str += consts.TAB + "hydrate(data) {" + consts.LN
	for _, field := range fields {
		str += consts.TAB + consts.TAB + fmt.Sprintf("if ( data.%s ) { this.%s = data.%s }", field, field, field) + consts.LN
	}
	str += consts.TAB + "}" + consts.LN
	str += consts.LN

	str += consts.TAB + fmt.Sprintf("static from(data) { return new %s().hydrate(data) }", name) + consts.LN

	str += "}" + consts.LN

	return str
}

func JSGetClassFromTransformationFields(name string, fields map[string]string) string {
	str := fmt.Sprintf("export class %s {", name) + consts.LN
	for field := range fields {
		str += consts.TAB + field + consts.LN
	}
	str += consts.LN

	str += consts.TAB + "constructor("
	for field := range fields {
		str += field + ","
	}
	str = strings.TrimSuffix(str, ",")
	str += ") {" + consts.LN
	for field := range fields {
		str += consts.TAB + consts.TAB + fmt.Sprintf("this.%s = %s", field, field) + consts.LN
	}
	str += consts.TAB + "}" + consts.LN
	str += consts.LN

	str += consts.TAB + "hydrate(data) {" + consts.LN
	for field, fieldTransformation := range fields {
		str += consts.TAB + consts.TAB + fmt.Sprintf("if ( data.%s ) { this.%s = %s }", field, field, fieldTransformation) + consts.LN
	}
	str += consts.TAB + "}" + consts.LN
	str += consts.LN

	str += consts.TAB + fmt.Sprintf("static from(data) { return new %s().hydrate(data) }", name) + consts.LN

	str += "}" + consts.LN

	return str
}
