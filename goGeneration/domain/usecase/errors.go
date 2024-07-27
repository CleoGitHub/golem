package usecase

import (
	"errors"
	"strings"
)

var (
	// ErrUnknownType is returned when the type is unknown
	ErrUnknownType = errors.New("unknown type ‘{type}’")

	// ErrDefaultFiedlRedefined is returned when the default field is redefined
	ErrDefaultFiedlRedefined = errors.New("default field '{defaultField}' is redefined")

	// ErrDefaultFiedlRedefined is returned when the default field is redefined
	ErrFiedlRedefined = errors.New("field '{field}' is already defined")

	// ErrRelationModelNotFound is returned when the relation model is not found
	ErrModelNotFound = errors.New("model '{model}' not found")

	// ErrRelationModelNotFound is returned when the relation model is not found
	ErrRelationModelNotFound = errors.New("relation with model '{model}' not found")
)

func NewErrUnknownType(t string) error {
	return errors.New(strings.Replace(ErrUnknownType.Error(), "{type}", t, 1))
}

func NewErrDefaultFiedlRedefined(defaultField string) error {
	return errors.New(
		strings.Replace(ErrDefaultFiedlRedefined.Error(), "{defaultField}", defaultField, 1),
	)
}

func NewErrFiedlRedefined(field string) error {
	return errors.New(
		strings.Replace(ErrDefaultFiedlRedefined.Error(), "{field}", field, 1),
	)
}

func NewErrModelNotFound(model string) error {
	str := strings.Replace(ErrModelNotFound.Error(), "{model}", model, 1)
	return errors.New(str)
}

func NewErrRelationModelNotFound(model string) error {
	str := strings.Replace(ErrRelationModelNotFound.Error(), "{model}", model, 1)
	return errors.New(str)
}
