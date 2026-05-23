package mysql

import (
	"context"

	"github.com/divord97/ccc/internal/domain/telephony"
	"github.com/jmoiron/sqlx"
)

type RoutingRuleRepo struct{ db *sqlx.DB }

func NewRoutingRuleRepo(db *sqlx.DB) *RoutingRuleRepo { return &RoutingRuleRepo{db: db} }

func (r *RoutingRuleRepo) Create(ctx context.Context, rule *telephony.RoutingRule) error {
	_, err := r.db.NamedExecContext(ctx,
		`INSERT INTO routing_rules (id, tenant_id, name, match_type, match_value, sip_trunk_id, priority, is_active, created_at, updated_at)
		 VALUES (:id, :tenant_id, :name, :match_type, :match_value, :sip_trunk_id, :priority, :is_active, :created_at, :updated_at)`, rule)
	return err
}

func (r *RoutingRuleRepo) GetByID(ctx context.Context, id int64) (*telephony.RoutingRule, error) {
	var rule telephony.RoutingRule
	if err := r.db.GetContext(ctx, &rule, `SELECT * FROM routing_rules WHERE id = ?`, id); err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *RoutingRuleRepo) Update(ctx context.Context, rule *telephony.RoutingRule) error {
	_, err := r.db.NamedExecContext(ctx,
		`UPDATE routing_rules SET name=:name, match_type=:match_type, match_value=:match_value,
		 sip_trunk_id=:sip_trunk_id, priority=:priority, is_active=:is_active, updated_at=:updated_at WHERE id=:id`, rule)
	return err
}

func (r *RoutingRuleRepo) ListActive(ctx context.Context, tenantID int64) ([]*telephony.RoutingRule, error) {
	var rules []*telephony.RoutingRule
	err := r.db.SelectContext(ctx, &rules, `SELECT * FROM routing_rules WHERE tenant_id = ? AND is_active = 1 ORDER BY priority`, tenantID)
	return rules, err
}

func (r *RoutingRuleRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]*telephony.RoutingRule, int64, error) {
	var total int64
	_ = r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM routing_rules WHERE tenant_id = ?`, tenantID)
	var rules []*telephony.RoutingRule
	err := r.db.SelectContext(ctx, &rules, `SELECT * FROM routing_rules WHERE tenant_id = ? ORDER BY priority LIMIT ? OFFSET ?`, tenantID, limit, offset)
	return rules, total, err
}

func (r *RoutingRuleRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM routing_rules WHERE id = ?`, id)
	return err
}
