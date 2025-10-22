package main

import (
	"log"
	"net"

	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/repository"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/transport"
	"google.golang.org/grpc"

	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
)

const (
	port string = ":50051"
)

func main() {
	log.Printf("starting server on port %v", port)
	orderRepository := repository.NewOrderRepository()

	srv := transport.NewOrderServer(orderRepository)

	grpcServer := grpc.NewServer()

	api.RegisterOrderServiceServer(grpcServer, srv)

	lis, _ := net.Listen("tcp", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error with starting server: %e", err)
	}
}
