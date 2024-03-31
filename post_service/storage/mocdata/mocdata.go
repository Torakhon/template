package mocdata

import (
	"context"
	"github.com/jmoiron/sqlx"
	pb "post_service/genproto/post"
)

type PostRepo struct {
	db *sqlx.DB
}

// MocNewPostRepo ...
func MocNewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (p *PostRepo) Create(ctx context.Context, req *pb.CreateReq) (*pb.Post, error) {
	return &pb.Post{
		Id:        "f3e3a146-72dd-4d64-a3a1-4672dddd64b7",
		Title:     "Test Title",
		Content:   "New content",
		UserId:    "f381e475-71a5-4a27-81e4-7571a58a27a4",
		Category:  "Test category",
		Likes:     10,
		Dislikes:  2,
		CreatedAt: "2003-04-04",
		UpdatedAt: "2005-04-04",
		DeletedAt: "2005-08-09",
		Comments:  nil,
	}, nil
}

func (p *PostRepo) GetPost(ctx context.Context, req *pb.GetReq) (*pb.Post, error) {
	return &pb.Post{
		Id:        "f3e3a146-72dd-4d64-a3a1-4672dddd64b7",
		Title:     "Test Title",
		Content:   "New content",
		UserId:    "f381e475-71a5-4a27-81e4-7571a58a27a4",
		Category:  "Test category",
		Likes:     10,
		Dislikes:  2,
		CreatedAt: "2003-04-04",
		UpdatedAt: "2005-04-04",
		DeletedAt: "2005-08-09",
		Comments:  nil,
	}, nil
}

func (p *PostRepo) SearchPost(ctx context.Context, req *pb.SearchReq) (*pb.PostsRes, error) {
	return &pb.PostsRes{
		Posts: []*pb.Post{
			{
				Id:        "8470dc46-64c3-4978-b0dc-4664c3897837",
				Title:     "First Post",
				Content:   "Content 1",
				UserId:    "7c61d8b4-654g-4683-92c9-b3554g267278",
				Category:  "Technology",
				Likes:     15,
				Dislikes:  3,
				CreatedAt: "2022-01-01",
				UpdatedAt: "2022-01-02",
				DeletedAt: "",
				Comments:  nil,
			},
			{
				Id:        "af61ed23-b841-43f5-a1ed-23b84133f5be",
				Title:     "Second Post",
				Content:   "Content 2",
				UserId:    "7c61d8b4-654g-4683-92c9-b3554g267278",
				Category:  "Travel",
				Likes:     20,
				Dislikes:  1,
				CreatedAt: "2022-02-01",
				UpdatedAt: "2022-02-02",
				DeletedAt: "",
				Comments:  nil,
			},
		},
	}, nil
}

func (p *PostRepo) UpdatePost(ctx context.Context, req *pb.UpdatePostReq) (*pb.Post, error) {
	return &pb.Post{
		Id:        "f3e3a146-72dd-4d64-a3a1-4672dddd64b7",
		Title:     "Test Title",
		Content:   "New content",
		UserId:    "f381e475-71a5-4a27-81e4-7571a58a27a4",
		Category:  "Test category",
		Likes:     10,
		Dislikes:  2,
		CreatedAt: "2003-04-04",
		UpdatedAt: "2005-04-04",
		DeletedAt: "2005-08-09",
		Comments:  nil,
	}, nil
}

func (p *PostRepo) DeletePost(ctx context.Context, req *pb.DeletePostReq) (*pb.DeletePostRes, error) {
	return &pb.DeletePostRes{
		Status: true,
	}, nil
}

func (p *PostRepo) PostClickLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error) {
	return &pb.PostLike{
		Like: true,
	}, nil
}

func (p *PostRepo) PostClickDisLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error) {
	return &pb.PostLike{
		Like: false,
	}, nil
}

func (p *PostRepo) Views(ctx context.Context, req *pb.ViewReq) (*pb.ViewRes, error) {
	return nil, nil
}
