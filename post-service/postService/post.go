package postservice

import (
	"context"
	"post-service/internal/entity/post"
	"post-service/internal/logger"
	"post-service/internal/service"
	pb "post-service/protos/postProto/postproto"
	userpb "post-service/protos/userProto/userproto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostGrpc struct {
	pb.UnimplementedPostServiceServer
	svc        service.PostService
	userClient userpb.UserServiceClient
}

func NewPostGrpc(svc service.PostService, userConn *userpb.UserServiceClient) *PostGrpc {
	return &PostGrpc{
		svc:        svc,
		userClient: *userConn,
	}
}

func (s *PostGrpc) verifyUsername(ctx context.Context, username string) error {
	_, err := s.userClient.GetUsersbyUsername(ctx, &userpb.GetbyUsernameReq{Username: username})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			logger.Logger.Printf("verifyUsername: Foydalanuvchi topilmadi: username=%s", username)
			return status.Error(codes.NotFound, "User not found")
		}
		logger.Logger.Printf("verifyUsername: User service xatosi: %v", err)
		return status.Error(codes.Internal, "Internal server error")
	}
	logger.Logger.Printf("verifyUsername: Foydalanuvchi tasdiqlandi: username=%s", username)
	return nil
}

func (s *PostGrpc) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.PostResponse, error) {
	logger.Logger.Printf("gRPC CreatePost chaqirildi: username=%s, title=%s", req.Username, req.Title)

	if err := s.verifyUsername(ctx, req.Username); err != nil {
		return nil, err
	}

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
	logger.Logger.Printf("gRPC ListPosts chaqirildi:  page=%d, limit=%d", req.Page, req.Limit)


	listReq := post.ListPostsRequest{
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
			Id:       p.ID,
			Username: p.Username,
			Title:    p.Title,
			Content:  p.Content,
			Tags:     p.Tags,
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
		Title:   resp.Title,
		Content: resp.Content, 
		Tags:    resp.Tags,    
		Username: resp.Username, 
	}, nil
}
func (s *PostGrpc) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	logger.Logger.Printf("gRPC DeletePost chaqirildi: id=%s", req.Id)

	postt, err := s.svc.GetPost(post.GetPostRequest{ID: req.Id})
	if err != nil {
		logger.Logger.Printf("DeletePost: Postni olishda xato: %v", err)
		return nil, err
	}
	if postt == nil {
		logger.Logger.Printf("DeletePost: Post topilmadi: id=%s", req.Id)
		return nil, status.Error(codes.NotFound, "Post not found")
	}
	if err := s.verifyUsername(ctx, postt.Username); err != nil {
		return nil, err
	}

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
