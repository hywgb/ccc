package ai

import "errors"

var (
	ErrAnnotationTaskNotFound  = errors.New("annotation task not found")
	ErrAnnotationTaskCompleted = errors.New("annotation task already completed")
	ErrAnnotationTaskCancelled = errors.New("annotation task is cancelled")
	ErrAnnotationResultExists  = errors.New("annotation result already exists for this item")
	ErrLLMModelConfigNotFound  = errors.New("LLM model config not found")
	ErrNoDefaultLLMModel       = errors.New("no default LLM model configured")
	ErrInvalidAnnotationType   = errors.New("invalid annotation type")
	ErrInvalidProviderType     = errors.New("invalid LLM provider type")
)
