package main

import (
	postclients "api-gateway/internal/clients/postClients"
	userclients "api-gateway/internal/clients/userClients"
	"api-gateway/internal/config"
	"api-gateway/internal/logger"

	"api-gateway/internal/https/router"
	"log"
)

func main() {
	logger.Init()

	c := config.Configuration()
	postConn := postclients.DialPostClients()
	userConn := userclients.DialUserClients()
	r := router.SetupRouter(userConn, postConn)

	log.Printf("API Gateway %s portida ishga tushdi", c.ApiGateway.Port)
	if err := r.Run(c.ApiGateway.Port); err != nil {
		log.Fatalf("HTTP server xatosi: %v", err)
	}
}

