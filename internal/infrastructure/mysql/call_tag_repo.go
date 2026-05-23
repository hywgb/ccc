package mysql

import (
	"context"

	"github.com/divord97/ccc/internal/domain/integration"
	"github.com/jmoiron/sqlx"
)

type CallTagAssignmentRepo struct{ db *sqlx.DB }

func NewCallTagAssignmentRepo(db *sqlx.DB) *CallTagAssignmentRepo {
	return &CallTagAssignmentRepo{db: db}
}

func (r *CallTagAssignmentRepo) Create(ctx context.Context, a *integration.CallTagAssignment) error {
	_, err := r.db.NamedExecContext(ctx,
		`INSERT INTO call_tag_assignments (id, tenant_id, call_id, tag_id, tag_name, created_by, created_at)
		 VALUES (:id, :tenant_id, :call_id, :tag_id, :tag_name, :created_by, :created_at)`, a)
	return err
}

func (r *CallTagAssignmentRepo) ListByCallID(ctx context.Context, callID int64) ([]*integration.CallTagAssignment, error) {
	var tags []*integration.CallTagAssignment
	err := r.db.SelectContext(ctx, &tags, `SELECT * FROM call_tag_assignments WHERE call_id = ? ORDER BY created_at`, callID)
	return tags, err
}

func (r *CallTagAssignmentRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM call_tag_assignments WHERE id = ?`, id)
	return err
}
