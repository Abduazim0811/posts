package userclients

import (
	"api-gateway/internal/config"
	"api-gateway/internal/protos/userProto/userproto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialUserClients() userproto.UserServiceClient {
	c := config.Configuration()

	conn, err := grpc.NewClient(c.User.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("post clients error:", err)
	}

	return userproto.NewUserServiceClient(conn)
}
