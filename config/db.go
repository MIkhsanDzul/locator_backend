package config

import (
	"fmt"
	"locator-backend/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=ikhsan24 dbname=tracker port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek ke database: ", err)
	}

	fmt.Println(" Berhasil konek ke database")
}

func SaveLocationToPostgres(username string, timestamp int64, formatted string, latitude, longitude float64) error {
	var count int64
	err := DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		newUser := model.User{Username: username}
		if err := DB.Create(&newUser).Error; err != nil {
			return err
		}
	}

	query := `
		INSERT INTO locations (username, timestamp, timestamp_formatted, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (username) DO UPDATE SET
		timestamp = EXCLUDED.timestamp,
		timestamp_formatted = EXCLUDED.timestamp_formatted,
		latitude = EXCLUDED.latitude,
		longitude = EXCLUDED.longitude
	`
	err = DB.Exec(query, username, timestamp, formatted, latitude, longitude).Error
	return err
}
