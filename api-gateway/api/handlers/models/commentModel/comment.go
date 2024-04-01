package commentModel

type Comments struct {
	Comments []Comment `json:"comments"`
}

type Comment struct {
	CommentID string ` json:"comment_id"`
	PostID    string ` json:"post_id"`
	UserID    string ` json:"user_id"`
	Content   string ` json:"content"`
	Likes     int64  ` json:"likes"`
}

type CreateReq struct {
	CommentID string ` json:"-"`
	PostID    string ` json:"post_id"`
	UserID    string ` json:"-"`
	Content   string ` json:"content"`
}

type UpdateComment struct {
	CommentID  string ` json:"comment_id"`
	UserID     string ` json:"-"`
	NewContent string ` json:"new_content"`
}

type Status struct {
	Status bool `json:"status"`
}
