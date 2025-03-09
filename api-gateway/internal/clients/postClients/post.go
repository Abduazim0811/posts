package postclients

import (
	"api-gateway/internal/config"
	"api-gateway/internal/protos/postProto/postproto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialPostClients() postproto.PostServiceClient {
	c := config.Configuration()

	conn, err := grpc.NewClient(c.Post.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("post clients error:", err)
	}

	return postproto.NewPostServiceClient(conn)
}