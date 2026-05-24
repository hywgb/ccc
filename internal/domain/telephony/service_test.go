package telephony

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutingService_MatchRule_CLIPrefix(t *testing.T) {
	repo := NewMockRoutingRuleRepo()
	svc := NewRoutingService(repo)
	ctx := context.Background()

	_ = repo.Create(ctx, &RoutingRule{ID: 1, TenantID: 1, Name: "China Mobile", MatchType: "prefix", MatchValue: "+86", SIPTrunkID: 100, Priority: 10, IsActive: true})
	_ = repo.Create(ctx, &RoutingRule{ID: 2, TenantID: 1, Name: "US", MatchType: "prefix", MatchValue: "+1", SIPTrunkID: 200, Priority: 10, IsActive: true})

	rule, err := svc.MatchRule(ctx, 1, "+8613800138000")
	require.NoError(t, err)
	assert.Equal(t, int64(100), rule.SIPTrunkID)

	rule, err = svc.MatchRule(ctx, 1, "+12025551234")
	require.NoError(t, err)
	assert.Equal(t, int64(200), rule.SIPTrunkID)
}

func TestRoutingService_MatchRule_Priority(t *testing.T) {
	repo := NewMockRoutingRuleRepo()
	svc := NewRoutingService(repo)
	ctx := context.Background()

	_ = repo.Create(ctx, &RoutingRule{ID: 1, TenantID: 1, MatchType: "prefix", MatchValue: "+86", SIPTrunkID: 100, Priority: 20, IsActive: true})
	_ = repo.Create(ctx, &RoutingRule{ID: 2, TenantID: 1, MatchType: "prefix", MatchValue: "+86", SIPTrunkID: 200, Priority: 5, IsActive: true})

	rule, err := svc.MatchRule(ctx, 1, "+8613800138000")
	require.NoError(t, err)
	assert.Equal(t, int64(200), rule.SIPTrunkID, "should pick lower priority number (higher priority)")
}

func TestRoutingService_MatchRule_TimeOfDay(t *testing.T) {
	repo := NewMockRoutingRuleRepo()
	svc := NewRoutingService(repo)
	ctx := context.Background()

	_ = repo.Create(ctx, &RoutingRule{ID: 1, TenantID: 1, MatchType: "time_of_day", MatchValue: "00:00-23:59", SIPTrunkID: 300, Priority: 10, IsActive: true})

	rule, err := svc.MatchRule(ctx, 1, "+8613800138000")
	require.NoError(t, err)
	assert.Equal(t, int64(300), rule.SIPTrunkID)
}

func TestRoutingService_MatchRule_NoMatch(t *testing.T) {
	repo := NewMockRoutingRuleRepo()
	svc := NewRoutingService(repo)
	ctx := context.Background()

	_ = repo.Create(ctx, &RoutingRule{ID: 1, TenantID: 1, MatchType: "prefix", MatchValue: "+44", SIPTrunkID: 100, Priority: 10, IsActive: true})

	_, err := svc.MatchRule(ctx, 1, "+8613800138000")
	assert.ErrorIs(t, err, ErrNoMatchingRoute)
}

func TestRoutingService_MatchRule_InactiveSkipped(t *testing.T) {
	repo := NewMockRoutingRuleRepo()
	svc := NewRoutingService(repo)
	ctx := context.Background()

	_ = repo.Create(ctx, &RoutingRule{ID: 1, TenantID: 1, MatchType: "prefix", MatchValue: "+86", SIPTrunkID: 100, Priority: 10, IsActive: false})

	_, err := svc.MatchRule(ctx, 1, "+8613800138000")
	assert.ErrorIs(t, err, ErrNoMatchingRoute)
}

func TestCLIPolicyService_SelectCLI_Fixed(t *testing.T) {
	policyRepo := NewMockCLIPolicyRepo()
	phoneRepo := NewMockPhoneNumberRepo()
	svc := NewCLIPolicyService(policyRepo, phoneRepo)
	ctx := context.Background()

	_ = phoneRepo.Create(ctx, &PhoneNumber{ID: 10, TenantID: 1, Number: "+861380001"})
	fixedID := int64(10)
	_ = policyRepo.Create(ctx, &CLIPolicy{ID: 1, TenantID: 1, Strategy: CLIStrategyFixed, FixedNumberID: &fixedID, IsDefault: true})

	pn, err := svc.SelectCLI(ctx, 1, nil, "+8613900139000")
	require.NoError(t, err)
	assert.Equal(t, "+861380001", pn.Number)
}

func TestCLIPolicyService_SelectCLI_RoundRobin(t *testing.T) {
	policyRepo := NewMockCLIPolicyRepo()
	phoneRepo := NewMockPhoneNumberRepo()
	svc := NewCLIPolicyService(policyRepo, phoneRepo)
	ctx := context.Background()

	_ = phoneRepo.Create(ctx, &PhoneNumber{ID: 10, TenantID: 1, Number: "+86A"})
	_ = phoneRepo.Create(ctx, &PhoneNumber{ID: 20, TenantID: 1, Number: "+86B"})
	policyID := int64(1)
	_ = policyRepo.Create(ctx, &CLIPolicy{ID: 1, TenantID: 1, Strategy: CLIStrategyRoundRobin, NumberPoolIDs: "10,20", IsDefault: true})

	pn1, err := svc.SelectCLI(ctx, 1, &policyID, "+8613900139000")
	require.NoError(t, err)

	pn2, err := svc.SelectCLI(ctx, 1, &policyID, "+8613900139000")
	require.NoError(t, err)

	// Round robin should alternate between numbers
	assert.NotEqual(t, pn1.Number, pn2.Number, "round robin should rotate")
}

func TestCLIPolicyService_SelectCLI_NoPolicy(t *testing.T) {
	policyRepo := NewMockCLIPolicyRepo()
	phoneRepo := NewMockPhoneNumberRepo()
	svc := NewCLIPolicyService(policyRepo, phoneRepo)
	ctx := context.Background()

	_, err := svc.SelectCLI(ctx, 1, nil, "+8613900139000")
	assert.ErrorIs(t, err, ErrNoCLIPolicy)
}

func TestCLIPolicyService_SelectCLI_MatchArea(t *testing.T) {
	policyRepo := NewMockCLIPolicyRepo()
	phoneRepo := NewMockPhoneNumberRepo()
	svc := NewCLIPolicyService(policyRepo, phoneRepo)
	ctx := context.Background()

	_ = phoneRepo.Create(ctx, &PhoneNumber{ID: 10, TenantID: 1, Number: "+021XXXX"})
	_ = phoneRepo.Create(ctx, &PhoneNumber{ID: 20, TenantID: 1, Number: "+010XXXX"})
	policyID := int64(1)
	_ = policyRepo.Create(ctx, &CLIPolicy{ID: 1, TenantID: 1, Strategy: CLIStrategyMatchArea, NumberPoolIDs: "10,20"})

	pn, err := svc.SelectCLI(ctx, 1, &policyID, "+010123456")
	require.NoError(t, err)
	assert.Equal(t, "+010XXXX", pn.Number, "should match area code")
}

func TestTrunkHealthCheck_OPTIONSKeepalive(t *testing.T) {
	trunkRepo := NewMockPhoneNumberRepo() // just need a SIPTrunkRepository
	_ = trunkRepo
	svc := NewTrunkHealthService(&mockSIPTrunkRepo{}, NewMockSIPTrunkGroupRepo())

	// Initially healthy
	h := svc.GetHealthStatus(100)
	assert.Equal(t, TrunkStatusActive, h.Status)

	// 1 failure — still active
	svc.RecordHealthCheck(100, false)
	h = svc.GetHealthStatus(100)
	assert.Equal(t, TrunkStatusActive, h.Status)
	assert.Equal(t, 1, h.FailCount)

	// 2 failures — still active
	svc.RecordHealthCheck(100, false)
	h = svc.GetHealthStatus(100)
	assert.Equal(t, TrunkStatusActive, h.Status)

	// 3 failures — down
	svc.RecordHealthCheck(100, false)
	h = svc.GetHealthStatus(100)
	assert.Equal(t, TrunkStatusDown, h.Status)

	// Success — recovered
	svc.RecordHealthCheck(100, true)
	h = svc.GetHealthStatus(100)
	assert.Equal(t, TrunkStatusActive, h.Status)
	assert.Equal(t, 0, h.FailCount)
}

func TestTrunkFailover_AutoSwitch(t *testing.T) {
	trunkRepo := &mockSIPTrunkRepo{
		trunks: map[int64]*SIPTrunk{
			1: {ID: 1, Name: "Primary", Status: "active"},
			2: {ID: 2, Name: "Backup", Status: "active"},
		},
	}
	groupRepo := NewMockSIPTrunkGroupRepo()
	ctx := context.Background()

	_ = groupRepo.Create(ctx, &SIPTrunkGroup{ID: 10, TenantID: 1, Name: "HA Group", Strategy: "priority"})
	_ = groupRepo.AddMember(ctx, &SIPTrunkGroupMember{ID: 1, GroupID: 10, SIPTrunkID: 1, Priority: 1})
	_ = groupRepo.AddMember(ctx, &SIPTrunkGroupMember{ID: 2, GroupID: 10, SIPTrunkID: 2, Priority: 2})

	svc := NewTrunkHealthService(trunkRepo, groupRepo)

	// Both healthy — should pick primary (trunk 1)
	trunk, err := svc.SelectHealthyTrunk(ctx, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), trunk.ID)

	// Mark primary as down
	svc.RecordHealthCheck(1, false)
	svc.RecordHealthCheck(1, false)
	svc.RecordHealthCheck(1, false)

	// Should failover to backup (trunk 2)
	trunk, err = svc.SelectHealthyTrunk(ctx, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), trunk.ID)

	// Mark both as down
	svc.RecordHealthCheck(2, false)
	svc.RecordHealthCheck(2, false)
	svc.RecordHealthCheck(2, false)

	_, err = svc.SelectHealthyTrunk(ctx, 10)
	assert.ErrorIs(t, err, ErrNoHealthyTrunk)
}

// mockSIPTrunkRepo implements SIPTrunkRepository for tests.
type mockSIPTrunkRepo struct {
	trunks map[int64]*SIPTrunk
}

func (r *mockSIPTrunkRepo) Create(_ context.Context, t *SIPTrunk) error {
	if r.trunks == nil {
		r.trunks = make(map[int64]*SIPTrunk)
	}
	r.trunks[t.ID] = t
	return nil
}

func (r *mockSIPTrunkRepo) GetByID(_ context.Context, id int64) (*SIPTrunk, error) {
	if r.trunks == nil {
		return nil, nil
	}
	return r.trunks[id], nil
}

func (r *mockSIPTrunkRepo) Update(_ context.Context, t *SIPTrunk) error {
	r.trunks[t.ID] = t
	return nil
}

func (r *mockSIPTrunkRepo) List(_ context.Context, _ int64, _, _ int) ([]*SIPTrunk, int64, error) {
	return nil, 0, nil
}

func (r *mockSIPTrunkRepo) ListAll(_ context.Context, _, _ int) ([]*SIPTrunk, int64, error) {
	return nil, 0, nil
}
