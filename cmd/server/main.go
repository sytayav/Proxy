package main

import (
	"Proxy/pkg/api"
	"Proxy/pkg/thumbnail"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
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
