package DataBase

import (
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Connect() *gorm.DB {
	once.Do(func() {
		dsn := os.Getenv("DB_DSN")
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("[ERROR] database connection error: %v", err)
			return
		}
	})
	return db
}
