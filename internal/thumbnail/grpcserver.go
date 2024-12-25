package thumbnail

import (
	"Proxy/api"
	"Proxy/pkg/utils"
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"

	// Importing SQLite3 driver for database interaction
	_ "github.com/mattn/go-sqlite3"
)

// GRPCServer ...
type GRPCServer struct {
	api.UnimplementedThumbnailServiceServer         // Встраиваем неинициализированный сервер
	db                                      *sql.DB // Доступ к базе данных
}

// NewGRPCServer создает новый экземпляр сервера с доступом к базе данных
func NewGRPCServer(db *sql.DB) *GRPCServer {
	return &GRPCServer{db: db}
}

// Реализация метода mustEmbedUnimplementedThumbnailServiceServer
func (s *GRPCServer) mustEmbedUnimplementedThumbnailServiceServer() {
	// Этот метод может остаться пустым
}

// Метод DownloadThumbnail
func (s *GRPCServer) DownloadThumbnail(ctx context.Context, req *api.ThumbnailRequest) (*api.ThumbnailResponse, error) {
	log.Printf("Received request for video URL: %s", req.VideoUrl)

	if req == nil || req.VideoUrl == "" {
		return nil, status.Errorf(codes.InvalidArgument, "request or video URL is nil")
	}

	// Проверяем кэш
	var imageData []byte
	var cacheStatus string
	err := s.db.QueryRow("SELECT image_data, cache_status FROM thumbnails WHERE video_url = ?", req.VideoUrl).Scan(&imageData, &cacheStatus)

	if err == nil {
		// Если данные найдены в кэше
		log.Printf("Cache hit for %s", req.VideoUrl)
	} else if err == sql.ErrNoRows {
		// Если данных нет в кэше, получаем миниатюру
		log.Printf("Cache miss for %s, fetching new thumbnail", req.VideoUrl)
		videoID, err := utils.ExtractVideoID(req.VideoUrl)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid video URL: %v", err)
		}
		thumbnailURL := "https://img.youtube.com/vi/" + videoID + "/maxresdefault.jpg"

		imageData, err = utils.FetchImage(thumbnailURL)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to fetch thumbnail: %v", err)
		}

		cacheStatus = "new"

		// Сохраняем в кэш
		_, err = s.db.Exec("INSERT INTO thumbnails (video_url, image_data, cache_status) VALUES (?, ?, ?)", req.VideoUrl, imageData, cacheStatus)
		if err != nil {
			log.Printf("Error inserting into cache: %v", err)
		}
	} else {
		return nil, status.Errorf(codes.Internal, "failed to query cache: %v", err)
	}

	// Возвращаем данные
	res := &api.ThumbnailResponse{
		ImageData:   imageData,
		VideoUrl:    req.VideoUrl,
		CacheStatus: cacheStatus,
	}
	return res, nil
}
