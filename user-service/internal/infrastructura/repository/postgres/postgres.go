package postgres

import (
	"database/sql"
	"fmt"
	"time"
	"user-service/internal/entity/users"
	"user-service/internal/infrastructura/repository"
	"user-service/internal/logger"

	"github.com/Masterminds/squirrel"
)

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) repository.UserRepository {
	return &UserPostgres{db: db}
}

func (u *UserPostgres) CreateUsers(req users.SignUpReq) (*users.User, error) {
	logger.Logger.Printf("CreateUsers boshlandi: username=%s, email=%s", req.Username, req.Email)

	sqlQuery, args, err := squirrel.
		Insert("users").
		Columns("username", "email", "password_hash", "full_name").
		Values(req.Username, req.Email, req.Password, req.Fullname).
		Suffix("RETURNING id, email, username, full_name, created_at, updated_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("CreateUsers: SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	var user users.User
	err = u.db.QueryRow(sqlQuery, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Fullname,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		logger.Logger.Printf("CreateUsers: Foydalanuvchi qo`shishda xato: %v", err)
		return nil, err
	}

	logger.Logger.Printf("CreateUsers muvaffaqiyatli yakunlandi: id=%s", user.ID)
	return &user, nil
}

func (u *UserPostgres) SignInUsers(req users.SignINReq) (string, error) {
	logger.Logger.Printf("SignInUsers boshlandi: email=%s", req.Email)
	var password string
	sqlQuery, args, err := squirrel.
		Select("password_hash").
		From("users").
		Where(squirrel.Eq{"email": req.Email}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("SignInUsers: SQL so`rov yaratishda xato: %v", err)
		return "", err
	}
	err = u.db.QueryRow(sqlQuery, args...).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Printf("SignInUsers: Foydalanuvchi topilmadi: email=%s", req.Email)
			return "", nil
		}
		logger.Logger.Printf("SignInUsers: Ma`lumot olishda xato: %v", err)
		return "", err
	}

	

	logger.Logger.Println("SignInUsers muvaffaqiyatli yakunlandi:")
	return password, nil
}

func (u *UserPostgres) GetUsersById(req users.UsersbyId) (*users.User, error) {
	logger.Logger.Printf("GetUsersById boshlandi: id=%s", req.ID)

	sqlQuery, args, err := squirrel.
		Select("id", "email", "username", "full_name", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("GetUsersById: SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	var user users.User
	err = u.db.QueryRow(sqlQuery, args...).Scan(
		&user.ID, &user.Email, &user.Username, &user.Fullname, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Printf("GetUsersById: Foydalanuvchi topilmadi: id=%s", req.ID)
			return nil, nil
		}
		logger.Logger.Printf("GetUsersById: Ma’lumot olishda xato: %v", err)
		return nil, err
	}

	logger.Logger.Printf("GetUsersById muvaffaqiyatli yakunlandi: id=%s", user.ID)
	return &user, nil
}

func (u *UserPostgres) GetUsersbyUsername(req users.UsersbyUsername) (*users.User, error) {
	logger.Logger.Printf("GetUsersbyUsername boshlandi: username=%s", req.Username)

	sqlQuery, args, err := squirrel.
		Select("id", "email", "username", "full_name", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"username": req.Username}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("GetUsersbyUsername: SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	var user users.User
	err = u.db.QueryRow(sqlQuery, args...).Scan(
		&user.ID, &user.Email, &user.Username, &user.Fullname, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Printf("GetUsersbyUsername: Foydalanuvchi topilmadi: username=%s", req.Username)
			return nil, nil
		}
		logger.Logger.Printf("GetUsersbyUsername: Ma`lumot olishda xato: %v", err)
		return nil, err
	}

	logger.Logger.Printf("GetUsersbyUsername muvaffaqiyatli yakunlandi: username=%s", user.Username)
	return &user, nil
}

func (u *UserPostgres) GetUsers() (*users.ListUsersRes, error) {
	logger.Logger.Println("GetUsers boshlandi")

	sqlQuery, args, err := squirrel.
		Select("id", "email", "username", "full_name", "created_at", "updated_at").
		From("users").
		OrderBy("created_at DESC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("GetUsers: SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	rows, err := u.db.Query(sqlQuery, args...)
	if err != nil {
		logger.Logger.Printf("GetUsers: Ma`lumot olishda xato: %v", err)
		return nil, err
	}
	defer rows.Close()

	var usersList []users.User
	for rows.Next() {
		var user users.User
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Fullname, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			logger.Logger.Printf("GetUsers: Satrdan o`qishda xato: %v", err)
			return nil, err
		}
		usersList = append(usersList, user)
	}

	logger.Logger.Printf("GetUsers muvaffaqiyatli yakunlandi: jami=%d", len(usersList))
	return &users.ListUsersRes{
		Users: usersList,
	}, nil
}

func (u *UserPostgres) UpdateUsers(req users.UpdateReq) (*users.UpdateRes, error) {
	logger.Logger.Printf("UpdateUsers boshlandi: id=%s", req.ID)

	updateBuilder := squirrel.Update("users").Where(squirrel.Eq{"id": req.ID})

	if req.Email != "" {
		updateBuilder = updateBuilder.Set("email", req.Email)
	}
	if req.Username != "" {
		updateBuilder = updateBuilder.Set("username", req.Username)
	}
	if req.Password != "" {
		updateBuilder = updateBuilder.Set("password_hash", req.Password)
	}
	if req.Fullname != "" {
		updateBuilder = updateBuilder.Set("full_name", req.Fullname)
	}
	updateBuilder = updateBuilder.Set("updated_at", time.Now().Format(time.RFC3339))

	sqlQuery, args, err := updateBuilder.
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("UpdateUsers: SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	result, err := u.db.Exec(sqlQuery, args...)
	if err != nil {
		logger.Logger.Printf("UpdateUsers: Yangilashda xato: %v", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		logger.Logger.Printf("UpdateUsers: Foydalanuvchi topilmadi yoki yangilanmadi: id=%s", req.ID)
		return &users.UpdateRes{Message: "User not found or no changes made"}, nil
	}

	logger.Logger.Printf("UpdateUsers muvaffaqiyatli yakunlandi: id=%s", req.ID)
	return &users.UpdateRes{Message: "User updated successfully"}, nil
}

func (u *UserPostgres) DeleteUsers(req users.UsersbyId) (*users.UpdateRes, error) {
	logger.Logger.Printf("DeleteUsers boshlandi: id=%s", req.ID)

	sqlQuery, args, err := squirrel.
		Delete("users").
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("DeleteUsers: SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	result, err := u.db.Exec(sqlQuery, args...)
	if err != nil {
		logger.Logger.Printf("DeleteUsers: O`chirishda xato: %v", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		logger.Logger.Printf("DeleteUsers: Foydalanuvchi topilmadi: id=%s", req.ID)
		return nil, fmt.Errorf("users not found")
	}

	logger.Logger.Printf("DeleteUsers muvaffaqiyatli yakunlandi: id=%s", req.ID)
	return &users.UpdateRes{Message: "User deleted successfully"}, nil
}