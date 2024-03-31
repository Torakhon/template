package storage

import (
	"github.com/jmoiron/sqlx"
	"post_service/storage/postgres"
	"post_service/storage/repo"
)

// IStorage ...
type IStorage interface {
	Post() repo.PostStorageI
}

type Pg struct {
	db       *sqlx.DB
	postRepo repo.PostStorageI
}

// NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *Pg {
	return &Pg{
		db:       db,
		postRepo: postgres.NewPostRepo(db),
	}
}

func (s Pg) Post() repo.PostStorageI {
	return s.postRepo
}
