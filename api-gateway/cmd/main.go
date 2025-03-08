package main

import (
	postclients "api-gateway/internal/clients/postClients"
	"api-gateway/internal/config"
	"api-gateway/internal/logger"

	"api-gateway/internal/https/router"
	"log"
)

func main() {
	logger.Init()

	c := config.Configuration()
	conn := postclients.DialPostClients()
	r := router.SetupRouter(conn)

	log.Printf("API Gateway %s portida ishga tushdi", c.ApiGateway.Port)
	if err := r.Run(c.ApiGateway.Port); err != nil {
		log.Fatalf("HTTP server xatosi: %v", err)
	}
}
