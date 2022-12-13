package database

import (
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database configuration struct
type DatabaseConfig struct {
	Driver   string
	Host     string
	Username string
	Password string
	Port     int
	Database string
}

type Database struct {
	*gorm.DB
}

func New(cfg *DatabaseConfig) (db *Database, err error) {
	// database configuration
	dbstring := "user=" + cfg.Username + " password=" + cfg.Password + " dbname=" + cfg.Database + " host=" + cfg.Host + " port=" + strconv.Itoa(cfg.Port) + " TimeZone=UTC"

	// database connection
	dbconn, err := gorm.Open(postgres.Open(dbstring), &gorm.Config{})
	if err != nil {
		return
	}

	// return
	db = &Database{dbconn}
	return
}
