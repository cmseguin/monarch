package utils

import (
	"database/sql"
	"errors"

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

func ConnectToDatabase(driver string, connection string) (db *sql.DB, err error) {
	db, err = sql.Open(driver, connection)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateMigrationTable(db *sql.DB) (err error) {
	// Check if the db driver is mysql
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

	return err
}

func CreateMigrationEntry(db *sql.DB, name string) (err error) {
	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
	case *sqlite.Driver:
		_, err = db.Exec("INSERT INTO migrations (key) VALUES (?)", name)
	case *pq.Driver:
		_, err = db.Exec("INSERT INTO migrations (key) VALUES ($1)", name)
	}

	return err
}

func ApplyMigration(db *sql.DB, name string) (err error) {
	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
	case *sqlite.Driver:
		_, err = db.Exec("UPDATE migrations SET is_applied = true, updated_at = NOW() WHERE name = ?", name)
	case *pq.Driver:
		_, err = db.Exec("UPDATE migrations SET is_applied = true, updated_at = NOW() WHERE name = $1", name)
	}

	return err
}

func RollbackMigration(db *sql.DB, name string) (err error) {
	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
	case *sqlite.Driver:
		_, err = db.Exec("UPDATE migrations SET is_applied = false, updated_at = NOW() WHERE name = ?", name)
	case *pq.Driver:
		_, err = db.Exec("UPDATE migrations SET is_applied = false, updated_at = NOW() WHERE name = $1", name)
	}

	return err
}

func GetMigrationStatus(db *sql.DB, name string) (isApplied bool, err error) {
	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
	case *sqlite.Driver:
		err = db.QueryRow("SELECT is_applied FROM migrations WHERE name = ?", name).Scan(&isApplied)
	case *pq.Driver:
		err = db.QueryRow("SELECT is_applied FROM migrations WHERE name = $1", name).Scan(&isApplied)
	}

	return isApplied, err
}

func GetMigrationsFromDatabase(db *sql.DB, applied bool) (migrations []string, err error) {
	var rows *sql.Rows

	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
	case *sqlite.Driver:
		rows, err = db.Query("SELECT key FROM migrations WHERE is_applied = ?", applied)
	case *pq.Driver:
		rows, err = db.Query("SELECT key FROM migrations WHERE is_applied = $1", applied)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var name string

		err = rows.Scan(&name)

		if err != nil {
			return nil, err
		}

		migrations = append(migrations, name)
	}

	return migrations, nil
}

func ExecuteMigration(db *sql.DB, sql string) (err error) {
	_, err = db.Exec(sql)

	return err
}

func InitDb(cmd *cobra.Command) (db *sql.DB, err error) {
	var supportedDrivers = []string{"mysql", "postgres", "sqlite"}

	connection := GetStringArg(cmd, "connection", "MONARCH_CONNECTION_STRING", "")

	if connection == "" {
		return nil, errors.New("connection string is required")
	}

	driver := GetStringArg(cmd, "driver", "MONARCH_DRIVER", "")

	if driver == "" {
		return nil, errors.New("driver is required")
	}

	// Check if the driver is supported
	supported := IsDriverSupported(driver, supportedDrivers)

	if !supported {
		return nil, errors.New("driver not supported")
	}

	// Try to connect to the database
	db, err = ConnectToDatabase(driver, connection)

	if err != nil {
		return nil, errors.New("error connecting to database")
	}

	return db, err
}
