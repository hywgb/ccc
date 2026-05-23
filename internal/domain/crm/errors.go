package crm

import "errors"

var (
	ErrCustomerNotFound    = errors.New("customer not found")
	ErrDuplicatePhone      = errors.New("phone number already exists for another customer")
	ErrNoPrimaryPhone      = errors.New("at least one primary phone required")
	ErrMultiplePrimary     = errors.New("only one primary phone allowed")
	ErrInvalidPhoneType    = errors.New("invalid phone type")
	ErrInvalidLevel        = errors.New("invalid customer level")
	ErrInvalidFieldType    = errors.New("invalid custom field type")
	ErrInvalidEntityType   = errors.New("invalid entity type")
	ErrRequiredFieldEmpty  = errors.New("required custom field is empty")
	ErrFieldNotFound       = errors.New("custom field definition not found")
)
