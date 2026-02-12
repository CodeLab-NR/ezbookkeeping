package datastore

import (
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// getMssqlConnectionString returns the connection string for MSSQL Server or Azure SQL DB
// Supports both basic authentication and Azure AD Service Principal authentication
func getMssqlConnectionString(dbConfig *settings.DatabaseConfig) (string, error) {
	// Validate configuration
	if err := ValidateMSSQLConfig(dbConfig); err != nil {
		return "", err
	}

	// Build appropriate connection string based on auth method
	return BuildMSSQLConnectionString(dbConfig)
}

// GetMssqlDriverName returns the driver name for MSSQL
func GetMssqlDriverName() string {
	return "mssql"
}

// ConfigureMssqlConnectionPool configures connection pool settings for MSSQL
// Returns the configured max idle connections, max open connections, and connection lifetime
func ConfigureMssqlConnectionPool(dbConfig *settings.DatabaseConfig) (int, int, int) {
	maxIdleConns := int(dbConfig.MaxIdleConnection)
	maxOpenConns := int(dbConfig.MaxOpenConnection)
	connMaxLifetime := int(dbConfig.ConnectionMaxLifeTime)

	// Azure-specific connection pooling configuration
	if dbConfig.DatabaseType == settings.AzureSqlDbType {
		// Override with Azure-specific settings if configured
		if dbConfig.AzureMaxIdleConns > 0 {
			maxIdleConns = dbConfig.AzureMaxIdleConns
		}
		if dbConfig.AzureMaxOpenConns > 0 {
			maxOpenConns = dbConfig.AzureMaxOpenConns
		}
		if dbConfig.AzureConnMaxLifetime > 0 {
			connMaxLifetime = dbConfig.AzureConnMaxLifetime
		}

		// Set reasonable defaults for Azure if not specified
		if maxIdleConns == 0 {
			maxIdleConns = 10
		}
		if maxOpenConns == 0 {
			maxOpenConns = 100
		}
		if connMaxLifetime == 0 {
			connMaxLifetime = 3600 // 1 hour
		}
	}

	return maxIdleConns, maxOpenConns, connMaxLifetime
}

// Example usage in datastore_container.go initializeDatabase function:
// ============================================================
// Add this to the initializeDatabase function after the existing type checks:
//
// } else if dbConfig.DatabaseType == settings.SqlServerDbType {
//    connStr, err = getMssqlConnectionString(dbConfig)
// } else if dbConfig.DatabaseType == settings.AzureSqlDbType {
//    connStr, err = getMssqlConnectionString(dbConfig)
//
// Then update the xorm.NewEngineGroup call to handle MSSQL driver:
//
// driverName := dbConfig.DatabaseType
// if dbConfig.DatabaseType == settings.SqlServerDbType || dbConfig.DatabaseType == settings.AzureSqlDbType {
//    driverName = GetMssqlDriverName()
// }
// engineGroup, err := xorm.NewEngineGroup(driverName, connStrs, xorm.RoundRobinPolicy())
//
// And update connection pool configuration:
//
// if dbConfig.DatabaseType == settings.SqlServerDbType || dbConfig.DatabaseType == settings.AzureSqlDbType {
//    maxIdle, maxOpen, connLifetime := ConfigureMssqlConnectionPool(dbConfig)
//    engineGroup.SetMaxIdleConns(maxIdle)
//    engineGroup.SetMaxOpenConns(maxOpen)
//    engineGroup.SetConnMaxLifetime(time.Duration(connLifetime) * time.Second)
// } else {
//    engineGroup.SetMaxIdleConns(int(dbConfig.MaxIdleConnection))
//    engineGroup.SetMaxOpenConns(int(dbConfig.MaxOpenConnection))
//    engineGroup.SetConnMaxLifetime(time.Duration(dbConfig.ConnectionMaxLifeTime) * time.Second)
// }
