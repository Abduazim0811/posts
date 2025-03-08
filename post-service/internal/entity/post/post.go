package post

type Post struct {
    ID        string    `json:"id,omitempty"`
    Username    string    `json:"user_id,omitempty"`
    Title     string   `json:"title,omitempty"`
    Content   string   `json:"content,omitempty"`
    CreatedAt string   `json:"created_at,omitempty"`
    UpdatedAt string   `json:"updated_at,omitempty"`
    Tags      []string `json:"tags,omitempty"`
}

type CreatePostRequest struct {
    Username  string    `json:"user_id,omitempty"`
    Title   string   `json:"title,omitempty"`
    Content string   `json:"content,omitempty"`
    Tags    []string `json:"tags,omitempty"`
}

type GetPostRequest struct {
    ID string `json:"id,omitempty"`
}

type ListPostsRequest struct {
    Username string  `json:"user_id,omitempty"`
    Page   int32 `json:"page,omitempty"`
    Limit  int32 `json:"limit,omitempty"`
}

type UpdatePostRequest struct {
    ID      string    `json:"id,omitempty"`
    Title   string   `json:"title,omitempty"`
    Content string   `json:"content,omitempty"`
    Tags    []string `json:"tags,omitempty"`
}

type DeletePostRequest struct {
    ID string `json:"id,omitempty"`
}

type PostResponse struct {
    ID string `json:"id,omitempty"`
    Message string `json:"message,omitempty"`
}

type GetPostResponse struct {
	ID       string
	Username string 
	Title    string
	Content  string
	Tags     []string
}

type ListPostsResponse struct {
    Posts   []*GetPostResponse `json:"posts,omitempty"`
    Total   int32   `json:"total,omitempty"`
}

type DeletePostResponse struct {
    Message string `json:"message,omitempty"`
}