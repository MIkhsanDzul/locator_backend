package config

import (
	"fmt"
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

// SaveLocationToPostgres menyimpan timestamp ke database PostgreSQL
func SaveLocationToPostgres(username string, timestamp int64, formatted string) error {
	return DB.Exec(`
		INSERT INTO locations (username, timestamp, timestamp_formatted)
		VALUES ($1, $2, $3)
		ON CONFLICT (username) DO UPDATE SET
			timestamp = excluded.timestamp,
			timestamp_formatted = excluded.timestamp_formatted
	`, username, timestamp, formatted).Error
}

