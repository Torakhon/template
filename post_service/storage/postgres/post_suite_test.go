package postgres

//
//import (
//	"github.com/stretchr/testify/suite"
//	"post_service/config"
//	pb "post_service/genproto/post"
//	"post_service/pkg/db"
//	"post_service/storage/repo"
//	"testing"
//)
//
//type PostRepositoryTestSite struct {
//	suite.Suite
//	CleanUpFunc func()
//	Repository  repo.PostStorageI
//}
//
//func (s *PostRepositoryTestSite) SetupSuite() {
//	pgPool, cleanUp := db.ConnectDBForSuite(config.Load())
//	s.Repository = NewPostRepo(pgPool)
//	s.CleanUpFunc = cleanUp
//}
//
//func (s *PostRepositoryTestSite) TestUserCRUD() {
//	post := &pb.CreateRequest{
//		Title:    "asasas",
//		ImageUrl: "sasas",
//		UserId:   "805209d0-5ad8-4187-ad3c-26f1ced97784",
//	}
//
//	reqPost, err := s.Repository.Create(post)
//	s.Suite.NotNil(reqPost)
//	s.Suite.NoError(err)
//
//	getPost, err := s.Repository.GetPost(reqPost.Id)
//	s.Suite.NotNil(getPost)
//	s.Suite.NoError(err)
//}
//
//func (s *PostRepositoryTestSite) TearDownSuite() {
//	s.CleanUpFunc()
//}
//
//func TestUserRepositoryTestSuite(t *testing.T) {
//	suite.Run(t, new(PostRepositoryTestSite))
//}
