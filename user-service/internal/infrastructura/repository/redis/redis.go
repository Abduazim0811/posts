package redis

import (
	"context"
	"encoding/json"
	"time"
	"user-service/internal/entity/users"
	"user-service/internal/logger"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewUserRedis(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) SetUser(ctx context.Context, key string, user *users.User, expiration time.Duration) error {
	logger.Logger.Printf("SetUser: Redis`ga keshlash boshlandi: key=%s", key)

	data, err := json.Marshal(user)
	if err != nil {
		logger.Logger.Printf("SetUser: JSON marshal qilishda xato: %v", err)
		return err
	}

	err = c.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		logger.Logger.Printf("SetUser: Redis`ga yozishda xato: %v", err)
		return err
	}

	logger.Logger.Printf("SetUser: Redis`ga muvaffaqiyatli keshlandi: key=%s", key)
	return nil
}

func (c *Cache) GetUser(ctx context.Context, key string) (*users.User, error) {
	logger.Logger.Printf("GetUser: Redis`dan olish boshlandi: key=%s", key)

	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		logger.Logger.Printf("GetUser: Redis`da topilmadi: key=%s", key)
		return nil, nil
	}
	if err != nil {
		logger.Logger.Printf("GetUser: Redis`dan olishda xato: %v", err)
		return nil, err
	}

	var user users.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		logger.Logger.Printf("GetUser: JSON unmarshal qilishda xato: %v", err)
		return nil, err
	}

	logger.Logger.Printf("GetUser: Redis`dan muvaffaqiyatli olindi: id=%s", user.ID)
	return &user, nil
}

func (c *Cache) DeleteUser(ctx context.Context, key string) error {
	logger.Logger.Printf("DeleteUser: Redis`dan o‘chirish boshlandi: key=%s", key)

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		logger.Logger.Printf("DeleteUser: Redis`dan o‘chirishda xato: %v", err)
		return err
	}

	logger.Logger.Printf("DeleteUser: Redis`dan muvaffaqiyatli o‘chirildi: key=%s", key)
	return nil
}

func (c *Cache) SetUsers(ctx context.Context, key string, usersRes *users.ListUsersRes, expiration time.Duration) error {
	logger.Logger.Printf("SetUsers: Redis`ga keshlash boshlandi: key=%s", key)

	data, err := json.Marshal(usersRes)
	if err != nil {
		logger.Logger.Printf("SetUsers: JSON marshal qilishda xato: %v", err)
		return err
	}

	err = c.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		logger.Logger.Printf("SetUsers: Redis`ga yozishda xato: %v", err)
		return err
	}

	logger.Logger.Printf("SetUsers: Redis`ga muvaffaqiyatli keshlandi: key=%s, jami=%d", key, len(usersRes.Users))
	return nil
}

func (c *Cache) GetUsers(ctx context.Context, key string) (*users.ListUsersRes, error) {
	logger.Logger.Printf("GetUsers: Redis`dan olish boshlandi: key=%s", key)

	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		logger.Logger.Printf("GetUsers: Redis`da topilmadi: key=%s", key)
		return nil, nil
	}
	if err != nil {
		logger.Logger.Printf("GetUsers: Redis`dan olishda xato: %v", err)
		return nil, err
	}

	var usersRes users.ListUsersRes
	err = json.Unmarshal(data, &usersRes)
	if err != nil {
		logger.Logger.Printf("GetUsers: JSON unmarshal qilishda xato: %v", err)
		return nil, err
	}

	logger.Logger.Printf("GetUsers: Redis`dan muvaffaqiyatli olindi: key=%s, jami=%d", key, len(usersRes.Users))
	return &usersRes, nil
}