package telephony

import (
	"context"
	"testing"
	"time"

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

// Prevent unused import for time in tests
var _ = time.Now
