package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBProduct struct {
	ID sql.NullString
}

type PingedStock struct {
	ID sql.NullString
}

func Connect() *gorm.DB {
	databaseDSN := "doadmin:qb0v2sk62i1hrmex@tcp(db-asos-do-user-9189913-0.b.db.ondigitalocean.com:25060)/flatspotdb"
	db, err := gorm.Open(mysql.Open(databaseDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database %s", err))
	}

	_ = db.AutoMigrate(&DBProduct{})
	_ = db.AutoMigrate(&PingedStock{})

	return db

}

