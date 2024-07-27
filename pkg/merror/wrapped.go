package merror

import "fmt"

type wrappedError struct {
	src     error
	details string
}

func NewWrappedError(err error, details string) error {
	return &wrappedError{
		src:     err,
		details: details,
	}
}

func (e *wrappedError) Error() string {
	return fmt.Sprintf("%s, %s", e.src.Error(), e.details)
}

func (e *wrappedError) Unwrap() error {
	return e.src
}
