package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(postgres.Open("host=127.0.0.1 port=5432 user=postgres password=root dbname=gosquash sslmode=disable"), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		panic(err)
	}
}
