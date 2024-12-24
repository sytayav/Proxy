package main

import (
	"Proxy/internal/db"
	"Proxy/pkg/api"
	"Proxy/pkg/thumbnail"
	"context"
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func GetThumbnail(videoUrl string, wg *sync.WaitGroup, async bool, database *sql.DB) {
	if wg != nil {
		defer wg.Done() // Уменьшаем счётчик ожидания по завершении работы функции
	}

	// Устанавливаем соединение с сервером
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewThumbnailServiceClient(conn)

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.DownloadThumbnail(ctx, &api.ThumbnailRequest{VideoUrl: videoUrl})
	if err != nil {
		log.Fatalf("could not download thumbnail: %v", err)
		return
	}

	if res == nil {
		log.Printf("No thumbnail found for %s", videoUrl)
		return
	}

	// Сохраняем изображение в файл
	err = saveImageToFile(videoUrl, res.ImageData)
	if err != nil {
		log.Printf("Ошибка при сохранении изображения: %v", err)
		return
	}

	log.Printf("Превью для %s загружено и сохранено, статус кэша: %s\n", videoUrl, res.CacheStatus)

	db.ChangeStatus(res, &videoUrl, database)
}

func saveImageToFile(videoUrl string, imageData []byte) error {
	// Извлекаем ID видео из URL для использования в качестве имени файла
	videoID, err := thumbnail.ExtractVideoID(videoUrl)
	if err != nil {
		return err
	}

	// Определяем путь для сохранения файла
	fileName := videoID + ".jpg"
	filePath := filepath.Join("thumbnails", fileName)

	// Создаем директорию, если она не существует
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	// Записываем данные в файл
	err = ioutil.WriteFile(filePath, imageData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Инициализация базы данных
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()

	// Настройка сервера gRPC
	go func() {
		lis, err := net.Listen("tcp", ":8084")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		api.RegisterThumbnailServiceServer(grpcServer, thumbnail.NewGRPCServer(database))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	asyncFlag := flag.Bool("async", false, "Download thumbnails asynchronously")
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("At least one video URL is required")
	}

	videoUrls := flag.Args()
	if *asyncFlag {
		var wg sync.WaitGroup // Создаём WaitGroup для отслеживания завершения горутин
		for _, videoUrl := range videoUrls {
			wg.Add(1)                                            // Увеличиваем счётчик
			go GetThumbnail(videoUrl, &wg, *asyncFlag, database) // Запускаем загрузку в горутине
		}
		wg.Wait() // Ожидаем завершения всех загрузок
	} else {
		for _, videoUrl := range videoUrls {
			GetThumbnail(videoUrl, nil, *asyncFlag, database) // Синхронная загрузка
		}
	}
}
