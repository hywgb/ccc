package integration

import (
	"context"
	"sync"
)

type MockDNCRepo struct {
	mu      sync.RWMutex
	entries map[int64]*DNCEntry
}

func NewMockDNCRepo() *MockDNCRepo {
	return &MockDNCRepo{entries: make(map[int64]*DNCEntry)}
}

func (r *MockDNCRepo) Create(_ context.Context, entry *DNCEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[entry.ID] = entry
	return nil
}

func (r *MockDNCRepo) GetByNumber(_ context.Context, tenantID int64, number string) (*DNCEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, e := range r.entries {
		if e.TenantID == tenantID && e.Number == number {
			return e, nil
		}
	}
	return nil, nil
}

func (r *MockDNCRepo) List(_ context.Context, tenantID int64, offset, limit int) ([]*DNCEntry, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*DNCEntry
	for _, e := range r.entries {
		if e.TenantID == tenantID {
			filtered = append(filtered, e)
		}
	}
	total := int64(len(filtered))
	if offset >= len(filtered) {
		return nil, total, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func (r *MockDNCRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, id)
	return nil
}

func (r *MockDNCRepo) CheckNumbers(_ context.Context, tenantID int64, numbers []string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	numSet := make(map[string]bool)
	for _, e := range r.entries {
		if e.TenantID == tenantID {
			numSet[e.Number] = true
		}
	}
	var blocked []string
	for _, n := range numbers {
		if numSet[n] {
			blocked = append(blocked, n)
		}
	}
	return blocked, nil
}

type MockCallTagRepo struct {
	mu   sync.RWMutex
	tags map[int64]*CallTagAssignment
}

func NewMockCallTagRepo() *MockCallTagRepo {
	return &MockCallTagRepo{tags: make(map[int64]*CallTagAssignment)}
}

func (r *MockCallTagRepo) Create(_ context.Context, a *CallTagAssignment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tags[a.ID] = a
	return nil
}

func (r *MockCallTagRepo) ListByCallID(_ context.Context, callID int64) ([]*CallTagAssignment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*CallTagAssignment
	for _, a := range r.tags {
		if a.CallID == callID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *MockCallTagRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tags, id)
	return nil
}
