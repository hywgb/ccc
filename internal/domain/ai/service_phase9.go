package ai

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/divord97/ccc/pkg/snowflake"
)

// DigitalEmployeeService manages digital employee (AI bot) entities and scenes.
type DigitalEmployeeService struct {
	employees DigitalEmployeeRepository
	scenes    DigitalEmployeeSceneRepository
}

func NewDigitalEmployeeService(employees DigitalEmployeeRepository, scenes DigitalEmployeeSceneRepository) *DigitalEmployeeService {
	return &DigitalEmployeeService{employees: employees, scenes: scenes}
}

type CreateDigitalEmployeeInput struct {
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
}

func (s *DigitalEmployeeService) Create(ctx context.Context, in CreateDigitalEmployeeInput) (*DigitalEmployee, error) {
	now := time.Now()
	de := &DigitalEmployee{
		ID:          snowflake.NextID(),
		TenantID:    in.TenantID,
		Name:        in.Name,
		Description: in.Description,
		AvatarURL:   in.AvatarURL,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.employees.Create(ctx, de); err != nil {
		return nil, err
	}
	return de, nil
}

func (s *DigitalEmployeeService) GetByID(ctx context.Context, id int64) (*DigitalEmployee, error) {
	de, err := s.employees.GetByID(ctx, id)
	if err != nil || de == nil {
		return nil, ErrDigitalEmployeeNotFound
	}
	return de, nil
}

func (s *DigitalEmployeeService) Update(ctx context.Context, de *DigitalEmployee) error {
	de.UpdatedAt = time.Now()
	return s.employees.Update(ctx, de)
}

func (s *DigitalEmployeeService) List(ctx context.Context, tenantID int64) ([]*DigitalEmployee, error) {
	return s.employees.List(ctx, tenantID)
}

type CreateSceneInput struct {
	DigitalEmployeeID  int64  `json:"digital_employee_id"`
	TenantID           int64  `json:"tenant_id"`
	Name               string `json:"name"`
	Intents            string `json:"intents"`
	FAQs               string `json:"faqs"`
	TransferSkillGroup *int64 `json:"transfer_skill_group"`
}

func (s *DigitalEmployeeService) CreateScene(ctx context.Context, in CreateSceneInput) (*DigitalEmployeeScene, error) {
	if in.Intents != "" && !json.Valid([]byte(in.Intents)) {
		return nil, ErrInvalidIntentConfig
	}
	now := time.Now()
	scene := &DigitalEmployeeScene{
		ID:                 snowflake.NextID(),
		DigitalEmployeeID:  in.DigitalEmployeeID,
		TenantID:           in.TenantID,
		Name:               in.Name,
		Intents:            in.Intents,
		FAQs:               in.FAQs,
		TransferSkillGroup: in.TransferSkillGroup,
		Status:             SceneStatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := s.scenes.Create(ctx, scene); err != nil {
		return nil, err
	}
	return scene, nil
}

func (s *DigitalEmployeeService) PublishScene(ctx context.Context, sceneID int64) (*DigitalEmployeeScene, error) {
	scene, err := s.scenes.GetByID(ctx, sceneID)
	if err != nil || scene == nil {
		return nil, ErrSceneNotFound
	}
	if scene.Status == SceneStatusPublished {
		return nil, ErrSceneAlreadyPublished
	}
	scene.Status = SceneStatusPublished
	scene.UpdatedAt = time.Now()
	if err := s.scenes.Update(ctx, scene); err != nil {
		return nil, err
	}
	return scene, nil
}

func (s *DigitalEmployeeService) ListScenes(ctx context.Context, digitalEmployeeID int64) ([]*DigitalEmployeeScene, error) {
	return s.scenes.List(ctx, digitalEmployeeID)
}

// IntentMatchResult represents the result of intent matching.
type IntentMatchResult struct {
	Matched    bool   `json:"matched"`
	IntentName string `json:"intent_name,omitempty"`
	Response   string `json:"response,omitempty"`
	Transfer   bool   `json:"transfer"`
}

// IntentConfig represents a single intent configuration entry.
type IntentConfig struct {
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
	Response string   `json:"response"`
	Transfer bool     `json:"transfer"`
}

// MatchIntent checks user input against a scene's intent configuration.
func (s *DigitalEmployeeService) MatchIntent(ctx context.Context, sceneID int64, userInput string) (*IntentMatchResult, error) {
	scene, err := s.scenes.GetByID(ctx, sceneID)
	if err != nil || scene == nil {
		return nil, ErrSceneNotFound
	}

	var intents []IntentConfig
	if scene.Intents != "" {
		if err := json.Unmarshal([]byte(scene.Intents), &intents); err != nil {
			return &IntentMatchResult{Matched: false}, nil
		}
	}

	for _, intent := range intents {
		for _, kw := range intent.Keywords {
			if containsIgnoreCase(userInput, kw) {
				return &IntentMatchResult{
					Matched:    true,
					IntentName: intent.Name,
					Response:   intent.Response,
					Transfer:   intent.Transfer,
				}, nil
			}
		}
	}

	return &IntentMatchResult{Matched: false}, nil
}

// containsIgnoreCase checks if s contains substr (case-insensitive, Unicode-safe).
func containsIgnoreCase(s, substr string) bool {
	return len(substr) > 0 && strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// QualityInspectionService manages QA rules, schemes, and results.
type QualityInspectionService struct {
	rules   QARuleRepository
	schemes QASchemeRepository
	results QAResultRepository
}

func NewQualityInspectionService(rules QARuleRepository, schemes QASchemeRepository, results QAResultRepository) *QualityInspectionService {
	return &QualityInspectionService{rules: rules, schemes: schemes, results: results}
}

type CreateQARuleInput struct {
	TenantID int64  `json:"tenant_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Config   string `json:"config"`
	Severity string `json:"severity"`
}

func (s *QualityInspectionService) CreateRule(ctx context.Context, in CreateQARuleInput) (*QARule, error) {
	if !IsValidQARuleType(in.Type) {
		return nil, ErrInvalidQARuleType
	}
	if in.Severity == "" {
		in.Severity = "warning"
	}
	now := time.Now()
	rule := &QARule{
		ID:        snowflake.NextID(),
		TenantID:  in.TenantID,
		Name:      in.Name,
		Type:      in.Type,
		Config:    in.Config,
		Severity:  in.Severity,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.rules.Create(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *QualityInspectionService) GetRule(ctx context.Context, id int64) (*QARule, error) {
	r, err := s.rules.GetByID(ctx, id)
	if err != nil || r == nil {
		return nil, ErrQARuleNotFound
	}
	return r, nil
}

func (s *QualityInspectionService) UpdateRule(ctx context.Context, rule *QARule) error {
	rule.UpdatedAt = time.Now()
	return s.rules.Update(ctx, rule)
}

func (s *QualityInspectionService) DeleteRule(ctx context.Context, id int64) error {
	return s.rules.Delete(ctx, id)
}

func (s *QualityInspectionService) ListRules(ctx context.Context, tenantID int64) ([]*QARule, error) {
	return s.rules.List(ctx, tenantID)
}

type CreateQASchemeInput struct {
	TenantID  int64              `json:"tenant_id"`
	Name      string             `json:"name"`
	RuleIDs   []SchemeRuleWeight `json:"rule_ids"`
	IsDefault bool               `json:"is_default"`
}

func (s *QualityInspectionService) CreateScheme(ctx context.Context, in CreateQASchemeInput) (*QAScheme, error) {
	ruleJSON, err := json.Marshal(in.RuleIDs)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	scheme := &QAScheme{
		ID:        snowflake.NextID(),
		TenantID:  in.TenantID,
		Name:      in.Name,
		RuleIDs:   string(ruleJSON),
		IsDefault: in.IsDefault,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.schemes.Create(ctx, scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}

func (s *QualityInspectionService) GetScheme(ctx context.Context, id int64) (*QAScheme, error) {
	scheme, err := s.schemes.GetByID(ctx, id)
	if err != nil || scheme == nil {
		return nil, ErrQASchemeNotFound
	}
	return scheme, nil
}

func (s *QualityInspectionService) UpdateScheme(ctx context.Context, scheme *QAScheme) error {
	scheme.UpdatedAt = time.Now()
	return s.schemes.Update(ctx, scheme)
}

func (s *QualityInspectionService) DeleteScheme(ctx context.Context, id int64) error {
	return s.schemes.Delete(ctx, id)
}

func (s *QualityInspectionService) ListSchemes(ctx context.Context, tenantID int64) ([]*QAScheme, error) {
	return s.schemes.List(ctx, tenantID)
}

// RunInspection runs quality inspection on a transcript using a scheme.
func (s *QualityInspectionService) RunInspection(ctx context.Context, tenantID, callID, schemeID int64, transcript string) (*QAResult, error) {
	if transcript == "" {
		return nil, ErrEmptyTranscript
	}

	scheme, err := s.schemes.GetByID(ctx, schemeID)
	if err != nil || scheme == nil {
		return nil, ErrQASchemeNotFound
	}

	var weights []SchemeRuleWeight
	if err := json.Unmarshal([]byte(scheme.RuleIDs), &weights); err != nil {
		return nil, err
	}

	ruleIDs := make([]int64, len(weights))
	for i, w := range weights {
		ruleIDs[i] = w.RuleID
	}
	rules, err := s.rules.ListByIDs(ctx, ruleIDs)
	if err != nil {
		return nil, err
	}

	ruleMap := make(map[int64]*QARule, len(rules))
	for _, r := range rules {
		ruleMap[r.ID] = r
	}

	var ruleResults []QARuleResult
	var totalScore, totalWeight float64

	for _, w := range weights {
		rule, ok := ruleMap[w.RuleID]
		if !ok {
			continue
		}
		rr := evaluateRule(rule, transcript)
		ruleResults = append(ruleResults, rr)
		totalScore += rr.Score * w.Weight
		totalWeight += w.Weight
	}

	finalScore := float64(0)
	if totalWeight > 0 {
		finalScore = totalScore / totalWeight
	}

	detailsJSON, _ := json.Marshal(ruleResults)

	now := time.Now()
	result := &QAResult{
		ID:        snowflake.NextID(),
		TenantID:  tenantID,
		CallID:    callID,
		SchemeID:  schemeID,
		Score:     finalScore,
		Details:   string(detailsJSON),
		Status:    QAResultStatusCompleted,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.results.Create(ctx, result); err != nil {
		return nil, err
	}
	return result, nil
}

// evaluateRule evaluates a single QA rule against a transcript.
func evaluateRule(rule *QARule, transcript string) QARuleResult {
	rr := QARuleResult{
		RuleID:   rule.ID,
		RuleName: rule.Name,
		RuleType: rule.Type,
	}

	switch rule.Type {
	case QARuleTypeKeyword:
		rr = evaluateKeywordRule(rule, transcript)
	case QARuleTypeSilence:
		rr.Passed = true
		rr.Score = 100
		rr.Detail = "silence check passed (stub)"
	case QARuleTypeSpeed:
		rr.Passed = true
		rr.Score = 100
		rr.Detail = "speed check passed (stub)"
	case QARuleTypeLLM:
		rr.Passed = true
		rr.Score = 80
		rr.Detail = "LLM analysis passed (stub)"
	default:
		rr.Passed = true
		rr.Score = 100
		rr.Detail = "rule check passed (stub)"
	}
	return rr
}

// KeywordRuleConfig represents the JSON config for a keyword rule.
type KeywordRuleConfig struct {
	Keywords []string `json:"keywords"`
	Require  string   `json:"require"` // "present" or "absent"
}

func evaluateKeywordRule(rule *QARule, transcript string) QARuleResult {
	rr := QARuleResult{
		RuleID:   rule.ID,
		RuleName: rule.Name,
		RuleType: rule.Type,
	}

	var cfg KeywordRuleConfig
	if err := json.Unmarshal([]byte(rule.Config), &cfg); err != nil {
		rr.Passed = false
		rr.Score = 0
		rr.Detail = "invalid keyword config"
		return rr
	}

	for _, kw := range cfg.Keywords {
		found := containsIgnoreCase(transcript, kw)
		if cfg.Require == "absent" && found {
			rr.Passed = false
			rr.Score = 0
			rr.Detail = "forbidden keyword found: " + kw
			return rr
		}
		if cfg.Require == "present" && found {
			rr.Passed = true
			rr.Score = 100
			rr.Detail = "required keyword found: " + kw
			return rr
		}
	}

	if cfg.Require == "present" {
		rr.Passed = false
		rr.Score = 0
		rr.Detail = "required keyword not found"
	} else {
		rr.Passed = true
		rr.Score = 100
		rr.Detail = "no forbidden keywords found"
	}
	return rr
}

// Appeal submits an appeal for a QA result.
func (s *QualityInspectionService) Appeal(ctx context.Context, resultID int64, note string) (*QAResult, error) {
	result, err := s.results.GetByID(ctx, resultID)
	if err != nil || result == nil {
		return nil, ErrQAResultNotFound
	}
	if result.Status != QAResultStatusCompleted {
		return nil, ErrQAResultNotAppealable
	}
	result.Status = QAResultStatusAppealed
	result.AppealNote = note
	result.UpdatedAt = time.Now()
	if err := s.results.Update(ctx, result); err != nil {
		return nil, err
	}
	return result, nil
}

// Review completes a review of an appealed QA result.
func (s *QualityInspectionService) Review(ctx context.Context, resultID int64, reviewerID int64, note string, newScore float64) (*QAResult, error) {
	result, err := s.results.GetByID(ctx, resultID)
	if err != nil || result == nil {
		return nil, ErrQAResultNotFound
	}
	if result.Status != QAResultStatusAppealed {
		return nil, ErrQAResultNotReviewable
	}
	result.Status = QAResultStatusReviewed
	result.ReviewerID = &reviewerID
	result.ReviewNote = note
	result.Score = newScore
	result.UpdatedAt = time.Now()
	if err := s.results.Update(ctx, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *QualityInspectionService) GetResult(ctx context.Context, id int64) (*QAResult, error) {
	r, err := s.results.GetByID(ctx, id)
	if err != nil || r == nil {
		return nil, ErrQAResultNotFound
	}
	return r, nil
}

func (s *QualityInspectionService) ListResults(ctx context.Context, tenantID int64, offset, limit int) ([]*QAResult, error) {
	return s.results.List(ctx, tenantID, offset, limit)
}

// ASRHotwordsService manages ASR custom hotword vocabularies.
type ASRHotwordsService struct {
	repo ASRHotwordsRepository
}

func NewASRHotwordsService(repo ASRHotwordsRepository) *ASRHotwordsService {
	return &ASRHotwordsService{repo: repo}
}

type CreateASRHotwordsInput struct {
	TenantID int64  `json:"tenant_id"`
	Name     string `json:"name"`
	Words    string `json:"words"`
}

func (s *ASRHotwordsService) Create(ctx context.Context, in CreateASRHotwordsInput) (*ASRHotwords, error) {
	now := time.Now()
	h := &ASRHotwords{
		ID:        snowflake.NextID(),
		TenantID:  in.TenantID,
		Name:      in.Name,
		Words:     in.Words,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, h); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *ASRHotwordsService) GetByID(ctx context.Context, id int64) (*ASRHotwords, error) {
	h, err := s.repo.GetByID(ctx, id)
	if err != nil || h == nil {
		return nil, ErrASRHotwordsNotFound
	}
	return h, nil
}

func (s *ASRHotwordsService) Update(ctx context.Context, h *ASRHotwords) error {
	h.UpdatedAt = time.Now()
	return s.repo.Update(ctx, h)
}

func (s *ASRHotwordsService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *ASRHotwordsService) List(ctx context.Context, tenantID int64) ([]*ASRHotwords, error) {
	return s.repo.List(ctx, tenantID)
}

// PerformanceScorecardService manages agent performance scorecards.
type PerformanceScorecardService struct {
	repo PerformanceScorecardRepository
}

func NewPerformanceScorecardService(repo PerformanceScorecardRepository) *PerformanceScorecardService {
	return &PerformanceScorecardService{repo: repo}
}

type GenerateScorecardInput struct {
	TenantID        int64   `json:"tenant_id"`
	AgentID         int64   `json:"agent_id"`
	Period          string  `json:"period"`
	TotalCalls      int     `json:"total_calls"`
	AvgHandleTime   float64 `json:"avg_handle_time"`
	AvgQAScore      float64 `json:"avg_qa_score"`
	CSATScore       float64 `json:"csat_score"`
	FirstCallResolv float64 `json:"first_call_resolution"`
	Adherence       float64 `json:"adherence"`
}

func (s *PerformanceScorecardService) Generate(ctx context.Context, in GenerateScorecardInput) (*PerformanceScorecard, error) {
	// Weighted overall score: QA 30%, CSAT 30%, FCR 20%, Adherence 20%
	overall := in.AvgQAScore*0.3 + in.CSATScore*0.3 + in.FirstCallResolv*0.2 + in.Adherence*0.2

	sc := &PerformanceScorecard{
		ID:              snowflake.NextID(),
		TenantID:        in.TenantID,
		AgentID:         in.AgentID,
		Period:          in.Period,
		TotalCalls:      in.TotalCalls,
		AvgHandleTime:   in.AvgHandleTime,
		AvgQAScore:      in.AvgQAScore,
		CSATScore:       in.CSATScore,
		FirstCallResolv: in.FirstCallResolv,
		Adherence:       in.Adherence,
		OverallScore:    overall,
		CreatedAt:       time.Now(),
	}
	if err := s.repo.Create(ctx, sc); err != nil {
		return nil, err
	}
	return sc, nil
}

func (s *PerformanceScorecardService) List(ctx context.Context, tenantID int64, period string) ([]*PerformanceScorecard, error) {
	return s.repo.List(ctx, tenantID, period)
}
