package integration

import "time"

type DNCEntry struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  int64     `db:"tenant_id" json:"tenant_id"`
	Number    string    `db:"number" json:"number"`
	Reason    string    `db:"reason" json:"reason"`
	Source    string    `db:"source" json:"source"` // manual, import, api
	ExpiresAt *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CallTagAssignment struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  int64     `db:"tenant_id" json:"tenant_id"`
	CallID    int64     `db:"call_id" json:"call_id"`
	TagID     int64     `db:"tag_id" json:"tag_id"`
	TagName   string    `db:"tag_name" json:"tag_name"`
	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
