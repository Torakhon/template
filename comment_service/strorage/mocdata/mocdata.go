package mocdata

import (
	pb "comment_service/genproto/comment"
	"context"
	"github.com/jmoiron/sqlx"
)

type CommentRepo struct {
	db *sqlx.DB
}

// NewCommentRepo ...
func NewCommentRepo(db *sqlx.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, req *pb.CreateReq) (*pb.Comment, error) {
	return &pb.Comment{
		CommentId: "49814fbd-6338-4f97-814f-bd63385f9798",
		PostId:    "41103e64-2ec8-474d-903e-642ec8174db0",
		UserId:    "8190dcbb-fac0-4c43-90dc-bbfac0cc432d",
		Content:   "Test content ",
		Likes:     12,
	}, nil
}

func (r *CommentRepo) GetCommentsByPostId(ctx context.Context, req *pb.GetByPostIdReq) (*pb.GetByIdCommentsRes, error) {
	return &pb.GetByIdCommentsRes{
		Comments: []*pb.Comment{
			{
				CommentId: "b9d3d075-414a-4b1a-93d0-75414abb1a4f",
				PostId:    "5e51b8a2-543f-4572-91b8-a2543f157276",
				UserId:    "7c61d8b4-654g-4683-92c9-b3554g267278",
				Content:   "First comment content",
				Likes:     10,
			},
			{
				CommentId: "7604480f-9df3-4e0d-8448-0f9df3ae0dcd",
				PostId:    "5e51b8a2-543f-4572-91b8-a2543f157276",
				UserId:    "8d72e9c5-654g-4683-92c9-b3554g267279",
				Content:   "Second comment content",
				Likes:     5,
			},
		},
	}, nil
}

func (r *CommentRepo) GetCommentsByOwnerId(ctx context.Context, req *pb.GetByOwnerIdReq) (*pb.GetByIdCommentsRes, error) {
	return &pb.GetByIdCommentsRes{
		Comments: []*pb.Comment{
			{
				CommentId: "b9d3d075-414a-4b1a-93d0-75414abb1a4f",
				PostId:    "5e51b8a2-543f-4572-91b8-a2543f157276",
				UserId:    "7c61d8b4-654g-4683-92c9-b3554g267278",
				Content:   "First comment content",
				Likes:     10,
			},
			{
				CommentId: "7604480f-9df3-4e0d-8448-0f9df3ae0dcd",
				PostId:    "5e51b8a2-543f-4572-91b8-a2543f157276",
				UserId:    "8d72e9c5-654g-4683-92c9-b3554g267279",
				Content:   "Second comment content",
				Likes:     5,
			},
		},
	}, nil
}

func (r *CommentRepo) UpdateComment(ctx context.Context, req *pb.UpdateCommentReq) (*pb.Comment, error) {
	return &pb.Comment{
		CommentId: "49814fbd-6338-4f97-814f-bd63385f9798",
		PostId:    "41103e64-2ec8-474d-903e-642ec8174db0",
		UserId:    "8190dcbb-fac0-4c43-90dc-bbfac0cc432d",
		Content:   "Test content Update",
		Likes:     20,
	}, nil
}

func (r *CommentRepo) DeleteComment(ctx context.Context, req *pb.DeleteCommentReq) (*pb.DeleteRes, error) {
	return &pb.DeleteRes{Status: true}, nil
}

func (r *CommentRepo) CommentClickLike(ctx context.Context, req *pb.ClickReq) (*pb.CommentLike, error) {
	return &pb.CommentLike{Like: true}, nil
}
