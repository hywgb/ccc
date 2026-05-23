package mysql

import (
	"context"
	"database/sql"

	"github.com/divord97/ccc/internal/domain/ai"
	"github.com/jmoiron/sqlx"
)

// --- AnnotationTaskRepo ---

type AnnotationTaskRepo struct{ db *sqlx.DB }

func NewAnnotationTaskRepo(db *sqlx.DB) *AnnotationTaskRepo {
	return &AnnotationTaskRepo{db: db}
}

func (r *AnnotationTaskRepo) Create(ctx context.Context, t *ai.AnnotationTask) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO annotation_tasks (id, tenant_id, name, type, dataset_id, assignee_id, total_items, labeled_items, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.TenantID, t.Name, t.Type, t.DatasetID, t.AssigneeID, t.TotalItems, t.LabeledItems, t.Status, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *AnnotationTaskRepo) GetByID(ctx context.Context, id int64) (*ai.AnnotationTask, error) {
	var t ai.AnnotationTask
	if err := r.db.GetContext(ctx, &t, "SELECT * FROM annotation_tasks WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ai.ErrAnnotationTaskNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *AnnotationTaskRepo) Update(ctx context.Context, t *ai.AnnotationTask) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE annotation_tasks SET name=?, assignee_id=?, labeled_items=?, status=?, updated_at=? WHERE id=?`,
		t.Name, t.AssigneeID, t.LabeledItems, t.Status, t.UpdatedAt, t.ID)
	return err
}

func (r *AnnotationTaskRepo) List(ctx context.Context, tenantID int64) ([]ai.AnnotationTask, error) {
	var out []ai.AnnotationTask
	if err := r.db.SelectContext(ctx, &out, "SELECT * FROM annotation_tasks WHERE tenant_id = ? ORDER BY created_at DESC", tenantID); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AnnotationTaskRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM annotation_tasks WHERE id = ?", id)
	return err
}

// --- AnnotationResultRepo ---

type AnnotationResultRepo struct{ db *sqlx.DB }

func NewAnnotationResultRepo(db *sqlx.DB) *AnnotationResultRepo {
	return &AnnotationResultRepo{db: db}
}

func (r *AnnotationResultRepo) Create(ctx context.Context, res *ai.AnnotationResult) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO annotation_results (id, task_id, tenant_id, item_index, raw_text, label, metadata, annotator_id, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		res.ID, res.TaskID, res.TenantID, res.ItemIndex, res.RawText, res.Label, res.Metadata, res.AnnotatorID, res.CreatedAt)
	return err
}

func (r *AnnotationResultRepo) ListByTaskID(ctx context.Context, taskID int64) ([]ai.AnnotationResult, error) {
	var out []ai.AnnotationResult
	if err := r.db.SelectContext(ctx, &out, "SELECT * FROM annotation_results WHERE task_id = ? ORDER BY item_index", taskID); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AnnotationResultRepo) GetByTaskAndIndex(ctx context.Context, taskID int64, itemIndex int) (*ai.AnnotationResult, error) {
	var res ai.AnnotationResult
	if err := r.db.GetContext(ctx, &res, "SELECT * FROM annotation_results WHERE task_id = ? AND item_index = ?", taskID, itemIndex); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

// --- LLMModelConfigRepo ---

type LLMModelConfigRepo struct{ db *sqlx.DB }

func NewLLMModelConfigRepo(db *sqlx.DB) *LLMModelConfigRepo {
	return &LLMModelConfigRepo{db: db}
}

func (r *LLMModelConfigRepo) Create(ctx context.Context, c *ai.LLMModelConfig) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO llm_model_configs (id, tenant_id, name, provider_type, endpoint, api_key, model_name, is_default, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.TenantID, c.Name, c.ProviderType, c.Endpoint, c.APIKey, c.ModelName, c.IsDefault, c.IsActive, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *LLMModelConfigRepo) GetByID(ctx context.Context, id int64) (*ai.LLMModelConfig, error) {
	var c ai.LLMModelConfig
	if err := r.db.GetContext(ctx, &c, "SELECT * FROM llm_model_configs WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ai.ErrLLMModelConfigNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *LLMModelConfigRepo) Update(ctx context.Context, c *ai.LLMModelConfig) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE llm_model_configs SET name=?, provider_type=?, endpoint=?, api_key=?, model_name=?, is_default=?, is_active=?, updated_at=? WHERE id=?`,
		c.Name, c.ProviderType, c.Endpoint, c.APIKey, c.ModelName, c.IsDefault, c.IsActive, c.UpdatedAt, c.ID)
	return err
}

func (r *LLMModelConfigRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM llm_model_configs WHERE id = ?", id)
	return err
}

func (r *LLMModelConfigRepo) List(ctx context.Context, tenantID int64) ([]ai.LLMModelConfig, error) {
	var out []ai.LLMModelConfig
	if err := r.db.SelectContext(ctx, &out, "SELECT * FROM llm_model_configs WHERE tenant_id = ? ORDER BY created_at DESC", tenantID); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *LLMModelConfigRepo) GetDefault(ctx context.Context, tenantID int64) (*ai.LLMModelConfig, error) {
	var c ai.LLMModelConfig
	if err := r.db.GetContext(ctx, &c, "SELECT * FROM llm_model_configs WHERE tenant_id = ? AND is_default = 1 AND is_active = 1 LIMIT 1", tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ai.ErrNoDefaultLLMModel
		}
		return nil, err
	}
	return &c, nil
}
