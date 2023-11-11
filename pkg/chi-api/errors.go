package chi_api

import (
	"fmt"
)

type ErrUnsupportedType struct {
	Type string
}

func (err *ErrUnsupportedType) Error() string {
	return fmt.Sprintf("unsupported type:%s", err.Type)
}
