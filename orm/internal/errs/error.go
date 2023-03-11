package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("orm: Only one level of pointer support")
)

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm: unsupported expresion type %v", expr)
}

func NewErrUnknownField(name string) error {
	return fmt.Errorf("orm: invalid field %s", name)
}

func NewErrInvalidTagContent(pair string) error {
	return fmt.Errorf("orm: invalid tag content %s", pair)
}
