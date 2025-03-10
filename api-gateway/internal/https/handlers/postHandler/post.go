package posthandler

import (
	"context"
	"net/http"
	"time"

	"api-gateway/internal/logger"
	pb "api-gateway/internal/protos/postProto/postproto"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	ClientPost pb.PostServiceClient
}

func NewPostHandler(client pb.PostServiceClient) *PostHandler {
	return &PostHandler{
		ClientPost: client,
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post
// @Tags post
// @Accept json
// @Produce json
// @Param post body postproto.CreatePostRequest true "Create post request body"
// @Success 200 {object} postproto.PostResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req pb.CreatePostRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Logger.Printf("CreatePost: JSON bog`lashda xato: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("CreatePost: So`rov qabul qilindi: username=%s, title=%s", req.Username, req.Title)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.ClientPost.CreatePost(ctx, &req)
	if err != nil {
		logger.Logger.Printf("CreatePost: gRPC xatosi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("CreatePost: Muvaffaqiyatli yakunlandi: id=%s", res.Id)
	c.JSON(http.StatusOK, res)
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Get details of a specific post by its ID
// @Tags post
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} postproto.GetPostResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logger.Logger.Println("GetPost: ID parametri bo`sh")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	logger.Logger.Printf("GetPost: So`rov qabul qilindi: id=%s", id)

	req := &pb.GetPostRequest{Id: id}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.ClientPost.GetPost(ctx, req)
	if err != nil {
		logger.Logger.Printf("GetPost: gRPC xatosi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("GetPost: Muvaffaqiyatli yakunlandi: id=%s", res.Id)
	c.JSON(http.StatusOK, res)
}

// ListPosts godoc
// @Summary Get all posts
// @Description Get a list of all posts
// @Tags post
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} postproto.ListPostsResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	var req pb.ListPostsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Logger.Printf("ListPosts: Query parametrlarni bog`lashda xato: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("ListPosts: So`rov qabul qilindi: page=%d, limit=%d",req.Page, req.Limit)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.ClientPost.ListPosts(ctx, &req)
	if err != nil {
		logger.Logger.Printf("ListPosts: gRPC xatosi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("ListPosts: Muvaffaqiyatli yakunlandi: total=%d", res.Total)
	c.JSON(http.StatusOK, res)
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update an existing post by its ID
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body postproto.UpdatePostRequest true "Post update request body"
// @Success 200 {object} postproto.PostResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        logger.Logger.Println("UpdatePost: ID parametri bo`sh")
        c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
        return
    }

    logger.Logger.Printf("UpdatePost: So`rov qabul qilindi: id=%s", id)

    // Body dan ma'lumotlarni olish
    var req pb.UpdatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Logger.Printf("UpdatePost: JSON bog`lashda xato: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Id ni path dan olingan qiymat bilan to'ldirish
    req.Id = id // Path dan kelgan id ni o'rnatish

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()

    res, err := h.ClientPost.UpdatePost(ctx, &req)
    if err != nil {
        logger.Logger.Printf("UpdatePost: gRPC xatosi: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    logger.Logger.Printf("UpdatePost: Muvaffaqiyatli yakunlandi: id=%s", res.Id)
    c.JSON(http.StatusOK, res)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete an existing post by its ID
// @Tags post
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} postproto.DeletePostResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logger.Logger.Println("DeletePost: ID parametri bo`sh")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	logger.Logger.Printf("DeletePost: So`rov qabul qilindi: id=%s", id)

	req := &pb.DeletePostRequest{Id: id}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.ClientPost.DeletePost(ctx, req)
	if err != nil {
		logger.Logger.Printf("DeletePost: gRPC xatosi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("DeletePost: Muvaffaqiyatli yakunlandi: message=%s", res.Message)
	c.JSON(http.StatusOK, res)
}