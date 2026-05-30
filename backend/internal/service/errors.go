package service

import (
	"errors"

	"gorm.io/gorm"
)

// wrapNotFound converts gorm.ErrRecordNotFound to the given domain error.
// Returns the original error if it's not a "not found" error.
func wrapNotFound(err error, domainErr error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domainErr
	}
	return err
}
