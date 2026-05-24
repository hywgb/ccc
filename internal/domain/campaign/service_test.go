package campaign

import (
	"context"
	"testing"

	"github.com/divord97/ccc/pkg/snowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	_ = snowflake.Init(1)
}

func newTestService() *CampaignService {
	return NewCampaignService(NewMockCampaignRepo(), NewMockCampaignCaseRepo(), nil)
}

func TestCampaignService_Create_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, err := svc.Create(ctx, CreateCampaignInput{
		TenantID:     1,
		Name:         "Test Campaign",
		DialingMode:  DialingModePredictive,
		SkillGroupID: 10,
	})

	require.NoError(t, err)
	assert.Equal(t, CampaignStatusDraft, c.Status)
	assert.Equal(t, DialingModePredictive, c.DialingMode)
	assert.Equal(t, "Test Campaign", c.Name)
	assert.NotZero(t, c.ID)
}

func TestCampaignService_Create_InvalidMode(t *testing.T) {
	svc := newTestService()

	_, err := svc.Create(context.Background(), CreateCampaignInput{
		TenantID:     1,
		Name:         "Bad",
		DialingMode:  "invalid",
		SkillGroupID: 10,
	})
	assert.ErrorIs(t, err, ErrInvalidDialingMode)
}

func TestCampaignService_Start_WithValidCases(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, _ := svc.Create(ctx, CreateCampaignInput{
		TenantID: 1, Name: "C", DialingMode: DialingModePreview, SkillGroupID: 10,
	})

	_ = svc.ImportCases(ctx, c.ID, []CaseInput{
		{PhoneNumber: "+8613800001111", CustomerName: "张三"},
		{PhoneNumber: "+8613800002222", CustomerName: "李四"},
	})

	started, err := svc.Start(ctx, c.ID)
	require.NoError(t, err)
	assert.Equal(t, CampaignStatusRunning, started.Status)
	assert.NotNil(t, started.StartedAt)
	assert.Equal(t, 2, started.TotalCases)
}

func TestCampaignService_Start_NoCases(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, _ := svc.Create(ctx, CreateCampaignInput{
		TenantID: 1, Name: "Empty", DialingMode: DialingModePreview, SkillGroupID: 10,
	})

	_, err := svc.Start(ctx, c.ID)
	assert.ErrorIs(t, err, ErrCampaignNoCases)
}

func TestCampaignService_Pause_Running(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, _ := svc.Create(ctx, CreateCampaignInput{
		TenantID: 1, Name: "C", DialingMode: DialingModeProgressive, SkillGroupID: 10,
	})
	_ = svc.ImportCases(ctx, c.ID, []CaseInput{{PhoneNumber: "+86138"}})
	_, _ = svc.Start(ctx, c.ID)

	paused, err := svc.Pause(ctx, c.ID)
	require.NoError(t, err)
	assert.Equal(t, CampaignStatusPaused, paused.Status)
}

func TestCampaignService_Abort(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, _ := svc.Create(ctx, CreateCampaignInput{
		TenantID: 1, Name: "C", DialingMode: DialingModePower, SkillGroupID: 10,
	})
	_ = svc.ImportCases(ctx, c.ID, []CaseInput{{PhoneNumber: "+86138"}})
	_, _ = svc.Start(ctx, c.ID)

	aborted, err := svc.Abort(ctx, c.ID)
	require.NoError(t, err)
	assert.Equal(t, CampaignStatusAborted, aborted.Status)
}

func TestCampaignService_ImportCases_DNCFilter(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, _ := svc.Create(ctx, CreateCampaignInput{
		TenantID: 1, Name: "C", DialingMode: DialingModePredictive, SkillGroupID: 10,
	})

	err := svc.ImportCases(ctx, c.ID, []CaseInput{
		{PhoneNumber: "+8613800001111"},
		{PhoneNumber: ""},
		{PhoneNumber: "+8613800003333"},
	})
	require.NoError(t, err)

	cases, total, err := svc.ListCases(ctx, c.ID, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, cases, 2)
}

func TestCampaignService_DialingMode_Predictive(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, _ := svc.Create(ctx, CreateCampaignInput{
		TenantID: 1, Name: "Predictive", DialingMode: DialingModePredictive, SkillGroupID: 10,
	})

	assert.Equal(t, DialingModePredictive, c.DialingMode)
	assert.Equal(t, 1.5, c.RatioMultiplier)
	assert.Equal(t, 3.0, c.MaxAbandonRate)
}

func TestCampaignService_DialingMode_Preview(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, _ := svc.Create(ctx, CreateCampaignInput{
		TenantID: 1, Name: "Preview", DialingMode: DialingModePreview, SkillGroupID: 10,
	})

	assert.Equal(t, DialingModePreview, c.DialingMode)
	assert.Equal(t, 30, c.PreviewTimeoutSec)
}

func TestCampaignService_DialingMode_Progressive(t *testing.T) {
	svc := newTestService()

	c, _ := svc.Create(context.Background(), CreateCampaignInput{
		TenantID: 1, Name: "Progressive", DialingMode: DialingModeProgressive, SkillGroupID: 10,
	})

	assert.Equal(t, DialingModeProgressive, c.DialingMode)
	assert.Equal(t, 1.0, c.RatioMultiplier)
}

func TestCampaignService_DialingMode_Power(t *testing.T) {
	svc := newTestService()

	c, _ := svc.Create(context.Background(), CreateCampaignInput{
		TenantID: 1, Name: "Power", DialingMode: DialingModePower, SkillGroupID: 10,
	})

	assert.Equal(t, DialingModePower, c.DialingMode)
	assert.Equal(t, 3.0, c.RatioMultiplier)
}

func TestCampaignService_RetryLogic(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, _ := svc.Create(ctx, CreateCampaignInput{
		TenantID: 1, Name: "Retry", DialingMode: DialingModePredictive, SkillGroupID: 10,
	})
	_ = svc.ImportCases(ctx, c.ID, []CaseInput{{PhoneNumber: "+86138"}})
	_, _ = svc.Start(ctx, c.ID)

	cs, _ := svc.GetNextCase(ctx, c.ID)
	require.NotNil(t, cs)

	updated, err := svc.MarkCaseFailed(ctx, cs.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, updated.AttemptCount)
	assert.Equal(t, CaseStatusPending, updated.Status)
	assert.NotNil(t, updated.NextAttemptAt)
}
