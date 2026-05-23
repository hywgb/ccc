package integration

import "context"

type DNCRepository interface {
	Create(ctx context.Context, entry *DNCEntry) error
	GetByNumber(ctx context.Context, tenantID int64, number string) (*DNCEntry, error)
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*DNCEntry, int64, error)
	Delete(ctx context.Context, id int64) error
	CheckNumbers(ctx context.Context, tenantID int64, numbers []string) ([]string, error)
}

type CallTagAssignmentRepository interface {
	Create(ctx context.Context, a *CallTagAssignment) error
	ListByCallID(ctx context.Context, callID int64) ([]*CallTagAssignment, error)
	Delete(ctx context.Context, id int64) error
}
