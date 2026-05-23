package mysql

import (
	"context"

	"github.com/divord97/ccc/internal/domain/telephony"
	"github.com/jmoiron/sqlx"
)

type CLIPolicyRepo struct{ db *sqlx.DB }

func NewCLIPolicyRepo(db *sqlx.DB) *CLIPolicyRepo { return &CLIPolicyRepo{db: db} }

func (r *CLIPolicyRepo) Create(ctx context.Context, p *telephony.CLIPolicy) error {
	_, err := r.db.NamedExecContext(ctx,
		`INSERT INTO cli_policies (id, tenant_id, name, strategy, fixed_number_id, number_pool_ids, is_default, created_at, updated_at)
		 VALUES (:id, :tenant_id, :name, :strategy, :fixed_number_id, :number_pool_ids, :is_default, :created_at, :updated_at)`, p)
	return err
}

func (r *CLIPolicyRepo) GetByID(ctx context.Context, id int64) (*telephony.CLIPolicy, error) {
	var p telephony.CLIPolicy
	if err := r.db.GetContext(ctx, &p, `SELECT * FROM cli_policies WHERE id = ?`, id); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *CLIPolicyRepo) GetDefault(ctx context.Context, tenantID int64) (*telephony.CLIPolicy, error) {
	var p telephony.CLIPolicy
	if err := r.db.GetContext(ctx, &p, `SELECT * FROM cli_policies WHERE tenant_id = ? AND is_default = 1 LIMIT 1`, tenantID); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *CLIPolicyRepo) Update(ctx context.Context, p *telephony.CLIPolicy) error {
	_, err := r.db.NamedExecContext(ctx,
		`UPDATE cli_policies SET name=:name, strategy=:strategy, fixed_number_id=:fixed_number_id,
		 number_pool_ids=:number_pool_ids, is_default=:is_default, updated_at=:updated_at WHERE id=:id`, p)
	return err
}

func (r *CLIPolicyRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]*telephony.CLIPolicy, int64, error) {
	var total int64
	_ = r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM cli_policies WHERE tenant_id = ?`, tenantID)
	var policies []*telephony.CLIPolicy
	err := r.db.SelectContext(ctx, &policies, `SELECT * FROM cli_policies WHERE tenant_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, tenantID, limit, offset)
	return policies, total, err
}
