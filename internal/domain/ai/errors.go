package ai

import "errors"

var (
	ErrArticleNotFound  = errors.New("knowledge article not found")
	ErrCategoryNotFound = errors.New("knowledge category not found")
	ErrScriptNotFound   = errors.New("agent script not found")
	ErrTemplateNotFound = errors.New("session info template not found")
	ErrInvalidScript    = errors.New("invalid script content JSON")
)
