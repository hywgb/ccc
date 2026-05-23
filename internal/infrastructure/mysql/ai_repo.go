package mysql

import (
	"context"
	"database/sql"

	"github.com/divord97/ccc/internal/domain/ai"
	"github.com/jmoiron/sqlx"
)

type KnowledgeCategoryRepo struct{ db *sqlx.DB }

func NewKnowledgeCategoryRepo(db *sqlx.DB) *KnowledgeCategoryRepo {
	return &KnowledgeCategoryRepo{db: db}
}

func (r *KnowledgeCategoryRepo) Create(ctx context.Context, c *ai.KnowledgeCategory) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO knowledge_categories (id, tenant_id, name, parent_id) VALUES (?,?,?,?)`,
		c.ID, c.TenantID, c.Name, c.ParentID)
	return err
}

func (r *KnowledgeCategoryRepo) List(ctx context.Context, tenantID int64) ([]*ai.KnowledgeCategory, error) {
	var result []*ai.KnowledgeCategory
	err := r.db.SelectContext(ctx, &result, `SELECT * FROM knowledge_categories WHERE tenant_id=?`, tenantID)
	return result, err
}

type KnowledgeArticleRepo struct{ db *sqlx.DB }

func NewKnowledgeArticleRepo(db *sqlx.DB) *KnowledgeArticleRepo {
	return &KnowledgeArticleRepo{db: db}
}

func (r *KnowledgeArticleRepo) Create(ctx context.Context, a *ai.KnowledgeArticle) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO knowledge_articles (id, tenant_id, category_id, title, content, tags, status, created_at, updated_at)
		 VALUES (?,?,?,?,?,?,?,?,?)`,
		a.ID, a.TenantID, a.CategoryID, a.Title, a.Content, a.Tags, a.Status, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *KnowledgeArticleRepo) GetByID(ctx context.Context, id int64) (*ai.KnowledgeArticle, error) {
	var a ai.KnowledgeArticle
	if err := r.db.GetContext(ctx, &a, `SELECT * FROM knowledge_articles WHERE id=?`, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *KnowledgeArticleRepo) Update(ctx context.Context, a *ai.KnowledgeArticle) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE knowledge_articles SET category_id=?, title=?, content=?, tags=?, status=?, updated_at=? WHERE id=?`,
		a.CategoryID, a.Title, a.Content, a.Tags, a.Status, a.UpdatedAt, a.ID)
	return err
}

func (r *KnowledgeArticleRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]*ai.KnowledgeArticle, error) {
	var result []*ai.KnowledgeArticle
	err := r.db.SelectContext(ctx, &result,
		`SELECT * FROM knowledge_articles WHERE tenant_id=? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		tenantID, limit, offset)
	return result, err
}

func (r *KnowledgeArticleRepo) Search(ctx context.Context, tenantID int64, query string, limit int) ([]*ai.KnowledgeArticle, error) {
	var result []*ai.KnowledgeArticle
	q := "%" + query + "%"
	err := r.db.SelectContext(ctx, &result,
		`SELECT * FROM knowledge_articles WHERE tenant_id=? AND status='published'
		 AND (title LIKE ? OR content LIKE ? OR tags LIKE ?) LIMIT ?`,
		tenantID, q, q, q, limit)
	return result, err
}

type AgentScriptRepo struct{ db *sqlx.DB }

func NewAgentScriptRepo(db *sqlx.DB) *AgentScriptRepo { return &AgentScriptRepo{db: db} }

func (r *AgentScriptRepo) Create(ctx context.Context, s *ai.AgentScript) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO agent_scripts (id, tenant_id, name, content, is_active, created_at, updated_at)
		 VALUES (?,?,?,?,?,?,?)`,
		s.ID, s.TenantID, s.Name, s.Content, s.IsActive, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *AgentScriptRepo) GetByID(ctx context.Context, id int64) (*ai.AgentScript, error) {
	var s ai.AgentScript
	if err := r.db.GetContext(ctx, &s, `SELECT * FROM agent_scripts WHERE id=?`, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *AgentScriptRepo) Update(ctx context.Context, s *ai.AgentScript) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE agent_scripts SET name=?, content=?, is_active=?, updated_at=? WHERE id=?`,
		s.Name, s.Content, s.IsActive, s.UpdatedAt, s.ID)
	return err
}

func (r *AgentScriptRepo) List(ctx context.Context, tenantID int64) ([]*ai.AgentScript, error) {
	var result []*ai.AgentScript
	err := r.db.SelectContext(ctx, &result, `SELECT * FROM agent_scripts WHERE tenant_id=?`, tenantID)
	return result, err
}

type SessionInfoTemplateRepo struct{ db *sqlx.DB }

func NewSessionInfoTemplateRepo(db *sqlx.DB) *SessionInfoTemplateRepo {
	return &SessionInfoTemplateRepo{db: db}
}

func (r *SessionInfoTemplateRepo) Create(ctx context.Context, t *ai.SessionInfoTemplate) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO session_info_templates (id, tenant_id, name, fields, is_default, created_at, updated_at)
		 VALUES (?,?,?,?,?,?,?)`,
		t.ID, t.TenantID, t.Name, t.Fields, t.IsDefault, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *SessionInfoTemplateRepo) GetByID(ctx context.Context, id int64) (*ai.SessionInfoTemplate, error) {
	var t ai.SessionInfoTemplate
	if err := r.db.GetContext(ctx, &t, `SELECT * FROM session_info_templates WHERE id=?`, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *SessionInfoTemplateRepo) Update(ctx context.Context, t *ai.SessionInfoTemplate) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE session_info_templates SET name=?, fields=?, is_default=?, updated_at=? WHERE id=?`,
		t.Name, t.Fields, t.IsDefault, t.UpdatedAt, t.ID)
	return err
}

func (r *SessionInfoTemplateRepo) List(ctx context.Context, tenantID int64) ([]*ai.SessionInfoTemplate, error) {
	var result []*ai.SessionInfoTemplate
	err := r.db.SelectContext(ctx, &result, `SELECT * FROM session_info_templates WHERE tenant_id=?`, tenantID)
	return result, err
}
