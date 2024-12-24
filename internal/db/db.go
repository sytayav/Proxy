package db

import (
	"Proxy/api"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
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

func ChangeStatus(res *api.ThumbnailResponse, videoUrl *string, database *sql.DB) {
	// Обновляем статус кэша, если это новое превью
	if res.CacheStatus == "new" {
		// Здесь должен быть код для обновления статуса в базе данных
		_, err := database.Exec("UPDATE thumbnails SET cache_status = ? WHERE video_url = ?", "hit", &videoUrl)
		if err != nil {
			log.Printf("Ошибка при обновлении статуса кэша: %v", err)
		}
	}
}
