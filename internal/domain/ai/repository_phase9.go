package ai

import "context"

type DigitalEmployeeRepository interface {
	Create(ctx context.Context, de *DigitalEmployee) error
	GetByID(ctx context.Context, id int64) (*DigitalEmployee, error)
	Update(ctx context.Context, de *DigitalEmployee) error
	List(ctx context.Context, tenantID int64) ([]*DigitalEmployee, error)
}

type DigitalEmployeeSceneRepository interface {
	Create(ctx context.Context, s *DigitalEmployeeScene) error
	GetByID(ctx context.Context, id int64) (*DigitalEmployeeScene, error)
	Update(ctx context.Context, s *DigitalEmployeeScene) error
	List(ctx context.Context, digitalEmployeeID int64) ([]*DigitalEmployeeScene, error)
}

type QARuleRepository interface {
	Create(ctx context.Context, r *QARule) error
	GetByID(ctx context.Context, id int64) (*QARule, error)
	Update(ctx context.Context, r *QARule) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, tenantID int64) ([]*QARule, error)
	ListByIDs(ctx context.Context, ids []int64) ([]*QARule, error)
}

type QASchemeRepository interface {
	Create(ctx context.Context, s *QAScheme) error
	GetByID(ctx context.Context, id int64) (*QAScheme, error)
	Update(ctx context.Context, s *QAScheme) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, tenantID int64) ([]*QAScheme, error)
}

type QAResultRepository interface {
	Create(ctx context.Context, r *QAResult) error
	GetByID(ctx context.Context, id int64) (*QAResult, error)
	Update(ctx context.Context, r *QAResult) error
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*QAResult, error)
	ListByCallID(ctx context.Context, callID int64) ([]*QAResult, error)
}

type ASRHotwordsRepository interface {
	Create(ctx context.Context, h *ASRHotwords) error
	GetByID(ctx context.Context, id int64) (*ASRHotwords, error)
	Update(ctx context.Context, h *ASRHotwords) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, tenantID int64) ([]*ASRHotwords, error)
}

type PerformanceScorecardRepository interface {
	Create(ctx context.Context, s *PerformanceScorecard) error
	GetByID(ctx context.Context, id int64) (*PerformanceScorecard, error)
	List(ctx context.Context, tenantID int64, period string) ([]*PerformanceScorecard, error)
	GetByAgentAndPeriod(ctx context.Context, agentID int64, period string) (*PerformanceScorecard, error)
}
