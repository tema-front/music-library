package migrations

import (
	"log"
	"music-library/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.Song{})

	if err != nil {
		log.Fatalf("couldn't migrate: %v", err)
	}
}