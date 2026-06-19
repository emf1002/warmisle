// Package pkg provides common utilities for warmisle.
package pkg

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance.
var DB *gorm.DB

// InitDatabase initializes the database connection.
func InitDatabase(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "sqlite",
		DSN:        dbPath + "?_journal_mode=WAL&_busy_timeout=5000&_cache=shared",
	}), &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Warn),
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	return err
}
