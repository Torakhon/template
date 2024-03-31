package storage

import (
	"comment_service/strorage/postgres"
	"comment_service/strorage/repo"
	"github.com/jmoiron/sqlx"
)

// IStorage ...
type IStorage interface {
	Comment() repo.CommentStorageI
}

type Pg struct {
	db          *sqlx.DB
	commentRepo repo.CommentStorageI
}

// NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *Pg {
	return &Pg{
		db:          db,
		commentRepo: postgres.NewCommentRepo(db),
	}
}

func (s Pg) Comment() repo.CommentStorageI {
	return s.commentRepo
}
