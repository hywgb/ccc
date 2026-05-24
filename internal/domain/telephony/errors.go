package telephony

import "errors"

var (
	ErrCarrierNotFound     = errors.New("carrier not found")
	ErrSIPTrunkNotFound    = errors.New("sip trunk not found")
	ErrPhoneNumberNotFound = errors.New("phone number not found")
	ErrPhoneNumberExists   = errors.New("phone number already exists")
	ErrNoMatchingRoute     = errors.New("no matching routing rule")
	ErrNoCLIPolicy         = errors.New("no CLI policy found")
	ErrNoCLINumber         = errors.New("no CLI number available")
	ErrTrunkGroupNotFound  = errors.New("trunk group not found")
	ErrNoHealthyTrunk      = errors.New("no healthy trunk available")
	ErrTrunkDown           = errors.New("trunk is down")
	ErrInvalidPhoneFormat  = errors.New("phone number must start with + and contain only digits (E.164)")
)
