package consts

import "github.com/cleogithub/golem/goGeneration/domain/model"

// CommonPkgs contains all following standart packages:
// context
// time
// json
// uuid
// errors
// io
// fmt
// strings

var CommonPkgs = map[string]*model.GoPkg{
	"context": {
		Alias:     "context",
		ShortName: "context",
		FullName:  "context",
	},
	"time": {
		Alias:     "time",
		ShortName: "time",
		FullName:  "time",
	},
	"json": {
		Alias:     "json",
		ShortName: "json",
		FullName:  "encoding/json",
	},
	"reflect": {
		Alias:     "reflect",
		ShortName: "reflect",
		FullName:  "reflect",
	},
	"uuid": {
		Alias:     "uuid",
		ShortName: "uuid",
		FullName:  "github.com/google/uuid",
	},
	"errors": {
		Alias:     "errors",
		ShortName: "errors",
		FullName:  "errors",
	},
	"io": {
		Alias:     "io",
		ShortName: "io",
		FullName:  "io",
	},
	"fmt": {
		Alias:     "fmt",
		ShortName: "fmt",
		FullName:  "fmt",
	},
	"strings": {
		Alias:     "strings",
		ShortName: "strings",
		FullName:  "strings",
	},
	"gorm": {
		Alias:     "gorm",
		ShortName: "gorm",
		FullName:  "gorm.io/gorm",
	},
	"gorm/clause": {
		Alias:     "clause",
		ShortName: "clause",
		FullName:  "gorm.io/gorm/clause",
	},
	"http": {
		Alias:     "http",
		ShortName: "http",
		FullName:  "net/http",
	},
	"merror": {
		Alias:     "merror",
		ShortName: "merror",
		FullName:  "github.com/cleogithub/golem-common/pkg/merror",
	},
	"router": {
		Alias:     "router",
		ShortName: "router",
		FullName:  "github.com/cleogithub/golem-common/pkg/router",
	},
	"slices": {
		Alias:     "slices",
		ShortName: "slices",
		FullName:  "slices",
	},
	"httpclient": {
		Alias:     "httpclient",
		ShortName: "httpclient",
		FullName:  "github.com/cleogithub/golem-common/pkg/httpclient",
	},
}
