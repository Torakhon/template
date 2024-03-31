package postgres

import (
	"comment_service/config"
	pb "comment_service/genproto/comment"
	"comment_service/pkg/db"
	"comment_service/strorage/repo"
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CommentRepositoryTestSite struct {
	suite.Suite
	CleanUpFunc func()
	Repository  repo.CommentStorageI
}

func (s *CommentRepositoryTestSite) SetupSuite() {
	pgPool, cleanUp := db.ConnectDBForSuite(config.Load())
	s.Repository = NewCommentRepo(pgPool)
	s.CleanUpFunc = cleanUp
}

func (s *CommentRepositoryTestSite) TestCommentCRUD() {
	comment := &pb.CreateReq{
		CommentId: "b3560590-6b18-4dfb-9605-906b183dfbda",
		PostId:    "0b9555c0-9dc1-4685-9555-c09dc1468532",
		UserId:    "3c6b1db1-d869-4bf0-ab1d-b1d8693bf09a",
		Content:   "Test content 1",
	}

	comment2 := &pb.CreateReq{
		CommentId: "0607835b-21f6-4797-8783-5b21f6d79705",
		PostId:    "0b9555c0-9dc1-4685-9555-c09dc1468532",
		UserId:    "3c6b1db1-d869-4bf0-ab1d-b1d8693bf09a",
		Content:   "Test content 2",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(time.Second*4))
	defer cancel()

	// Testing Create method
	com, err := s.Repository.Create(ctx, comment)

	s.Suite.NotNil(com)
	s.Suite.NoError(err)
	s.Suite.Equal(comment.CommentId, com.CommentId)
	s.Suite.Equal(comment.PostId, com.PostId)
	s.Suite.Equal(comment.UserId, com.UserId)
	s.Suite.Equal(comment.Content, com.Content)

	// Testing Create method
	com2, err := s.Repository.Create(ctx, comment2)

	s.Suite.NotNil(com2)
	s.Suite.NoError(err)
	s.Suite.Equal(comment2.CommentId, com2.CommentId)
	s.Suite.Equal(comment2.PostId, com2.PostId)
	s.Suite.Equal(comment2.UserId, com2.UserId)
	s.Suite.Equal(comment2.Content, com2.Content)

	// Testing Get Comment By Post ID method
	comments, err := s.Repository.GetCommentsByPostId(ctx, &pb.GetByPostIdReq{
		PostId: com.PostId,
		Limit:  10,
		Page:   1,
	})
	s.Suite.NotNil(comments)
	s.Suite.NoError(err)

	// Testing Get Comment By Owner ID method
	comments, err = s.Repository.GetCommentsByOwnerId(ctx, &pb.GetByOwnerIdReq{
		OwnerId: com.UserId,
		Limit:   10,
		Page:    1,
	})
	s.Suite.NotNil(comments)
	s.Suite.NoError(err)

	// Testing Update Comment method
	com, err = s.Repository.UpdateComment(ctx, &pb.UpdateCommentReq{
		CommentId:  com2.CommentId,
		UserId:     com2.UserId,
		NewContent: "Update content",
	})
	s.Suite.NotNil(com)
	s.Suite.NoError(err)
	s.Suite.Equal(com2.CommentId, com.CommentId)
	s.Suite.Equal(com2.PostId, com.PostId)
	s.Suite.Equal(com2.UserId, com.UserId)
	s.Suite.Equal("Update content", com.Content)

	// Testing Comment Click Like method
	CommentClickLike, err := s.Repository.CommentClickLike(ctx, &pb.ClickReq{
		CommentId: com.CommentId,
		UserId:    com.UserId,
	})
	s.Suite.Equal(CommentClickLike.Like, true)
	s.Suite.NotNil(CommentClickLike)
	s.Suite.NoError(err)

	// Testing Comment remove the like method
	CommentClickLike, err = s.Repository.CommentClickLike(ctx, &pb.ClickReq{
		CommentId: com.CommentId,
		UserId:    com.UserId,
	})
	s.Suite.Equal(CommentClickLike.Like, false)
	s.Suite.NotNil(CommentClickLike)
	s.Suite.NoError(err)

	// Testing Comment delete method
	deleteRes, err := s.Repository.DeleteComment(ctx, &pb.DeleteCommentReq{
		CommentId: com.CommentId,
	})
	s.Suite.Equal(deleteRes.Status, true)
	s.Suite.NoError(err)

	deleteRes, err = s.Repository.DeleteComment(ctx, &pb.DeleteCommentReq{
		CommentId: com2.CommentId,
	})
	s.Suite.Equal(deleteRes.Status, true)
	s.Suite.NoError(err)

}

func (s *CommentRepositoryTestSite) TearDownSuite() {
	s.CleanUpFunc()
}

func TestCommentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CommentRepositoryTestSite))
}
