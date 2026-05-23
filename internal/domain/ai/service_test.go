package ai

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

func newTestKnowledgeService() *KnowledgeService {
	return NewKnowledgeService(
		NewMockKnowledgeCategoryRepo(),
		NewMockKnowledgeArticleRepo(),
	)
}

func newTestScriptService() *AgentScriptService {
	return NewAgentScriptService(NewMockAgentScriptRepo())
}

func TestKnowledgeService_CreateArticle_Published(t *testing.T) {
	svc := newTestKnowledgeService()
	ctx := context.Background()

	a, err := svc.CreateArticle(ctx, CreateArticleInput{
		TenantID: 1,
		Title:    "如何重置密码",
		Content:  "步骤1: 点击忘记密码...",
		Tags:     "密码,重置,常见问题",
		Status:   "published",
	})
	require.NoError(t, err)
	assert.NotZero(t, a.ID)
	assert.Equal(t, "published", a.Status)
}

func TestKnowledgeService_Search_MatchTitle(t *testing.T) {
	svc := newTestKnowledgeService()
	ctx := context.Background()

	_, _ = svc.CreateArticle(ctx, CreateArticleInput{
		TenantID: 1, Title: "如何退款", Content: "联系客服申请退款", Tags: "退款", Status: "published",
	})
	_, _ = svc.CreateArticle(ctx, CreateArticleInput{
		TenantID: 1, Title: "配送时间", Content: "一般3-5个工作日", Tags: "配送", Status: "published",
	})
	_, _ = svc.CreateArticle(ctx, CreateArticleInput{
		TenantID: 1, Title: "退款政策", Content: "30天无理由退款", Tags: "退款,政策", Status: "published",
	})

	results, err := svc.Search(ctx, 1, "退款", 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 2)
}

func TestKnowledgeService_Search_DraftNotReturned(t *testing.T) {
	svc := newTestKnowledgeService()
	ctx := context.Background()

	_, _ = svc.CreateArticle(ctx, CreateArticleInput{
		TenantID: 1, Title: "草稿文章", Content: "内容", Status: "draft",
	})

	results, err := svc.Search(ctx, 1, "草稿", 10)
	require.NoError(t, err)
	assert.Len(t, results, 0)
}

func TestKnowledgeService_CreateCategory(t *testing.T) {
	svc := newTestKnowledgeService()
	ctx := context.Background()

	cat, err := svc.CreateCategory(ctx, 1, "常见问题")
	require.NoError(t, err)
	assert.NotZero(t, cat.ID)

	cats, err := svc.ListCategories(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, cats, 1)
}

func TestAgentScriptService_Create_ValidJSON(t *testing.T) {
	svc := newTestScriptService()
	ctx := context.Background()

	s, err := svc.Create(ctx, CreateScriptInput{
		TenantID: 1,
		Name:     "欢迎话术",
		Content:  `{"steps":[{"text":"您好，请问有什么可以帮您？"},{"text":"好的，我来为您查看"}]}`,
		IsActive: true,
	})
	require.NoError(t, err)
	assert.NotZero(t, s.ID)
	assert.True(t, s.IsActive)
}

func TestAgentScriptService_Create_InvalidJSON_Error(t *testing.T) {
	svc := newTestScriptService()

	_, err := svc.Create(context.Background(), CreateScriptInput{
		TenantID: 1, Name: "Bad", Content: `{invalid}`,
	})
	assert.ErrorIs(t, err, ErrInvalidScript)
}

func TestAgentScriptService_Update(t *testing.T) {
	svc := newTestScriptService()
	ctx := context.Background()

	s, _ := svc.Create(ctx, CreateScriptInput{
		TenantID: 1, Name: "V1",
		Content: `{"steps":[]}`, IsActive: true,
	})

	s.Name = "V2"
	s.IsActive = false
	err := svc.Update(ctx, s)
	require.NoError(t, err)

	updated, _ := svc.GetByID(ctx, s.ID)
	assert.Equal(t, "V2", updated.Name)
	assert.False(t, updated.IsActive)
}
