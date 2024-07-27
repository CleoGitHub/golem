package gopkgmanager

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrUnknownPkg = fmt.Errorf("unknown package {pkg}")
)

func NewErrUnknownPkg(pkg string) error {
	return errors.New(strings.ReplaceAll(ErrUnknownPkg.Error(), "{pkg}", pkg))
}
