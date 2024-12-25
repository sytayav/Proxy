package main

import (
	"Proxy/api"
	"Proxy/internal/db"
	"Proxy/internal/thumbnail"
	"Proxy/pkg/utils"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"

	// Importing SQLite3 driver for database interaction
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Инициализация базы данных
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
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
		var wg sync.WaitGroup
		for _, videoUrl := range videoUrls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				err := utils.GetThumbnail(url, nil, database)
				if err != nil {
					log.Printf("Ошибка при обработке %s: %v", url, err)
				}
			}(videoUrl)
		}
		wg.Wait()
	} else {
		for _, videoUrl := range videoUrls {
			err := utils.GetThumbnail(videoUrl, nil, database)
			if err != nil {
				log.Printf("Ошибка при обработке %s: %v", videoUrl, err)
			}
		}
	}
}
