package main

import (
	"gorm-management-users/models"

	"gorm.io/gorm"
)

func migrateDatabase(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate database schema")
	}
}
