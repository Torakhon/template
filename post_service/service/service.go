package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	pb "post_service/genproto/post"
	l "post_service/pkg/logger"
	grpcClient "post_service/service/grpc_client"
	"post_service/storage"
)

// PostService ...
type PostService struct {
	storage    storage.IStorage
	logger     l.Logger
	grpcClient grpcClient.IserviveManager
}

// NewPostService ...
func NewPostService(db *sqlx.DB, log l.Logger, grpcClient grpcClient.IserviveManager) *PostService {
	return &PostService{
		storage:    storage.NewStoragePg(db),
		logger:     log,
		grpcClient: grpcClient,
	}
}

func (p *PostService) Create(ctx context.Context, req *pb.CreateReq) (*pb.Post, error) {
	return p.storage.Post().Create(ctx, req)
}

func (p *PostService) GetPost(ctx context.Context, req *pb.GetReq) (*pb.Post, error) {
	return p.storage.Post().GetPost(ctx, req)
}

func (p *PostService) SearchPost(ctx context.Context, req *pb.SearchReq) (*pb.PostsRes, error) {
	return p.storage.Post().SearchPost(ctx, req)
}

func (p *PostService) UpdatePost(ctx context.Context, req *pb.UpdatePostReq) (*pb.Post, error) {
	return p.storage.Post().UpdatePost(ctx, req)
}

func (p *PostService) DeletePost(ctx context.Context, req *pb.DeletePostReq) (*pb.DeletePostRes, error) {
	return p.storage.Post().DeletePost(ctx, req)
}

func (p *PostService) PostClickLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error) {
	return p.storage.Post().PostClickLike(ctx, req)
}

func (p *PostService) PostClickDisLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error) {
	return p.storage.Post().PostClickDisLike(ctx, req)
}

func (p *PostService) Views(ctx context.Context, req *pb.ViewReq) (*pb.ViewRes, error) {
	return p.storage.Post().Views(ctx, req)
}

//func (p *PostService) Create(ctx context.Context, req *pb.CreateReq) (*pb.Post, error) {
//	return p.storage.Post().Create(ctx, req)
//}
//
//func (p *PostService) GetPost(ctx context.Context, req *pb.GetReq) (*pb.Post, error) {
//	postModel, err := p.storage.Post().GetPost(ctx, req)
//	if err != nil {
//		return nil, err
//	}
//
//	_, err = p.storage.Post().Views(ctx, &pb.ViewReq{
//		PostId: postModel.Id,
//		UserId: postModel.UserId,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	comments, err := p.grpcClient.CommentService().GetCommentsByPostId(ctx, &commentPb.GetByPostIdReq{
//		PostId: postModel.Id,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	for _, comment := range comments.Comments {
//		comm := &pb.Comment{
//			CommentId: comment.CommentId,
//			PostId:    comment.PostId,
//			UserId:    comment.UserId,
//			Content:   comment.Content,
//			Likes:     comment.Likes,
//			Dislikes:  comment.Dislikes,
//		}
//		postModel.Comments = append(postModel.Comments, comm)
//	}
//
//	return postModel, nil
//}
//
//func (p *PostService) SearchPost(ctx context.Context, req *pb.SearchReq) (*pb.PostsRes, error) {
//	posts, err := p.storage.Post().SearchPost(ctx, req)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, postModel := range posts.Posts {
//		comments, err := p.grpcClient.CommentService().GetCommentsByPostId(ctx, &commentPb.GetByPostIdReq{
//			PostId: postModel.Id,
//		})
//		if err != nil {
//			return nil, err
//		}
//
//		_, err = p.storage.Post().Views(ctx, &pb.ViewReq{
//			PostId: postModel.Id,
//			UserId: postModel.UserId,
//		})
//		if err != nil {
//			return nil, err
//		}
//
//		var com []*pb.Comment
//		for _, comment := range comments.Comments {
//			var comm pb.Comment
//			comm.PostId = comment.PostId
//			comm.UserId = comment.UserId
//			comm.Content = comment.Content
//			comm.CommentId = comment.CommentId
//			com = append(com, &comm)
//		}
//
//	}
//
//	return posts, nil
//}
//
//func (p *PostService) UpdatePost(ctx context.Context, req *pb.UpdatePostReq) (*pb.Post, error) {
//	return p.storage.Post().UpdatePost(ctx, req)
//}
//
//func (p *PostService) DeletePost(ctx context.Context, req *pb.DeletePostReq) (*pb.DeletePostRes, error) {
//	return p.storage.Post().DeletePost(ctx, req)
//}
//
//func (p *PostService) PostClickLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error) {
//	return p.storage.Post().PostClickLike(ctx, req)
//}
//
//func (p *PostService) PostClickDisLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error) {
//	return p.storage.Post().PostClickDisLike(ctx, req)
//}
//
//func (p *PostService) Views(ctx context.Context, req *pb.ViewReq) (*pb.ViewRes, error) {
//	return p.storage.Post().Views(ctx, req)
//}
