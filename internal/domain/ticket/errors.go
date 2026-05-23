package ticket

import "errors"

var (
	ErrTicketNotFound      = errors.New("ticket not found")
	ErrTemplateNotFound    = errors.New("ticket template not found")
	ErrCategoryNotFound    = errors.New("ticket category not found")
	ErrInvalidTransition   = errors.New("invalid ticket status transition")
	ErrTemplateNotPublished = errors.New("template is not published")
	ErrAlreadyPublished    = errors.New("template is already published")
	ErrInvalidPriority     = errors.New("invalid ticket priority")
	ErrInvalidFlowGraph    = errors.New("invalid flow graph JSON")
)
