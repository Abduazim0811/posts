package service

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/entity/users"
	"user-service/internal/infrastructura/repository"
	"user-service/internal/infrastructura/repository/redis"
	"user-service/internal/logger"
	"user-service/internal/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      repository.UserRepository
	cache     *redis.Cache
}

func NewUserService(repo repository.UserRepository, cache *redis.Cache) *UserService {
	return &UserService{
		repo:      repo,
		cache:     cache,
	}
}

// SignUp yangi foydalanuvchi yaratadi va token qaytaradi
func (s *UserService) SignUp(ctx context.Context, req users.SignUpReq) (*users.Response, error) {
	logger.Logger.Printf("SignUp boshlandi: email=%s, username=%s", req.Email, req.Username)

	// Parolni hash qilish
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Printf("SignUp: Parolni hash qilishda xato: %v", err)
		return nil, err
	}

	resp, err := s.repo.CreateUsers(users.SignUpReq{
		Email:    req.Email,
		Password: string(passwordHash),
		Fullname: req.Fullname,
		Username: req.Username,
	})
	if err != nil {
		logger.Logger.Printf("SignUp: Foydalanuvchi yaratishda xato: %v", err)
		return nil, err
	}

	cacheKey := fmt.Sprintf("user:%s", req.Email)
	getUserPost := &users.User{
		ID:        resp.ID,
		Email:     resp.Email,
		Fullname:  resp.Fullname,
		Username:  resp.Username,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
	err = s.cache.SetUser(ctx, cacheKey, getUserPost, 10*time.Minute)
	if err != nil {
		logger.Logger.Printf("SignUp: Keshga saqlashda xato: %v", err)
	}

	token, err := jwt.GenerateToken(req.Email)
	if err != nil {
		logger.Logger.Printf("SignUp: Token yaratishda xato: %v", err)
		return nil, err
	}

	logger.Logger.Printf("SignUp muvaffaqiyatli yakunlandi: email=%s, id=%s", req.Email, resp.ID)
	return &users.Response{Token: token}, nil
}

func (s *UserService) SignIn(ctx context.Context, req users.SignINReq) (*users.Response, error) {
	logger.Logger.Printf("SignIn boshlandi: email=%s", req.Email)

	storedPassword, err := s.repo.SignInUsers(users.SignINReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Logger.Printf("SignIn: Foydalanuvchi ma'lumotlarini olishda xato: %v", err)
		return nil, err
	}
	if storedPassword == "" {
		logger.Logger.Printf("SignIn: Foydalanuvchi topilmadi: email=%s", req.Email)
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		logger.Logger.Printf("SignIn: Parol mos emas: email=%s", req.Email)
		return nil, fmt.Errorf("invalid credentials")
	}


	token, err := jwt.GenerateToken(req.Email)
	if err != nil {
		logger.Logger.Printf("SignIn: Token yaratishda xato: %v", err)
		return nil, err
	}

	logger.Logger.Printf("SignIn muvaffaqiyatli yakunlandi: email=%s", req.Email)
	return &users.Response{Token: token}, nil
}

func (s *UserService) GetUsersById(ctx context.Context, req users.UsersbyId) (*users.User, error) {
	logger.Logger.Printf("GetUsersById boshlandi: id=%s", req.ID)

	cacheKey := fmt.Sprintf("user:%s", req.ID)
	cachedUser, err := s.cache.GetUser(ctx, cacheKey)
	if err != nil {
		logger.Logger.Printf("GetUsersById: Redis`dan olishda xato: %v", err)
	}
	if cachedUser != nil {
		logger.Logger.Printf("GetUsersById: Redis`dan olindi: id=%s", cachedUser.ID)
		return cachedUser, nil
	}

	// Postgres’dan olish
	user, err := s.repo.GetUsersById(req)
	if err != nil {
		logger.Logger.Printf("GetUsersById: Postgres’dan olishda xato: %v", err)
		return nil, err
	}
	if user == nil {
		logger.Logger.Printf("GetUsersById: Foydalanuvchi topilmadi: id=%s", req.ID)
		return nil, fmt.Errorf("user not found")
	}

	// Redis’ga keshlash
	err = s.cache.SetUser(ctx, cacheKey, user, 10*time.Minute)
	if err != nil {
		logger.Logger.Printf("GetUsersById: Keshga saqlashda xato: %v", err)
	}

	logger.Logger.Printf("GetUsersById muvaffaqiyatli yakunlandi: id=%s", user.ID)
	return user, nil
}

func (s *UserService) GetUsers(ctx context.Context) (*users.ListUsersRes, error) {
	logger.Logger.Println("GetUsers boshlandi")

	cacheKey := "users:all"
	cachedUsers, err := s.cache.GetUsers(ctx, cacheKey)
	if err != nil {
		logger.Logger.Printf("GetUsers: Redisdan olishda xato: %v", err)
	}
	if cachedUsers != nil {
		logger.Logger.Printf("GetUsers: Redisdan olindi: jami=%d", len(cachedUsers.Users))
		return cachedUsers, nil
	}

	res, err := s.repo.GetUsers()
	if err != nil {
		logger.Logger.Printf("GetUsers: Postgresdan olishda xato: %v", err)
		return nil, err
	}

	err = s.cache.SetUsers(ctx, cacheKey, res, 10*time.Minute)
	if err != nil {
		logger.Logger.Printf("GetUsers: Keshga saqlashda xato: %v", err)
	}

	logger.Logger.Printf("GetUsers muvaffaqiyatli yakunlandi: jami=%d", len(res.Users))
	return res, nil
}

func (s *UserService) UpdateUsers(ctx context.Context, req users.UpdateReq) (*users.UpdateRes, error) {
	logger.Logger.Printf("UpdateUsers boshlandi: id=%s", req.ID)

	var password string
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Logger.Printf("UpdateUsers: Parolni hash qilishda xato: %v", err)
			return nil, err
		}
		password = string(hash)
	} else {
		password = req.Password
	}

	res, err := s.repo.UpdateUsers(users.UpdateReq{
		ID:       req.ID,
		Email:    req.Email,
		Username: req.Username,
		Password: password,
		Fullname: req.Fullname,
	})
	if err != nil {
		logger.Logger.Printf("UpdateUsers: Postgres`da yangilashda xato: %v", err)
		return nil, err
	}

	cacheKey := fmt.Sprintf("user:%s", req.ID)
	user, err := s.repo.GetUsersById(users.UsersbyId{ID: req.ID})
	if err == nil && user != nil {
		err = s.cache.SetUser(ctx, cacheKey, user, 10*time.Minute)
		if err != nil {
			logger.Logger.Printf("UpdateUsers: Keshga saqlashda xato: %v", err)
		}
	}
	err = s.cache.DeleteUser(ctx, "users:all")
	if err != nil {
		logger.Logger.Printf("UpdateUsers: Ro`yxat keshini o`chirishda xato: %v", err)
	}

	logger.Logger.Printf("UpdateUsers muvaffaqiyatli yakunlandi: id=%s", req.ID)
	return res, nil
}

func (s *UserService) DeleteUsers(ctx context.Context, req users.UsersbyId) (*users.UpdateRes, error) {
	logger.Logger.Printf("DeleteUsers boshlandi: id=%s", req.ID)

	res, err := s.repo.DeleteUsers(req)
	if err != nil {
		logger.Logger.Printf("DeleteUsers: Postgres`da o`chirishda xato: %v", err)
		return nil, err
	}

	cacheKey := fmt.Sprintf("user:%s", req.ID)
	err = s.cache.DeleteUser(ctx, cacheKey)
	if err != nil {
		logger.Logger.Printf("DeleteUsers: Foydalanuvchi keshini o`chirishda xato: %v", err)
	}
	err = s.cache.DeleteUser(ctx, "users:all")
	if err != nil {
		logger.Logger.Printf("DeleteUsers: Ro`yxat keshini o`chirishda xato: %v", err)
	}

	logger.Logger.Printf("DeleteUsers muvaffaqiyatli yakunlandi: id=%s", req.ID)
	return res, nil
}
