package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// Инициализация базы данных
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "cache.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS thumbnails (video_url TEXT PRIMARY KEY, image_data BLOB, cache_status TEXT)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}
