package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"user-service/internal/config"
	"user-service/internal/connection"
	"user-service/internal/infrastructura/repository/postgres"
	cache "user-service/internal/infrastructura/repository/redis"
	"user-service/internal/logger"
	"user-service/internal/service"
	pb "user-service/protos/userProto/userproto"
	userservice "user-service/userService"

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

	repo := postgres.NewUserPostgres(conn.DB)
	cache := cache.NewUserRedis(conn.Redis)
	svc := service.NewUserService(repo, cache)

	grpcServer := grpc.NewServer()
	userGrpc := userservice.NewUserGrpc(*svc)
	pb.RegisterUserServiceServer(grpcServer, &userGrpc)

	lis, err := net.Listen(c.User.Host, c.User.Port)
	if err != nil {
		logger.Logger.Fatalf("Portni tinglashda xato: %v", err)
	}

	logger.Logger.Printf("gRPC server %s portida ishga tushdi", c.User.Port)

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