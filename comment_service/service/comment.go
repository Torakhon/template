package service

import (
	commentPb "comment_service/genproto/comment"
	l "comment_service/pkg/logger"
	grpcClient "comment_service/service/grpc_client"
	"comment_service/strorage"
	"context"
	"github.com/jmoiron/sqlx"
)

// CommentService ...
type CommentService struct {
	storage    storage.IStorage
	logger     l.Logger
	grpcClient grpcClient.IserviveManager
}

// NewCommentService ...
func NewCommentService(db *sqlx.DB, log l.Logger, grpcClient grpcClient.IserviveManager) *CommentService {
	return &CommentService{
		storage:    storage.NewStoragePg(db),
		logger:     log,
		grpcClient: grpcClient,
	}
}

//         mongo DB

//// UserService ...
//type UserService struct {
//	storage    storage.IStorage
//	logger     l.Logger
//	grpcClient grpcClient.IServiceManager
//}
//
//// NewUserService ...
//func NewUserService(db *mongo.Client, log l.Logger, grpcClient grpcClient.IServiceManager) *UserService {
//	return &UserService{
//		storage:    storage.NewStoragePg(db),
//		logger:     log,
//		grpcClient: grpcClient,
//	}
//}

func (r *CommentService) Create(ctx context.Context, req *commentPb.CreateReq) (*commentPb.Comment, error) {
	return r.storage.Comment().Create(ctx, req)
}

func (r *CommentService) GetCommentsByPostId(ctx context.Context, req *commentPb.GetByPostIdReq) (*commentPb.GetByIdCommentsRes, error) {
	return r.storage.Comment().GetCommentsByPostId(ctx, req)
}

func (r *CommentService) GetCommentsByOwnerId(ctx context.Context, req *commentPb.GetByOwnerIdReq) (*commentPb.GetByIdCommentsRes, error) {
	return r.storage.Comment().GetCommentsByOwnerId(ctx, req)
}

func (r *CommentService) UpdateComment(ctx context.Context, req *commentPb.UpdateCommentReq) (*commentPb.Comment, error) {
	return r.storage.Comment().UpdateComment(ctx, req)
}

func (r *CommentService) DeleteComment(ctx context.Context, req *commentPb.DeleteCommentReq) (*commentPb.DeleteRes, error) {
	return r.storage.Comment().DeleteComment(ctx, req)
}

func (r *CommentService) CommentClickLike(ctx context.Context, req *commentPb.ClickReq) (*commentPb.CommentLike, error) {
	return r.storage.Comment().CommentClickLike(ctx, req)
}
