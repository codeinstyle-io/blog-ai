package utils

import (
	"github.com/glebarez/go-sqlite"
)

const UNIQUE_CONSTRAINT = 2067

func IsConstraintError(err error) bool {
	castedErr := err.(*sqlite.Error)
	return castedErr.Code() == UNIQUE_CONSTRAINT
}
