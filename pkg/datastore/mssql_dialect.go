package datastore

import (
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// SetSavePointMssql sets a save point in the current transaction for MSSQL
// MSSQL uses SAVE TRANSACTION instead of SAVEPOINT
func SetSavePointMssql(sess interface{ Exec(sql string, args ...interface{}) (interface{}, error) }, savePointName string) error {
	// MSSQL syntax: SAVE TRANSACTION savePointName
	_, err := sess.Exec("SAVE TRANSACTION [" + savePointName + "]")
	return err
}

// RollbackToSavePointMssql rolls back to the specified save point in the current transaction for MSSQL
// MSSQL uses ROLLBACK TRANSACTION instead of ROLLBACK TO SAVEPOINT
func RollbackToSavePointMssql(sess interface{ Exec(sql string, args ...interface{}) (interface{}, error) }, savePointName string) error {
	// MSSQL syntax: ROLLBACK TRANSACTION savePointName
	_, err := sess.Exec("ROLLBACK TRANSACTION [" + savePointName + "]")
	return err
}

// GetDatabaseTypeDriver returns the driver name for a given database type
func GetDatabaseTypeDriver(dbType string) string {
	switch dbType {
	case settings.SqlServerDbType, settings.AzureSqlDbType:
		return "mssql"
	case settings.MySqlDbType:
		return settings.MySqlDbType
	case settings.PostgresDbType:
		return settings.PostgresDbType
	case settings.Sqlite3DbType:
		return settings.Sqlite3DbType
	default:
		return dbType
	}
}

// IsMssqlDatabase returns true if the database type is SQL Server or Azure SQL DB
func IsMssqlDatabase(dbType string) bool {
	return dbType == settings.SqlServerDbType || dbType == settings.AzureSqlDbType
}

// SupportsTransactionSavepoints returns true if the database type supports transaction savepoints
// PostgreSQL and MSSQL support SAVEPOINT syntax
func SupportsTransactionSavepoints(dbType string) bool {
	return dbType == settings.PostgresDbType || IsMssqlDatabase(dbType)
}

// GetDateTimeFormat returns the SQL datetime format string for the given database type
// Useful for constructing datetime values in SQL queries
func GetDateTimeFormat(dbType string) string {
	switch dbType {
	case settings.PostgresDbType:
		return "2006-01-02 15:04:05"
	case settings.MySqlDbType:
		return "2006-01-02 15:04:05"
	case settings.Sqlite3DbType:
		return "2006-01-02 15:04:05"
	case settings.SqlServerDbType, settings.AzureSqlDbType:
		return "2006-01-02T15:04:05"
	default:
		return "2006-01-02 15:04:05"
	}
}

// SQL Syntax Differences Reference:
// ================================
//
// SAVEPOINT/TRANSACTION HANDLING:
// - PostgreSQL: SAVEPOINT name; ROLLBACK TO SAVEPOINT name;
// - MSSQL/Azure: SAVE TRANSACTION name; ROLLBACK TRANSACTION name;
// - MySQL: SAVEPOINT name; ROLLBACK TO SAVEPOINT name;
// - SQLite: SAVEPOINT name; ROLLBACK TO SAVEPOINT name;
//
// AUTO INCREMENT:
// - PostgreSQL: SERIAL or BIGSERIAL
// - MSSQL/Azure: IDENTITY(1, 1)
// - MySQL: AUTO_INCREMENT
// - SQLite: AUTOINCREMENT
//
// DATETIME FUNCTIONS:
// - PostgreSQL: CURRENT_TIMESTAMP, NOW()
// - MSSQL/Azure: GETDATE(), CURRENT_TIMESTAMP
// - MySQL: CURRENT_TIMESTAMP, NOW()
// - SQLite: datetime('now'), CURRENT_TIMESTAMP
//
// STRING FUNCTIONS:
// - PostgreSQL: CONCAT(), ||
// - MSSQL/Azure: +, CONCAT()
// - MySQL: CONCAT()
// - SQLite: ||
//
// LIMIT OFFSET:
// - PostgreSQL: LIMIT n OFFSET m
// - MSSQL/Azure: OFFSET m ROWS FETCH NEXT n ROWS ONLY
// - MySQL: LIMIT m, n (offset m, fetch n)
// - SQLite: LIMIT n OFFSET m
//
// TABLE NAME ESCAPING:
// - PostgreSQL: "table_name"
// - MSSQL/Azure: [table_name] or "table_name"
// - MySQL: `table_name` or "table_name"
// - SQLite: [table_name] or "table_name"
