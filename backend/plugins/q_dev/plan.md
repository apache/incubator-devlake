# Feature Design: User Display Name Integration

## Overview
Enhance the Q Developer plugin to display user-friendly names instead of UUIDs in dashboards by integrating with AWS IAM Identity Center to fetch and store user display names.

## Current State
- User data from S3 contains UUID-based user identifiers
- Dashboards show cryptic UUIDs instead of readable user names
- No mapping between UUIDs and human-readable names

## Proposed Solution

### 1. AWS IAM Identity Center Integration
**New AWS Service Integration:**
- Add AWS Identity Store SDK support alongside existing S3 integration
- Implement `identitystore.DescribeUser` API calls to fetch user details
- Extend connection configuration to include Identity Store ID and region

### 2. Database Schema Changes
**Modified Table: `_tool_q_dev_user_data`**
```sql
-- Add new column to existing user data table
ALTER TABLE _tool_q_dev_user_data 
ADD COLUMN display_name VARCHAR(255) NULL;
```

**Modified Table: `_tool_q_dev_user_metrics`**
```sql
-- Add new column to existing user metrics table
ALTER TABLE _tool_q_dev_user_metrics 
ADD COLUMN display_name VARCHAR(255) NULL;
```

### 3. Configuration Updates
**Extended Connection Model:**
```go
type QDevConn struct {
    // Existing fields...
    AccessKeyId      string `json:"accessKeyId"`
    SecretAccessKey  string `json:"secretAccessKey"`
    Region          string `json:"region"`           // S3 region
    Bucket          string `json:"bucket"`
    RateLimitPerHour int   `json:"rateLimitPerHour"`
    
    // New fields for IAM Identity Center
    IdentityStoreId     string `json:"identityStoreId"`     // Required for user name resolution
    IdentityStoreRegion string `json:"identityStoreRegion"` // IAM IDC region (may differ from S3)
}
```

### 4. Data Pipeline Enhancement
**Modified 3-Stage Process:**
1. **Collection**: `collectQDevS3Files` - Unchanged
2. **Extraction**: `extractQDevS3Data` - Enhanced to resolve user display names during data insertion
3. **Conversion**: `convertQDevUserMetrics` - Enhanced to include display names in aggregated metrics

### 5. Implementation Components

**New Files:**
- `tasks/identity_client.go` - AWS Identity Store client wrapper

**Modified Files:**
- `models/connection.go` - Add IdentityStoreId and IdentityStoreRegion fields
- `models/user_data.go` - Add DisplayName field
- `models/user_metrics.go` - Add DisplayName field
- `impl/impl.go` - Update table registration for schema changes
- `tasks/s3_data_extractor.go` - Integrate user name resolution during data extraction
- `tasks/user_metrics_converter.go` - Include display names in metrics aggregation

### 6. Model Updates
**Enhanced User Data Model:**
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

**Enhanced User Metrics Model:**
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

### 7. API Integration Details
**AWS Identity Store Operations:**
```go
// New client in tasks/identity_client.go
type QDevIdentityClient struct {
    IdentityStore *identitystore.IdentityStore
    StoreId       string
    Region        string
}

func NewQDevIdentityClient(connection *models.QDevConnection) (*QDevIdentityClient, error) {
    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String(connection.IdentityStoreRegion), // Separate region
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

// Core operation with caching
func (client *QDevIdentityClient) ResolveUserDisplayName(userId string) (string, error) {
    input := &identitystore.DescribeUserInput{
        IdentityStoreId: aws.String(client.StoreId),
        UserId:         aws.String(userId),
    }
    
    result, err := client.IdentityStore.DescribeUser(input)
    if err != nil {
        return userId, err // Fallback to UUID on error
    }
    
    if result.DisplayName != nil {
        return *result.DisplayName, nil
    }
    
    return userId, nil // Fallback to UUID if no display name
}
```

### 8. Enhanced Data Extraction Logic
**Modified S3 Data Extractor:**
```go
// In tasks/s3_data_extractor.go
func (extractor *QDevS3DataExtractor) extractUserData(csvData [][]string) error {
    // Create identity client for user name resolution
    identityClient, err := NewQDevIdentityClient(extractor.connection)
    if err != nil {
        // Log warning but continue without name resolution
        extractor.logger.Warn("Failed to create identity client, using UUIDs", err)
    }
    
    for _, row := range csvData {
        userData := &models.QDevUserData{
            // ... populate existing fields
            UserId: row[userIdColumn],
        }
        
        // Resolve display name if identity client is available
        if identityClient != nil {
            displayName, err := identityClient.ResolveUserDisplayName(userData.UserId)
            if err != nil {
                extractor.logger.Warn("Failed to resolve display name for user", userData.UserId, err)
                userData.DisplayName = userData.UserId // Fallback to UUID
            } else {
                userData.DisplayName = displayName
            }
        } else {
            userData.DisplayName = userData.UserId // Fallback to UUID
        }
        
        // Save to database
        // ...
    }
}
```

### 9. Error Handling & Resilience
- **Graceful Degradation**: If Identity Center is unavailable, store UUID as display name
- **Regional Flexibility**: Support different regions for S3 and IAM Identity Center
- **Rate Limiting**: Respect AWS API limits for Identity Store operations
- **Retry Logic**: Handle transient failures with exponential backoff
- **Caching**: Implement in-memory cache during extraction to avoid duplicate API calls

### 10. Configuration Validation
**Enhanced Connection Testing:**
```go
func TestConnection(connection *models.QDevConnection) error {
    // Test S3 access (existing)
    if err := testS3Access(connection); err != nil {
        return err
    }
    
    // Test Identity Store access (new)
    if connection.IdentityStoreId != "" {
        if err := testIdentityStoreAccess(connection); err != nil {
            return fmt.Errorf("Identity Store access failed: %w", err)
        }
    }
    
    return nil
}
```

### 11. Migration Strategy
**Database Migration:**
```go
// In models/migrationscripts/
func (*addDisplayNameFields) Up(basicRes context.BasicRes) errors.Error {
    return basicRes.GetDal().AutoMigrate(
        &models.QDevUserData{},
        &models.QDevUserMetrics{},
    )
}
```

## Benefits
1. **Improved UX**: Dashboards show meaningful user names instead of UUIDs
2. **Regional Flexibility**: S3 and IAM Identity Center can be in different regions
3. **Simplified Schema**: No additional tables, just enhanced existing ones
4. **Better Analytics**: Easier to identify usage patterns by actual users
5. **Backward Compatibility**: Existing data remains functional

## Implementation Priority
1. **Phase 1**: Database schema updates and model changes
2. **Phase 2**: Identity Center client and connection configuration
3. **Phase 3**: Enhanced data extraction with name resolution
4. **Phase 4**: Dashboard updates and testing

This approach maintains data consistency by storing display names directly with user data while providing the flexibility to handle different AWS regions for S3 and IAM Identity Center services.
