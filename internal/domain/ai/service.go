package ai

import (
	"context"
	"encoding/json"
	"time"

	"github.com/divord97/ccc/pkg/snowflake"
)

// KnowledgeService manages knowledge base categories and articles.
type KnowledgeService struct {
	categories KnowledgeCategoryRepository
	articles   KnowledgeArticleRepository
}

func NewKnowledgeService(categories KnowledgeCategoryRepository, articles KnowledgeArticleRepository) *KnowledgeService {
	return &KnowledgeService{categories: categories, articles: articles}
}

type CreateArticleInput struct {
	TenantID   int64  `json:"tenant_id"`
	CategoryID *int64 `json:"category_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Tags       string `json:"tags"`
	Status     string `json:"status"`
}

func (s *KnowledgeService) CreateArticle(ctx context.Context, in CreateArticleInput) (*KnowledgeArticle, error) {
	if in.Status == "" {
		in.Status = "draft"
	}
	now := time.Now()
	a := &KnowledgeArticle{
		ID:         snowflake.NextID(),
		TenantID:   in.TenantID,
		CategoryID: in.CategoryID,
		Title:      in.Title,
		Content:    in.Content,
		Tags:       in.Tags,
		Status:     in.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.articles.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *KnowledgeService) GetArticle(ctx context.Context, id int64) (*KnowledgeArticle, error) {
	a, err := s.articles.GetByID(ctx, id)
	if err != nil || a == nil {
		return nil, ErrArticleNotFound
	}
	return a, nil
}

func (s *KnowledgeService) UpdateArticle(ctx context.Context, a *KnowledgeArticle) error {
	a.UpdatedAt = time.Now()
	return s.articles.Update(ctx, a)
}

func (s *KnowledgeService) ListArticles(ctx context.Context, tenantID int64, offset, limit int) ([]*KnowledgeArticle, error) {
	return s.articles.List(ctx, tenantID, offset, limit)
}

func (s *KnowledgeService) Search(ctx context.Context, tenantID int64, query string, limit int) ([]*KnowledgeArticle, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.articles.Search(ctx, tenantID, query, limit)
}

func (s *KnowledgeService) CreateCategory(ctx context.Context, tenantID int64, name string) (*KnowledgeCategory, error) {
	c := &KnowledgeCategory{
		ID:       snowflake.NextID(),
		TenantID: tenantID,
		Name:     name,
	}
	if err := s.categories.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *KnowledgeService) ListCategories(ctx context.Context, tenantID int64) ([]*KnowledgeCategory, error) {
	return s.categories.List(ctx, tenantID)
}

// AgentScriptService manages agent script templates.
type AgentScriptService struct {
	scripts AgentScriptRepository
}

func NewAgentScriptService(scripts AgentScriptRepository) *AgentScriptService {
	return &AgentScriptService{scripts: scripts}
}

type CreateScriptInput struct {
	TenantID int64  `json:"tenant_id"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	IsActive bool   `json:"is_active"`
}

func (s *AgentScriptService) Create(ctx context.Context, in CreateScriptInput) (*AgentScript, error) {
	if in.Content != "" && !json.Valid([]byte(in.Content)) {
		return nil, ErrInvalidScript
	}

	now := time.Now()
	script := &AgentScript{
		ID:        snowflake.NextID(),
		TenantID:  in.TenantID,
		Name:      in.Name,
		Content:   in.Content,
		IsActive:  in.IsActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.scripts.Create(ctx, script); err != nil {
		return nil, err
	}
	return script, nil
}

func (s *AgentScriptService) GetByID(ctx context.Context, id int64) (*AgentScript, error) {
	script, err := s.scripts.GetByID(ctx, id)
	if err != nil || script == nil {
		return nil, ErrScriptNotFound
	}
	return script, nil
}

func (s *AgentScriptService) Update(ctx context.Context, script *AgentScript) error {
	script.UpdatedAt = time.Now()
	return s.scripts.Update(ctx, script)
}

func (s *AgentScriptService) List(ctx context.Context, tenantID int64) ([]*AgentScript, error) {
	return s.scripts.List(ctx, tenantID)
}

// SessionInfoTemplateService manages session info templates.
type SessionInfoTemplateService struct {
	templates SessionInfoTemplateRepository
}

func NewSessionInfoTemplateService(templates SessionInfoTemplateRepository) *SessionInfoTemplateService {
	return &SessionInfoTemplateService{templates: templates}
}

func (s *SessionInfoTemplateService) Create(ctx context.Context, tenantID int64, name, fields string, isDefault bool) (*SessionInfoTemplate, error) {
	now := time.Now()
	t := &SessionInfoTemplate{
		ID:        snowflake.NextID(),
		TenantID:  tenantID,
		Name:      name,
		Fields:    fields,
		IsDefault: isDefault,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.templates.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *SessionInfoTemplateService) GetByID(ctx context.Context, id int64) (*SessionInfoTemplate, error) {
	t, err := s.templates.GetByID(ctx, id)
	if err != nil || t == nil {
		return nil, ErrTemplateNotFound
	}
	return t, nil
}

func (s *SessionInfoTemplateService) Update(ctx context.Context, t *SessionInfoTemplate) error {
	t.UpdatedAt = time.Now()
	return s.templates.Update(ctx, t)
}

func (s *SessionInfoTemplateService) List(ctx context.Context, tenantID int64) ([]*SessionInfoTemplate, error) {
	return s.templates.List(ctx, tenantID)
}
