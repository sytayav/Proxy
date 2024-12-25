package utils

import (
	"Proxy/api"
	"Proxy/internal/db"
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func GetThumbnail(videoUrl string, wg *sync.WaitGroup, database *sql.DB) error {
	if wg != nil {
		defer wg.Done() // Уменьшаем счётчик ожидания по завершении работы функции
	}

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Устанавливаем соединение с сервером
	conn, err := grpc.DialContext(ctx, ":9098", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewThumbnailServiceClient(conn)

	// Создаем контекст с таймаутом для RPC-вызова
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer rpcCancel()

	res, err := c.DownloadThumbnail(rpcCtx, &api.ThumbnailRequest{VideoUrl: videoUrl})
	if err != nil {
		return fmt.Errorf("could not download thumbnail: %v", err)
	}

	// Сохраняем изображение в файл
	err = SaveImageToFile(videoUrl, res.ImageData)
	if err != nil {
		return fmt.Errorf("Ошибка при сохранении изображения: %v", err)
	}

	log.Printf("Превью для %s загружено и сохранено, статус кэша: %s\n", videoUrl, res.CacheStatus)

	err = db.ChangeStatus(res, &videoUrl, database)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса в базе данных: %w", err)
	}

	return nil
}

// Функция для извлечения ID видео из URL
func ExtractVideoID(videoURL string) (string, error) {
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

func SaveImageToFile(videoUrl string, imageData []byte) error {
	// Извлекаем ID видео из URL для использования в качестве имени файла
	videoID, err := ExtractVideoID(videoUrl)
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

// Функция для получения изображения по URL
func FetchImage(url string) ([]byte, error) {
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
