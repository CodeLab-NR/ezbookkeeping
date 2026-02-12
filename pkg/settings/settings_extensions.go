package settings

// New database type constants for SQL Server and Azure SQL DB
const (
	SqlServerDbType string = "sqlserver"
	AzureSqlDbType  string = "azuresqldb"
)

// Azure SQL DB authentication method types
const (
	AzureAuthPassword        string = "password"
	AzureAuthServicePrincipal string = "service_principal"
)

// ExtendedDatabaseConfig represents extended database configuration for SQL Server and Azure SQL DB
type ExtendedDatabaseConfig struct {
	// Azure SQL DB specific fields
	AzureAuthMethod   string // "password" or "service_principal"
	AzureTenantID     string
	AzureClientID     string
	AzureClientSecret string

	// Connection pooling settings
	AzureMaxIdleConns   int
	AzureMaxOpenConns   int
	AzureConnMaxLifetime int
}

// Update the existing DatabaseConfig struct in setting.go to include:
// The following fields should be added to the DatabaseConfig struct:
// 
// Extended configuration for Azure SQL DB and SQL Server
// AzureAuthMethod   string // "password" (default) or "service_principal"
// AzureTenantID     string
// AzureClientID     string
// AzureClientSecret string
// AzureMaxIdleConns   int // overrides MaxIdleConnection if set
// AzureMaxOpenConns   int // overrides MaxOpenConnection if set
// AzureConnMaxLifetime int // overrides ConnectionMaxLifeTime if set
