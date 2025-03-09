package userclients

import (
	"log"
	"post-service/internal/config"
	"post-service/protos/userProto/userproto"

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
