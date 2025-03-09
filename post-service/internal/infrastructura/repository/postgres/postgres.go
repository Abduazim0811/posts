package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"post-service/internal/entity/post"
	"post-service/internal/infrastructura/repository"
	"post-service/internal/logger"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

type PostPostgres struct {
	db *sql.DB
}

func NewPostPostgres(db *sql.DB) repository.PostRepository {
	return &PostPostgres{db: db}
}

func (p *PostPostgres) CreatePost(req post.CreatePostRequest) (*post.PostResponse, error) {
	logger.Logger.Printf("CreatePost boshlandi: username=%s, title=%s", req.Username, req.Title)

	sqlQuery, args, err := squirrel.
		Insert("posts").
		Columns("username", "title", "content", "tags").
		Values(req.Username, req.Title, req.Content, pq.Array(req.Tags)).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	var res post.PostResponse
	err = p.db.QueryRow(sqlQuery, args...).Scan(&res.ID)
	if err != nil {
		logger.Logger.Printf("Postni saqlashda xato: %v", err)
		return nil, err
	}

	res.Message = "Post successfully created"
	logger.Logger.Printf("Post muvaffaqiyatli yaratildi: id=%s", res.ID)
	return &res, nil
}

func (p *PostPostgres) GetPost(req post.GetPostRequest) (*post.GetPostResponse, error) {
	logger.Logger.Printf("GetPost boshlandi: id=%s", req.ID)

	sqlQuery, args, err := squirrel.
		Select("id", "username", "title", "content", "tags").
		From("posts").
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	var id, username, title, content string
	var tags []string
	err = p.db.QueryRow(sqlQuery, args...).Scan(&id, &username, &title, &content, pq.Array(&tags))
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Printf("Post topilmadi: id=%s", req.ID)
			return nil, nil
		}
		logger.Logger.Printf("Postni olishda xato: %v", err)
		return nil, err
	}

	logger.Logger.Printf("Post muvaffaqiyatli olindi: id=%s", id)
	return &post.GetPostResponse{
		ID:       id,
		Username: username,
		Title:    title,
		Content:  content,
		Tags:     tags,
	}, nil
}

func (p *PostPostgres) ListPosts(req post.ListPostsRequest) (*post.ListPostsResponse, error) {
	logger.Logger.Printf("ListPosts boshlandi: page=%d, limit=%d", req.Page, req.Limit)

	offset := (req.Page - 1) * req.Limit

	sqlQuery, args, err := squirrel.
		Select("id", "username", "title", "content", "tags").
		From("posts").
		OrderBy("created_at DESC").
		Limit(uint64(req.Limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	rows, err := p.db.Query(sqlQuery, args...)
	if err != nil {
		logger.Logger.Printf("Postlarni olishda xato: %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []*post.GetPostResponse
	for rows.Next() {
		var id, username, title, content string
		var tags []string
		err := rows.Scan(&id, &username, &title, &content, &tags)
		if err != nil {
			logger.Logger.Printf("Postlarni skanerlashda xato: %v", err)
			return nil, err
		}

		posts = append(posts, &post.GetPostResponse{
			ID:       id,
			Username: username,
			Title:    title,
			Content:  content,
			Tags:     tags,
		})
	}

	// Umumiy postlar sonini olish
	countQuery, countArgs, err := squirrel.
		Select("COUNT(*)").
		From("posts").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("Count SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	var total int32
	err = p.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		logger.Logger.Printf("Umumiy postlar sonini olishda xato: %v", err)
		return nil, err
	}

	logger.Logger.Printf("Postlar muvaffaqiyatli olindi: total=%d", total)
	return &post.ListPostsResponse{
		Posts: posts,
		Total: total,
	}, nil
}


func (p *PostPostgres) UpdatePost(req post.UpdatePostRequest) (*post.PostResponse, error) {
	logger.Logger.Printf("UpdatePost boshlandi: id=%s", req.ID)

	updateBuilder := squirrel.Update("posts").Where(squirrel.Eq{"id": req.ID})

	if req.Title != "" {
		updateBuilder = updateBuilder.Set("title", req.Title)
	}
	if req.Content != "" {
		updateBuilder = updateBuilder.Set("content", req.Content)
	}
	if len(req.Tags) > 0 {
		updateBuilder = updateBuilder.Set("tags", pq.Array(req.Tags))

	}

	updateBuilder = updateBuilder.Set("updated_at", time.Now().Format(time.RFC3339))

	sqlQuery, args, err := updateBuilder.
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("UpdatePost: SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	result, err := p.db.Exec(sqlQuery, args...)
	if err != nil {
		logger.Logger.Printf("UpdatePost: Yangilashda xato: %v", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		logger.Logger.Printf("UpdatePost: Post topilmadi yoki yangilanmadi: id=%s", req.ID)
		return &post.PostResponse{Message: "Post not found or no changes made"}, nil
	}

	updateBuilder = squirrel.Update("posts").Where(squirrel.Eq{"id": req.ID})
	sqlQuery, args, err = updateBuilder.
		Suffix("RETURNING id, username, title, content, tags").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("UpdatePost: RETURNING so`rov yaratishda xato: %v", err)
		return nil, err
	}

	var resp post.PostResponse
	var tagsInterface interface{}
	err = p.db.QueryRow(sqlQuery, args...).Scan(&resp.ID, &resp.Username, &resp.Title, &resp.Content, &tagsInterface)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Printf("UpdatePost: Post topilmadi: id=%s", req.ID)
			return &post.PostResponse{Message: "Post not found"}, nil
		}
		logger.Logger.Printf("UpdatePost: Ma'lumot olishda xato: %v", err)
		return nil, err
	}

	switch v := tagsInterface.(type) {
	case []string:
		resp.Tags = v
	case []byte:
		if err := json.Unmarshal(v, &resp.Tags); err != nil {
			logger.Logger.Printf("UpdatePost: Tagsni JSON dan olishda xato: %v", err)
			return nil, err
		}
	default:
		logger.Logger.Printf("UpdatePost: Noto'g'ri tags turi: %T", v)
		return nil, fmt.Errorf("unsupported tags type")
	}

	resp.Message = "Post updated successfully"
	logger.Logger.Printf("UpdatePost muvaffaqiyatli yakunlandi: id=%s", resp.ID)
	return &resp, nil
}

func (p *PostPostgres) DeletePost(req post.DeletePostRequest) (*post.DeletePostResponse, error) {
	logger.Logger.Printf("DeletePost boshlandi: id=%s", req.ID)

	sqlQuery, args, err := squirrel.
		Delete("posts").
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.Logger.Printf("SQL so`rov yaratishda xato: %v", err)
		return nil, err
	}

	result, err := p.db.Exec(sqlQuery, args...)
	if err != nil {
		logger.Logger.Printf("Postni o`chirishda xato: %v", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Logger.Printf("RowsAffected olishda xato: %v", err)
		return nil, err
	}
	if rowsAffected == 0 {
		logger.Logger.Printf("Post topilmadi: id=%s", req.ID)
		return &post.DeletePostResponse{
			Message: "Post not found",
		}, nil
	}

	logger.Logger.Printf("Post muvaffaqiyatli o`chirildi: id=%s", req.ID)
	return &post.DeletePostResponse{
		Message: "Post deleted successfully",
	}, nil
}
