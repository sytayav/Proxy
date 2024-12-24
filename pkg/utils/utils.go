package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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
