package ai

import (
	"context"
	"sync"
)

// --- MockAnnotationTaskRepo ---

type MockAnnotationTaskRepo struct {
	mu    sync.RWMutex
	tasks map[int64]*AnnotationTask
	seq   int64
}

func NewMockAnnotationTaskRepo() *MockAnnotationTaskRepo {
	return &MockAnnotationTaskRepo{tasks: make(map[int64]*AnnotationTask)}
}

func (m *MockAnnotationTaskRepo) Create(_ context.Context, t *AnnotationTask) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seq++
	t.ID = m.seq
	clone := *t
	m.tasks[t.ID] = &clone
	return nil
}

func (m *MockAnnotationTaskRepo) GetByID(_ context.Context, id int64) (*AnnotationTask, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.tasks[id]
	if !ok {
		return nil, ErrAnnotationTaskNotFound
	}
	clone := *t
	return &clone, nil
}

func (m *MockAnnotationTaskRepo) Update(_ context.Context, t *AnnotationTask) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	clone := *t
	m.tasks[t.ID] = &clone
	return nil
}

func (m *MockAnnotationTaskRepo) List(_ context.Context, tenantID int64) ([]AnnotationTask, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []AnnotationTask
	for _, t := range m.tasks {
		if t.TenantID == tenantID {
			out = append(out, *t)
		}
	}
	return out, nil
}

func (m *MockAnnotationTaskRepo) Delete(_ context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tasks, id)
	return nil
}

// --- MockAnnotationResultRepo ---

type MockAnnotationResultRepo struct {
	mu      sync.RWMutex
	results map[int64]*AnnotationResult
	seq     int64
}

func NewMockAnnotationResultRepo() *MockAnnotationResultRepo {
	return &MockAnnotationResultRepo{results: make(map[int64]*AnnotationResult)}
}

func (m *MockAnnotationResultRepo) Create(_ context.Context, r *AnnotationResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seq++
	r.ID = m.seq
	clone := *r
	m.results[r.ID] = &clone
	return nil
}

func (m *MockAnnotationResultRepo) ListByTaskID(_ context.Context, taskID int64) ([]AnnotationResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []AnnotationResult
	for _, r := range m.results {
		if r.TaskID == taskID {
			out = append(out, *r)
		}
	}
	return out, nil
}

func (m *MockAnnotationResultRepo) GetByTaskAndIndex(_ context.Context, taskID int64, itemIndex int) (*AnnotationResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, r := range m.results {
		if r.TaskID == taskID && r.ItemIndex == itemIndex {
			clone := *r
			return &clone, nil
		}
	}
	return nil, nil
}

// --- MockLLMModelConfigRepo ---

type MockLLMModelConfigRepo struct {
	mu      sync.RWMutex
	configs map[int64]*LLMModelConfig
	seq     int64
}

func NewMockLLMModelConfigRepo() *MockLLMModelConfigRepo {
	return &MockLLMModelConfigRepo{configs: make(map[int64]*LLMModelConfig)}
}

func (m *MockLLMModelConfigRepo) Create(_ context.Context, c *LLMModelConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seq++
	c.ID = m.seq
	clone := *c
	m.configs[c.ID] = &clone
	return nil
}

func (m *MockLLMModelConfigRepo) GetByID(_ context.Context, id int64) (*LLMModelConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.configs[id]
	if !ok {
		return nil, ErrLLMModelConfigNotFound
	}
	clone := *c
	return &clone, nil
}

func (m *MockLLMModelConfigRepo) Update(_ context.Context, c *LLMModelConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	clone := *c
	m.configs[c.ID] = &clone
	return nil
}

func (m *MockLLMModelConfigRepo) Delete(_ context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.configs, id)
	return nil
}

func (m *MockLLMModelConfigRepo) List(_ context.Context, tenantID int64) ([]LLMModelConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []LLMModelConfig
	for _, c := range m.configs {
		if c.TenantID == tenantID {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (m *MockLLMModelConfigRepo) GetDefault(_ context.Context, tenantID int64) (*LLMModelConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, c := range m.configs {
		if c.TenantID == tenantID && c.IsDefault && c.IsActive {
			clone := *c
			return &clone, nil
		}
	}
	return nil, ErrNoDefaultLLMModel
}
