package test

import (
	"Proxy/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestExtractVideoID
func TestExtractVideoID(t *testing.T) {
	tests := []struct {
		url      string
		expected string
		err      error
	}{
		{"https://www.youtube.com/watch?v=pYZigRVSOXM", "pYZigRVSOXM", nil},
		{"https://youtu.be/USAI8p0DdA0", "USAI8p0DdA0", nil},
		{"invalid-url", "", status.Errorf(codes.InvalidArgument, "unable to extract video ID")},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			id, err := utils.ExtractVideoID(tt.url)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, id)
			}
		})
	}
}

// TestSaveImageToFile
func TestSaveImageToFile(t *testing.T) {
	videoURL := "https://www.youtube.com/watch?v=pYZigRVSOXM"
	imageData := []byte{1, 2, 3}

	err := utils.SaveImageToFile(videoURL, imageData)
	require.NoError(t, err)

	// Проверяем, что файл был создан
	fileName := "thumbnails/pYZigRVSOXM.jpg"
	file, err := os.Stat(fileName)
	require.NoError(t, err)
	require.True(t, file.Size() > 0)

	// Удаляем файл после теста
	os.Remove(fileName)
}

// TestFetchImage
func TestFetchImage(t *testing.T) {
	// Создание тестового HTTP сервера
	imageData := []byte{1, 2, 3}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(imageData)
	}))
	defer ts.Close()

	// Тестирование FetchImage
	data, err := utils.FetchImage(ts.URL)
	require.NoError(t, err)
	require.Equal(t, imageData, data)
}

// TestFetchImage_BadRequest
func TestFetchImage_BadRequest(t *testing.T) {
	// Создание тестового HTTP сервера, который возвращает статус 500
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	// Тестирование FetchImage на основании статуса ошибки
	data, err := utils.FetchImage(ts.URL)
	require.Error(t, err)
	require.Nil(t, data)
}
