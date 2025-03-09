package main

import (
	"net"
	"os"
	"os/signal"
	userclients "post-service/internal/clients/userClients"
	"post-service/internal/config"
	"post-service/internal/connection"
	"post-service/internal/infrastructura/repository/postgres"
	cache "post-service/internal/infrastructura/repository/redis"
	"post-service/internal/logger"
	"post-service/internal/service"
	postservice "post-service/postService"
	pb "post-service/protos/postProto/postproto"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	logger.Init()
	c := config.Configuration()
	conn, err := connection.NewConnection()
	if err != nil {
		logger.Logger.Fatalf("Ulanishlarda xato: %v", err)
	}
	defer conn.Close()

	repo := postgres.NewPostPostgres(conn.DB)
	cache := cache.NewCache(conn.Redis)
	svc := service.NewPostService(repo, cache)

	userConn := userclients.DialUserClients()

	grpcServer := grpc.NewServer()
	postGrpc := postservice.NewPostGrpc(*svc, &userConn)
	pb.RegisterPostServiceServer(grpcServer, postGrpc)

	lis, err := net.Listen(c.Post.Host, c.Post.Port)
	if err != nil {
		logger.Logger.Fatalf("Portni tinglashda xato: %v", err)
	}

	logger.Logger.Printf("gRPC server %s portida ishga tushdi", c.Post.Port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Logger.Fatalf("gRPC serverni ishga tushirishda xato: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Logger.Println("Server to`xtatilmoqda...")
	grpcServer.GracefulStop()
	logger.Logger.Println("Server muvaffaqiyatli to`xtatildi")
}