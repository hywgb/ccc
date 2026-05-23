package aianalysis

import (
	"context"
	"os"
	"testing"

	"github.com/divord97/ccc/internal/infrastructure/llm"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestService() *Service {
	logger := zerolog.New(os.Stdout)
	return NewService(llm.NewStubProvider(), logger)
}

func TestRealtimeTranscription_StreamPush(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	// Summarize doubles as transcription stream output processing
	result, err := svc.GenerateSummary(ctx, 1001, "客户：你好，我想查询我的订单状态。坐席：好的，请提供您的订单号。")
	require.NoError(t, err)
	assert.Contains(t, result.Summary, "[summary]")
}

func TestRealtimeAssist_ScriptRecommend(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	scripts := []string{"感谢您的来电，请问有什么可以帮助您？", "请稍等，我为您查询一下。", "感谢您的耐心等待。"}
	rec, err := svc.RecommendScript(ctx, 1002, "客户：我要查快递", scripts)
	require.NoError(t, err)
	assert.Equal(t, scripts[0], rec.Recommended)
}

func TestAutoFormFill_ExtractFields(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	result, err := svc.AutoFillTicket(ctx, 1003, "客户：我的订单号是12345，商品有质量问题需要退换。")
	require.NoError(t, err)
	assert.NotEmpty(t, result.Fields)
	assert.Contains(t, result.Fields, "subject")
}

func TestSessionTagAnalysis_Classification(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	result, err := svc.ExtractTags(ctx, 1004, "客户咨询退款流程和退货地址")
	require.NoError(t, err)
	assert.NotEmpty(t, result.Tags)
}

func TestAISatisfaction_Prediction(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	result, err := svc.PredictSatisfaction(ctx, 1005, "服务非常满意，谢谢你的帮助！")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, result.Score, float64(1))
	assert.LessOrEqual(t, result.Score, float64(5))
}

func TestSentimentAnalysis_Classification(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	result, err := svc.AnalyzeSentiment(ctx, 1006, "这个服务太差了，我非常不满意！")
	require.NoError(t, err)
	assert.NotEmpty(t, result.Label)
	assert.Greater(t, result.Confidence, float64(0))
}

func TestIVRAnalysis_PathSummary(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	result, err := svc.AnalyzeIVRPath(ctx, 1007, "welcome→menu→queue→agent")
	require.NoError(t, err)
	assert.Contains(t, result.Analysis, "ivr-analysis")
}

func TestCompletionScore_Judgement(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	result, err := svc.JudgeCompletion(ctx, 1008, "问题已经解决，客户表示满意。")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, result.Score, float64(1))
	assert.LessOrEqual(t, result.Score, float64(5))
}

func TestPostCallActions_Extract(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	result, err := svc.ExtractPostCallActions(ctx, 1009, "需要在24小时内给客户回电确认退款进度")
	require.NoError(t, err)
	assert.NotEmpty(t, result.Actions)
}

func TestBatchTagAnalysis(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	result, err := svc.RunBatchTagAnalysis(ctx, BatchTagAnalysisInput{
		TenantID:    1,
		CallIDs:     []int64{1, 2, 3},
		Transcripts: []string{"退款相关", "投诉相关", "咨询相关"},
	})
	require.NoError(t, err)
	assert.Len(t, result.Results, 3)
}

func TestHotwordAnalysis(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	result, err := svc.AnalyzeHotwords(ctx, []string{
		"退款 退款 投诉",
		"退款 咨询 查询",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.Hotwords)
}
