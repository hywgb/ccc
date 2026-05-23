package ai

import "context"

// AnnotationTaskRepository manages annotation task persistence.
type AnnotationTaskRepository interface {
	Create(ctx context.Context, task *AnnotationTask) error
	GetByID(ctx context.Context, id int64) (*AnnotationTask, error)
	Update(ctx context.Context, task *AnnotationTask) error
	List(ctx context.Context, tenantID int64) ([]AnnotationTask, error)
	Delete(ctx context.Context, id int64) error
}

// AnnotationResultRepository manages annotation result persistence.
type AnnotationResultRepository interface {
	Create(ctx context.Context, result *AnnotationResult) error
	ListByTaskID(ctx context.Context, taskID int64) ([]AnnotationResult, error)
	GetByTaskAndIndex(ctx context.Context, taskID int64, itemIndex int) (*AnnotationResult, error)
}

// LLMModelConfigRepository manages LLM model configurations.
type LLMModelConfigRepository interface {
	Create(ctx context.Context, config *LLMModelConfig) error
	GetByID(ctx context.Context, id int64) (*LLMModelConfig, error)
	Update(ctx context.Context, config *LLMModelConfig) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, tenantID int64) ([]LLMModelConfig, error)
	GetDefault(ctx context.Context, tenantID int64) (*LLMModelConfig, error)
}
