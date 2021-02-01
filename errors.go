package stickers

import (
	"fmt"
)

// ErrFetchingLabelsForUser
type ErrFetchingLabelsForUser struct {
	User string
}

// Error
func (e ErrFetchingLabelsForUser) Error() string {
	return fmt.Sprintf("could not fetch labels for user '%s'", e.User)
}

// ErrFilterNotFound
type ErrFilterNotFound struct {
	Filter string
	User   string
}

// ErrFilterNotFound
func (e ErrFilterNotFound) Error() string {
	return fmt.Sprintf("could not fetch filter '%s' for user '%s'", e.Filter, e.User)
}

// ErrLabelNotFound
type ErrLabelNotFound struct {
	Label string
	User  string
}

// Error
func (e ErrLabelNotFound) Error() string {
	return fmt.Sprintf("could not fetch label '%s' for user '%s'", e.Label, e.User)
}

// ErrFetchingFiltersForUser
type ErrFetchingFiltersForUser struct {
	User string
}

// Error
func (e ErrFetchingFiltersForUser) Error() string {
	return fmt.Sprintf("could not fetch filters for user '%s'", e.User)
}

// ErrMissingServiceAccountFile
type ErrMissingServiceAccountFile struct{}

// Error
func (e ErrMissingServiceAccountFile) Error() string {
	return "missing service account file, set with `SetServiceAccountFile(filename)`"
}

// ErrInvalidServiceAccountFile
type ErrInvalidServiceAccountFile struct {
	File string
}

// Error
func (e ErrInvalidServiceAccountFile) Error() string {
	return fmt.Sprintf("invalid service account file, `%s`", e.File)
}

// ErrLabelMismatch
type ErrLabelMismatch struct {
	Label string
}

// ErrLabelMismatch
func (e ErrLabelMismatch) Error() string {
	return fmt.Sprintf("label mismatch for `%s`", e.Label)
}

// ErrFilterMismatch
type ErrFilterMismatch struct {
	Filter string
}

// ErrLabelMismatch
func (e ErrFilterMismatch) Error() string {
	return fmt.Sprintf("filter mismatch for `%s`", e.Filter)
}
