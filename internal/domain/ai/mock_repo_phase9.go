package ai

import (
	"context"
	"sync"
)

// --- MockDigitalEmployeeRepo ---

type MockDigitalEmployeeRepo struct {
	mu   sync.RWMutex
	data map[int64]*DigitalEmployee
}

func NewMockDigitalEmployeeRepo() *MockDigitalEmployeeRepo {
	return &MockDigitalEmployeeRepo{data: make(map[int64]*DigitalEmployee)}
}

func (r *MockDigitalEmployeeRepo) Create(_ context.Context, de *DigitalEmployee) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[de.ID] = de
	return nil
}

func (r *MockDigitalEmployeeRepo) GetByID(_ context.Context, id int64) (*DigitalEmployee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if de, ok := r.data[id]; ok {
		return de, nil
	}
	return nil, nil
}

func (r *MockDigitalEmployeeRepo) Update(_ context.Context, de *DigitalEmployee) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[de.ID] = de
	return nil
}

func (r *MockDigitalEmployeeRepo) List(_ context.Context, tenantID int64) ([]*DigitalEmployee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*DigitalEmployee
	for _, de := range r.data {
		if de.TenantID == tenantID {
			result = append(result, de)
		}
	}
	return result, nil
}

// --- MockDigitalEmployeeSceneRepo ---

type MockDigitalEmployeeSceneRepo struct {
	mu   sync.RWMutex
	data map[int64]*DigitalEmployeeScene
}

func NewMockDigitalEmployeeSceneRepo() *MockDigitalEmployeeSceneRepo {
	return &MockDigitalEmployeeSceneRepo{data: make(map[int64]*DigitalEmployeeScene)}
}

func (r *MockDigitalEmployeeSceneRepo) Create(_ context.Context, s *DigitalEmployeeScene) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[s.ID] = s
	return nil
}

func (r *MockDigitalEmployeeSceneRepo) GetByID(_ context.Context, id int64) (*DigitalEmployeeScene, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.data[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (r *MockDigitalEmployeeSceneRepo) Update(_ context.Context, s *DigitalEmployeeScene) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[s.ID] = s
	return nil
}

func (r *MockDigitalEmployeeSceneRepo) List(_ context.Context, digitalEmployeeID int64) ([]*DigitalEmployeeScene, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*DigitalEmployeeScene
	for _, s := range r.data {
		if s.DigitalEmployeeID == digitalEmployeeID {
			result = append(result, s)
		}
	}
	return result, nil
}

// --- MockQARuleRepo ---

type MockQARuleRepo struct {
	mu   sync.RWMutex
	data map[int64]*QARule
}

func NewMockQARuleRepo() *MockQARuleRepo {
	return &MockQARuleRepo{data: make(map[int64]*QARule)}
}

func (r *MockQARuleRepo) Create(_ context.Context, rule *QARule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[rule.ID] = rule
	return nil
}

func (r *MockQARuleRepo) GetByID(_ context.Context, id int64) (*QARule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if rule, ok := r.data[id]; ok {
		return rule, nil
	}
	return nil, nil
}

func (r *MockQARuleRepo) Update(_ context.Context, rule *QARule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[rule.ID] = rule
	return nil
}

func (r *MockQARuleRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}

func (r *MockQARuleRepo) List(_ context.Context, tenantID int64) ([]*QARule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*QARule
	for _, rule := range r.data {
		if rule.TenantID == tenantID {
			result = append(result, rule)
		}
	}
	return result, nil
}

func (r *MockQARuleRepo) ListByIDs(_ context.Context, ids []int64) ([]*QARule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	idSet := make(map[int64]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}
	var result []*QARule
	for _, rule := range r.data {
		if idSet[rule.ID] {
			result = append(result, rule)
		}
	}
	return result, nil
}

// --- MockQASchemeRepo ---

type MockQASchemeRepo struct {
	mu   sync.RWMutex
	data map[int64]*QAScheme
}

func NewMockQASchemeRepo() *MockQASchemeRepo {
	return &MockQASchemeRepo{data: make(map[int64]*QAScheme)}
}

func (r *MockQASchemeRepo) Create(_ context.Context, s *QAScheme) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[s.ID] = s
	return nil
}

func (r *MockQASchemeRepo) GetByID(_ context.Context, id int64) (*QAScheme, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.data[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (r *MockQASchemeRepo) Update(_ context.Context, s *QAScheme) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[s.ID] = s
	return nil
}

func (r *MockQASchemeRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}

func (r *MockQASchemeRepo) List(_ context.Context, tenantID int64) ([]*QAScheme, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*QAScheme
	for _, s := range r.data {
		if s.TenantID == tenantID {
			result = append(result, s)
		}
	}
	return result, nil
}

// --- MockQAResultRepo ---

type MockQAResultRepo struct {
	mu   sync.RWMutex
	data map[int64]*QAResult
}

func NewMockQAResultRepo() *MockQAResultRepo {
	return &MockQAResultRepo{data: make(map[int64]*QAResult)}
}

func (r *MockQAResultRepo) Create(_ context.Context, res *QAResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[res.ID] = res
	return nil
}

func (r *MockQAResultRepo) GetByID(_ context.Context, id int64) (*QAResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if res, ok := r.data[id]; ok {
		return res, nil
	}
	return nil, nil
}

func (r *MockQAResultRepo) Update(_ context.Context, res *QAResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[res.ID] = res
	return nil
}

func (r *MockQAResultRepo) List(_ context.Context, tenantID int64, offset, limit int) ([]*QAResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var all []*QAResult
	for _, res := range r.data {
		if res.TenantID == tenantID {
			all = append(all, res)
		}
	}
	if offset >= len(all) {
		return nil, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}

func (r *MockQAResultRepo) ListByCallID(_ context.Context, callID int64) ([]*QAResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*QAResult
	for _, res := range r.data {
		if res.CallID == callID {
			result = append(result, res)
		}
	}
	return result, nil
}

// --- MockASRHotwordsRepo ---

type MockASRHotwordsRepo struct {
	mu   sync.RWMutex
	data map[int64]*ASRHotwords
}

func NewMockASRHotwordsRepo() *MockASRHotwordsRepo {
	return &MockASRHotwordsRepo{data: make(map[int64]*ASRHotwords)}
}

func (r *MockASRHotwordsRepo) Create(_ context.Context, h *ASRHotwords) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[h.ID] = h
	return nil
}

func (r *MockASRHotwordsRepo) GetByID(_ context.Context, id int64) (*ASRHotwords, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if h, ok := r.data[id]; ok {
		return h, nil
	}
	return nil, nil
}

func (r *MockASRHotwordsRepo) Update(_ context.Context, h *ASRHotwords) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[h.ID] = h
	return nil
}

func (r *MockASRHotwordsRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}

func (r *MockASRHotwordsRepo) List(_ context.Context, tenantID int64) ([]*ASRHotwords, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*ASRHotwords
	for _, h := range r.data {
		if h.TenantID == tenantID {
			result = append(result, h)
		}
	}
	return result, nil
}

// --- MockPerformanceScorecardRepo ---

type MockPerformanceScorecardRepo struct {
	mu   sync.RWMutex
	data map[int64]*PerformanceScorecard
}

func NewMockPerformanceScorecardRepo() *MockPerformanceScorecardRepo {
	return &MockPerformanceScorecardRepo{data: make(map[int64]*PerformanceScorecard)}
}

func (r *MockPerformanceScorecardRepo) Create(_ context.Context, s *PerformanceScorecard) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[s.ID] = s
	return nil
}

func (r *MockPerformanceScorecardRepo) GetByID(_ context.Context, id int64) (*PerformanceScorecard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.data[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (r *MockPerformanceScorecardRepo) List(_ context.Context, tenantID int64, period string) ([]*PerformanceScorecard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*PerformanceScorecard
	for _, s := range r.data {
		if s.TenantID == tenantID && (period == "" || s.Period == period) {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *MockPerformanceScorecardRepo) GetByAgentAndPeriod(_ context.Context, agentID int64, period string) (*PerformanceScorecard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.data {
		if s.AgentID == agentID && s.Period == period {
			return s, nil
		}
	}
	return nil, nil
}
