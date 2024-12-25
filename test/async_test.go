package test

import (
	"Proxy/api"
	"Proxy/internal/db"
	"Proxy/internal/thumbnail"
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

const bufSize2 = 1024 * 1024

var l *bufconn.Listener

func init() {
	l = bufconn.Listen(bufSize2)
}

func bufDialer2(context.Context, string) (net.Conn, error) {
	return l.Dial()
}

func TestAsyncDownload(t *testing.T) {
	var wg sync.WaitGroup
	database, err := db.InitDB()
	require.NoError(t, err)
	defer database.Close()

	// Настройка и запуск gRPC сервера
	s := grpc.NewServer()
	api.RegisterThumbnailServiceServer(s, thumbnail.NewGRPCServer(database))

	go func() {
		if err = s.Serve(l); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	videoUrls := []string{
		"https://www.youtube.com/watch?v=pYZigRVSOXM",
		"https://www.youtube.com/watch?v=USAI8p0DdA0",
	}

	for _, url := range videoUrls {
		wg.Add(1)
		go func(videoUrl string) {
			defer wg.Done()
			// Устанавливаем соединение с сервером
			conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer2), grpc.WithInsecure())
			require.NoError(t, err)
			defer conn.Close()

			c := api.NewThumbnailServiceClient(conn)

			// Вызываем метод загрузки
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err = c.DownloadThumbnail(ctx, &api.ThumbnailRequest{VideoUrl: videoUrl})
			require.NoError(t, err)
		}(url)
	}
	wg.Wait()
}
