package router

import (
	posthandler "api-gateway/internal/https/handlers/postHandler"
	pb "api-gateway/internal/protos/post_proto/postproto"
	middleware "api-gateway/internal/rateLimiting"
	"time"
	_ "api-gateway/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)
// @title Post Service API
// @version 1.0
// @description This is an API Gateway for managing posts in a blogging system.
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter the token in the format `Bearer {token}`
// @host localhost:8080
// @BasePath /
func SetupRouter(postClient pb.PostServiceClient) *gin.Engine {
	r := gin.Default()

	url := ginSwagger.URL("swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	rateLimiter := middleware.NewRateLimiter(10, time.Minute)
	r.Use(rateLimiter.Limit())

	postHandler := posthandler.NewPostHandler(postClient)

	r.POST("/posts", postHandler.CreatePost)   
	r.GET("/posts/:id", postHandler.GetPost)      
	r.GET("/posts", postHandler.ListPosts)         
	r.PUT("/posts/:id", postHandler.UpdatePost)    
	r.DELETE("/posts/:id", postHandler.DeletePost) 

	return r
}
