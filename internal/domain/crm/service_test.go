package crm

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

func newTestService() *CustomerService {
	return NewCustomerService(
		NewMockCustomerRepo(),
		NewMockCustomerPhoneRepo(),
		NewMockInteractionRepo(),
		NewMockCustomFieldRepo(),
	)
}

func TestCustomerService_Create_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, err := svc.Create(ctx, CreateCustomerInput{
		TenantID: 1,
		Name:     "张三",
		Email:    "zhangsan@example.com",
		Level:    "normal",
		Phones: []PhoneInput{
			{PhoneType: "mobile", Number: "+8613800001111", IsPrimary: true},
		},
	})

	require.NoError(t, err)
	assert.NotZero(t, c.ID)
	assert.Equal(t, "张三", c.Name)
	assert.Equal(t, "normal", c.Level)
}

func TestCustomerService_Create_WithMultiPhones(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, err := svc.Create(ctx, CreateCustomerInput{
		TenantID: 1,
		Name:     "李四",
		Level:    "vip",
		Phones: []PhoneInput{
			{PhoneType: "mobile", Number: "+8613800001111", IsPrimary: true},
			{PhoneType: "landline", Number: "021-12345678", IsPrimary: false},
			{PhoneType: "backup", Number: "+8613900009999", IsPrimary: false},
		},
	})

	require.NoError(t, err)
	phones, err := svc.phones.ListByCustomer(ctx, c.ID)
	require.NoError(t, err)
	assert.Len(t, phones, 3)

	var primaryCount int
	for _, p := range phones {
		if p.IsPrimary {
			primaryCount++
		}
	}
	assert.Equal(t, 1, primaryCount)
}

func TestCustomerService_Create_NoPrimary_Error(t *testing.T) {
	svc := newTestService()

	_, err := svc.Create(context.Background(), CreateCustomerInput{
		TenantID: 1,
		Name:     "王五",
		Level:    "normal",
		Phones: []PhoneInput{
			{PhoneType: "mobile", Number: "+8613800001111", IsPrimary: false},
		},
	})
	assert.ErrorIs(t, err, ErrNoPrimaryPhone)
}

func TestCustomerService_Create_MultiplePrimary_Error(t *testing.T) {
	svc := newTestService()

	_, err := svc.Create(context.Background(), CreateCustomerInput{
		TenantID: 1,
		Name:     "赵六",
		Level:    "normal",
		Phones: []PhoneInput{
			{PhoneType: "mobile", Number: "+8613800001111", IsPrimary: true},
			{PhoneType: "landline", Number: "021-12345678", IsPrimary: true},
		},
	})
	assert.ErrorIs(t, err, ErrMultiplePrimary)
}

func TestCustomerService_Create_InvalidLevel_Error(t *testing.T) {
	svc := newTestService()

	_, err := svc.Create(context.Background(), CreateCustomerInput{
		TenantID: 1,
		Name:     "Test",
		Level:    "diamond",
		Phones: []PhoneInput{
			{PhoneType: "mobile", Number: "+8613800001111", IsPrimary: true},
		},
	})
	assert.ErrorIs(t, err, ErrInvalidLevel)
}

func TestCustomerService_Create_InvalidPhoneType_Error(t *testing.T) {
	svc := newTestService()

	_, err := svc.Create(context.Background(), CreateCustomerInput{
		TenantID: 1,
		Name:     "Test",
		Level:    "normal",
		Phones: []PhoneInput{
			{PhoneType: "fax", Number: "+8613800001111", IsPrimary: true},
		},
	})
	assert.ErrorIs(t, err, ErrInvalidPhoneType)
}

func TestCustomerService_FindByPhone_AnyPhoneMatch(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateCustomerInput{
		TenantID: 1,
		Name:     "多号码客户",
		Level:    "normal",
		Phones: []PhoneInput{
			{PhoneType: "mobile", Number: "+8613800001111", IsPrimary: true},
			{PhoneType: "landline", Number: "021-99998888", IsPrimary: false},
		},
	})
	require.NoError(t, err)

	// Find by primary phone
	found, err := svc.FindByPhone(ctx, 1, "+8613800001111")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.ID, found.ID)

	// Find by secondary phone
	found2, err := svc.FindByPhone(ctx, 1, "021-99998888")
	require.NoError(t, err)
	require.NotNil(t, found2)
	assert.Equal(t, created.ID, found2.ID)

	// Not found
	found3, err := svc.FindByPhone(ctx, 1, "+8613800009999")
	require.NoError(t, err)
	assert.Nil(t, found3)
}

func TestCustomerService_RecordInteraction(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	c, _ := svc.Create(ctx, CreateCustomerInput{
		TenantID: 1, Name: "客户A", Level: "normal",
		Phones: []PhoneInput{{PhoneType: "mobile", Number: "+8613800001111", IsPrimary: true}},
	})

	err := svc.RecordInteraction(ctx, RecordInteractionInput{
		CustomerID: c.ID,
		TenantID:   1,
		Channel:    "call",
		Direction:  "inbound",
		Summary:    "首次咨询产品",
		AgentName:  "坐席小王",
	})
	require.NoError(t, err)

	interactions, err := svc.ListInteractions(ctx, c.ID, 0, 10)
	require.NoError(t, err)
	assert.Len(t, interactions, 1)
	assert.Equal(t, "call", interactions[0].Channel)
	assert.Equal(t, "首次咨询产品", interactions[0].Summary)
}

func TestCustomerService_CustomFieldDefinition(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	err := svc.CreateFieldDefinition(ctx, CustomFieldDefinition{
		TenantID:   1,
		EntityType: "customer",
		FieldName:  "industry",
		FieldType:  "select",
		Options:    `["IT","Finance","Healthcare"]`,
		IsRequired: true,
		SortOrder:  1,
	})
	require.NoError(t, err)

	err = svc.CreateFieldDefinition(ctx, CustomFieldDefinition{
		TenantID:   1,
		EntityType: "customer",
		FieldName:  "contract_date",
		FieldType:  "date",
		SortOrder:  2,
	})
	require.NoError(t, err)

	fields, err := svc.ListFieldDefinitions(ctx, 1, "customer")
	require.NoError(t, err)
	assert.Len(t, fields, 2)
}

func TestCustomerService_CustomFieldDefinition_InvalidType(t *testing.T) {
	svc := newTestService()

	err := svc.CreateFieldDefinition(context.Background(), CustomFieldDefinition{
		TenantID:   1,
		EntityType: "customer",
		FieldName:  "test",
		FieldType:  "blob",
	})
	assert.ErrorIs(t, err, ErrInvalidFieldType)
}

func TestCustomerService_CustomFieldDefinition_InvalidEntity(t *testing.T) {
	svc := newTestService()

	err := svc.CreateFieldDefinition(context.Background(), CustomFieldDefinition{
		TenantID:   1,
		EntityType: "order",
		FieldName:  "test",
		FieldType:  "text",
	})
	assert.ErrorIs(t, err, ErrInvalidEntityType)
}

func TestCustomerService_BatchImport_CSV(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	records := []CreateCustomerInput{
		{TenantID: 1, Name: "导入客户A", Level: "normal", Phones: []PhoneInput{{PhoneType: "mobile", Number: "+8613800001111", IsPrimary: true}}},
		{TenantID: 1, Name: "导入客户B", Level: "vip", Phones: []PhoneInput{{PhoneType: "mobile", Number: "+8613800002222", IsPrimary: true}}},
		{TenantID: 1, Name: "导入客户C", Level: "svip", Phones: []PhoneInput{{PhoneType: "mobile", Number: "+8613800003333", IsPrimary: true}}},
	}

	result, err := svc.BatchImport(ctx, records)
	require.NoError(t, err)
	assert.Equal(t, 3, result.Success)
	assert.Equal(t, 0, result.Failed)

	// Verify all imported
	customers, err := svc.customers.List(ctx, 1, 0, 100)
	require.NoError(t, err)
	assert.Len(t, customers, 3)
}

func TestCustomerService_BatchImport_PartialFailure(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	records := []CreateCustomerInput{
		{TenantID: 1, Name: "Good", Level: "normal", Phones: []PhoneInput{{PhoneType: "mobile", Number: "+8613800001111", IsPrimary: true}}},
		{TenantID: 1, Name: "Bad", Level: "diamond", Phones: []PhoneInput{{PhoneType: "mobile", Number: "+8613800002222", IsPrimary: true}}}, // invalid level
		{TenantID: 1, Name: "Good2", Level: "vip", Phones: []PhoneInput{{PhoneType: "mobile", Number: "+8613800003333", IsPrimary: true}}},
	}

	result, err := svc.BatchImport(ctx, records)
	require.NoError(t, err)
	assert.Equal(t, 2, result.Success)
	assert.Equal(t, 1, result.Failed)
	assert.Len(t, result.Errors, 1)
}
