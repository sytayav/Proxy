package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Инициализация базы данных
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "cache.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Создаём таблицу для кеширования данных.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS thumbnails (video_url TEXT PRIMARY KEY, image_data BLOB, cache_status TEXT)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return db
}
