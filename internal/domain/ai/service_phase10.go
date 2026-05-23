package ai

import (
	"context"
	"time"

	"github.com/divord97/ccc/pkg/snowflake"
)

// AnnotationService manages annotation tasks and results.
type AnnotationService struct {
	taskRepo   AnnotationTaskRepository
	resultRepo AnnotationResultRepository
}

func NewAnnotationService(taskRepo AnnotationTaskRepository, resultRepo AnnotationResultRepository) *AnnotationService {
	return &AnnotationService{taskRepo: taskRepo, resultRepo: resultRepo}
}

func (s *AnnotationService) CreateTask(ctx context.Context, in CreateAnnotationTaskInput) (*AnnotationTask, error) {
	if !isValidAnnotationType(in.Type) {
		return nil, ErrInvalidAnnotationType
	}
	task := &AnnotationTask{
		ID:           snowflake.NextID(),
		TenantID:     in.TenantID,
		Name:         in.Name,
		Type:         in.Type,
		DatasetID:    in.DatasetID,
		AssigneeID:   in.AssigneeID,
		TotalItems:   in.TotalItems,
		LabeledItems: 0,
		Status:       AnnotationTaskStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *AnnotationService) GetTask(ctx context.Context, id int64) (*AnnotationTask, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *AnnotationService) ListTasks(ctx context.Context, tenantID int64) ([]AnnotationTask, error) {
	return s.taskRepo.List(ctx, tenantID)
}

func (s *AnnotationService) StartTask(ctx context.Context, id int64) (*AnnotationTask, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task.Status == AnnotationTaskStatusCompleted {
		return nil, ErrAnnotationTaskCompleted
	}
	if task.Status == AnnotationTaskStatusCancelled {
		return nil, ErrAnnotationTaskCancelled
	}
	task.Status = AnnotationTaskStatusInProgress
	task.UpdatedAt = time.Now()
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *AnnotationService) CompleteTask(ctx context.Context, id int64) (*AnnotationTask, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	task.Status = AnnotationTaskStatusCompleted
	task.UpdatedAt = time.Now()
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *AnnotationService) CancelTask(ctx context.Context, id int64) (*AnnotationTask, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task.Status == AnnotationTaskStatusCompleted {
		return nil, ErrAnnotationTaskCompleted
	}
	task.Status = AnnotationTaskStatusCancelled
	task.UpdatedAt = time.Now()
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *AnnotationService) SubmitAnnotation(ctx context.Context, in SubmitAnnotationInput) (*AnnotationResult, error) {
	task, err := s.taskRepo.GetByID(ctx, in.TaskID)
	if err != nil {
		return nil, err
	}
	if task.Status == AnnotationTaskStatusCompleted {
		return nil, ErrAnnotationTaskCompleted
	}
	if task.Status == AnnotationTaskStatusCancelled {
		return nil, ErrAnnotationTaskCancelled
	}

	existing, _ := s.resultRepo.GetByTaskAndIndex(ctx, in.TaskID, in.ItemIndex)
	if existing != nil {
		return nil, ErrAnnotationResultExists
	}

	result := &AnnotationResult{
		ID:          snowflake.NextID(),
		TaskID:      in.TaskID,
		TenantID:    in.TenantID,
		ItemIndex:   in.ItemIndex,
		RawText:     in.RawText,
		Label:       in.Label,
		Metadata:    in.Metadata,
		AnnotatorID: in.AnnotatorID,
		CreatedAt:   time.Now(),
	}
	if err := s.resultRepo.Create(ctx, result); err != nil {
		return nil, err
	}

	task.LabeledItems++
	task.UpdatedAt = time.Now()
	if task.LabeledItems >= task.TotalItems {
		task.Status = AnnotationTaskStatusCompleted
	} else if task.Status == AnnotationTaskStatusPending {
		task.Status = AnnotationTaskStatusInProgress
	}
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *AnnotationService) ListResults(ctx context.Context, taskID int64) ([]AnnotationResult, error) {
	return s.resultRepo.ListByTaskID(ctx, taskID)
}

// LLMGatewayService manages multi-model LLM configurations.
type LLMGatewayService struct {
	configRepo LLMModelConfigRepository
}

func NewLLMGatewayService(configRepo LLMModelConfigRepository) *LLMGatewayService {
	return &LLMGatewayService{configRepo: configRepo}
}

func (s *LLMGatewayService) CreateConfig(ctx context.Context, in CreateLLMModelConfigInput) (*LLMModelConfig, error) {
	if !isValidProviderType(in.ProviderType) {
		return nil, ErrInvalidProviderType
	}
	config := &LLMModelConfig{
		ID:           snowflake.NextID(),
		TenantID:     in.TenantID,
		Name:         in.Name,
		ProviderType: in.ProviderType,
		Endpoint:     in.Endpoint,
		APIKey:       in.APIKey,
		ModelName:    in.ModelName,
		IsDefault:    in.IsDefault,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.configRepo.Create(ctx, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (s *LLMGatewayService) GetConfig(ctx context.Context, id int64) (*LLMModelConfig, error) {
	return s.configRepo.GetByID(ctx, id)
}

func (s *LLMGatewayService) ListConfigs(ctx context.Context, tenantID int64) ([]LLMModelConfig, error) {
	return s.configRepo.List(ctx, tenantID)
}

func (s *LLMGatewayService) UpdateConfig(ctx context.Context, config *LLMModelConfig) error {
	config.UpdatedAt = time.Now()
	return s.configRepo.Update(ctx, config)
}

func (s *LLMGatewayService) DeleteConfig(ctx context.Context, id int64) error {
	return s.configRepo.Delete(ctx, id)
}

func (s *LLMGatewayService) GetDefault(ctx context.Context, tenantID int64) (*LLMModelConfig, error) {
	return s.configRepo.GetDefault(ctx, tenantID)
}

func isValidAnnotationType(t AnnotationType) bool {
	switch t {
	case AnnotationTypeIntent, AnnotationTypeEntity, AnnotationTypeSentiment, AnnotationTypeFAQ:
		return true
	}
	return false
}

func isValidProviderType(t string) bool {
	switch t {
	case "tongyi", "openai", "bailian", "self_hosted":
		return true
	}
	return false
}
