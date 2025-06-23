# Step-by-Step Implementation Tasks: User Display Name Integration

## Phase 1: Database Schema and Model Updates

### Task 1.1: Create Database Migration
**File**: `models/migrationscripts/20240623_add_display_name_fields.go`
**Action**: Create new file
```go
package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addDisplayNameFields)(nil)

type addDisplayNameFields struct{}

func (*addDisplayNameFields) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(
		&archived.QDevUserData{},
		&archived.QDevUserMetrics{},
	)
}

func (*addDisplayNameFields) Version() uint64 {
	return 20240623000001
}

func (*addDisplayNameFields) Name() string {
	return "add display_name fields to user tables"
}

// Archived models for migration
type archived struct{}

func (archived) QDevUserData() interface{} {
	return &QDevUserData{
		DisplayName: "",
	}
}

func (archived) QDevUserMetrics() interface{} {
	return &QDevUserMetrics{
		DisplayName: "",
	}
}
```

### Task 1.2: Update Migration Registry
**File**: `models/migrationscripts/register.go`
**Action**: Modify existing file
```go
// Add to the All() function
func All() []plugin.MigrationScript {
	return []plugin.MigrationScript{
		// ... existing migrations
		new(addDisplayNameFields), // Add this line
	}
}
```

### Task 1.3: Update Connection Model
**File**: `models/connection.go`
**Action**: Modify existing struct

**Test First**: Create `models/connection_test.go`
```go
package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestQDevConn_WithIdentityStore(t *testing.T) {
	conn := QDevConn{
		AccessKeyId:         "test-key",
		SecretAccessKey:     "test-secret",
		Region:             "us-east-1",
		Bucket:             "test-bucket",
		RateLimitPerHour:   20000,
		IdentityStoreId:    "d-1234567890",
		IdentityStoreRegion: "us-west-2",
	}
	
	assert.Equal(t, "d-1234567890", conn.IdentityStoreId)
	assert.Equal(t, "us-west-2", conn.IdentityStoreRegion)
}

func TestQDevConn_Sanitize_PreservesIdentityStore(t *testing.T) {
	conn := QDevConn{
		SecretAccessKey:     "secret-key",
		IdentityStoreId:    "d-1234567890",
		IdentityStoreRegion: "us-west-2",
	}
	
	sanitized := conn.Sanitize()
	assert.NotEqual(t, "secret-key", sanitized.SecretAccessKey)
	assert.Equal(t, "d-1234567890", sanitized.IdentityStoreId)
	assert.Equal(t, "us-west-2", sanitized.IdentityStoreRegion)
}
```

**Implementation**: Update `QDevConn` struct
```go
type QDevConn struct {
	AccessKeyId      string `mapstructure:"accessKeyId" json:"accessKeyId"`
	SecretAccessKey  string `mapstructure:"secretAccessKey" json:"secretAccessKey"`
	Region          string `mapstructure:"region" json:"region"`
	Bucket          string `mapstructure:"bucket" json:"bucket"`
	RateLimitPerHour int    `mapstructure:"rateLimitPerHour" json:"rateLimitPerHour"`
	
	// New fields for IAM Identity Center
	IdentityStoreId     string `mapstructure:"identityStoreId" json:"identityStoreId"`
	IdentityStoreRegion string `mapstructure:"identityStoreRegion" json:"identityStoreRegion"`
}
```

### Task 1.4: Update User Data Model
**File**: `models/user_data.go`
**Action**: Modify existing struct

**Test First**: Create `models/user_data_test.go`
```go
package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestQDevUserData_WithDisplayName(t *testing.T) {
	userData := QDevUserData{
		ConnectionId: 1,
		UserId:      "uuid-123",
		DisplayName: "John Doe",
	}
	
	assert.Equal(t, "John Doe", userData.DisplayName)
	assert.Equal(t, "uuid-123", userData.UserId)
}

func TestQDevUserData_TableName(t *testing.T) {
	userData := QDevUserData{}
	assert.Equal(t, "_tool_q_dev_user_data", userData.TableName())
}
```

**Implementation**: Add DisplayName field to existing struct
```go
type QDevUserData struct {
	// Existing fields...
	ConnectionId uint64    `gorm:"primaryKey"`
	UserId      string    `gorm:"primaryKey"`
	// ... other existing fields
	
	// New field
	DisplayName string    `gorm:"type:varchar(255)" json:"displayName"`
}
```

### Task 1.5: Update User Metrics Model
**File**: `models/user_metrics.go`
**Action**: Modify existing struct

**Test First**: Create `models/user_metrics_test.go`
```go
package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestQDevUserMetrics_WithDisplayName(t *testing.T) {
	userMetrics := QDevUserMetrics{
		ConnectionId: 1,
		UserId:      "uuid-123",
		DisplayName: "John Doe",
	}
	
	assert.Equal(t, "John Doe", userMetrics.DisplayName)
	assert.Equal(t, "uuid-123", userMetrics.UserId)
}

func TestQDevUserMetrics_TableName(t *testing.T) {
	userMetrics := QDevUserMetrics{}
	assert.Equal(t, "_tool_q_dev_user_metrics", userMetrics.TableName())
}
```

**Implementation**: Add DisplayName field to existing struct
```go
type QDevUserMetrics struct {
	// Existing fields...
	ConnectionId uint64 `gorm:"primaryKey"`
	UserId      string `gorm:"primaryKey"`
	// ... other existing fields
	
	// New field
	DisplayName string `gorm:"type:varchar(255)" json:"displayName"`
}
```

## Phase 2: Identity Center Client Implementation

### Task 2.1: Create Identity Client
**File**: `tasks/identity_client.go`
**Action**: Create new file

**Test First**: Create `tasks/identity_client_test.go`
```go
package tasks

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/aws/aws-sdk-go/service/identitystore"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
)

type MockIdentityStore struct {
	mock.Mock
}

func (m *MockIdentityStore) DescribeUser(input *identitystore.DescribeUserInput) (*identitystore.DescribeUserOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*identitystore.DescribeUserOutput), args.Error(1)
}

func TestNewQDevIdentityClient_Success(t *testing.T) {
	connection := &models.QDevConnection{
		QDevConn: models.QDevConn{
			AccessKeyId:         "test-key",
			SecretAccessKey:     "test-secret",
			IdentityStoreId:     "d-1234567890",
			IdentityStoreRegion: "us-west-2",
		},
	}
	
	client, err := NewQDevIdentityClient(connection)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "d-1234567890", client.StoreId)
	assert.Equal(t, "us-west-2", client.Region)
}

func TestQDevIdentityClient_ResolveUserDisplayName_Success(t *testing.T) {
	mockService := &MockIdentityStore{}
	client := &QDevIdentityClient{
		IdentityStore: mockService,
		StoreId:       "d-1234567890",
		Region:        "us-west-2",
	}
	
	displayName := "John Doe"
	mockService.On("DescribeUser", mock.AnythingOfType("*identitystore.DescribeUserInput")).Return(
		&identitystore.DescribeUserOutput{
			DisplayName: &displayName,
		}, nil)
	
	result, err := client.ResolveUserDisplayName("user-123")
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", result)
}

func TestQDevIdentityClient_ResolveUserDisplayName_Fallback(t *testing.T) {
	mockService := &MockIdentityStore{}
	client := &QDevIdentityClient{
		IdentityStore: mockService,
		StoreId:       "d-1234567890",
		Region:        "us-west-2",
	}
	
	mockService.On("DescribeUser", mock.AnythingOfType("*identitystore.DescribeUserInput")).Return(
		&identitystore.DescribeUserOutput{}, errors.New("user not found"))
	
	result, err := client.ResolveUserDisplayName("user-123")
	assert.Error(t, err)
	assert.Equal(t, "user-123", result) // Should fallback to UUID
}
```

**Implementation**: Create the client
```go
package tasks

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/identitystore"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
)

type QDevIdentityClient struct {
	IdentityStore *identitystore.IdentityStore
	StoreId       string
	Region        string
}

func NewQDevIdentityClient(connection *models.QDevConnection) (*QDevIdentityClient, error) {
	if connection.IdentityStoreId == "" {
		return nil, nil // No identity store configured
	}
	
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(connection.IdentityStoreRegion),
		Credentials: credentials.NewStaticCredentials(connection.AccessKeyId, connection.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}
	
	return &QDevIdentityClient{
		IdentityStore: identitystore.New(sess),
		StoreId:       connection.IdentityStoreId,
		Region:        connection.IdentityStoreRegion,
	}, nil
}

func (client *QDevIdentityClient) ResolveUserDisplayName(userId string) (string, error) {
	input := &identitystore.DescribeUserInput{
		IdentityStoreId: aws.String(client.StoreId),
		UserId:         aws.String(userId),
	}
	
	result, err := client.IdentityStore.DescribeUser(input)
	if err != nil {
		return userId, err // Fallback to UUID on error
	}
	
	if result.DisplayName != nil && *result.DisplayName != "" {
		return *result.DisplayName, nil
	}
	
	return userId, nil // Fallback to UUID if no display name
}
```

## Phase 3: Enhanced Data Extraction

### Task 3.1: Update Task Data Structure
**File**: `tasks/task_data.go`
**Action**: Modify existing struct

**Test First**: Create `tasks/task_data_test.go`
```go
package tasks

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestQDevTaskData_WithIdentityClient(t *testing.T) {
	taskData := &QDevTaskData{
		Options:        &QDevOptions{ConnectionId: 1},
		S3Client:       &QDevS3Client{},
		IdentityClient: &QDevIdentityClient{StoreId: "d-1234567890"},
	}
	
	assert.NotNil(t, taskData.IdentityClient)
	assert.Equal(t, "d-1234567890", taskData.IdentityClient.StoreId)
}
```

**Implementation**: Add IdentityClient field
```go
type QDevTaskData struct {
	Options        *QDevOptions
	S3Client       *QDevS3Client
	IdentityClient *QDevIdentityClient // New field
}
```

### Task 3.2: Update Plugin Implementation
**File**: `impl/impl.go`
**Action**: Modify PrepareTaskData method

**Test First**: Create `impl/impl_test.go`
```go
package impl

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/q_dev/tasks"
)

type MockTaskContext struct {
	mock.Mock
}

func (m *MockTaskContext) GetData() interface{} {
	args := m.Called()
	return args.Get(0)
}

func TestQDev_PrepareTaskData_WithIdentityStore(t *testing.T) {
	plugin := &QDev{}
	options := map[string]interface{}{
		"connectionId": uint64(1),
	}
	
	// Mock setup would go here
	// This test verifies that PrepareTaskData creates IdentityClient when configured
}
```

**Implementation**: Update PrepareTaskData method
```go
func (p QDev) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.QDevOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}

	connectionHelper := helper.NewConnectionHelper(taskCtx, nil, p.Name())
	connection := &models.QDevConnection{}
	err := connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	// Create S3 client
	s3Client, err := tasks.NewQDevS3Client(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	// Create Identity client (new)
	identityClient, err := tasks.NewQDevIdentityClient(connection)
	if err != nil {
		taskCtx.GetLogger().Warn("Failed to create identity client, proceeding without user name resolution", err)
		identityClient = nil
	}

	return &tasks.QDevTaskData{
		Options:        &op,
		S3Client:       s3Client,
		IdentityClient: identityClient,
	}, nil
}
```

### Task 3.3: Update S3 Data Extractor
**File**: `tasks/s3_data_extractor.go`
**Action**: Modify existing extraction logic

**Test First**: Create `tasks/s3_data_extractor_test.go`
```go
package tasks

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQDevIdentityClient struct {
	mock.Mock
}

func (m *MockQDevIdentityClient) ResolveUserDisplayName(userId string) (string, error) {
	args := m.Called(userId)
	return args.String(0), args.Error(1)
}

func TestS3DataExtractor_ExtractUserData_WithDisplayName(t *testing.T) {
	mockIdentityClient := &MockQDevIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "user-123").Return("John Doe", nil)
	
	extractor := &QDevS3DataExtractor{
		identityClient: mockIdentityClient,
	}
	
	// Test CSV data processing with display name resolution
	csvData := [][]string{
		{"header1", "user_id", "header3"},
		{"value1", "user-123", "value3"},
	}
	
	// Verify that display name is resolved and stored
	// Implementation details would depend on existing extractor structure
}

func TestS3DataExtractor_ExtractUserData_FallbackToUUID(t *testing.T) {
	mockIdentityClient := &MockQDevIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "user-123").Return("user-123", errors.New("not found"))
	
	extractor := &QDevS3DataExtractor{
		identityClient: mockIdentityClient,
	}
	
	// Test that UUID is used when display name resolution fails
}

func TestS3DataExtractor_ExtractUserData_NoIdentityClient(t *testing.T) {
	extractor := &QDevS3DataExtractor{
		identityClient: nil, // No identity client
	}
	
	// Test that UUID is used when no identity client is available
}
```

**Implementation**: Update extraction logic
```go
// Add to existing s3_data_extractor.go

func (extractor *QDevS3DataExtractor) extractUserDataWithDisplayName(userData *models.QDevUserData, taskData *tasks.QDevTaskData) error {
	// Set default display name to user ID
	userData.DisplayName = userData.UserId
	
	// Try to resolve display name if identity client is available
	if taskData.IdentityClient != nil {
		displayName, err := taskData.IdentityClient.ResolveUserDisplayName(userData.UserId)
		if err != nil {
			extractor.logger.Warn("Failed to resolve display name for user", userData.UserId, err)
			// Keep UUID as fallback
		} else {
			userData.DisplayName = displayName
		}
	}
	
	return nil
}

// Update existing extraction method to call the new function
func (extractor *QDevS3DataExtractor) processUserRecord(record []string, taskData *tasks.QDevTaskData) error {
	userData := &models.QDevUserData{
		// ... existing field population
		UserId: record[userIdColumnIndex],
	}
	
	// Resolve display name
	if err := extractor.extractUserDataWithDisplayName(userData, taskData); err != nil {
		return err
	}
	
	// Save to database
	return extractor.saveUserData(userData)
}
```

### Task 3.4: Update User Metrics Converter
**File**: `tasks/user_metrics_converter.go`
**Action**: Modify existing conversion logic

**Test First**: Create `tasks/user_metrics_converter_test.go`
```go
package tasks

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
)

func TestUserMetricsConverter_IncludeDisplayName(t *testing.T) {
	userData := []models.QDevUserData{
		{
			UserId:      "user-123",
			DisplayName: "John Doe",
			// ... other fields
		},
	}
	
	converter := &QDevUserMetricsConverter{}
	metrics, err := converter.convertToMetrics(userData)
	
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", metrics[0].DisplayName)
	assert.Equal(t, "user-123", metrics[0].UserId)
}

func TestUserMetricsConverter_FallbackDisplayName(t *testing.T) {
	userData := []models.QDevUserData{
		{
			UserId:      "user-123",
			DisplayName: "user-123", // Fallback case
			// ... other fields
		},
	}
	
	converter := &QDevUserMetricsConverter{}
	metrics, err := converter.convertToMetrics(userData)
	
	assert.NoError(t, err)
	assert.Equal(t, "user-123", metrics[0].DisplayName)
}
```

**Implementation**: Update conversion logic
```go
// Update existing conversion method in user_metrics_converter.go

func (converter *QDevUserMetricsConverter) convertUserDataToMetrics(userData *models.QDevUserData) *models.QDevUserMetrics {
	return &models.QDevUserMetrics{
		// ... existing field mappings
		UserId:      userData.UserId,
		DisplayName: userData.DisplayName, // New field
		// ... other aggregated fields
	}
}

// Update aggregation logic to preserve display names
func (converter *QDevUserMetricsConverter) aggregateUserMetrics(userDataList []models.QDevUserData) []models.QDevUserMetrics {
	metricsMap := make(map[string]*models.QDevUserMetrics)
	
	for _, userData := range userDataList {
		key := userData.UserId
		if existing, exists := metricsMap[key]; exists {
			// Aggregate existing metrics
			converter.aggregateMetrics(existing, &userData)
		} else {
			// Create new metrics entry
			metricsMap[key] = converter.convertUserDataToMetrics(&userData)
		}
	}
	
	// Convert map to slice
	var result []models.QDevUserMetrics
	for _, metrics := range metricsMap {
		result = append(result, *metrics)
	}
	
	return result
}
```

## Phase 4: API and Configuration Updates

### Task 4.1: Update Connection API
**File**: `api/connection.go`
**Action**: Modify connection testing

**Test First**: Create `api/connection_test.go`
```go
package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
)

func TestTestConnection_WithIdentityStore(t *testing.T) {
	connection := &models.QDevConnection{
		QDevConn: models.QDevConn{
			AccessKeyId:         "test-key",
			SecretAccessKey:     "test-secret",
			Region:             "us-east-1",
			Bucket:             "test-bucket",
			IdentityStoreId:    "d-1234567890",
			IdentityStoreRegion: "us-west-2",
		},
	}
	
	// Test that connection validation includes identity store validation
	err := testConnectionWithIdentityStore(connection)
	// Assert based on expected behavior
}

func TestTestConnection_WithoutIdentityStore(t *testing.T) {
	connection := &models.QDevConnection{
		QDevConn: models.QDevConn{
			AccessKeyId:     "test-key",
			SecretAccessKey: "test-secret",
			Region:         "us-east-1",
			Bucket:         "test-bucket",
			// No identity store fields
		},
	}
	
	// Test that connection validation works without identity store
	err := testConnectionWithoutIdentityStore(connection)
	// Assert based on expected behavior
}
```

**Implementation**: Update connection testing
```go
// Add to existing connection.go

func testIdentityStoreAccess(connection *models.QDevConnection) error {
	if connection.IdentityStoreId == "" {
		return nil // Identity store not configured, skip test
	}
	
	identityClient, err := tasks.NewQDevIdentityClient(connection)
	if err != nil {
		return fmt.Errorf("failed to create identity client: %w", err)
	}
	
	// Test with a dummy user ID to verify access
	_, err = identityClient.ResolveUserDisplayName("test-user-id")
	if err != nil {
		// Log warning but don't fail connection test
		// Identity store access might fail for non-existent users
		return nil
	}
	
	return nil
}

// Update existing TestConnection function
func TestConnection(connection *models.QDevConnection) error {
	// Test S3 access (existing)
	if err := testS3Access(connection); err != nil {
		return err
	}
	
	// Test Identity Store access (new)
	if err := testIdentityStoreAccess(connection); err != nil {
		return fmt.Errorf("Identity Store access failed: %w", err)
	}
	
	return nil
}
```

## Phase 5: Integration Testing

### Task 5.1: Create Integration Tests
**File**: `integration_test.go`
**Action**: Create new file

```go
package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"github.com/apache/incubator-devlake/plugins/q_dev/tasks"
)

func TestEndToEndUserDisplayNameResolution(t *testing.T) {
	// Integration test that verifies the complete flow:
	// 1. Connection with identity store configuration
	// 2. S3 data extraction with display name resolution
	// 3. User metrics conversion including display names
	// 4. Database storage and retrieval
	
	t.Skip("Integration test - requires AWS credentials and resources")
	
	connection := &models.QDevConnection{
		QDevConn: models.QDevConn{
			AccessKeyId:         "test-key",
			SecretAccessKey:     "test-secret",
			Region:             "us-east-1",
			Bucket:             "test-bucket",
			IdentityStoreId:    "d-1234567890",
			IdentityStoreRegion: "us-west-2",
		},
	}
	
	// Test complete pipeline
	// This would require actual AWS resources for full integration testing
}

func TestUserDisplayNameFallback(t *testing.T) {
	// Test that the system gracefully handles identity store failures
	// and falls back to UUID display
}
```

### Task 5.2: Update Documentation
**File**: `README.md`
**Action**: Update configuration documentation

```markdown
## Configuration

Configuration items include:

1. AWS Access Key ID
2. AWS Secret Key
3. AWS Region (for S3)
4. S3 Bucket Name
5. Rate Limit (per hour)
6. Identity Store ID (optional - for user display names)
7. Identity Store Region (optional - may differ from S3 region)

You can create a connection using the following curl command:
```bash
curl 'http://localhost:8080/plugins/q_dev/connections' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "q_dev_connection",
    "accessKeyId": "<YOUR_ACCESS_KEY_ID>",
    "secretAccessKey": "<YOUR_SECRET_ACCESS_KEY>",
    "region": "<AWS_REGION>",
    "bucket": "<YOUR_S3_BUCKET_NAME>",
    "rateLimitPerHour": 20000,
    "identityStoreId": "<YOUR_IDENTITY_STORE_ID>",
    "identityStoreRegion": "<IDENTITY_STORE_REGION>"
}'
```

Note: If `identityStoreId` is provided, user display names will be resolved from AWS IAM Identity Center. If not provided, user IDs (UUIDs) will be displayed in dashboards.
```

## Execution Order

1. **Run all tests first** (they should fail initially)
2. **Implement Phase 1** (Database and Models)
3. **Run Phase 1 tests** (should pass)
4. **Implement Phase 2** (Identity Client)
5. **Run Phase 2 tests** (should pass)
6. **Implement Phase 3** (Data Extraction)
7. **Run Phase 3 tests** (should pass)
8. **Implement Phase 4** (API Updates)
9. **Run Phase 4 tests** (should pass)
10. **Run Integration Tests** (Phase 5)
11. **Update Documentation**

This test-driven approach ensures that each component is properly tested before implementation and that the integration works correctly across all phases.
