package ai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Annotation Service Tests ---

func TestAnnotation_CreateTask(t *testing.T) {
	svc := NewAnnotationService(NewMockAnnotationTaskRepo(), NewMockAnnotationResultRepo())
	task, err := svc.CreateTask(context.Background(), CreateAnnotationTaskInput{
		TenantID:   1,
		Name:       "Intent labeling batch 1",
		Type:       AnnotationTypeIntent,
		DatasetID:  "ds-001",
		TotalItems: 100,
	})
	require.NoError(t, err)
	assert.Equal(t, "Intent labeling batch 1", task.Name)
	assert.Equal(t, AnnotationTaskStatusPending, task.Status)
	assert.Equal(t, 0, task.LabeledItems)
	assert.Equal(t, 100, task.TotalItems)
}

func TestAnnotation_InvalidType(t *testing.T) {
	svc := NewAnnotationService(NewMockAnnotationTaskRepo(), NewMockAnnotationResultRepo())
	_, err := svc.CreateTask(context.Background(), CreateAnnotationTaskInput{
		TenantID:   1,
		Name:       "Bad type",
		Type:       "invalid_type",
		DatasetID:  "ds-001",
		TotalItems: 10,
	})
	assert.ErrorIs(t, err, ErrInvalidAnnotationType)
}

func TestAnnotation_TaskLifecycle(t *testing.T) {
	svc := NewAnnotationService(NewMockAnnotationTaskRepo(), NewMockAnnotationResultRepo())
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, CreateAnnotationTaskInput{
		TenantID: 1, Name: "Lifecycle test", Type: AnnotationTypeSentiment,
		DatasetID: "ds-002", TotalItems: 5,
	})
	require.NoError(t, err)
	assert.Equal(t, AnnotationTaskStatusPending, task.Status)

	task, err = svc.StartTask(ctx, task.ID)
	require.NoError(t, err)
	assert.Equal(t, AnnotationTaskStatusInProgress, task.Status)

	task, err = svc.CompleteTask(ctx, task.ID)
	require.NoError(t, err)
	assert.Equal(t, AnnotationTaskStatusCompleted, task.Status)

	// Cannot start a completed task
	_, err = svc.StartTask(ctx, task.ID)
	assert.ErrorIs(t, err, ErrAnnotationTaskCompleted)
}

func TestAnnotation_CancelTask(t *testing.T) {
	svc := NewAnnotationService(NewMockAnnotationTaskRepo(), NewMockAnnotationResultRepo())
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, CreateAnnotationTaskInput{
		TenantID: 1, Name: "Cancel test", Type: AnnotationTypeFAQ,
		DatasetID: "ds-003", TotalItems: 10,
	})
	require.NoError(t, err)

	task, err = svc.CancelTask(ctx, task.ID)
	require.NoError(t, err)
	assert.Equal(t, AnnotationTaskStatusCancelled, task.Status)

	// Cannot submit to cancelled task
	_, err = svc.SubmitAnnotation(ctx, SubmitAnnotationInput{
		TaskID: task.ID, TenantID: 1, ItemIndex: 0, RawText: "hello",
		Label: "greeting", AnnotatorID: 1,
	})
	assert.ErrorIs(t, err, ErrAnnotationTaskCancelled)
}

func TestAnnotation_SubmitAndAutoComplete(t *testing.T) {
	svc := NewAnnotationService(NewMockAnnotationTaskRepo(), NewMockAnnotationResultRepo())
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, CreateAnnotationTaskInput{
		TenantID: 1, Name: "Auto-complete", Type: AnnotationTypeEntity,
		DatasetID: "ds-004", TotalItems: 2,
	})
	require.NoError(t, err)

	// Submit first item — should transition to in_progress
	_, err = svc.SubmitAnnotation(ctx, SubmitAnnotationInput{
		TaskID: task.ID, TenantID: 1, ItemIndex: 0, RawText: "item1",
		Label: "person", AnnotatorID: 1,
	})
	require.NoError(t, err)

	task, _ = svc.GetTask(ctx, task.ID)
	assert.Equal(t, AnnotationTaskStatusInProgress, task.Status)
	assert.Equal(t, 1, task.LabeledItems)

	// Submit second item — should auto-complete
	_, err = svc.SubmitAnnotation(ctx, SubmitAnnotationInput{
		TaskID: task.ID, TenantID: 1, ItemIndex: 1, RawText: "item2",
		Label: "location", AnnotatorID: 1,
	})
	require.NoError(t, err)

	task, _ = svc.GetTask(ctx, task.ID)
	assert.Equal(t, AnnotationTaskStatusCompleted, task.Status)
	assert.Equal(t, 2, task.LabeledItems)
}

func TestAnnotation_DuplicateSubmit(t *testing.T) {
	svc := NewAnnotationService(NewMockAnnotationTaskRepo(), NewMockAnnotationResultRepo())
	ctx := context.Background()

	task, _ := svc.CreateTask(ctx, CreateAnnotationTaskInput{
		TenantID: 1, Name: "Dup test", Type: AnnotationTypeIntent,
		DatasetID: "ds-005", TotalItems: 10,
	})

	_, err := svc.SubmitAnnotation(ctx, SubmitAnnotationInput{
		TaskID: task.ID, TenantID: 1, ItemIndex: 0, RawText: "text",
		Label: "intent_a", AnnotatorID: 1,
	})
	require.NoError(t, err)

	// Duplicate submit for same item index
	_, err = svc.SubmitAnnotation(ctx, SubmitAnnotationInput{
		TaskID: task.ID, TenantID: 1, ItemIndex: 0, RawText: "text",
		Label: "intent_b", AnnotatorID: 1,
	})
	assert.ErrorIs(t, err, ErrAnnotationResultExists)
}

func TestAnnotation_ListResults(t *testing.T) {
	svc := NewAnnotationService(NewMockAnnotationTaskRepo(), NewMockAnnotationResultRepo())
	ctx := context.Background()

	task, _ := svc.CreateTask(ctx, CreateAnnotationTaskInput{
		TenantID: 1, Name: "List test", Type: AnnotationTypeIntent,
		DatasetID: "ds-006", TotalItems: 10,
	})

	for i := 0; i < 3; i++ {
		svc.SubmitAnnotation(ctx, SubmitAnnotationInput{
			TaskID: task.ID, TenantID: 1, ItemIndex: i, RawText: "text",
			Label: "label", AnnotatorID: 1,
		})
	}

	results, err := svc.ListResults(ctx, task.ID)
	require.NoError(t, err)
	assert.Len(t, results, 3)
}

// --- LLM Gateway Service Tests ---

func TestLLMGateway_CreateConfig(t *testing.T) {
	svc := NewLLMGatewayService(NewMockLLMModelConfigRepo())
	config, err := svc.CreateConfig(context.Background(), CreateLLMModelConfigInput{
		TenantID:     1,
		Name:         "Tongyi Qianwen",
		ProviderType: "tongyi",
		Endpoint:     "https://dashscope.aliyuncs.com",
		APIKey:       "sk-test",
		ModelName:    "qwen-max",
		IsDefault:    true,
	})
	require.NoError(t, err)
	assert.Equal(t, "Tongyi Qianwen", config.Name)
	assert.True(t, config.IsDefault)
	assert.True(t, config.IsActive)
}

func TestLLMGateway_InvalidProviderType(t *testing.T) {
	svc := NewLLMGatewayService(NewMockLLMModelConfigRepo())
	_, err := svc.CreateConfig(context.Background(), CreateLLMModelConfigInput{
		TenantID: 1, Name: "Bad", ProviderType: "chatgpt", ModelName: "gpt-4",
	})
	assert.ErrorIs(t, err, ErrInvalidProviderType)
}

func TestLLMGateway_GetDefault(t *testing.T) {
	svc := NewLLMGatewayService(NewMockLLMModelConfigRepo())
	ctx := context.Background()

	svc.CreateConfig(ctx, CreateLLMModelConfigInput{
		TenantID: 1, Name: "Model A", ProviderType: "tongyi",
		ModelName: "qwen-turbo", IsDefault: false,
	})
	svc.CreateConfig(ctx, CreateLLMModelConfigInput{
		TenantID: 1, Name: "Model B", ProviderType: "bailian",
		ModelName: "qwen-max", IsDefault: true,
	})

	def, err := svc.GetDefault(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "Model B", def.Name)
}

func TestLLMGateway_NoDefault(t *testing.T) {
	svc := NewLLMGatewayService(NewMockLLMModelConfigRepo())
	_, err := svc.GetDefault(context.Background(), 999)
	assert.ErrorIs(t, err, ErrNoDefaultLLMModel)
}


