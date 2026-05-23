package telephony

import (
	"context"
	"sync"
)

type MockRoutingRuleRepo struct {
	mu    sync.RWMutex
	rules map[int64]*RoutingRule
}

func NewMockRoutingRuleRepo() *MockRoutingRuleRepo {
	return &MockRoutingRuleRepo{rules: make(map[int64]*RoutingRule)}
}

func (r *MockRoutingRuleRepo) Create(_ context.Context, rule *RoutingRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.ID] = rule
	return nil
}

func (r *MockRoutingRuleRepo) GetByID(_ context.Context, id int64) (*RoutingRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.rules[id], nil
}

func (r *MockRoutingRuleRepo) Update(_ context.Context, rule *RoutingRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.ID] = rule
	return nil
}

func (r *MockRoutingRuleRepo) ListActive(_ context.Context, tenantID int64) ([]*RoutingRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*RoutingRule
	for _, rule := range r.rules {
		if rule.TenantID == tenantID && rule.IsActive {
			result = append(result, rule)
		}
	}
	return result, nil
}

func (r *MockRoutingRuleRepo) List(_ context.Context, tenantID int64, offset, limit int) ([]*RoutingRule, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*RoutingRule
	for _, rule := range r.rules {
		if rule.TenantID == tenantID {
			filtered = append(filtered, rule)
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

func (r *MockRoutingRuleRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rules, id)
	return nil
}

type MockCLIPolicyRepo struct {
	mu       sync.RWMutex
	policies map[int64]*CLIPolicy
}

func NewMockCLIPolicyRepo() *MockCLIPolicyRepo {
	return &MockCLIPolicyRepo{policies: make(map[int64]*CLIPolicy)}
}

func (r *MockCLIPolicyRepo) Create(_ context.Context, p *CLIPolicy) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.policies[p.ID] = p
	return nil
}

func (r *MockCLIPolicyRepo) GetByID(_ context.Context, id int64) (*CLIPolicy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.policies[id], nil
}

func (r *MockCLIPolicyRepo) GetDefault(_ context.Context, tenantID int64) (*CLIPolicy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.policies {
		if p.TenantID == tenantID && p.IsDefault {
			return p, nil
		}
	}
	return nil, nil
}

func (r *MockCLIPolicyRepo) Update(_ context.Context, p *CLIPolicy) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.policies[p.ID] = p
	return nil
}

func (r *MockCLIPolicyRepo) List(_ context.Context, tenantID int64, offset, limit int) ([]*CLIPolicy, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*CLIPolicy
	for _, p := range r.policies {
		if p.TenantID == tenantID {
			filtered = append(filtered, p)
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

type MockPhoneNumberRepo struct {
	mu      sync.RWMutex
	numbers map[int64]*PhoneNumber
}

func NewMockPhoneNumberRepo() *MockPhoneNumberRepo {
	return &MockPhoneNumberRepo{numbers: make(map[int64]*PhoneNumber)}
}

func (r *MockPhoneNumberRepo) Create(_ context.Context, p *PhoneNumber) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.numbers[p.ID] = p
	return nil
}

func (r *MockPhoneNumberRepo) GetByID(_ context.Context, id int64) (*PhoneNumber, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.numbers[id], nil
}

func (r *MockPhoneNumberRepo) GetByNumber(_ context.Context, tenantID int64, number string) (*PhoneNumber, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.numbers {
		if p.TenantID == tenantID && p.Number == number {
			return p, nil
		}
	}
	return nil, nil
}

func (r *MockPhoneNumberRepo) Update(_ context.Context, p *PhoneNumber) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.numbers[p.ID] = p
	return nil
}

func (r *MockPhoneNumberRepo) List(_ context.Context, tenantID int64, offset, limit int) ([]*PhoneNumber, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*PhoneNumber
	for _, p := range r.numbers {
		if p.TenantID == tenantID {
			filtered = append(filtered, p)
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
