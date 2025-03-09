package router

import (
	_ "api-gateway/docs"
	posthandler "api-gateway/internal/https/handlers/postHandler"
	userhandler "api-gateway/internal/https/handlers/userHandler"
	"api-gateway/internal/pkg/jwt"
	pb "api-gateway/internal/protos/postProto/postproto"
	"api-gateway/internal/protos/userProto/userproto"
	middleware "api-gateway/internal/rateLimiting"
	"time"

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
func SetupRouter(userClient userproto.UserServiceClient, postClient pb.PostServiceClient) *gin.Engine {
	r := gin.Default()

	url := ginSwagger.URL("swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	rateLimiter := middleware.NewRateLimiter(10, time.Minute)
	r.Use(rateLimiter.Limit())

	userHandler := userhandler.NewUserHandler(userClient)
	postHandler := posthandler.NewPostHandler(postClient)

	r.POST("/signup", userHandler.SignUp)
	r.POST("/signin", userHandler.SignIn)

	userGroup := r.Group("/users").Use(jwt.AuthMiddleware())
	{
		userGroup.GET("/:id", userHandler.GetUsersbyId)
		userGroup.GET("", userHandler.GetUsers)
		userGroup.PUT("/:id", userHandler.UpdateUsers)
		userGroup.DELETE("/:id", userHandler.DeleteUsers)
	}

	postGroup := r.Group("/posts").Use(jwt.AuthMiddleware())
	{
		postGroup.POST("", postHandler.CreatePost)
		postGroup.GET("/:id", postHandler.GetPost)
		postGroup.GET("", postHandler.ListPosts)
		postGroup.PUT("/:id", postHandler.UpdatePost)
		postGroup.DELETE("/:id", postHandler.DeletePost)
	}

	return r
}
