package mysql

import (
	"context"

	"github.com/divord97/ccc/internal/domain/integration"
	"github.com/jmoiron/sqlx"
)

type DNCRepo struct{ db *sqlx.DB }

func NewDNCRepo(db *sqlx.DB) *DNCRepo { return &DNCRepo{db: db} }

func (r *DNCRepo) Create(ctx context.Context, entry *integration.DNCEntry) error {
	_, err := r.db.NamedExecContext(ctx,
		`INSERT INTO dnc_list (id, tenant_id, number, reason, source, expires_at, created_at)
		 VALUES (:id, :tenant_id, :number, :reason, :source, :expires_at, :created_at)`, entry)
	return err
}

func (r *DNCRepo) GetByNumber(ctx context.Context, tenantID int64, number string) (*integration.DNCEntry, error) {
	var entry integration.DNCEntry
	if err := r.db.GetContext(ctx, &entry, `SELECT * FROM dnc_list WHERE tenant_id = ? AND number = ? LIMIT 1`, tenantID, number); err != nil {
		return nil, nil
	}
	return &entry, nil
}

func (r *DNCRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]*integration.DNCEntry, int64, error) {
	var total int64
	_ = r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM dnc_list WHERE tenant_id = ?`, tenantID)
	var entries []*integration.DNCEntry
	err := r.db.SelectContext(ctx, &entries, `SELECT * FROM dnc_list WHERE tenant_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, tenantID, limit, offset)
	return entries, total, err
}

func (r *DNCRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM dnc_list WHERE id = ?`, id)
	return err
}

func (r *DNCRepo) CheckNumbers(ctx context.Context, tenantID int64, numbers []string) ([]string, error) {
	if len(numbers) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(`SELECT number FROM dnc_list WHERE tenant_id = ? AND number IN (?) AND (expires_at IS NULL OR expires_at > NOW())`, tenantID, numbers)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var blocked []string
	err = r.db.SelectContext(ctx, &blocked, query, args...)
	return blocked, err
}
