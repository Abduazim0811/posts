package users

type User struct {
	ID        string
	Email     string
	Username  string
	Password  string 
	Fullname  string
	CreatedAt string 
	UpdatedAt string
}

type SignUpReq struct {
	Email    string
	Username string
	Password string
	Fullname string
}

type UsersbyUsername struct{
	Username string
}

type SignUpRes struct{
	ID		string
}
type Response struct {
	Token string
}

type SignINReq struct {
	Email    string
	Password string
}

type UsersbyId struct {
	ID string
}

type ListUsersRes struct {
	Users []User
}

type UpdateReq struct {
	ID       string
	Email    string
	Username string
	Password string
	Fullname string
}

type UpdateRes struct {
	Message string
}