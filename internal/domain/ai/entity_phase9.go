package ai

import "time"

// DigitalEmployee represents an AI bot / digital agent.
type DigitalEmployee struct {
	ID          int64     `db:"id" json:"id"`
	TenantID    int64     `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	AvatarURL   string    `db:"avatar_url" json:"avatar_url"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// DigitalEmployeeScene represents a conversation scenario for a digital employee.
type DigitalEmployeeScene struct {
	ID                 int64     `db:"id" json:"id"`
	DigitalEmployeeID  int64     `db:"digital_employee_id" json:"digital_employee_id"`
	TenantID           int64     `db:"tenant_id" json:"tenant_id"`
	Name               string    `db:"name" json:"name"`
	Intents            string    `db:"intents" json:"intents"`       // JSON array of intent configs
	FAQs               string    `db:"faqs" json:"faqs"`             // JSON array of FAQ entries
	TransferSkillGroup *int64    `db:"transfer_skill_group" json:"transfer_skill_group,omitempty"`
	Status             string    `db:"status" json:"status"` // draft, published
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

// QARule represents a quality assurance inspection rule.
type QARule struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  int64     `db:"tenant_id" json:"tenant_id"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"` // keyword, regex, silence, speed, interruption, energy, duration, entity, role, abnormal_hangup, llm
	Config    string    `db:"config" json:"config"` // JSON rule config
	Severity  string    `db:"severity" json:"severity"` // info, warning, critical
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// QAScheme represents a quality assurance inspection scheme (a collection of rules).
type QAScheme struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  int64     `db:"tenant_id" json:"tenant_id"`
	Name      string    `db:"name" json:"name"`
	RuleIDs   string    `db:"rule_ids" json:"rule_ids"` // JSON array of rule IDs with weights
	IsDefault bool      `db:"is_default" json:"is_default"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// SchemeRuleWeight represents a rule with its weight within a scheme.
type SchemeRuleWeight struct {
	RuleID int64   `json:"rule_id"`
	Weight float64 `json:"weight"` // 0-100
}

// QAResult represents the result of a quality inspection on a call.
type QAResult struct {
	ID         int64     `db:"id" json:"id"`
	TenantID   int64     `db:"tenant_id" json:"tenant_id"`
	CallID     int64     `db:"call_id" json:"call_id"`
	SchemeID   int64     `db:"scheme_id" json:"scheme_id"`
	Score      float64   `db:"score" json:"score"` // 0-100
	Details    string    `db:"details" json:"details"` // JSON array of rule results
	Status     string    `db:"status" json:"status"` // completed, appealed, reviewed
	AppealNote string    `db:"appeal_note" json:"appeal_note,omitempty"`
	ReviewNote string    `db:"review_note" json:"review_note,omitempty"`
	ReviewerID *int64    `db:"reviewer_id" json:"reviewer_id,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// ASRHotwords represents a custom ASR hotword vocabulary.
type ASRHotwords struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  int64     `db:"tenant_id" json:"tenant_id"`
	Name      string    `db:"name" json:"name"`
	Words     string    `db:"words" json:"words"` // JSON array of hotword entries
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// PerformanceScorecard represents an agent performance scorecard.
type PerformanceScorecard struct {
	ID              int64     `db:"id" json:"id"`
	TenantID        int64     `db:"tenant_id" json:"tenant_id"`
	AgentID         int64     `db:"agent_id" json:"agent_id"`
	Period          string    `db:"period" json:"period"` // YYYY-MM
	TotalCalls      int       `db:"total_calls" json:"total_calls"`
	AvgHandleTime   float64   `db:"avg_handle_time" json:"avg_handle_time"`
	AvgQAScore      float64   `db:"avg_qa_score" json:"avg_qa_score"`
	CSATScore       float64   `db:"csat_score" json:"csat_score"`
	FirstCallResolv float64   `db:"first_call_resolution" json:"first_call_resolution"` // percentage
	Adherence       float64   `db:"adherence" json:"adherence"` // schedule adherence %
	OverallScore    float64   `db:"overall_score" json:"overall_score"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// QARuleResult represents a single rule's result within a QA inspection.
type QARuleResult struct {
	RuleID   int64   `json:"rule_id"`
	RuleName string  `json:"rule_name"`
	RuleType string  `json:"rule_type"`
	Passed   bool    `json:"passed"`
	Score    float64 `json:"score"`
	Detail   string  `json:"detail"`
}

// Valid QA rule types.
const (
	QARuleTypeKeyword         = "keyword"
	QARuleTypeRegex           = "regex"
	QARuleTypeSilence         = "silence"
	QARuleTypeSpeed           = "speed"
	QARuleTypeInterruption    = "interruption"
	QARuleTypeEnergy          = "energy"
	QARuleTypeDuration        = "duration"
	QARuleTypeEntity          = "entity"
	QARuleTypeRole            = "role"
	QARuleTypeAbnormalHangup  = "abnormal_hangup"
	QARuleTypeLLM             = "llm"
)

var validQARuleTypes = map[string]bool{
	QARuleTypeKeyword: true, QARuleTypeRegex: true, QARuleTypeSilence: true,
	QARuleTypeSpeed: true, QARuleTypeInterruption: true, QARuleTypeEnergy: true,
	QARuleTypeDuration: true, QARuleTypeEntity: true, QARuleTypeRole: true,
	QARuleTypeAbnormalHangup: true, QARuleTypeLLM: true,
}

func IsValidQARuleType(t string) bool {
	return validQARuleTypes[t]
}

const (
	QAResultStatusCompleted = "completed"
	QAResultStatusAppealed  = "appealed"
	QAResultStatusReviewed  = "reviewed"
)

const (
	SceneStatusDraft     = "draft"
	SceneStatusPublished = "published"
)
