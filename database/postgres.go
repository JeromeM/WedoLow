package database

import (
	"users-service/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(url string) *PostgresDB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&model.User{})

	return &PostgresDB{db: db}
}
