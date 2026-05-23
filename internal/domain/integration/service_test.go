package integration

import (
	"context"
	"testing"
	"time"

	"github.com/divord97/ccc/pkg/snowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	_ = snowflake.Init(1)
}

func TestDNCService_CheckDNC_Blocked(t *testing.T) {
	repo := NewMockDNCRepo()
	svc := NewDNCService(repo)
	ctx := context.Background()

	_ = svc.AddEntry(ctx, &DNCEntry{TenantID: 1, Number: "+8613800138000", Reason: "customer request", Source: "manual"})

	err := svc.CheckDNC(ctx, 1, "+8613800138000")
	assert.ErrorIs(t, err, ErrDNCBlocked)
}

func TestDNCService_CheckDNC_NotBlocked(t *testing.T) {
	repo := NewMockDNCRepo()
	svc := NewDNCService(repo)
	ctx := context.Background()

	err := svc.CheckDNC(ctx, 1, "+8613800138000")
	assert.NoError(t, err)
}

func TestDNCService_CheckDNC_Expired(t *testing.T) {
	repo := NewMockDNCRepo()
	svc := NewDNCService(repo)
	ctx := context.Background()

	past := time.Now().Add(-24 * time.Hour)
	_ = svc.AddEntry(ctx, &DNCEntry{TenantID: 1, Number: "+8613800138000", Reason: "temp", Source: "api", ExpiresAt: &past})

	err := svc.CheckDNC(ctx, 1, "+8613800138000")
	assert.NoError(t, err, "expired DNC should not block")
}

func TestDNCService_CheckBatch(t *testing.T) {
	repo := NewMockDNCRepo()
	svc := NewDNCService(repo)
	ctx := context.Background()

	_ = svc.AddEntry(ctx, &DNCEntry{TenantID: 1, Number: "+86A", Source: "manual"})
	_ = svc.AddEntry(ctx, &DNCEntry{TenantID: 1, Number: "+86C", Source: "manual"})

	blocked, err := svc.CheckBatch(ctx, 1, []string{"+86A", "+86B", "+86C"})
	require.NoError(t, err)
	assert.Len(t, blocked, 2)
	assert.Contains(t, blocked, "+86A")
	assert.Contains(t, blocked, "+86C")
}

func TestCallTagService_AssignTag(t *testing.T) {
	repo := NewMockCallTagRepo()
	svc := NewCallTagService(repo)
	ctx := context.Background()

	err := svc.AssignTag(ctx, &CallTagAssignment{TenantID: 1, CallID: 100, TagID: 1, TagName: "VIP"})
	require.NoError(t, err)

	tags, err := svc.GetCallTags(ctx, 100)
	require.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, "VIP", tags[0].TagName)
}

func TestCallTagService_AssignTag_Duplicate(t *testing.T) {
	repo := NewMockCallTagRepo()
	svc := NewCallTagService(repo)
	ctx := context.Background()

	_ = svc.AssignTag(ctx, &CallTagAssignment{TenantID: 1, CallID: 100, TagID: 1, TagName: "VIP"})
	err := svc.AssignTag(ctx, &CallTagAssignment{TenantID: 1, CallID: 100, TagID: 1, TagName: "VIP"})
	assert.ErrorIs(t, err, ErrTagAlreadyExists)
}

func TestCallTagService_RemoveTag(t *testing.T) {
	repo := NewMockCallTagRepo()
	svc := NewCallTagService(repo)
	ctx := context.Background()

	_ = svc.AssignTag(ctx, &CallTagAssignment{TenantID: 1, CallID: 100, TagID: 1, TagName: "VIP"})
	tags, _ := svc.GetCallTags(ctx, 100)
	require.Len(t, tags, 1)

	_ = svc.RemoveTag(ctx, tags[0].ID)
	tags, _ = svc.GetCallTags(ctx, 100)
	assert.Len(t, tags, 0)
}
