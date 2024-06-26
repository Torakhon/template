package storage

import (
	"github.com/jmoiron/sqlx"
	"user_service/storage/postgres"
	"user_service/storage/repo"
)

// IStorage ...
type IStorage interface {
	User() repo.UserStorageI
}

type Pg struct {
	db       *sqlx.DB
	userRepo repo.UserStorageI
}

// NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *Pg {
	return &Pg{
		db:       db,
		userRepo: postgres.NewUserRepo(db),
	}
}

func (s Pg) User() repo.UserStorageI {
	return s.userRepo
}

//// IStorage ...
//type IStorage interface {
//	User() repo.UserStorageI
//}
//
//type Pg struct {
//	db       *mongo.Client
//	userRepo repo.UserStorageI
//}
//
//// NewStoragePg ...
//func NewStoragePg(db *mongo.Client) *Pg {
//	return &Pg{
//		db:       db,
//		userRepo: mongoDBdatabase.NewUserRepo(db,"v1","users"),
//	}
//}
//
//func (s Pg) User() repo.UserStorageI {
//	return s.userRepo
//}
