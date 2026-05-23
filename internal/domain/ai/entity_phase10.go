package ai

import "time"

// --- Annotation Management ---

type AnnotationTaskStatus string

const (
	AnnotationTaskStatusPending    AnnotationTaskStatus = "pending"
	AnnotationTaskStatusInProgress AnnotationTaskStatus = "in_progress"
	AnnotationTaskStatusCompleted  AnnotationTaskStatus = "completed"
	AnnotationTaskStatusCancelled  AnnotationTaskStatus = "cancelled"
)

type AnnotationType string

const (
	AnnotationTypeIntent    AnnotationType = "intent"
	AnnotationTypeEntity    AnnotationType = "entity"
	AnnotationTypeSentiment AnnotationType = "sentiment"
	AnnotationTypeFAQ       AnnotationType = "faq"
)

// AnnotationTask represents a data labeling task for digital employee training.
type AnnotationTask struct {
	ID           int64                `db:"id" json:"id"`
	TenantID     int64                `db:"tenant_id" json:"tenant_id"`
	Name         string               `db:"name" json:"name"`
	Type         AnnotationType       `db:"type" json:"type"`
	DatasetID    string               `db:"dataset_id" json:"dataset_id"`
	AssigneeID   *int64               `db:"assignee_id" json:"assignee_id,omitempty"`
	TotalItems   int                  `db:"total_items" json:"total_items"`
	LabeledItems int                  `db:"labeled_items" json:"labeled_items"`
	Status       AnnotationTaskStatus `db:"status" json:"status"`
	CreatedAt    time.Time            `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `db:"updated_at" json:"updated_at"`
}

// AnnotationResult represents a single labeling result within a task.
type AnnotationResult struct {
	ID         int64     `db:"id" json:"id"`
	TaskID     int64     `db:"task_id" json:"task_id"`
	TenantID   int64     `db:"tenant_id" json:"tenant_id"`
	ItemIndex  int       `db:"item_index" json:"item_index"`
	RawText    string    `db:"raw_text" json:"raw_text"`
	Label      string    `db:"label" json:"label"`
	Metadata   string    `db:"metadata" json:"metadata"` // JSON
	AnnotatorID int64    `db:"annotator_id" json:"annotator_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// --- LLM Gateway ---

type LLMModelConfig struct {
	ID         int64     `db:"id" json:"id"`
	TenantID   int64     `db:"tenant_id" json:"tenant_id"`
	Name       string    `db:"name" json:"name"`
	ProviderType string  `db:"provider_type" json:"provider_type"` // tongyi, openai, bailian, self_hosted
	Endpoint   string    `db:"endpoint" json:"endpoint"`
	APIKey     string    `db:"api_key" json:"-"`
	ModelName  string    `db:"model_name" json:"model_name"`
	IsDefault  bool      `db:"is_default" json:"is_default"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// --- Inputs ---

type CreateAnnotationTaskInput struct {
	TenantID   int64          `json:"tenant_id"`
	Name       string         `json:"name"`
	Type       AnnotationType `json:"type"`
	DatasetID  string         `json:"dataset_id"`
	AssigneeID *int64         `json:"assignee_id,omitempty"`
	TotalItems int            `json:"total_items"`
}

type SubmitAnnotationInput struct {
	TaskID      int64  `json:"task_id"`
	TenantID    int64  `json:"tenant_id"`
	ItemIndex   int    `json:"item_index"`
	RawText     string `json:"raw_text"`
	Label       string `json:"label"`
	Metadata    string `json:"metadata"`
	AnnotatorID int64  `json:"annotator_id"`
}

type CreateLLMModelConfigInput struct {
	TenantID     int64  `json:"tenant_id"`
	Name         string `json:"name"`
	ProviderType string `json:"provider_type"`
	Endpoint     string `json:"endpoint"`
	APIKey       string `json:"api_key"`
	ModelName    string `json:"model_name"`
	IsDefault    bool   `json:"is_default"`
}
