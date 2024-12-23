package main

import (
	"Proxy/internal/db"
	"Proxy/pkg/api"
	"Proxy/pkg/thumbnail"
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

func GetThumbnail(videoUrl string, wg *sync.WaitGroup, async bool) {
	if wg != nil {
		defer wg.Done() // Уменьшаем счётчик ожидания по завершении работы функции
	}

	// Устанавливаем соединение с сервером
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewThumbnailServiceClient(conn)
	res, err := c.DownloadThumbnail(context.Background(), &api.ThumbnailRequest{VideoUrl: videoUrl})
	if err != nil {
		log.Fatalf("could not download thumbnail: %v", err)
		return
	}

	if res == nil {
		log.Printf("No thumbnail found for %s", videoUrl)
		return
	}

	// Выводим информацию о миниатюре и статусе кэширования
	log.Printf("Thumbnail for %s downloaded, data length: %d bytes, cache status: %s\n",
		videoUrl, len(res.ImageData), res.CacheStatus)
}

func main() {
	// Инициализация базы данных
	database := db.InitDB()
	if database == nil {
		log.Fatal("Database initialization failed")
		return
	}
	defer database.Close()

	// Настройка сервера gRPC
	go func() {
		lis, err := net.Listen("tcp", ":8080")
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
			wg.Add(1)                                  // Увеличиваем счётчик
			go GetThumbnail(videoUrl, &wg, *asyncFlag) // Запускаем загрузку в горутине
		}
		wg.Wait() // Ожидаем завершения всех загрузок
	} else {
		for _, videoUrl := range videoUrls {
			GetThumbnail(videoUrl, nil, *asyncFlag) // Синхронная загрузка
		}
	}
}
