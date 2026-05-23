package integration

import "errors"

var (
	ErrDNCBlocked       = errors.New("number is on DNC list")
	ErrDNCEntryNotFound = errors.New("DNC entry not found")
	ErrTagAlreadyExists = errors.New("tag already assigned to call")
)
