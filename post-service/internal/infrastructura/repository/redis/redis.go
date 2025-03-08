package cache

import (
	"context"
	"encoding/json"
	"post-service/internal/entity/post"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) SetPost(ctx context.Context, key string, post *post.GetPostResponse, expiration time.Duration) error {
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *Cache) GetPost(ctx context.Context, key string) (*post.GetPostResponse, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var p post.GetPostResponse
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Cache) DeletePost(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) SetPosts(ctx context.Context, key string, posts *post.ListPostsResponse, expiration time.Duration) error {
	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *Cache) GetPosts(ctx context.Context, key string) (*post.ListPostsResponse, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var posts post.ListPostsResponse
	err = json.Unmarshal(data, &posts)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}
