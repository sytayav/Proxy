package test

import (
	"Proxy/api"
	"Proxy/internal/db"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	// Удаляем файл базы данных, если он существует
	os.Remove("cache.db")

	// Инициализируем базу данных
	database, err := db.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Проверяем, создается ли таблица
	var exists int
	row := database.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='thumbnails'")
	err = row.Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check if table exists: %v", err)
	}

	if exists == 0 {
		t.Fatal("Table 'thumbnails' does not exist")
	}

	// Закрываем базу данных
	defer func() {
		if err := database.Close(); err != nil {
			t.Fatalf("Failed to close database: %v", err)
		}
	}()
}

func TestChangeStatus(t *testing.T) {
	// Удаляем файл базы данных, если он существует
	os.Remove("cache.db")

	// Инициализируем базу данных
	database, err := db.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Добавим тестовые данные в таблицу
	videoUrl := "http://example.com/video"
	_, err = database.Exec("INSERT INTO thumbnails (video_url, image_data, cache_status) VALUES (?, ?, ?)", videoUrl, nil, "new")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Создаем объект ThumbnailResponse с новым статусом кэша
	res := &api.ThumbnailResponse{CacheStatus: "new"}

	// Вызываем функцию ChangeStatus
	err = db.ChangeStatus(res, &videoUrl, database)
	if err != nil {
		t.Fatalf("ChangeStatus returned an error: %v", err)
	}

	// Проверяем, что статус обновился
	var cacheStatus string
	row := database.QueryRow("SELECT cache_status FROM thumbnails WHERE video_url = ?", videoUrl)
	err = row.Scan(&cacheStatus)
	if err != nil {
		t.Fatalf("Failed to query cache_status: %v", err)
	}

	if cacheStatus != "hit" {
		t.Errorf("Expected cache_status to be 'hit', got '%s'", cacheStatus)
	}

	// Удаляем тестовые данные
	_, err = database.Exec("DELETE FROM thumbnails WHERE video_url = ?", videoUrl)
	if err != nil {
		t.Fatalf("Failed to delete test data: %v", err)
	}
}
