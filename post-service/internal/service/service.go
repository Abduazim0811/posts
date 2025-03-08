package service

import (
	"context"
	"fmt"
	"post-service/internal/entity/post"
	"post-service/internal/infrastructura/repository"
	cache "post-service/internal/infrastructura/repository/redis"
	"post-service/internal/logger"
	"time"
)

type PostService struct {
	repo  repository.PostRepository
	cache *cache.Cache
}

func NewPostService(repo repository.PostRepository, cache *cache.Cache) *PostService {
	return &PostService{
		repo:  repo,
		cache: cache,
	}
}

func (s *PostService) CreatePost(req post.CreatePostRequest) (*post.PostResponse, error) {
	ctx := context.Background()
	resp, err := s.repo.CreatePost(req)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("post:%s", resp.ID)
	getPostResp := &post.GetPostResponse{
		ID:       resp.ID,
		Username: req.Username,
		Title:    req.Title,
		Content:  req.Content,
		Tags:     req.Tags,
	}
	err = s.cache.SetPost(ctx, cacheKey, getPostResp, 10*time.Minute)
	if err != nil {
		logger.Logger.Printf("Keshga saqlashda xato: %v", err)
	}

	return resp, nil
}

func (s *PostService) GetPost(req post.GetPostRequest) (*post.GetPostResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("post:%s", req.ID)

	cachedPost, err := s.cache.GetPost(ctx, cacheKey)
	if err != nil {
		logger.Logger.Printf("Keshdan olishda xato: %v", err)
	}
	if cachedPost != nil {
		logger.Logger.Printf("Post keshdan olindi: id=%s", req.ID)
		return cachedPost, nil
	}

	resp, err := s.repo.GetPost(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		err = s.cache.SetPost(ctx, cacheKey, resp, 10*time.Minute)
		if err != nil {
			logger.Logger.Printf("Keshga saqlashda xato: %v", err)
		}
	}

	return resp, nil
}

func (s *PostService) ListPosts(req post.ListPostsRequest) (*post.ListPostsResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("posts:username:%s:page:%d:limit:%d", req.Username, req.Page, req.Limit)

	cachedPosts, err := s.cache.GetPosts(ctx, cacheKey)
	if err != nil {
		logger.Logger.Printf("Keshdan olishda xato: %v", err)
	}
	if cachedPosts != nil {
		logger.Logger.Printf("Postlar keshdan olindi: username=%s", req.Username)
		return cachedPosts, nil
	}

	resp, err := s.repo.ListPosts(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		err = s.cache.SetPosts(ctx, cacheKey, resp, 10*time.Minute)
		if err != nil {
			logger.Logger.Printf("Keshga saqlashda xato: %v", err)
		}
	}

	return resp, nil
}

func (s *PostService) UpdatePost(req post.UpdatePostRequest) (*post.PostResponse, error) {
	ctx := context.Background()
	resp, err := s.repo.UpdatePost(req)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("post:%s", resp.ID)
	cachedPost, err := s.cache.GetPost(ctx, cacheKey)
	if err != nil {
		logger.Logger.Printf("Keshdan olishda xato: %v", err)
	}
	if cachedPost != nil {
		if req.Title != "" {
			cachedPost.Title = req.Title
		}
		if req.Content != "" {
			cachedPost.Content = req.Content
		}
		if len(req.Tags) > 0 {
			cachedPost.Tags = req.Tags
		}
		err = s.cache.SetPost(ctx, cacheKey, cachedPost, 10*time.Minute)
		if err != nil {
			logger.Logger.Printf("Keshni yangilashda xato: %v", err)
		}
	}

	return resp, nil
}

func (s *PostService) DeletePost(req post.DeletePostRequest) (*post.DeletePostResponse, error) {
	ctx := context.Background()
	resp, err := s.repo.DeletePost(req)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("post:%s", req.ID)
	err = s.cache.DeletePost(ctx, cacheKey)
	if err != nil {
		logger.Logger.Printf("Keshdan o`chirishda xato: %v", err)
	}


	return resp, nil
}