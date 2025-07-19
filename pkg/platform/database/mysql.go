package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenDatabase creates and returns a GORM DB instance for MySQL.
func OpenDatabase(dbURL string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Use mysql.Open for GORM with MySQL
	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the MySQL database!")
	return db, nil
}

// RunMigrations applies database migrations using Goose for MySQL.
func RunMigrations(dbURL, migrationsDir string) error {
	// Goose works with *sql.DB, use "mysql" driver
	sqlDB, err := sql.Open("mysql", dbURL)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	// Set dialect for MySQL
	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	log.Printf("Running migrations from: %s", migrationsDir)
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return err
	}
	log.Println("Migrations applied successfully!")
	return nil
}
