package repo

import (
	"context"
	pb "post_service/genproto/post"
)

// PostStorageI ...
type PostStorageI interface {
	Create(ctx context.Context, req *pb.CreateReq) (*pb.Post, error)
	GetPost(ctx context.Context, req *pb.GetReq) (*pb.Post, error)
	SearchPost(ctx context.Context, req *pb.SearchReq) (*pb.PostsRes, error)
	UpdatePost(ctx context.Context, req *pb.UpdatePostReq) (*pb.Post, error)
	DeletePost(ctx context.Context, req *pb.DeletePostReq) (*pb.DeletePostRes, error)
	PostClickLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error)
	PostClickDisLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error)
	Views(ctx context.Context, req *pb.ViewReq) (*pb.ViewRes, error)
}
