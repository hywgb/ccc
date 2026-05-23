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
)
