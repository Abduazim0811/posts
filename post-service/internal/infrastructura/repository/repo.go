package repository

import "post-service/internal/entity/post"

type PostRepository interface {
	CreatePost(req post.CreatePostRequest) (*post.PostResponse, error)
	GetPost(req post.GetPostRequest) (*post.GetPostResponse, error)
	ListPosts(req post.ListPostsRequest) (*post.ListPostsResponse, error)
	UpdatePost(req post.UpdatePostRequest) (*post.PostResponse, error)
	DeletePost(req post.DeletePostRequest) (*post.DeletePostResponse, error)
}
