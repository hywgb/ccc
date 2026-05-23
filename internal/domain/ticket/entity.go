package ticket

import (
	"time"
)

type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "open"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusPending    TicketStatus = "pending"
	TicketStatusResolved   TicketStatus = "resolved"
	TicketStatusClosed     TicketStatus = "closed"
)

type TicketCategory struct {
	ID       int64  `db:"id" json:"id"`
	TenantID int64  `db:"tenant_id" json:"tenant_id"`
	Name     string `db:"name" json:"name"`
	ParentID *int64 `db:"parent_id" json:"parent_id,omitempty"`
}

type TicketTemplate struct {
	ID           int64  `db:"id" json:"id"`
	TenantID     int64  `db:"tenant_id" json:"tenant_id"`
	Name         string `db:"name" json:"name"`
	CategoryID   *int64 `db:"category_id" json:"category_id,omitempty"`
	Fields       string `db:"fields" json:"fields"`       // JSON array of field definitions
	FlowGraph    string `db:"flow_graph" json:"flow_graph"` // JSON DAG for ticket flow
	OnlineStatus string `db:"online_status" json:"online_status"` // draft, published, offline
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Ticket struct {
	ID           int64        `db:"id" json:"id"`
	TenantID     int64        `db:"tenant_id" json:"tenant_id"`
	TemplateID   *int64       `db:"template_id" json:"template_id,omitempty"`
	CategoryID   *int64       `db:"category_id" json:"category_id,omitempty"`
	Title        string       `db:"title" json:"title"`
	Description  string       `db:"description" json:"description"`
	Status       TicketStatus `db:"status" json:"status"`
	Priority     string       `db:"priority" json:"priority"` // low, medium, high, urgent
	CustomerID   *int64       `db:"customer_id" json:"customer_id,omitempty"`
	AssigneeID   *int64       `db:"assignee_id" json:"assignee_id,omitempty"`
	CallID       *int64       `db:"call_id" json:"call_id,omitempty"`
	CustomData   string       `db:"custom_data" json:"custom_data"`
	CreatedAt    time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updated_at"`
	ResolvedAt   *time.Time   `db:"resolved_at" json:"resolved_at,omitempty"`
}

type TicketComment struct {
	ID        int64     `db:"id" json:"id"`
	TicketID  int64     `db:"ticket_id" json:"ticket_id"`
	AuthorID  int64     `db:"author_id" json:"author_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
