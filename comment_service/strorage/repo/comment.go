package repo

import (
	pb "comment_service/genproto/comment"
	"context"
)

// CommentStorageI  ...
type CommentStorageI interface {
	Create(ctx context.Context, req *pb.CreateReq) (*pb.Comment, error)
	GetCommentsByPostId(ctx context.Context, req *pb.GetByPostIdReq) (*pb.GetByIdCommentsRes, error)
	GetCommentsByOwnerId(ctx context.Context, req *pb.GetByOwnerIdReq) (*pb.GetByIdCommentsRes, error)
	UpdateComment(ctx context.Context, req *pb.UpdateCommentReq) (*pb.Comment, error)
	DeleteComment(ctx context.Context, req *pb.DeleteCommentReq) (*pb.DeleteRes, error)
	CommentClickLike(ctx context.Context, req *pb.ClickReq) (*pb.CommentLike, error)
}
