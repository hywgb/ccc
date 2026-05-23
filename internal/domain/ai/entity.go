package ai

import (
	"time"
)

type KnowledgeCategory struct {
	ID       int64  `db:"id" json:"id"`
	TenantID int64  `db:"tenant_id" json:"tenant_id"`
	Name     string `db:"name" json:"name"`
	ParentID *int64 `db:"parent_id" json:"parent_id,omitempty"`
}

type KnowledgeArticle struct {
	ID         int64     `db:"id" json:"id"`
	TenantID   int64     `db:"tenant_id" json:"tenant_id"`
	CategoryID *int64    `db:"category_id" json:"category_id,omitempty"`
	Title      string    `db:"title" json:"title"`
	Content    string    `db:"content" json:"content"`
	Tags       string    `db:"tags" json:"tags"` // comma-separated
	Status     string    `db:"status" json:"status"` // draft, published
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type AgentScript struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  int64     `db:"tenant_id" json:"tenant_id"`
	Name      string    `db:"name" json:"name"`
	Content   string    `db:"content" json:"content"` // JSON script flow
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type SessionInfoTemplate struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  int64     `db:"tenant_id" json:"tenant_id"`
	Name      string    `db:"name" json:"name"`
	Fields    string    `db:"fields" json:"fields"` // JSON array of field definitions
	IsDefault bool      `db:"is_default" json:"is_default"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
