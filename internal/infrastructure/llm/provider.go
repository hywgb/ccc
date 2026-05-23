package llm

import "context"

// Provider is the pluggable LLM provider interface for all AI operations.
type Provider interface {
	// Text operations (Phase 8)
	Correct(ctx context.Context, text string) (string, error)
	Expand(ctx context.Context, text string) (string, error)
	Optimize(ctx context.Context, text string) (string, error)

	// Phase 9 AI operations
	Summarize(ctx context.Context, transcript string) (string, error)
	AnalyzeSentiment(ctx context.Context, text string) (SentimentResult, error)
	ExtractTags(ctx context.Context, transcript string) ([]string, error)
	PredictSatisfaction(ctx context.Context, transcript string) (float64, error)
	AnalyzeIVRPath(ctx context.Context, ivrPath string) (string, error)
	JudgeCompletion(ctx context.Context, transcript string) (float64, error)
	ExtractPostCallActions(ctx context.Context, transcript string) ([]string, error)
	AutoFillTicket(ctx context.Context, transcript string) (map[string]string, error)
	RecommendScript(ctx context.Context, transcript string, scripts []string) (string, error)
	QAInspectLLM(ctx context.Context, transcript, prompt string) (float64, string, error)
}

// SentimentResult represents sentiment analysis output.
type SentimentResult struct {
	Label      string  `json:"label"` // positive, negative, neutral
	Confidence float64 `json:"confidence"`
}

// ASRProvider is the interface for speech-to-text services.
type ASRProvider interface {
	Transcribe(ctx context.Context, audioURL string) (string, error)
}

// TTSProvider is the interface for text-to-speech services.
type TTSProvider interface {
	Synthesize(ctx context.Context, text string, voice string) ([]byte, error)
}
