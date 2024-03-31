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
