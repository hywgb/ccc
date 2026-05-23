package crm

import "context"

type CustomerRepository interface {
	Create(ctx context.Context, c *Customer) error
	GetByID(ctx context.Context, id int64) (*Customer, error)
	Update(ctx context.Context, c *Customer) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*Customer, error)
	FindByPhone(ctx context.Context, tenantID int64, phone string) (*Customer, error)
}

type CustomerPhoneRepository interface {
	Create(ctx context.Context, p *CustomerPhone) error
	ListByCustomer(ctx context.Context, customerID int64) ([]*CustomerPhone, error)
	DeleteByCustomer(ctx context.Context, customerID int64) error
}

type CustomerInteractionRepository interface {
	Create(ctx context.Context, i *CustomerInteraction) error
	ListByCustomer(ctx context.Context, customerID int64, offset, limit int) ([]*CustomerInteraction, error)
}

type CustomFieldDefinitionRepository interface {
	Create(ctx context.Context, d *CustomFieldDefinition) error
	Update(ctx context.Context, d *CustomFieldDefinition) error
	Delete(ctx context.Context, id int64) error
	ListByEntity(ctx context.Context, tenantID int64, entityType string) ([]*CustomFieldDefinition, error)
}
