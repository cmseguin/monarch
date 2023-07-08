package utils

import (
	"database/sql"

	"github.com/cmseguin/khata"
	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/spf13/cobra"
	"modernc.org/sqlite"
)

func IsDriverSupported(driver string, supportedDriverList []string) bool {
	var supported bool = false

	for _, supportedDriver := range supportedDriverList {
		if driver == supportedDriver {
			supported = true
		}
	}

	return supported
}

func ConnectToDatabase(driver string, connection string) (*sql.DB, *khata.Khata) {
	db, err := sql.Open(driver, connection)

	if err != nil {
		return nil, khata.Wrap(err).SetExitCode(1).Explain("Could not connect to database")
	}

	err = db.Ping()

	if err != nil {
		return nil, khata.Wrap(err).SetExitCode(1).Explain("Could not ping database")
	}

	return db, nil
}

func CreateMigrationTable(db *sql.DB) *khata.Khata {
	var err error

	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
		_, err = db.Exec(`
				CREATE TABLE IF NOT EXISTS migrations 
				(
					id INT NOT NULL AUTO_INCREMENT, 
					key VARCHAR(255) NOT NULL,
					is_applied BOOLEAN NOT NULL DEFAULT FALSE,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					PRIMARY KEY (id)
				)
			`)
	case *pq.Driver:
		_, err = db.Exec(`
				CREATE TABLE IF NOT EXISTS migrations
				(
					id SERIAL PRIMARY KEY,
					key VARCHAR(255) NOT NULL,
					is_applied BOOLEAN NOT NULL DEFAULT FALSE,
					created_at TIMESTAMP NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMP NOT NULL DEFAULT NOW()
				)
			`)
	case *sqlite.Driver:
		_, err = db.Exec(`
				CREATE TABLE IF NOT EXISTS migrations
				(
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					key VARCHAR(255) NOT NULL,
					is_applied BOOLEAN NOT NULL DEFAULT FALSE,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				)
			`)
	}

	if err != nil {
		return khata.Wrap(err).SetExitCode(1).Explain("Could not create migrations table")
	}

	return nil
}

func CreateMigrationEntry(db *sql.DB, name string) *khata.Khata {
	var err error

	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
	case *sqlite.Driver:
		_, err = db.Exec("INSERT INTO migrations (key) VALUES (?)", name)
	case *pq.Driver:
		_, err = db.Exec("INSERT INTO migrations (key) VALUES ($1)", name)
	}

	if err != nil {
		return khata.Wrap(err).SetExitCode(1).Explain("Could not create migration entry")
	}

	return nil
}

func ApplyMigration(db *sql.DB, name string) *khata.Khata {
	var err error

	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
	case *sqlite.Driver:
		_, err = db.Exec("UPDATE migrations SET is_applied = true, updated_at = NOW() WHERE key = ?", name)
	case *pq.Driver:
		_, err = db.Exec("UPDATE migrations SET is_applied = true, updated_at = NOW() WHERE key = $1", name)
	}

	if err != nil {
		return khata.Wrap(err).SetExitCode(1).Explain("Could not apply migration")
	}

	return nil
}

func RollbackMigration(db *sql.DB, name string) *khata.Khata {
	var err error
	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
	case *sqlite.Driver:
		_, err = db.Exec("UPDATE migrations SET is_applied = false, updated_at = NOW() WHERE key = ?", name)
	case *pq.Driver:
		_, err = db.Exec("UPDATE migrations SET is_applied = false, updated_at = NOW() WHERE key = $1", name)
	}

	if err != nil {
		return khata.Wrap(err).SetExitCode(1).Explain("Could not rollback migration")
	}

	return nil
}

func GetMigrationStatus(db *sql.DB, name string) (bool, *khata.Khata) {
	var err error
	var isApplied bool

	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
	case *sqlite.Driver:
		err = db.QueryRow("SELECT is_applied FROM migrations WHERE key = ?", name).Scan(&isApplied)
	case *pq.Driver:
		err = db.QueryRow("SELECT is_applied FROM migrations WHERE key = $1", name).Scan(&isApplied)
	}

	if err != nil {
		return false, khata.Wrap(err).SetExitCode(1).Explain("Could not get migration status")
	}

	return isApplied, nil
}

func GetMigrationsFromDatabase(db *sql.DB, applied bool) ([]string, *khata.Khata) {
	var err error
	var migrations []string
	var rows *sql.Rows

	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
	case *sqlite.Driver:
		rows, err = db.Query("SELECT key FROM migrations WHERE is_applied = ?", applied)
	case *pq.Driver:
		rows, err = db.Query("SELECT key FROM migrations WHERE is_applied = $1", applied)
	}

	if err != nil {
		return nil, khata.Wrap(err).SetExitCode(1).Explain("Could not get migrations from database")
	}

	defer rows.Close()

	for rows.Next() {
		var name string

		err = rows.Scan(&name)

		if err != nil {
			return nil, khata.Wrap(err).SetExitCode(1).Explain("Could not scan migration row")
		}

		migrations = append(migrations, name)
	}

	return migrations, nil
}

func ExecuteMigration(db *sql.DB, sql string) *khata.Khata {
	_, err := db.Exec(sql)

	if err != nil {
		return khata.Wrap(err).SetExitCode(1).Explain("Could not execute migration")
	}

	return nil
}

func InitDb(cmd *cobra.Command) (*sql.DB, *khata.Khata) {
	var supportedDrivers = []string{"mysql", "postgres", "sqlite"}

	connection := GetStringArg(cmd, "connection", "MONARCH_CONNECTION_STRING", "")

	if connection == "" {
		return nil, khata.New("connection string is required").SetExitCode(1)
	}

	driver := GetStringArg(cmd, "driver", "MONARCH_DRIVER", "")

	if driver == "" {
		return nil, khata.New("driver is required").SetExitCode(1)
	}

	// Check if the driver is supported
	supported := IsDriverSupported(driver, supportedDrivers)

	if !supported {
		return nil, khata.New("driver not supported").SetExitCode(1)
	}

	// Try to connect to the database
	db, err := ConnectToDatabase(driver, connection)

	if err != nil {
		return nil, khata.Wrap(err).SetExitCode(1).Explain("Could not connect to database")
	}

	return db, err
}
