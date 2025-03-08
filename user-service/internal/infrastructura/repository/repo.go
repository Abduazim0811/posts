package repository

import "user-service/internal/entity/users"

type UserRepository interface {
	CreateUsers(req users.SignUpReq) (*users.User, error) 
	SignInUsers(req users.SignINReq) (string, error)
	GetUsersById(req users.UsersbyId) (*users.User, error)
	GetUsers() (*users.ListUsersRes, error)
	UpdateUsers(req users.UpdateReq) (*users.UpdateRes, error)
	DeleteUsers(req users.UsersbyId) (*users.UpdateRes, error)
}
