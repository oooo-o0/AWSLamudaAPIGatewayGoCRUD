package common

import (
	"errors"
	"strings"
)

func ConflictError(msg string) error {
	return errors.New("conflict: " + msg)
}

func IsConflictError(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "conflict:")
}
