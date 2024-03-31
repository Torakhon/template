package postgres

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"user-service/config"
	pb "user-service/genproto/user"
	"user-service/pkg/db"
	"user-service/storage/repo"
)

type UserRepositoryTestSite struct {
	suite.Suite
	CleanUpFunc func()
	Repository  repo.UserStorageI
}

func (s *UserRepositoryTestSite) SetupSuite() {
	pgPool, cleanUp := db.ConnectDBForSuite(config.Load())
	s.Repository = NewUserRepo(pgPool)
	s.CleanUpFunc = cleanUp
}

func (s *UserRepositoryTestSite) TestUserCRUD() {
	user := &pb.CreateRequest{
		Name:     "Husan",
		LastName: "Aliyev",
	}

	respUser, err := s.Repository.Create(user)
	s.Suite.NotNil(respUser)
	s.Suite.NoError(err)
	s.Suite.Equal(user.Name, respUser.Name)
	s.Suite.Equal(user.LastName, respUser.LastName)

	GetUser, err := s.Repository.GetUser(respUser.Uuid)
	s.Suite.NotNil(GetUser)
	s.Suite.NoError(err)

	user.Name = "Laziz"
	user.LastName = "Hasanov"

	UpdateUser, err := s.Repository.UpdateUser(&pb.UpRequest{
		Uuid:     GetUser.Uuid,
		Name:     "laziz",
		LastName: "Hasanov",
		UserName: "Ali",
		Email:    user.Email,
		Password: "11223344",
	})
	s.Suite.NotNil(UpdateUser)
	s.Suite.NoError(err)
	s.Suite.Equal(respUser.Name, UpdateUser.Name)
	s.Suite.Equal(respUser.LastName, UpdateUser.LastName)

	GetUsers, err := s.Repository.GetAllusers(1, 10)
	s.Suite.NotNil(GetUsers)
	s.Suite.NoError(err)

	err = s.Repository.DeleteUser(GetUser.Uuid)
	s.Suite.NoError(err)

}

func (s *UserRepositoryTestSite) TearDownSuite() {
	s.CleanUpFunc()
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSite))
}
