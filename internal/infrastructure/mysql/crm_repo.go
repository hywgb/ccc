package mysql

import (
	"context"
	"database/sql"

	"github.com/divord97/ccc/internal/domain/crm"
	"github.com/jmoiron/sqlx"
)

type CustomerRepo struct{ db *sqlx.DB }

func NewCustomerRepo(db *sqlx.DB) *CustomerRepo { return &CustomerRepo{db: db} }

func (r *CustomerRepo) Create(ctx context.Context, c *crm.Customer) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO customers (id, tenant_id, name, email, company, level, custom_data, created_at, updated_at)
		 VALUES (?,?,?,?,?,?,?,?,?)`,
		c.ID, c.TenantID, c.Name, c.Email, c.Company, c.Level, c.CustomData, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *CustomerRepo) GetByID(ctx context.Context, id int64) (*crm.Customer, error) {
	var c crm.Customer
	if err := r.db.GetContext(ctx, &c, `SELECT * FROM customers WHERE id=?`, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepo) Update(ctx context.Context, c *crm.Customer) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE customers SET name=?, email=?, company=?, level=?, custom_data=?, updated_at=? WHERE id=?`,
		c.Name, c.Email, c.Company, c.Level, c.CustomData, c.UpdatedAt, c.ID)
	return err
}

func (r *CustomerRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM customers WHERE id=?`, id)
	return err
}

func (r *CustomerRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]*crm.Customer, error) {
	var result []*crm.Customer
	err := r.db.SelectContext(ctx, &result,
		`SELECT * FROM customers WHERE tenant_id=? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		tenantID, limit, offset)
	return result, err
}

func (r *CustomerRepo) FindByPhone(ctx context.Context, tenantID int64, phone string) (*crm.Customer, error) {
	var c crm.Customer
	err := r.db.GetContext(ctx, &c,
		`SELECT c.* FROM customers c
		 JOIN customer_phones cp ON c.id = cp.customer_id
		 WHERE c.tenant_id=? AND cp.number=? LIMIT 1`, tenantID, phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

type CustomerPhoneRepo struct{ db *sqlx.DB }

func NewCustomerPhoneRepo(db *sqlx.DB) *CustomerPhoneRepo { return &CustomerPhoneRepo{db: db} }

func (r *CustomerPhoneRepo) Create(ctx context.Context, p *crm.CustomerPhone) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO customer_phones (id, customer_id, phone_type, number, is_primary) VALUES (?,?,?,?,?)`,
		p.ID, p.CustomerID, p.PhoneType, p.Number, p.IsPrimary)
	return err
}

func (r *CustomerPhoneRepo) ListByCustomer(ctx context.Context, customerID int64) ([]*crm.CustomerPhone, error) {
	var result []*crm.CustomerPhone
	err := r.db.SelectContext(ctx, &result, `SELECT * FROM customer_phones WHERE customer_id=?`, customerID)
	return result, err
}

func (r *CustomerPhoneRepo) DeleteByCustomer(ctx context.Context, customerID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM customer_phones WHERE customer_id=?`, customerID)
	return err
}

type InteractionRepo struct{ db *sqlx.DB }

func NewInteractionRepo(db *sqlx.DB) *InteractionRepo { return &InteractionRepo{db: db} }

func (r *InteractionRepo) Create(ctx context.Context, i *crm.CustomerInteraction) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO customer_interactions (id, customer_id, tenant_id, channel, direction, summary, call_id, ticket_id, agent_name, created_at)
		 VALUES (?,?,?,?,?,?,?,?,?,?)`,
		i.ID, i.CustomerID, i.TenantID, i.Channel, i.Direction, i.Summary, i.CallID, i.TicketID, i.AgentName, i.CreatedAt)
	return err
}

func (r *InteractionRepo) ListByCustomer(ctx context.Context, customerID int64, offset, limit int) ([]*crm.CustomerInteraction, error) {
	var result []*crm.CustomerInteraction
	err := r.db.SelectContext(ctx, &result,
		`SELECT * FROM customer_interactions WHERE customer_id=? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		customerID, limit, offset)
	return result, err
}

type CustomFieldRepo struct{ db *sqlx.DB }

func NewCustomFieldRepo(db *sqlx.DB) *CustomFieldRepo { return &CustomFieldRepo{db: db} }

func (r *CustomFieldRepo) Create(ctx context.Context, d *crm.CustomFieldDefinition) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO custom_field_definitions (id, tenant_id, entity_type, field_name, field_type, options, is_required, sort_order)
		 VALUES (?,?,?,?,?,?,?,?)`,
		d.ID, d.TenantID, d.EntityType, d.FieldName, d.FieldType, d.Options, d.IsRequired, d.SortOrder)
	return err
}

func (r *CustomFieldRepo) Update(ctx context.Context, d *crm.CustomFieldDefinition) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE custom_field_definitions SET field_name=?, field_type=?, options=?, is_required=?, sort_order=? WHERE id=?`,
		d.FieldName, d.FieldType, d.Options, d.IsRequired, d.SortOrder, d.ID)
	return err
}

func (r *CustomFieldRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM custom_field_definitions WHERE id=?`, id)
	return err
}

func (r *CustomFieldRepo) ListByEntity(ctx context.Context, tenantID int64, entityType string) ([]*crm.CustomFieldDefinition, error) {
	var result []*crm.CustomFieldDefinition
	err := r.db.SelectContext(ctx, &result,
		`SELECT * FROM custom_field_definitions WHERE tenant_id=? AND entity_type=? ORDER BY sort_order`,
		tenantID, entityType)
	return result, err
}
