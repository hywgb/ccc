package crm

import (
	"time"
)

type Customer struct {
	ID         int64     `db:"id" json:"id"`
	TenantID   int64     `db:"tenant_id" json:"tenant_id"`
	Name       string    `db:"name" json:"name"`
	Email      string    `db:"email" json:"email"`
	Company    string    `db:"company" json:"company"`
	Level      string    `db:"level" json:"level"` // normal, vip, svip
	CustomData string    `db:"custom_data" json:"custom_data"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type CustomerPhone struct {
	ID         int64  `db:"id" json:"id"`
	CustomerID int64  `db:"customer_id" json:"customer_id"`
	PhoneType  string `db:"phone_type" json:"phone_type"` // mobile, landline, backup
	Number     string `db:"number" json:"number"`
	IsPrimary  bool   `db:"is_primary" json:"is_primary"`
}

type CustomerInteraction struct {
	ID         int64     `db:"id" json:"id"`
	CustomerID int64     `db:"customer_id" json:"customer_id"`
	TenantID   int64     `db:"tenant_id" json:"tenant_id"`
	Channel    string    `db:"channel" json:"channel"` // call, ticket, im
	Direction  string    `db:"direction" json:"direction"`
	Summary    string    `db:"summary" json:"summary"`
	CallID     *int64    `db:"call_id" json:"call_id,omitempty"`
	TicketID   *int64    `db:"ticket_id" json:"ticket_id,omitempty"`
	AgentName  string    `db:"agent_name" json:"agent_name"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type CustomFieldDefinition struct {
	ID         int64  `db:"id" json:"id"`
	TenantID   int64  `db:"tenant_id" json:"tenant_id"`
	EntityType string `db:"entity_type" json:"entity_type"` // customer, ticket
	FieldName  string `db:"field_name" json:"field_name"`
	FieldType  string `db:"field_type" json:"field_type"` // text, number, select, date
	Options    string `db:"options" json:"options"`
	IsRequired bool   `db:"is_required" json:"is_required"`
	SortOrder  int    `db:"sort_order" json:"sort_order"`
}
