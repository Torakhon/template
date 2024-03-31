package postModel

import (
	"api-gateway/api/handlers/models/commentModel"
)

type Posts struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	Id        string                 `json:"id"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	UserId    string                 ` json:"user_id"`
	Category  string                 ` json:"category"`
	Likes     int64                  ` json:"likes"`
	Dislikes  int64                  ` json:"dislikes"`
	Views     int64                  ` json:"views"`
	CreatedAt string                 ` json:"created_at"`
	UpdatedAt string                 ` json:"updated_at"`
	DeletedAt string                 ` json:"deleted_at"`
	Comments  []commentModel.Comment `json:"comments"`
}

type GetUserWithPosts struct {
	UserId string `json:"user_id"`
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
}

type CreateReq struct {
	ID       string `json:"-"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	UserID   string `json:"-"`
	Category string `json:"category"`
}

type SearchReq struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Page  int32  `json:"page"`
	Limit int32  ` json:"limit"`
}

type UpdatePostReq struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type Status struct {
	Status bool `json:"'status'"`
}
