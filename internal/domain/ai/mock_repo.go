package ai

import (
	"context"
	"strings"
	"sync"
)

type MockKnowledgeCategoryRepo struct {
	mu   sync.RWMutex
	data map[int64]*KnowledgeCategory
}

func NewMockKnowledgeCategoryRepo() *MockKnowledgeCategoryRepo {
	return &MockKnowledgeCategoryRepo{data: make(map[int64]*KnowledgeCategory)}
}

func (r *MockKnowledgeCategoryRepo) Create(_ context.Context, c *KnowledgeCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[c.ID] = c
	return nil
}

func (r *MockKnowledgeCategoryRepo) List(_ context.Context, tenantID int64) ([]*KnowledgeCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*KnowledgeCategory
	for _, c := range r.data {
		if c.TenantID == tenantID {
			result = append(result, c)
		}
	}
	return result, nil
}

type MockKnowledgeArticleRepo struct {
	mu   sync.RWMutex
	data map[int64]*KnowledgeArticle
}

func NewMockKnowledgeArticleRepo() *MockKnowledgeArticleRepo {
	return &MockKnowledgeArticleRepo{data: make(map[int64]*KnowledgeArticle)}
}

func (r *MockKnowledgeArticleRepo) Create(_ context.Context, a *KnowledgeArticle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[a.ID] = a
	return nil
}

func (r *MockKnowledgeArticleRepo) GetByID(_ context.Context, id int64) (*KnowledgeArticle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if a, ok := r.data[id]; ok {
		return a, nil
	}
	return nil, nil
}

func (r *MockKnowledgeArticleRepo) Update(_ context.Context, a *KnowledgeArticle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[a.ID] = a
	return nil
}

func (r *MockKnowledgeArticleRepo) List(_ context.Context, tenantID int64, offset, limit int) ([]*KnowledgeArticle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*KnowledgeArticle
	for _, a := range r.data {
		if a.TenantID == tenantID {
			result = append(result, a)
		}
	}
	if offset >= len(result) {
		return nil, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

func (r *MockKnowledgeArticleRepo) Search(_ context.Context, tenantID int64, query string, limit int) ([]*KnowledgeArticle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	q := strings.ToLower(query)
	var result []*KnowledgeArticle
	for _, a := range r.data {
		if a.TenantID != tenantID || a.Status != "published" {
			continue
		}
		if strings.Contains(strings.ToLower(a.Title), q) || strings.Contains(strings.ToLower(a.Content), q) || strings.Contains(strings.ToLower(a.Tags), q) {
			result = append(result, a)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

type MockAgentScriptRepo struct {
	mu   sync.RWMutex
	data map[int64]*AgentScript
}

func NewMockAgentScriptRepo() *MockAgentScriptRepo {
	return &MockAgentScriptRepo{data: make(map[int64]*AgentScript)}
}

func (r *MockAgentScriptRepo) Create(_ context.Context, s *AgentScript) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[s.ID] = s
	return nil
}

func (r *MockAgentScriptRepo) GetByID(_ context.Context, id int64) (*AgentScript, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.data[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (r *MockAgentScriptRepo) Update(_ context.Context, s *AgentScript) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[s.ID] = s
	return nil
}

func (r *MockAgentScriptRepo) List(_ context.Context, tenantID int64) ([]*AgentScript, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*AgentScript
	for _, s := range r.data {
		if s.TenantID == tenantID {
			result = append(result, s)
		}
	}
	return result, nil
}

type MockSessionInfoTemplateRepo struct {
	mu   sync.RWMutex
	data map[int64]*SessionInfoTemplate
}

func NewMockSessionInfoTemplateRepo() *MockSessionInfoTemplateRepo {
	return &MockSessionInfoTemplateRepo{data: make(map[int64]*SessionInfoTemplate)}
}

func (r *MockSessionInfoTemplateRepo) Create(_ context.Context, t *SessionInfoTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[t.ID] = t
	return nil
}

func (r *MockSessionInfoTemplateRepo) GetByID(_ context.Context, id int64) (*SessionInfoTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if t, ok := r.data[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (r *MockSessionInfoTemplateRepo) Update(_ context.Context, t *SessionInfoTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[t.ID] = t
	return nil
}

func (r *MockSessionInfoTemplateRepo) List(_ context.Context, tenantID int64) ([]*SessionInfoTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*SessionInfoTemplate
	for _, t := range r.data {
		if t.TenantID == tenantID {
			result = append(result, t)
		}
	}
	return result, nil
}
