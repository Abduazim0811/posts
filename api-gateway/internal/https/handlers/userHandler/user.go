package userhandler

import (
	"api-gateway/internal/logger"
	"api-gateway/internal/protos/userProto/userproto"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	ClientUser userproto.UserServiceClient
}

func NewUserHandler(client userproto.UserServiceClient) *UserHandler {
	return &UserHandler{
		ClientUser: client,
	}
}

// SignUp godoc
// @Summary Foydalanuvchini ro`yxatdan o`tkazish
// @Description Yangi foydalanuvchi yaratadi va token qaytaradi
// @Tags users
// @Accept json
// @Produce json
// @Param request body userproto.SignUpReq true "Ro`yxatdan o`tish so`rovi"
// @Success 200 {object} userproto.Response 
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /signup [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	var req userproto.SignUpReq
	if err := c.BindJSON(&req); err != nil {
		logger.Logger.Printf("SignUp: JSON bog`lashda xato: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("SignUp: So`rov qabul qilindi: username=%s, email=%s", req.Username, req.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.ClientUser.SignUp(ctx, &req)
	if err != nil {
		logger.Logger.Printf("SignUp: gRPC xatosi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("SignUp: Muvaffaqiyatli yakunlandi: Token=%s", res.Token)
	c.JSON(http.StatusOK, res)
}

// SignIn godoc
// @Summary Foydalanuvchi tizimga kirishi
// @Description Foydalanuvchi email va parol bilan kiradi va token oladi
// @Tags users
// @Accept json
// @Produce json
// @Param request body userproto.SignINReq true "Kirish so`rovi"
// @Success 200 {object} userproto.Response "Muvaffaqiyatli kirish"
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Router /signin [post]
func (h *UserHandler) SignIn(c *gin.Context) {
	var req userproto.SignINReq
	if err := c.BindJSON(&req); err != nil {
		logger.Logger.Printf("SignIn: JSON bog`lashda xato: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("SignIn: So`rov qabul qilindi: email=%s", req.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.ClientUser.SignIn(ctx, &req)
	if err != nil {
		logger.Logger.Printf("SignIn: gRPC xatosi: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials: " + err.Error()})
		return
	}

	logger.Logger.Printf("SignIn: Muvaffaqiyatli yakunlandi: Token=%s", res.Token)
	c.JSON(http.StatusOK, res)
}

// GetUsersbyId godoc
// @Summary ID bo`yicha foydalanuvchi ma'lumotlarini olish
// @Description Berilgan ID bo`yicha foydalanuvchi ma'lumotlarini qaytaradi
// @Tags users
// @Produce json
// @Param id path string true "Foydalanuvchi ID"
// @Success 200 {object} userproto.User "Foydalanuvchi ma'lumotlari"
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /users/{id} [get]
func (h *UserHandler) GetUsersbyId(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logger.Logger.Printf("GetUsersbyId: ID kiritilmagan")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
		return
	}

	logger.Logger.Printf("GetUsersbyId: So`rov qabul qilindi: id=%s", id)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.ClientUser.GetUsersbyId(ctx, &userproto.UsersbyId{Id: id})
	if err != nil {
		logger.Logger.Printf("GetUsersbyId: gRPC xatosi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res == nil {
		logger.Logger.Printf("GetUsersbyId: Foydalanuvchi topilmadi: id=%s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	logger.Logger.Printf("GetUsersbyId: Muvaffaqiyatli yakunlandi: id=%s", res.Id)
	c.JSON(http.StatusOK, res)
}

// GetUsers godoc
// @Summary Barcha foydalanuvchilar ro`yxatini olish
// @Description Barcha foydalanuvchilar ro`yxatini qaytaradi
// @Tags users
// @Produce json
// @Success 200 {object} userproto.ListUsersRes "Foydalanuvchilar ro`yxati"
// @Failure 500 {object} string
// @Security Bearer
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	logger.Logger.Printf("GetUsers: So`rov qabul qilindi")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.ClientUser.GetUsers(ctx, &userproto.UsersEmpty{})
	if err != nil {
		logger.Logger.Printf("GetUsers: gRPC xatosi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("GetUsers: Muvaffaqiyatli yakunlandi: jami=%d", len(res.Users))
	c.JSON(http.StatusOK, res)
}

// UpdateUsers godoc
// @Summary Foydalanuvchi ma'lumotlarini yangilash
// @Description ID bo`yicha foydalanuvchi ma'lumotlarini yangilaydi
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Foydalanuvchi ID"
// @Param request body userproto.UpdateReq true "Yangilash so`rovi"
// @Success 200 {object} userproto.UpdateRes "Muvaffaqiyatli yangilandi"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUsers(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logger.Logger.Printf("UpdateUsers: ID kiritilmagan")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
		return
	}

	var req userproto.UpdateReq
	if err := c.BindJSON(&req); err != nil {
		logger.Logger.Printf("UpdateUsers: JSON bog`lashda xato: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = id

	logger.Logger.Printf("UpdateUsers: So`rov qabul qilindi: id=%s", req.Id)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.ClientUser.UpdateUsers(ctx, &req)
	if err != nil {
		logger.Logger.Printf("UpdateUsers: gRPC xatosi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("UpdateUsers: Muvaffaqiyatli yakunlandi: id=%s", req.Id)
	c.JSON(http.StatusOK, res)
}

// DeleteUsers godoc
// @Summary Foydalanuvchini o`chirish
// @Description ID bo`yicha foydalanuvchini o`chiradi
// @Tags users
// @Produce json
// @Param id path string true "Foydalanuvchi ID"
// @Success 200 {object} userproto.UpdateRes "Muvaffaqiyatli o`chirildi"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUsers(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logger.Logger.Printf("DeleteUsers: ID kiritilmagan")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
		return
	}

	logger.Logger.Printf("DeleteUsers: So`rov qabul qilindi: id=%s", id)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.ClientUser.DeleteUsers(ctx, &userproto.UsersbyId{Id: id})
	if err != nil {
		logger.Logger.Printf("DeleteUsers: gRPC xatosi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Printf("DeleteUsers: Muvaffaqiyatli yakunlandi: id=%s", id)
	c.JSON(http.StatusOK, res)
}