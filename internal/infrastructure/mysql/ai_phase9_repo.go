package mysql

import (
	"context"
	"database/sql"

	"github.com/divord97/ccc/internal/domain/ai"
	"github.com/jmoiron/sqlx"
)

// --- DigitalEmployeeRepo ---

type DigitalEmployeeRepo struct{ db *sqlx.DB }

func NewDigitalEmployeeRepo(db *sqlx.DB) *DigitalEmployeeRepo {
	return &DigitalEmployeeRepo{db: db}
}

func (r *DigitalEmployeeRepo) Create(ctx context.Context, de *ai.DigitalEmployee) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO digital_employees (id, tenant_id, name, description, avatar_url, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		de.ID, de.TenantID, de.Name, de.Description, de.AvatarURL, de.IsActive, de.CreatedAt, de.UpdatedAt)
	return err
}

func (r *DigitalEmployeeRepo) GetByID(ctx context.Context, id int64) (*ai.DigitalEmployee, error) {
	var de ai.DigitalEmployee
	if err := r.db.GetContext(ctx, &de, "SELECT * FROM digital_employees WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &de, nil
}

func (r *DigitalEmployeeRepo) Update(ctx context.Context, de *ai.DigitalEmployee) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE digital_employees SET name=?, description=?, avatar_url=?, is_active=?, updated_at=? WHERE id=?",
		de.Name, de.Description, de.AvatarURL, de.IsActive, de.UpdatedAt, de.ID)
	return err
}

func (r *DigitalEmployeeRepo) List(ctx context.Context, tenantID int64) ([]*ai.DigitalEmployee, error) {
	var result []*ai.DigitalEmployee
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM digital_employees WHERE tenant_id = ? ORDER BY created_at DESC", tenantID)
	return result, err
}

// --- DigitalEmployeeSceneRepo ---

type DigitalEmployeeSceneRepo struct{ db *sqlx.DB }

func NewDigitalEmployeeSceneRepo(db *sqlx.DB) *DigitalEmployeeSceneRepo {
	return &DigitalEmployeeSceneRepo{db: db}
}

func (r *DigitalEmployeeSceneRepo) Create(ctx context.Context, s *ai.DigitalEmployeeScene) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO digital_employee_scenes (id, digital_employee_id, tenant_id, name, intents, faqs, transfer_skill_group, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		s.ID, s.DigitalEmployeeID, s.TenantID, s.Name, s.Intents, s.FAQs, s.TransferSkillGroup, s.Status, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *DigitalEmployeeSceneRepo) GetByID(ctx context.Context, id int64) (*ai.DigitalEmployeeScene, error) {
	var s ai.DigitalEmployeeScene
	if err := r.db.GetContext(ctx, &s, "SELECT * FROM digital_employee_scenes WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *DigitalEmployeeSceneRepo) Update(ctx context.Context, s *ai.DigitalEmployeeScene) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE digital_employee_scenes SET name=?, intents=?, faqs=?, transfer_skill_group=?, status=?, updated_at=? WHERE id=?",
		s.Name, s.Intents, s.FAQs, s.TransferSkillGroup, s.Status, s.UpdatedAt, s.ID)
	return err
}

func (r *DigitalEmployeeSceneRepo) List(ctx context.Context, digitalEmployeeID int64) ([]*ai.DigitalEmployeeScene, error) {
	var result []*ai.DigitalEmployeeScene
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM digital_employee_scenes WHERE digital_employee_id = ? ORDER BY created_at DESC", digitalEmployeeID)
	return result, err
}

// --- QARuleRepo ---

type QARuleRepo struct{ db *sqlx.DB }

func NewQARuleRepo(db *sqlx.DB) *QARuleRepo {
	return &QARuleRepo{db: db}
}

func (r *QARuleRepo) Create(ctx context.Context, rule *ai.QARule) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO qa_rules (id, tenant_id, name, type, config, severity, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		rule.ID, rule.TenantID, rule.Name, rule.Type, rule.Config, rule.Severity, rule.IsActive, rule.CreatedAt, rule.UpdatedAt)
	return err
}

func (r *QARuleRepo) GetByID(ctx context.Context, id int64) (*ai.QARule, error) {
	var rule ai.QARule
	if err := r.db.GetContext(ctx, &rule, "SELECT * FROM qa_rules WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

func (r *QARuleRepo) Update(ctx context.Context, rule *ai.QARule) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE qa_rules SET name=?, type=?, config=?, severity=?, is_active=?, updated_at=? WHERE id=?",
		rule.Name, rule.Type, rule.Config, rule.Severity, rule.IsActive, rule.UpdatedAt, rule.ID)
	return err
}

func (r *QARuleRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM qa_rules WHERE id = ?", id)
	return err
}

func (r *QARuleRepo) List(ctx context.Context, tenantID int64) ([]*ai.QARule, error) {
	var result []*ai.QARule
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM qa_rules WHERE tenant_id = ? ORDER BY created_at DESC", tenantID)
	return result, err
}

func (r *QARuleRepo) ListByIDs(ctx context.Context, ids []int64) ([]*ai.QARule, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In("SELECT * FROM qa_rules WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var result []*ai.QARule
	err = r.db.SelectContext(ctx, &result, query, args...)
	return result, err
}

// --- QASchemeRepo ---

type QASchemeRepo struct{ db *sqlx.DB }

func NewQASchemeRepo(db *sqlx.DB) *QASchemeRepo {
	return &QASchemeRepo{db: db}
}

func (r *QASchemeRepo) Create(ctx context.Context, s *ai.QAScheme) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO qa_schemes (id, tenant_id, name, rule_ids, is_default, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		s.ID, s.TenantID, s.Name, s.RuleIDs, s.IsDefault, s.IsActive, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *QASchemeRepo) GetByID(ctx context.Context, id int64) (*ai.QAScheme, error) {
	var s ai.QAScheme
	if err := r.db.GetContext(ctx, &s, "SELECT * FROM qa_schemes WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *QASchemeRepo) Update(ctx context.Context, s *ai.QAScheme) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE qa_schemes SET name=?, rule_ids=?, is_default=?, is_active=?, updated_at=? WHERE id=?",
		s.Name, s.RuleIDs, s.IsDefault, s.IsActive, s.UpdatedAt, s.ID)
	return err
}

func (r *QASchemeRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM qa_schemes WHERE id = ?", id)
	return err
}

func (r *QASchemeRepo) List(ctx context.Context, tenantID int64) ([]*ai.QAScheme, error) {
	var result []*ai.QAScheme
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM qa_schemes WHERE tenant_id = ? ORDER BY created_at DESC", tenantID)
	return result, err
}

// --- QAResultRepo ---

type QAResultRepo struct{ db *sqlx.DB }

func NewQAResultRepo(db *sqlx.DB) *QAResultRepo {
	return &QAResultRepo{db: db}
}

func (r *QAResultRepo) Create(ctx context.Context, res *ai.QAResult) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO qa_results (id, tenant_id, call_id, scheme_id, score, details, status, appeal_note, review_note, reviewer_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		res.ID, res.TenantID, res.CallID, res.SchemeID, res.Score, res.Details, res.Status, res.AppealNote, res.ReviewNote, res.ReviewerID, res.CreatedAt, res.UpdatedAt)
	return err
}

func (r *QAResultRepo) GetByID(ctx context.Context, id int64) (*ai.QAResult, error) {
	var res ai.QAResult
	if err := r.db.GetContext(ctx, &res, "SELECT * FROM qa_results WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *QAResultRepo) Update(ctx context.Context, res *ai.QAResult) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE qa_results SET score=?, details=?, status=?, appeal_note=?, review_note=?, reviewer_id=?, updated_at=? WHERE id=?",
		res.Score, res.Details, res.Status, res.AppealNote, res.ReviewNote, res.ReviewerID, res.UpdatedAt, res.ID)
	return err
}

func (r *QAResultRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]*ai.QAResult, error) {
	var result []*ai.QAResult
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM qa_results WHERE tenant_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?", tenantID, limit, offset)
	return result, err
}

func (r *QAResultRepo) ListByCallID(ctx context.Context, callID int64) ([]*ai.QAResult, error) {
	var result []*ai.QAResult
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM qa_results WHERE call_id = ? ORDER BY created_at DESC", callID)
	return result, err
}

// --- ASRHotwordsRepo ---

type ASRHotwordsRepo struct{ db *sqlx.DB }

func NewASRHotwordsRepo(db *sqlx.DB) *ASRHotwordsRepo {
	return &ASRHotwordsRepo{db: db}
}

func (r *ASRHotwordsRepo) Create(ctx context.Context, h *ai.ASRHotwords) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO asr_hotwords (id, tenant_id, name, words, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		h.ID, h.TenantID, h.Name, h.Words, h.IsActive, h.CreatedAt, h.UpdatedAt)
	return err
}

func (r *ASRHotwordsRepo) GetByID(ctx context.Context, id int64) (*ai.ASRHotwords, error) {
	var h ai.ASRHotwords
	if err := r.db.GetContext(ctx, &h, "SELECT * FROM asr_hotwords WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &h, nil
}

func (r *ASRHotwordsRepo) Update(ctx context.Context, h *ai.ASRHotwords) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE asr_hotwords SET name=?, words=?, is_active=?, updated_at=? WHERE id=?",
		h.Name, h.Words, h.IsActive, h.UpdatedAt, h.ID)
	return err
}

func (r *ASRHotwordsRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM asr_hotwords WHERE id = ?", id)
	return err
}

func (r *ASRHotwordsRepo) List(ctx context.Context, tenantID int64) ([]*ai.ASRHotwords, error) {
	var result []*ai.ASRHotwords
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM asr_hotwords WHERE tenant_id = ? ORDER BY created_at DESC", tenantID)
	return result, err
}

// --- PerformanceScorecardRepo ---

type PerformanceScorecardRepo struct{ db *sqlx.DB }

func NewPerformanceScorecardRepo(db *sqlx.DB) *PerformanceScorecardRepo {
	return &PerformanceScorecardRepo{db: db}
}

func (r *PerformanceScorecardRepo) Create(ctx context.Context, s *ai.PerformanceScorecard) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO performance_scorecards (id, tenant_id, agent_id, period, total_calls, avg_handle_time, avg_qa_score, csat_score, first_call_resolution, adherence, overall_score, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		s.ID, s.TenantID, s.AgentID, s.Period, s.TotalCalls, s.AvgHandleTime, s.AvgQAScore, s.CSATScore, s.FirstCallResolv, s.Adherence, s.OverallScore, s.CreatedAt)
	return err
}

func (r *PerformanceScorecardRepo) GetByID(ctx context.Context, id int64) (*ai.PerformanceScorecard, error) {
	var s ai.PerformanceScorecard
	if err := r.db.GetContext(ctx, &s, "SELECT * FROM performance_scorecards WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *PerformanceScorecardRepo) List(ctx context.Context, tenantID int64, period string) ([]*ai.PerformanceScorecard, error) {
	var result []*ai.PerformanceScorecard
	if period != "" {
		err := r.db.SelectContext(ctx, &result, "SELECT * FROM performance_scorecards WHERE tenant_id = ? AND period = ? ORDER BY overall_score DESC", tenantID, period)
		return result, err
	}
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM performance_scorecards WHERE tenant_id = ? ORDER BY created_at DESC", tenantID)
	return result, err
}

func (r *PerformanceScorecardRepo) GetByAgentAndPeriod(ctx context.Context, agentID int64, period string) (*ai.PerformanceScorecard, error) {
	var s ai.PerformanceScorecard
	if err := r.db.GetContext(ctx, &s, "SELECT * FROM performance_scorecards WHERE agent_id = ? AND period = ?", agentID, period); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}
