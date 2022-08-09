package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	MaxIdleConnection = 10
	MaxOpenConnection = 10
	connectionPattern = "host=%s user=%s password=%s dbname=%s sslmode=disable"
)

// WritePostgresDB function for creating database connection for write-access
func WritePostgresDB() *sql.DB {
	return CreateDBConnection(fmt.Sprintf(connectionPattern,
		os.Getenv("WRITE_DB_HOST"), os.Getenv("WRITE_DB_USER"), os.Getenv("WRITE_DB_PASSWORD"), os.Getenv("WRITE_DB_NAME")))

}

// ReadPostgresDB function for creating database connection for read-access
func ReadPostgresDB() *sql.DB {
	return CreateDBConnection(fmt.Sprintf(connectionPattern,
		os.Getenv("READ_DB_HOST"), os.Getenv("READ_DB_USER"), os.Getenv("READ_DB_PASSWORD"), os.Getenv("READ_DB_NAME")))

}

// ReadCronPostgresDB function for creating database connection for read-access
func ReadCronPostgresDB() *sql.DB {
	return CreateDBConnection(fmt.Sprintf(connectionPattern,
		os.Getenv("READ_DB_CRON_HOST"), os.Getenv("READ_DB_CRON_USER"), os.Getenv("READ_DB_CRON_PASSWORD"), os.Getenv("READ_DB_CRON_NAME")))

}

// CreateDBConnection function for creating database connection
func CreateDBConnection(descriptor string) *sql.DB {
	db, err := sql.Open("postgres", descriptor)
	if err != nil {
		defer db.Close()
		return db
	}

	db.SetMaxIdleConns(MaxIdleConnection)
	db.SetMaxOpenConns(MaxOpenConnection)
	db.SetConnMaxLifetime(time.Second * 5)

	return db
}

// CloseDb function for closing database connection
func CloseDb(db *sql.DB) {
	if db != nil {
		db.Close()
		db = nil
	}
}
