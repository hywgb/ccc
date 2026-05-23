package aianalysis

import (
	"context"
	"encoding/json"

	"github.com/divord97/ccc/internal/infrastructure/llm"
	"github.com/rs/zerolog"
)

// Service provides AI analysis operations for calls and sessions.
type Service struct {
	llm    llm.Provider
	logger zerolog.Logger
}

func NewService(provider llm.Provider, logger zerolog.Logger) *Service {
	return &Service{llm: provider, logger: logger}
}

// SummaryResult holds the AI-generated summary.
type SummaryResult struct {
	CallID  int64  `json:"call_id"`
	Summary string `json:"summary"`
}

// GenerateSummary produces an AI summary of a call transcript.
func (s *Service) GenerateSummary(ctx context.Context, callID int64, transcript string) (*SummaryResult, error) {
	summary, err := s.llm.Summarize(ctx, transcript)
	if err != nil {
		s.logger.Error().Err(err).Int64("call_id", callID).Msg("ai: summary failed")
		return nil, err
	}
	return &SummaryResult{CallID: callID, Summary: summary}, nil
}

// SentimentResult holds the AI sentiment analysis output.
type SentimentResult struct {
	CallID     int64   `json:"call_id"`
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
}

// AnalyzeSentiment runs sentiment analysis on a transcript.
func (s *Service) AnalyzeSentiment(ctx context.Context, callID int64, transcript string) (*SentimentResult, error) {
	result, err := s.llm.AnalyzeSentiment(ctx, transcript)
	if err != nil {
		s.logger.Error().Err(err).Int64("call_id", callID).Msg("ai: sentiment failed")
		return nil, err
	}
	return &SentimentResult{CallID: callID, Label: result.Label, Confidence: result.Confidence}, nil
}

// TagResult holds AI-extracted tags.
type TagResult struct {
	CallID int64    `json:"call_id"`
	Tags   []string `json:"tags"`
}

// ExtractTags generates AI tags for a transcript.
func (s *Service) ExtractTags(ctx context.Context, callID int64, transcript string) (*TagResult, error) {
	tags, err := s.llm.ExtractTags(ctx, transcript)
	if err != nil {
		s.logger.Error().Err(err).Int64("call_id", callID).Msg("ai: tag extraction failed")
		return nil, err
	}
	return &TagResult{CallID: callID, Tags: tags}, nil
}

// SatisfactionResult holds AI satisfaction prediction.
type SatisfactionResult struct {
	CallID int64   `json:"call_id"`
	Score  float64 `json:"score"` // 1-5
}

// PredictSatisfaction predicts customer satisfaction from transcript.
func (s *Service) PredictSatisfaction(ctx context.Context, callID int64, transcript string) (*SatisfactionResult, error) {
	score, err := s.llm.PredictSatisfaction(ctx, transcript)
	if err != nil {
		s.logger.Error().Err(err).Int64("call_id", callID).Msg("ai: satisfaction prediction failed")
		return nil, err
	}
	return &SatisfactionResult{CallID: callID, Score: score}, nil
}

// IVRAnalysisResult holds AI IVR path analysis.
type IVRAnalysisResult struct {
	CallID   int64  `json:"call_id"`
	Analysis string `json:"analysis"`
}

// AnalyzeIVRPath analyzes a call's IVR navigation path.
func (s *Service) AnalyzeIVRPath(ctx context.Context, callID int64, ivrPath string) (*IVRAnalysisResult, error) {
	analysis, err := s.llm.AnalyzeIVRPath(ctx, ivrPath)
	if err != nil {
		s.logger.Error().Err(err).Int64("call_id", callID).Msg("ai: ivr analysis failed")
		return nil, err
	}
	return &IVRAnalysisResult{CallID: callID, Analysis: analysis}, nil
}

// CompletionResult holds AI completion score.
type CompletionResult struct {
	CallID int64   `json:"call_id"`
	Score  float64 `json:"score"` // 1-5
}

// JudgeCompletion judges whether the customer's issue was resolved.
func (s *Service) JudgeCompletion(ctx context.Context, callID int64, transcript string) (*CompletionResult, error) {
	score, err := s.llm.JudgeCompletion(ctx, transcript)
	if err != nil {
		s.logger.Error().Err(err).Int64("call_id", callID).Msg("ai: completion judgement failed")
		return nil, err
	}
	return &CompletionResult{CallID: callID, Score: score}, nil
}

// PostCallActionsResult holds extracted post-call action items.
type PostCallActionsResult struct {
	CallID  int64    `json:"call_id"`
	Actions []string `json:"actions"`
}

// ExtractPostCallActions extracts action items from a call transcript.
func (s *Service) ExtractPostCallActions(ctx context.Context, callID int64, transcript string) (*PostCallActionsResult, error) {
	actions, err := s.llm.ExtractPostCallActions(ctx, transcript)
	if err != nil {
		s.logger.Error().Err(err).Int64("call_id", callID).Msg("ai: post-call actions extraction failed")
		return nil, err
	}
	return &PostCallActionsResult{CallID: callID, Actions: actions}, nil
}

// AutoFillResult holds auto-filled ticket fields.
type AutoFillResult struct {
	CallID int64             `json:"call_id"`
	Fields map[string]string `json:"fields"`
}

// AutoFillTicket extracts information from transcript to auto-fill ticket fields.
func (s *Service) AutoFillTicket(ctx context.Context, callID int64, transcript string) (*AutoFillResult, error) {
	fields, err := s.llm.AutoFillTicket(ctx, transcript)
	if err != nil {
		s.logger.Error().Err(err).Int64("call_id", callID).Msg("ai: auto-fill failed")
		return nil, err
	}
	return &AutoFillResult{CallID: callID, Fields: fields}, nil
}

// ScriptRecommendation holds recommended agent script.
type ScriptRecommendation struct {
	CallID     int64  `json:"call_id"`
	Recommended string `json:"recommended"`
}

// RecommendScript finds the best matching script for the current conversation.
func (s *Service) RecommendScript(ctx context.Context, callID int64, transcript string, scripts []string) (*ScriptRecommendation, error) {
	rec, err := s.llm.RecommendScript(ctx, transcript, scripts)
	if err != nil {
		s.logger.Error().Err(err).Int64("call_id", callID).Msg("ai: script recommendation failed")
		return nil, err
	}
	return &ScriptRecommendation{CallID: callID, Recommended: rec}, nil
}

// BatchTagAnalysisInput holds input for batch tag analysis.
type BatchTagAnalysisInput struct {
	TenantID    int64    `json:"tenant_id"`
	CallIDs     []int64  `json:"call_ids"`
	Transcripts []string `json:"transcripts"`
}

// BatchTagResult holds batch tag analysis results.
type BatchTagResult struct {
	Results []TagResult `json:"results"`
}

// RunBatchTagAnalysis runs AI tag analysis on multiple transcripts.
func (s *Service) RunBatchTagAnalysis(ctx context.Context, in BatchTagAnalysisInput) (*BatchTagResult, error) {
	var results []TagResult
	for i, callID := range in.CallIDs {
		if i >= len(in.Transcripts) {
			break
		}
		tags, err := s.llm.ExtractTags(ctx, in.Transcripts[i])
		if err != nil {
			s.logger.Warn().Err(err).Int64("call_id", callID).Msg("ai: batch tag failed for call")
			continue
		}
		results = append(results, TagResult{CallID: callID, Tags: tags})
	}
	return &BatchTagResult{Results: results}, nil
}

// HotwordAnalysisResult holds hotword frequency analysis.
type HotwordAnalysisResult struct {
	Hotwords []HotwordEntry `json:"hotwords"`
}

// HotwordEntry is a single hotword with its frequency.
type HotwordEntry struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// AnalyzeHotwords extracts and counts frequently used words from transcripts.
func (s *Service) AnalyzeHotwords(_ context.Context, transcripts []string) (*HotwordAnalysisResult, error) {
	freq := make(map[string]int)
	for _, t := range transcripts {
		tags, _ := json.Marshal(t)
		_ = tags
		words := splitWords(t)
		for _, w := range words {
			if len(w) >= 2 {
				freq[w]++
			}
		}
	}
	var entries []HotwordEntry
	for w, c := range freq {
		entries = append(entries, HotwordEntry{Word: w, Count: c})
	}
	return &HotwordAnalysisResult{Hotwords: entries}, nil
}

// splitWords is a simple word splitter (for CJK, each char pair is a "word"; for ASCII, split by space).
func splitWords(s string) []string {
	var words []string
	current := ""
	for _, r := range s {
		if r == ' ' || r == ',' || r == '。' || r == '，' || r == '？' || r == '！' || r == '\n' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}
