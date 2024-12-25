package test

import (
	"Proxy/api"
	"Proxy/internal/thumbnail"
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestDownloadThumbnail(t *testing.T) {
	// Инициализация базы данных
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Создание таблицы для теста
	_, err = db.Exec("CREATE TABLE thumbnails (video_url TEXT PRIMARY KEY, image_data BLOB, cache_status TEXT)")
	require.NoError(t, err)

	// Инициализация gRPC-сервера с подключенной базой данных
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	api.RegisterThumbnailServiceServer(s, thumbnail.NewGRPCServer(db))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := api.NewThumbnailServiceClient(conn)
	req := &api.ThumbnailRequest{VideoUrl: "https://www.youtube.com/watch?v=pYZigRVSOXM"}
	res, err := client.DownloadThumbnail(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, "new", res.CacheStatus, "Expected cache status to be 'new'")
}
