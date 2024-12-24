package thumbnail

import (
	api2 "Proxy/api"
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // Обратите внимание на символ подчеркивания
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// GRPCServer ...
type GRPCServer struct {
	api2.UnimplementedThumbnailServiceServer         // Встраиваем неинициализированный сервер
	db                                       *sql.DB // Доступ к базе данных
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
func (s *GRPCServer) DownloadThumbnail(ctx context.Context, req *api2.ThumbnailRequest) (*api2.ThumbnailResponse, error) {
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
		videoID, err := ExtractVideoID(req.VideoUrl)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid video URL: %v", err)
		}
		thumbnailURL := "https://img.youtube.com/vi/" + videoID + "/maxresdefault.jpg"

		imageData, err = fetchImage(thumbnailURL)
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
	res := &api2.ThumbnailResponse{
		ImageData:   imageData,
		VideoUrl:    req.VideoUrl,
		CacheStatus: cacheStatus,
	}
	return res, nil
}

// Функция для извлечения ID видео из URL
func ExtractVideoID(videoURL string) (string, error) {
	// Простейший пример извлечения ID видео
	if strings.Contains(videoURL, "v=") {
		parts := strings.Split(videoURL, "v=")
		if len(parts) > 1 {
			if strings.Contains(parts[1], "&") {
				return strings.Split(parts[1], "&")[0], nil
			}
			return parts[1], nil
		}
	}
	if strings.Contains(videoURL, "youtu.be/") {
		parts := strings.Split(videoURL, "youtu.be/")
		if len(parts) > 1 {
			return parts[1], nil
		}
	}
	return "", status.Errorf(codes.InvalidArgument, "unable to extract video ID")
}

// Функция для получения изображения по URL
func fetchImage(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, status.Errorf(codes.Internal, "failed to fetch image, status: %s", response.Status)
	}

	return ioutil.ReadAll(response.Body)
}
