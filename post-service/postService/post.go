package postservice

import (
	"context"
	"post-service/internal/entity/post"
	"post-service/internal/service"
	pb "post-service/protos/postProto/postproto"
	"post-service/internal/logger"
)

type PostGrpc struct {
	pb.UnimplementedPostServiceServer
	svc service.PostService
}

func NewPostGrpc(svc service.PostService) *PostGrpc {
	return &PostGrpc{
		svc: svc,
	}
}

func (s *PostGrpc) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.PostResponse, error) {
	logger.Logger.Printf("gRPC CreatePost chaqirildi: username=%s, title=%s", req.Username, req.Title)

	createReq := post.CreatePostRequest{
		Username: req.Username,
		Title:    req.Title,
		Content:  req.Content,
		Tags:     req.Tags,
	}

	resp, err := s.svc.CreatePost(createReq)
	if err != nil {
		logger.Logger.Printf("CreatePost xatosi: %v", err)
		return nil, err
	}

	return &pb.PostResponse{
		Id:      resp.ID,
		Message: resp.Message,
	}, nil
}

func (s *PostGrpc) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	logger.Logger.Printf("gRPC GetPost chaqirildi: id=%s", req.Id)

	getReq := post.GetPostRequest{
		ID: req.Id,
	}

	resp, err := s.svc.GetPost(getReq)
	if err != nil {
		logger.Logger.Printf("GetPost xatosi: %v", err)
		return nil, err
	}
	if resp == nil {
		logger.Logger.Printf("Post topilmadi: id=%s", req.Id)
		return nil, nil
	}

	return &pb.GetPostResponse{
		Id:       resp.ID,
		Username: resp.Username,
		Title:    resp.Title,
		Content:  resp.Content,
		Tags:     resp.Tags,
	}, nil
}

func (s *PostGrpc) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	logger.Logger.Printf("gRPC ListPosts chaqirildi: username=%s, page=%d, limit=%d", req.Username, req.Page, req.Limit)

	listReq := post.ListPostsRequest{
		Username: req.Username,
		Page:     req.Page,
		Limit:    req.Limit,
	}

	resp, err := s.svc.ListPosts(listReq)
	if err != nil {
		logger.Logger.Printf("ListPosts xatosi: %v", err)
		return nil, err
	}

	pbPosts := make([]*pb.GetPostResponse, len(resp.Posts))
	for i, p := range resp.Posts {
		pbPosts[i] = &pb.GetPostResponse{
			Id:        p.ID,
			Username:  p.Username,
			Title:     p.Title,
			Content:   p.Content,
			Tags:      p.Tags,
		}
	}

	return &pb.ListPostsResponse{
		Posts: pbPosts,
		Total: resp.Total,
	}, nil
}

func (s *PostGrpc) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.PostResponse, error) {
	logger.Logger.Printf("gRPC UpdatePost chaqirildi: id=%s", req.Id)

	updateReq := post.UpdatePostRequest{
		ID:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Tags:    req.Tags,
	}

	resp, err := s.svc.UpdatePost(updateReq)
	if err != nil {
		logger.Logger.Printf("UpdatePost xatosi: %v", err)
		return nil, err
	}

	return &pb.PostResponse{
		Id:      resp.ID,
		Message: resp.Message,
	}, nil
}

func (s *PostGrpc) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	logger.Logger.Printf("gRPC DeletePost chaqirildi: id=%s", req.Id)

	deleteReq := post.DeletePostRequest{
		ID: req.Id,
	}

	resp, err := s.svc.DeletePost(deleteReq)
	if err != nil {
		logger.Logger.Printf("DeletePost xatosi: %v", err)
		return nil, err
	}

	return &pb.DeletePostResponse{
		Message: resp.Message,
	}, nil
}