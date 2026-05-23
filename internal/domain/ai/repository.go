package ai

import "context"

type KnowledgeCategoryRepository interface {
	Create(ctx context.Context, c *KnowledgeCategory) error
	List(ctx context.Context, tenantID int64) ([]*KnowledgeCategory, error)
}

type KnowledgeArticleRepository interface {
	Create(ctx context.Context, a *KnowledgeArticle) error
	GetByID(ctx context.Context, id int64) (*KnowledgeArticle, error)
	Update(ctx context.Context, a *KnowledgeArticle) error
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*KnowledgeArticle, error)
	Search(ctx context.Context, tenantID int64, query string, limit int) ([]*KnowledgeArticle, error)
}

type AgentScriptRepository interface {
	Create(ctx context.Context, s *AgentScript) error
	GetByID(ctx context.Context, id int64) (*AgentScript, error)
	Update(ctx context.Context, s *AgentScript) error
	List(ctx context.Context, tenantID int64) ([]*AgentScript, error)
}

type SessionInfoTemplateRepository interface {
	Create(ctx context.Context, t *SessionInfoTemplate) error
	GetByID(ctx context.Context, id int64) (*SessionInfoTemplate, error)
	Update(ctx context.Context, t *SessionInfoTemplate) error
	List(ctx context.Context, tenantID int64) ([]*SessionInfoTemplate, error)
}
