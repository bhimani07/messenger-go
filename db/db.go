package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var conn *gorm.DB = nil

func GetDbInstance() *gorm.DB {
	if conn == nil {
		dsn := "host=localhost user=postgres password=postgres dbname=messenger port=5432 sslmode=disable"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		conn = db
	}
	return conn
}
