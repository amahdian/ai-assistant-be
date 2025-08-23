package errs

import (
	"fmt"
)

const (
	DefaultEntryNotFoundMessage = "%q entry by id %d could not be found"
	InvalidSearchFieldMessage   = "invalid field %q for search condition"
	FailedToListItemsMessage    = "failed to list %q"
)

func NewInvalidSearchFieldErr(fieldName string) error {
	return fmt.Errorf(InvalidSearchFieldMessage, fieldName)
}

type EntryNotFoundErr struct {
	message string
}

func NewEntryNotFoundErr(message string) EntryNotFoundErr {
	return EntryNotFoundErr{message}
}

func (e EntryNotFoundErr) Error() string {
	return e.message
}
