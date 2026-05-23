package llm

import "context"

// StubProvider is a placeholder LLM provider that echoes input.
// Replace with Aliyun Tongyi or other provider in production.
type StubProvider struct{}

func NewStubProvider() *StubProvider { return &StubProvider{} }

func (p *StubProvider) Correct(_ context.Context, text string) (string, error) {
	return "[corrected] " + text, nil
}

func (p *StubProvider) Expand(_ context.Context, text string) (string, error) {
	return "[expanded] " + text, nil
}

func (p *StubProvider) Optimize(_ context.Context, text string) (string, error) {
	return "[optimized] " + text, nil
}

func (p *StubProvider) Summarize(_ context.Context, transcript string) (string, error) {
	if len(transcript) > 100 {
		return "[summary] " + transcript[:100] + "...", nil
	}
	return "[summary] " + transcript, nil
}

func (p *StubProvider) AnalyzeSentiment(_ context.Context, _ string) (SentimentResult, error) {
	return SentimentResult{Label: "neutral", Confidence: 0.85}, nil
}

func (p *StubProvider) ExtractTags(_ context.Context, _ string) ([]string, error) {
	return []string{"general", "inquiry"}, nil
}

func (p *StubProvider) PredictSatisfaction(_ context.Context, _ string) (float64, error) {
	return 4.0, nil // out of 5
}

func (p *StubProvider) AnalyzeIVRPath(_ context.Context, ivrPath string) (string, error) {
	return "[ivr-analysis] " + ivrPath, nil
}

func (p *StubProvider) JudgeCompletion(_ context.Context, _ string) (float64, error) {
	return 4.0, nil // 1-5 scale
}

func (p *StubProvider) ExtractPostCallActions(_ context.Context, _ string) ([]string, error) {
	return []string{"follow up with customer"}, nil
}

func (p *StubProvider) AutoFillTicket(_ context.Context, _ string) (map[string]string, error) {
	return map[string]string{
		"subject":     "Customer inquiry",
		"description": "Auto-filled from transcript",
	}, nil
}

func (p *StubProvider) RecommendScript(_ context.Context, _ string, scripts []string) (string, error) {
	if len(scripts) > 0 {
		return scripts[0], nil
	}
	return "", nil
}

func (p *StubProvider) QAInspectLLM(_ context.Context, _, _ string) (float64, string, error) {
	return 80, "LLM inspection passed (stub)", nil
}
