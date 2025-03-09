package userservice

import (
	"context"
	"user-service/internal/entity/users"
	"user-service/internal/logger"
	"user-service/internal/service"
	pb "user-service/protos/userProto/userproto"
)

type UserGrpc struct {
	pb.UnimplementedUserServiceServer
	svc service.UserService
}

func NewUserGrpc(svc service.UserService) UserGrpc {
	return UserGrpc{svc: svc}
}

func (u *UserGrpc) SignUp(ctx context.Context, req *pb.SignUpReq) (*pb.Response, error) {
	logger.Logger.Printf("gRPC SignUp chaqirildi: email=%s", req.Email)

	signupReq := users.SignUpReq{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
		Fullname: req.Fullname,
	}

	resp, err := u.svc.SignUp(ctx, signupReq)
	if err != nil {
		logger.Logger.Printf("SignUp xatosi: %v", err)
		return nil, err
	}

	return &pb.Response{
		Token: resp.Token,
	}, nil
}

func (u *UserGrpc) SignIn(ctx context.Context, req *pb.SignINReq) (*pb.Response, error) {
	logger.Logger.Printf("gRPC SignIn chaqirildi: email=%s", req.Email)

	signInReq := users.SignINReq{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := u.svc.SignIn(ctx, signInReq)
	if err != nil {
		logger.Logger.Printf("SignIn xatosi: %v", err)
		return nil, err
	}

	return &pb.Response{
		Token: resp.Token,
	}, nil
}

func (u *UserGrpc) GetUsersbyId(ctx context.Context, req *pb.UsersbyId) (*pb.User, error) {
	logger.Logger.Printf("gRPC GetUsersById chaqirildi: id=%s", req.Id)

	getReq := users.UsersbyId{
		ID: req.Id,
	}

	resp, err := u.svc.GetUsersById(ctx, getReq)
	if err != nil {
		logger.Logger.Printf("GetUsersById xatosi: %v", err)
		return nil, err
	}
	if resp == nil {
		logger.Logger.Printf("Foydalanuvchi topilmadi: id=%s", req.Id)
		return nil, nil
	}

	return &pb.User{
		Id:       resp.ID,
		Email:    resp.Email,
		Username: resp.Username,
		Fullname: resp.Fullname,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}, nil
}

func (u *UserGrpc) GetUsersbyUsername(ctx context.Context, req *pb.GetbyUsernameReq) (*pb.User, error){
	logger.Logger.Printf("gRPC GetUsersByUsername chaqirildi: username=%s", req.Username)

	getReq := users.UsersbyUsername{
		Username: req.Username,
	}

	resp, err := u.svc.GetUsersbyUsername(ctx, getReq)
	if err != nil {
		logger.Logger.Printf("GetUsersById xatosi: %v", err)
		return nil, err
	}
	if resp == nil {
		logger.Logger.Printf("Foydalanuvchi topilmadi: username=%s", req.Username)
		return nil, nil
	}

	return &pb.User{
		Id:       resp.ID,
		Email:    resp.Email,
		Username: resp.Username,
		Fullname: resp.Fullname,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}, nil
}

func (u *UserGrpc) GetUsers(ctx context.Context, req *pb.UsersEmpty) (*pb.ListUsersRes, error) {
	logger.Logger.Printf("gRPC GetUsers chaqirildi")

	

	resp, err := u.svc.GetUsers(ctx)
	if err != nil {
		logger.Logger.Printf("GetUsers xatosi: %v", err)
		return nil, err
	}

	pbUsers := make([]*pb.User, len(resp.Users))
	for i, user := range resp.Users {
		pbUsers[i] = &pb.User{
			Id:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Fullname: user.Fullname,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return &pb.ListUsersRes{
		Users: pbUsers,
	}, nil
}

func (u *UserGrpc) UpdateUsers(ctx context.Context, req *pb.UpdateReq) (*pb.UpdateRes, error) {
	logger.Logger.Printf("gRPC UpdateUsers chaqirildi: id=%s", req.Id)

	updateReq := users.UpdateReq{
		ID:       req.Id,
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Fullname: req.Fullname,
	}

	resp, err := u.svc.UpdateUsers(ctx, updateReq)
	if err != nil {
		logger.Logger.Printf("UpdateUsers xatosi: %v", err)
		return nil, err
	}

	return &pb.UpdateRes{
		Message: resp.Message,
	}, nil
}

func (u *UserGrpc) DeleteUsers(ctx context.Context, req *pb.UsersbyId) (*pb.UpdateRes, error) {
	logger.Logger.Printf("gRPC DeleteUsers chaqirildi: id=%s", req.Id)

	deleteReq := users.UsersbyId{
		ID: req.Id,
	}

	resp, err := u.svc.DeleteUsers(ctx, deleteReq)
	if err != nil {
		logger.Logger.Printf("DeleteUsers xatosi: %v", err)
		return nil, err
	}

	return &pb.UpdateRes{
		Message: resp.Message,
	}, nil
}