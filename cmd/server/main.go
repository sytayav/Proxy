package main

import (
	"Proxy/api"
	"Proxy/internal/thumbnail"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

// UnaryInterceptor устанавливает таймаут для каждого входящего запроса
func UnaryInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Создаем новый контекст с таймаутом
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Вызываем основной обработчик
		return handler(ctx, req)
	}
}

func main() {
	s := grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptor(15 * time.Second)))
	srv := &thumbnail.GRPCServer{}
	api.RegisterThumbnailServiceServer(s, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {

		log.Fatal(err)
	}
}
