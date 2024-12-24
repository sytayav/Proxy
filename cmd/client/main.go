package main

import (
	"Proxy/api"
	"Proxy/internal/db"
	"Proxy/internal/thumbnail"
	"Proxy/pkg/utils"
	"context"
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

func GetThumbnail(videoUrl string, wg *sync.WaitGroup, database *sql.DB) {
	if wg != nil {
		defer wg.Done() // Уменьшаем счётчик ожидания по завершении работы функции
	}

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Устанавливаем соединение с сервером
	conn, err := grpc.DialContext(ctx, ":9092", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewThumbnailServiceClient(conn)

	// Создаем контекст с таймаутом для RPC-вызова
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer rpcCancel()

	res, err := c.DownloadThumbnail(rpcCtx, &api.ThumbnailRequest{VideoUrl: videoUrl})
	if err != nil {
		log.Fatalf("could not download thumbnail: %v", err)
		return
	}

	// Сохраняем изображение в файл
	err = utils.SaveImageToFile(videoUrl, res.ImageData)
	if err != nil {
		log.Printf("Ошибка при сохранении изображения: %v", err)
		return
	}

	log.Printf("Превью для %s загружено и сохранено, статус кэша: %s\n", videoUrl, res.CacheStatus)

	db.ChangeStatus(res, &videoUrl, database)
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
			wg.Add(1)                                // Увеличиваем счётчик
			go GetThumbnail(videoUrl, &wg, database) // Запускаем загрузку в горутине
		}
		wg.Wait() // Ожидаем завершения всех загрузок
	} else {
		for _, videoUrl := range videoUrls {
			GetThumbnail(videoUrl, nil, database) // Синхронная загрузка
		}
	}
}
