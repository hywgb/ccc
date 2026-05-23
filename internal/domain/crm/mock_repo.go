package crm

import (
	"context"
	"sync"
)

// MockCustomerRepo is an in-memory CustomerRepository for testing.
type MockCustomerRepo struct {
	mu   sync.RWMutex
	data map[int64]*Customer
}

func NewMockCustomerRepo() *MockCustomerRepo {
	return &MockCustomerRepo{data: make(map[int64]*Customer)}
}

func (r *MockCustomerRepo) Create(_ context.Context, c *Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[c.ID] = c
	return nil
}

func (r *MockCustomerRepo) GetByID(_ context.Context, id int64) (*Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if c, ok := r.data[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (r *MockCustomerRepo) Update(_ context.Context, c *Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[c.ID] = c
	return nil
}

func (r *MockCustomerRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}

func (r *MockCustomerRepo) List(_ context.Context, tenantID int64, offset, limit int) ([]*Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Customer
	for _, c := range r.data {
		if c.TenantID == tenantID {
			result = append(result, c)
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

func (r *MockCustomerRepo) FindByPhone(_ context.Context, tenantID int64, phone string) (*Customer, error) {
	return nil, nil // phone lookup is handled via MockCustomerPhoneRepo
}

// MockCustomerPhoneRepo is an in-memory CustomerPhoneRepository for testing.
type MockCustomerPhoneRepo struct {
	mu   sync.RWMutex
	data map[int64]*CustomerPhone
	seq  int64
}

func NewMockCustomerPhoneRepo() *MockCustomerPhoneRepo {
	return &MockCustomerPhoneRepo{data: make(map[int64]*CustomerPhone)}
}

func (r *MockCustomerPhoneRepo) Create(_ context.Context, p *CustomerPhone) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	if p.ID == 0 {
		p.ID = r.seq
	}
	r.data[p.ID] = p
	return nil
}

func (r *MockCustomerPhoneRepo) ListByCustomer(_ context.Context, customerID int64) ([]*CustomerPhone, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*CustomerPhone
	for _, p := range r.data {
		if p.CustomerID == customerID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *MockCustomerPhoneRepo) DeleteByCustomer(_ context.Context, customerID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, p := range r.data {
		if p.CustomerID == customerID {
			delete(r.data, id)
		}
	}
	return nil
}

// FindCustomerByPhone looks up a customer ID from any of their phone numbers.
func (r *MockCustomerPhoneRepo) FindCustomerByPhone(phone string) int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.data {
		if p.Number == phone {
			return p.CustomerID
		}
	}
	return 0
}

// MockInteractionRepo is an in-memory CustomerInteractionRepository for testing.
type MockInteractionRepo struct {
	mu   sync.RWMutex
	data map[int64]*CustomerInteraction
	seq  int64
}

func NewMockInteractionRepo() *MockInteractionRepo {
	return &MockInteractionRepo{data: make(map[int64]*CustomerInteraction)}
}

func (r *MockInteractionRepo) Create(_ context.Context, i *CustomerInteraction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	if i.ID == 0 {
		i.ID = r.seq
	}
	r.data[i.ID] = i
	return nil
}

func (r *MockInteractionRepo) ListByCustomer(_ context.Context, customerID int64, offset, limit int) ([]*CustomerInteraction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*CustomerInteraction
	for _, i := range r.data {
		if i.CustomerID == customerID {
			result = append(result, i)
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

// MockCustomFieldRepo is an in-memory CustomFieldDefinitionRepository for testing.
type MockCustomFieldRepo struct {
	mu   sync.RWMutex
	data map[int64]*CustomFieldDefinition
	seq  int64
}

func NewMockCustomFieldRepo() *MockCustomFieldRepo {
	return &MockCustomFieldRepo{data: make(map[int64]*CustomFieldDefinition)}
}

func (r *MockCustomFieldRepo) Create(_ context.Context, d *CustomFieldDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	if d.ID == 0 {
		d.ID = r.seq
	}
	r.data[d.ID] = d
	return nil
}

func (r *MockCustomFieldRepo) Update(_ context.Context, d *CustomFieldDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[d.ID] = d
	return nil
}

func (r *MockCustomFieldRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}

func (r *MockCustomFieldRepo) ListByEntity(_ context.Context, tenantID int64, entityType string) ([]*CustomFieldDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*CustomFieldDefinition
	for _, d := range r.data {
		if d.TenantID == tenantID && d.EntityType == entityType {
			result = append(result, d)
		}
	}
	return result, nil
}
