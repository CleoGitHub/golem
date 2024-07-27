package merror

import (
	"fmt"
	"runtime"
)

type Merror struct {
	line     int
	file     string
	function string
	cause    error
}

func (m *Merror) Error() string {
	line := fmt.Sprintf("at %s:%d in %s\n", m.file, m.line, m.function)
	if _, ok := m.cause.(*Merror); !ok {
		return "Error happened, reason: '" + m.cause.Error() + "'. Stack: \n" + line
	} else {
		return m.cause.Error() + line
	}
}

func Stack(err error) error {
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc).Name()

	m := &Merror{
		line:     line,
		file:     file,
		function: function,
		cause:    err,
	}
	return m
}

func (m *Merror) Unwrap() error {
	return m.cause
}
