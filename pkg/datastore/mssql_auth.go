package datastore

import (
	"fmt"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// BuildMSSQLConnectionString builds a connection string for SQL Server or Azure SQL DB
// Supports both basic authentication and Azure AD Service Principal authentication
func BuildMSSQLConnectionString(dbConfig *settings.DatabaseConfig) (string, error) {
	if dbConfig == nil {
		return "", errs.ErrDatabaseIsNull
	}

	// Check if this is Azure SQL DB with Service Principal authentication
	if dbConfig.DatabaseType == settings.AzureSqlDbType &&
		dbConfig.AzureAuthMethod == settings.AzureAuthServicePrincipal {
		return buildAzureServicePrincipalConnectionString(dbConfig)
	}

	// For SQL Server and Azure SQL DB with basic auth
	return buildBasicMSSQLConnectionString(dbConfig)
}

// buildBasicMSSQLConnectionString builds connection string for SQL Server or Azure SQL DB with username/password
func buildBasicMSSQLConnectionString(dbConfig *settings.DatabaseConfig) (string, error) {
	if dbConfig.DatabaseHost == "" {
		return "", errs.ErrDatabaseHostInvalid
	}

	if dbConfig.DatabaseUser == "" {
		return "", fmt.Errorf("database user is required")
	}

	// For Azure SQL DB, append server name to username if not already present
	username := dbConfig.DatabaseUser
	if dbConfig.DatabaseType == settings.AzureSqlDbType && !strings.Contains(username, "@") {
		// Extract server name from host (e.g., "servername.database.windows.net" -> "servername")
		serverName := extractServerName(dbConfig.DatabaseHost)
		if serverName != "" {
			username = fmt.Sprintf("%s@%s", username, serverName)
		}
	}

	// Build connection string: server=<host>;user id=<user>;password=<pass>;database=<db>
	connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
		dbConfig.DatabaseHost,
		username,
		dbConfig.DatabasePassword,
		dbConfig.DatabaseName,
	)

	// Optional: Add port if specified in host (for SQL Server)
	// go-mssqldb will handle the default port 1433 if not specified

	return connStr, nil
}

// buildAzureServicePrincipalConnectionString builds connection string for Azure SQL DB with Service Principal
func buildAzureServicePrincipalConnectionString(dbConfig *settings.DatabaseConfig) (string, error) {
	if dbConfig.DatabaseHost == "" {
		return "", errs.ErrDatabaseHostInvalid
	}

	if dbConfig.AzureClientID == "" {
		return "", fmt.Errorf("azure_client_id is required for service principal authentication")
	}

	if dbConfig.AzureClientSecret == "" {
		return "", fmt.Errorf("azure_client_secret is required for service principal authentication")
	}

	// Service Principal authentication connection string for Azure SQL DB
	// Format: server=<host>;fedauth=ActiveDirectoryServicePrincipal;User ID=<client-id>;Password=<client-secret>;database=<db>
	connStr := fmt.Sprintf("server=%s;fedauth=ActiveDirectoryServicePrincipal;User ID=%s;Password=%s;database=%s",
		dbConfig.DatabaseHost,
		dbConfig.AzureClientID,
		dbConfig.AzureClientSecret,
		dbConfig.DatabaseName,
	)

	return connStr, nil
}

// extractServerName extracts server name from full host address
// Examples:
// "servername.database.windows.net" -> "servername"
// "localhost:1433" -> "localhost"
// "127.0.0.1" -> "127.0.0.1"
func extractServerName(host string) string {
	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// For Azure SQL DB, remove domain suffix
	if idx := strings.Index(host, ".database.windows.net"); idx != -1 {
		host = host[:idx]
	}

	return host
}

// ValidateMSSQLConfig validates the MSSQL configuration
func ValidateMSSQLConfig(dbConfig *settings.DatabaseConfig) error {
	if dbConfig == nil {
		return errs.ErrDatabaseIsNull
	}

	if dbConfig.DatabaseHost == "" {
		return errs.ErrDatabaseHostInvalid
	}

	if dbConfig.DatabaseName == "" {
		return fmt.Errorf("database name is required")
	}

	// For Service Principal auth, validate required fields
	if dbConfig.AzureAuthMethod == settings.AzureAuthServicePrincipal {
		if dbConfig.AzureClientID == "" {
			return fmt.Errorf("azure_client_id is required for service principal authentication")
		}
		if dbConfig.AzureClientSecret == "" {
			return fmt.Errorf("azure_client_secret is required for service principal authentication")
		}
		if dbConfig.AzureTenantID == "" {
			return fmt.Errorf("azure_tenant_id is required for service principal authentication")
		}
	} else {
		// For basic auth, validate user credentials
		if dbConfig.DatabaseUser == "" {
			return fmt.Errorf("database user is required for basic authentication")
		}
	}

	return nil
}
