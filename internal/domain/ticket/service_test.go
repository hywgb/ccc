package ticket

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

func newTestServices() (*TicketService, *TicketTemplateService) {
	categories := NewMockCategoryRepo()
	templates := NewMockTemplateRepo()
	tickets := NewMockTicketRepo()
	comments := NewMockCommentRepo()
	return NewTicketService(tickets, templates, comments),
		NewTicketTemplateService(templates, categories)
}

func TestTicketTemplateService_Publish_Online(t *testing.T) {
	_, tmplSvc := newTestServices()
	ctx := context.Background()

	tmpl, err := tmplSvc.Create(ctx, CreateTemplateInput{
		TenantID: 1,
		Name:     "投诉模板",
		Fields:   `[{"name":"reason","type":"select","options":["质量","服务","其他"]}]`,
	})
	require.NoError(t, err)
	assert.Equal(t, "draft", tmpl.OnlineStatus)

	published, err := tmplSvc.Publish(ctx, tmpl.ID)
	require.NoError(t, err)
	assert.Equal(t, "published", published.OnlineStatus)
}

func TestTicketTemplateService_Publish_AlreadyPublished(t *testing.T) {
	_, tmplSvc := newTestServices()
	ctx := context.Background()

	tmpl, _ := tmplSvc.Create(ctx, CreateTemplateInput{
		TenantID: 1, Name: "T1",
	})
	_, _ = tmplSvc.Publish(ctx, tmpl.ID)

	_, err := tmplSvc.Publish(ctx, tmpl.ID)
	assert.ErrorIs(t, err, ErrAlreadyPublished)
}

func TestTicketTemplateService_Offline(t *testing.T) {
	_, tmplSvc := newTestServices()
	ctx := context.Background()

	tmpl, _ := tmplSvc.Create(ctx, CreateTemplateInput{
		TenantID: 1, Name: "T1",
	})
	_, _ = tmplSvc.Publish(ctx, tmpl.ID)

	offline, err := tmplSvc.Offline(ctx, tmpl.ID)
	require.NoError(t, err)
	assert.Equal(t, "offline", offline.OnlineStatus)
}

func TestTicketTemplateService_FlowGraph_Validation(t *testing.T) {
	_, tmplSvc := newTestServices()
	ctx := context.Background()

	tmpl, _ := tmplSvc.Create(ctx, CreateTemplateInput{
		TenantID:  1,
		Name:      "有流程的模板",
		FlowGraph: `{"nodes":[{"id":"start"},{"id":"end"}],"edges":[{"from":"start","to":"end"}]}`,
	})
	require.NotNil(t, tmpl)
	assert.NotEmpty(t, tmpl.FlowGraph)

	// Invalid JSON should fail
	_, err := tmplSvc.Create(ctx, CreateTemplateInput{
		TenantID:  1,
		Name:      "坏流程",
		FlowGraph: `{invalid json}`,
	})
	assert.ErrorIs(t, err, ErrInvalidFlowGraph)
}

func TestTicketService_Create_FromTemplate(t *testing.T) {
	ticketSvc, tmplSvc := newTestServices()
	ctx := context.Background()

	tmpl, _ := tmplSvc.Create(ctx, CreateTemplateInput{
		TenantID: 1, Name: "T1",
	})
	_, _ = tmplSvc.Publish(ctx, tmpl.ID)

	tk, err := ticketSvc.Create(ctx, CreateTicketInput{
		TenantID:    1,
		TemplateID:  &tmpl.ID,
		Title:       "客户投诉001",
		Description: "产品质量问题",
		Priority:    "high",
	})
	require.NoError(t, err)
	assert.Equal(t, TicketStatusOpen, tk.Status)
	assert.Equal(t, "high", tk.Priority)
	assert.Equal(t, &tmpl.ID, tk.TemplateID)
}

func TestTicketService_Create_InvalidPriority(t *testing.T) {
	ticketSvc, _ := newTestServices()

	_, err := ticketSvc.Create(context.Background(), CreateTicketInput{
		TenantID: 1, Title: "T", Priority: "critical",
	})
	assert.ErrorIs(t, err, ErrInvalidPriority)
}

func TestTicketService_Create_UnpublishedTemplate_Error(t *testing.T) {
	ticketSvc, tmplSvc := newTestServices()
	ctx := context.Background()

	tmpl, _ := tmplSvc.Create(ctx, CreateTemplateInput{
		TenantID: 1, Name: "Draft Template",
	})

	_, err := ticketSvc.Create(ctx, CreateTicketInput{
		TenantID:   1,
		TemplateID: &tmpl.ID,
		Title:      "Test",
		Priority:   "low",
	})
	assert.ErrorIs(t, err, ErrTemplateNotPublished)
}

func TestTicketService_Assign_Success(t *testing.T) {
	ticketSvc, _ := newTestServices()
	ctx := context.Background()

	tk, _ := ticketSvc.Create(ctx, CreateTicketInput{
		TenantID: 1, Title: "需分配", Priority: "medium",
	})

	agentID := int64(100)
	assigned, err := ticketSvc.Assign(ctx, tk.ID, agentID)
	require.NoError(t, err)
	assert.Equal(t, &agentID, assigned.AssigneeID)
}

func TestTicketService_Transition_OpenToInProgress(t *testing.T) {
	ticketSvc, _ := newTestServices()
	ctx := context.Background()

	tk, _ := ticketSvc.Create(ctx, CreateTicketInput{
		TenantID: 1, Title: "T", Priority: "low",
	})

	updated, err := ticketSvc.Transition(ctx, tk.ID, TicketStatusInProgress)
	require.NoError(t, err)
	assert.Equal(t, TicketStatusInProgress, updated.Status)
}

func TestTicketService_Transition_InvalidState_Error(t *testing.T) {
	ticketSvc, _ := newTestServices()
	ctx := context.Background()

	tk, _ := ticketSvc.Create(ctx, CreateTicketInput{
		TenantID: 1, Title: "T", Priority: "low",
	})

	// Cannot go directly from open to closed
	_, err := ticketSvc.Transition(ctx, tk.ID, TicketStatusClosed)
	assert.ErrorIs(t, err, ErrInvalidTransition)
}

func TestTicketService_Transition_ResolveSetTimestamp(t *testing.T) {
	ticketSvc, _ := newTestServices()
	ctx := context.Background()

	tk, _ := ticketSvc.Create(ctx, CreateTicketInput{
		TenantID: 1, Title: "T", Priority: "low",
	})
	_, _ = ticketSvc.Transition(ctx, tk.ID, TicketStatusInProgress)

	resolved, err := ticketSvc.Transition(ctx, tk.ID, TicketStatusResolved)
	require.NoError(t, err)
	assert.Equal(t, TicketStatusResolved, resolved.Status)
	assert.NotNil(t, resolved.ResolvedAt)
}

func TestTicketService_AddComment(t *testing.T) {
	ticketSvc, _ := newTestServices()
	ctx := context.Background()

	tk, _ := ticketSvc.Create(ctx, CreateTicketInput{
		TenantID: 1, Title: "评论测试", Priority: "low",
	})

	err := ticketSvc.AddComment(ctx, tk.ID, 100, "已联系客户")
	require.NoError(t, err)

	comments, err := ticketSvc.ListComments(ctx, tk.ID)
	require.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, "已联系客户", comments[0].Content)
}
